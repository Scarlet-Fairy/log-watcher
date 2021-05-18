package service

import "time"

type Log struct {
	Id        string
	Timestamp time.Time
	Message   string
}
