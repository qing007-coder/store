package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       string `gorm:"primaryKey"`
	Account  string
	Password string
	Email    string
	Status   bool
}

type UserRole struct {
	ID   string `gorm:"primaryKey"`
	Type string
	V1   string
	V2   string
	V3   string
}

type ArticleJudgeRecord struct {
	gorm.Model
	ArticleID       string
	IsJudge         bool
	Result          bool
	AdministratorID string
	JudgeTime       time.Time
}
