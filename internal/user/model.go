package user

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email" binding:"required,email"`
	CreatedAt    string `json:"created_at"`
	PasswordHash string `json:"-"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type AuthResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

func UserToUserRespose(user *User) *UserResponse {
	return &UserResponse{ID: user.ID, Email: user.Email, CreatedAt: user.CreatedAt}
}
