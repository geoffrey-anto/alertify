package models

// Field names should start with an uppercase letter
type User struct {
	ID    int    `json:"id" xml:"id" form:"id"`
	Email string `json:"email" xml:"email" form:"email"`
	Pass  string `json:"pass" xml:"pass" form:"pass"`
	Fname string `json:"fname" xml:"fname" form:"fname"`
	Lname string `json:"lname" xml:"lname" form:"lname"`
}

type Reminder struct {
	ID               int    `json:"id" xml:"id" form:"id"`
	UserID           int    `json:"user_id" xml:"user_id" form:"user_id"`
	Name             string `json:"name" xml:"name" form:"name"`
	Status           string `json:"status" xml:"status" form:"status"`
	Description      string `json:"description" xml:"description" form:"description"`
	Category         string `json:"category" xml:"category" form:"category"`
	CreatedAt        string `json:"created_at" xml:"created_at" form:"created_at"`
	UpdatedAt        string `json:"updated_at" xml:"updated_at" form:"updated_at"`
	ReminderInterval string `json:"reminder_interval" xml:"reminder_interval" form:"reminder_interval"`
	ReminderEnd      string `json:"reminder_end" xml:"reminder_end" form:"reminder_end"`
}
