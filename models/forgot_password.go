package models

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}
