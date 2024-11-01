package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           string    `gorm:"primaryKey"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Account      string    `json:"account"`
	Nickname     string    `json:"nickname"`
	Introduction string    `json:"introduction"`
	Gender       string    `json:"gender"`
	Sign         string    `json:"sign"`
	Password     string    `json:"password"`
	Email        string    `json:"email"`
	Avatar       string    `json:"avatar"`
}

type ReceiverAddress struct {
	gorm.Model
	DetailedAddress string `json:"detailed_address"`
	ReceiverName    string `json:"receiver_name"`
	PhoneNumber     string `json:"phone_number"`
	Label           string `json:"label"`
	UserID          string `json:"user_id"`
}

type Favourites struct {
	gorm.Model
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"` //
	Category string `json:"category"`
}

type Footprint struct {
	gorm.Model
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
	Category string `json:"category"`
}

type Follow struct {
	gorm.Model
	UserID     string `json:"user_id"`
	MerchantID string `json:"merchant_id"`
}

type Order struct {
	ID                 string    `json:"id"`                   // 订单ID
	UserID             string    `json:"user_id"`              // 用户ID
	MerchandiseStyleID string    `json:"merchandise_style_id"` // 款式id
	MerchandiseID      string    `json:"merchandise_id"`       // 商品id
	TotalAmount        float32   `json:"total_amount"`         // 总金额
	Status             string    `json:"status"`               // 订单状态
	PaymentMethod      string    `json:"payment_method"`       // 支付方式
	ShippingAddress    string    `json:"shipping_address"`     // 收货地址
	CreatedAt          time.Time `json:"created_at"`           // 下单时间
	UpdatedAt          time.Time `json:"updated_at"`           // 更新时间
	Quantity           int       `json:"quantity"`             // 购买数量
	Price              float32   `json:"price"`                // 单价
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
	UserID string `json:"user_id"`
}

// MerchantRecord 商家行为记录
type MerchantRecord struct {
	ID     string    `json:"id"`
	Time   time.Time `json:"time"`
	UserID string    `json:"user_id"`
	Action string    `json:"action"`
	Source string    `json:"source"`
}
