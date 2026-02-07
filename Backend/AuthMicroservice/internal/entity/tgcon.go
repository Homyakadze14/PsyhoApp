package entity

import "time"

type TgConnection struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TgUserID  int       `json:"tg_user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
