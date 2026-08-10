package auth

type LoginInput struct {
	Username string `json:"Username" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

