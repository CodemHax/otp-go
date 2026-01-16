package models

type EmailModel struct {
	Email string `json:"email"`
}

type VerifyModel struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
