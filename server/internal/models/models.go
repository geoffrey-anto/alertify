package models

// Field names should start with an uppercase letter
type User struct {
	ID    int    `json:"id" xml:"id" form:"id"`
	Email string `json:"email" xml:"email" form:"email"`
	Pass  string `json:"pass" xml:"pass" form:"pass"`
	Fname string `json:"fname" xml:"fname" form:"fname"`
	Lname string `json:"lname" xml:"lname" form:"lname"`
}

type UserLogin struct {
	Email string `json:"email" xml:"email" form:"email"`
	Pass  string `json:"pass" xml:"pass" form:"pass"`
}
