package model

type User struct {
	Id       int64  `gorm:"primaryKey"`
	Username string `gorm:"type:varchar(33);unique;not null"`
	Password string `gorm:"type:varchar(33);not null"`
}
