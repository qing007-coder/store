package model

import "time"

type User struct {
	ID       string `gorm:"primaryKey"`
	Account  string
	Name     string
	Password string
	Email    string
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
