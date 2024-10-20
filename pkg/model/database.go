package model

type User struct {
	ID       string `gorm:"primaryKey"`
	Account  string
	Password string
	Email    string
}

type UserRole struct {
	ID   string `gorm:"primaryKey"`
	Type string
	V1   string
	V2   string
	V3   string
}

type Client struct {
	ID     string `gorm:"primaryKey"`
	Secret string
	Domain string
	Scope  string
	UserID string
}
