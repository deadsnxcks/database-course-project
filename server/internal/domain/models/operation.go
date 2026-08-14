package models

import "time"

type Operation struct {
	ID 			int64
	Title 		string
	CreatedAt	time.Time
}