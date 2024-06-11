package entity

import (
	"time"
)

type Result struct {
	TopCookies []string
}

type FileRow struct {
	Cookie    string
	TimeStamp time.Time
}
