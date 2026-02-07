package entity

import "time"

type AuthCode struct {
	Length int
	TTL    time.Duration
}
