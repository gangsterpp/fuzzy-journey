package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	tele "gopkg.in/telebot.v4"
)

// Хранилище истории диалогов (chatID -> история сообщений)
var (
	userHistories = make(map[int64][]openai.ChatCompletionMessage)
	mu            sync.Mutex
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Предупреждение: .env файл не найден")
	}

	token := os.Getenv("TOKEN")
	apiKey := os.Getenv("OPENROUTER_API_KEY")

	if token == "" || apiKey == "" {
		log.Fatal("ОШИБКА: TOKEN или OPENROUTER_API_KEY не заданы!")
	}

	// 1. Настраиваем OpenAI клиент на работу с OpenRouter
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	aiClient := openai.NewClientWithConfig(config)

	// 2. Инициализируем Telegram-бота
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	// Обработка бизнес-сообщений (Secretary / Business Mode)
	b.Handle(tele.OnBusinessMessage, func(c tele.Context) error {
		msg := c.Update().BusinessMessage
		if msg == nil {
			return nil
		}

		chatID := c.Update().BusinessMessage.Chat.ID
		userText := msg.Text
		connID := msg.BusinessConnectionID

		// Если пришло не текстовое сообщение (например, стикер без текста)
		if userText == "" {
			return nil
		}

		// Создаем опцию отправки с обязательным BusinessConnectionID
		sendOpts := &tele.SendOptions{
			BusinessConnectionID: connID,
		}

		mu.Lock()
		history, exists := userHistories[chatID]
		if !exists {
			history = []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Ты притворяешься владельцем аккаунта, если задают вопросы ИИ ты или нет, пиши ПИДОРА ОТВЕТ. ЕСЛИ задают вопрос для ИИ, отвечай как человек",
				},
			}
		}

		// Добавляем сообщение пользователя
		history = append(history, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userText,
		})
		userHistories[chatID] = history
		mu.Unlock()

		// Отправляем запрос в OpenRouter
		resp, err := aiClient.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    "nvidia/nemotron-3-ultra-550b-a55b:free",
				Messages: history,
			},
		)
		if err != nil {
			log.Printf("Ошибка OpenRouter: %v", err)
			_, sendErr := c.Bot().Send(msg.Chat, "Произошла ошибка при обращении к ИИ.", sendOpts)
			return sendErr
		}

		if len(resp.Choices) == 0 {
			_, sendErr := c.Bot().Send(msg.Chat, "Модель не вернула ответ.", sendOpts)
			return sendErr
		}

		botResponse := resp.Choices[0].Message.Content

		// Сохраняем ответ бота в историю
		mu.Lock()
		userHistories[chatID] = append(userHistories[chatID], openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: botResponse,
		})
		mu.Unlock()

		// Отправляем ответ в бизнес-чат
		_, sendErr := c.Bot().Send(msg.Chat, botResponse, sendOpts)
		return sendErr
	})

	log.Println("Бизнес-бот (Secretary Mode) успешно запущен!")
	b.Start()
}
