package auth

type AuthModel struct {
	Email    string `json:"email"`
	Password string `json:"password"  binding:"required,min=8"`
}
