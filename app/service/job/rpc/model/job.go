package model

import "time"

type Job struct {
	JobID          int64     `gorm:"primaryKey;column:job_id"`
	UserID         int64     `gorm:"index;column:user_id"`
	Status         int8      `gorm:"not null"`
	FileName       string    `gorm:"type:varchar(50);not null;column:file_name"`
	JobName        string    `gorm:"type:varchar(50);not null;column:job_name"`
	ConsequentFile string    `gorm:"type:varchar(50);column:consequent_file"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime"`
}
