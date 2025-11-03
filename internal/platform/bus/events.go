package bus

import "time"

type Topic string

const VerifyEmailTopic Topic = "verify_email"

type VerifyEmailEvent struct {
	Email    string    `json:"email"`
	ExpireAt time.Time `json:"expireAt"`
	Token    string    `json:"token"`
}

const ResetPasswordTopic Topic = "reset_password"

type ResetPasswordEvent struct {
	Email    string    `json:"email"`
	ExpireAt time.Time `json:"expireAt"`
	Token    string    `json:"token"`
}
