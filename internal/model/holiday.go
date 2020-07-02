package model

import "time"

type Holiday struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	HolidayType string    `db:"type" json:"type"`
	HolidayDate time.Time `db:"date" json:"date"`
}

type TimeUnit string

const (
	Epoch       TimeUnit = "Epoch"
	Month       TimeUnit = "Month"
	DateOfMonth TimeUnit = "DateOfMonth"
	Year        TimeUnit = "Year"
)