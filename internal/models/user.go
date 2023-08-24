package models

type User struct {
	Id        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Password  string `json:"pass"`
	ImagePath string `json:"image"`
}
