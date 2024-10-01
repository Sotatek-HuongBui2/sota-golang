package utils

import (
	"time"
	"vtcanteen/models"
)

type UserMetadatas struct {
	IsVerified       bool       `json:"is_verified"`
	VerificationCode string     `json:"verification_code"`
	ExpireAt         *time.Time `json:"expire_at"`
}

type RegisterResponse struct {
	User        *models.Users `json:"user"`
	MessageMail interface{}   `json:"message_mail"`
}
