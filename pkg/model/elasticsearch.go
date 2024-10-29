package model

import "time"

// Merchandise 商品
type Merchandise struct {
	ID          string   `json:"id"`           // 商品id
	Name        string   `json:"name"`         // 商品名
	Info        string   `json:"info"`         // 商品简介
	PictureList []string `json:"picture_list"` // 照片路径
	MerchantID  string   `json:"merchant_id"`  // 商家id
	Delivery    string   `json:"delivery"`     // 快递发货情况
	Category    string   `json:"category"`     // 商品分类
	Views       uint     `json:"views"`        // 浏览量
	SalesVolume uint     `json:"sales_volume"` // 购买量
}

// MerchandiseStyle 商品款式
type MerchandiseStyle struct {
	ID            string    `json:"id"`             // 款式id
	MerchandiseID string    `json:"merchandise_id"` // 商品id
	Name          string    `json:"name"`           // 款式名
	Info          string    `json:"info"`           // 简介
	Picture       string    `json:"picture"`        // 款式封面
	Price         float32   `json:"price"`          // 款式价格
	Stock         uint      `json:"stock"`          // 款式库存数量
	Status        string    `json:"status"`         // 款式状态（如在售、下架）
	CreatedAt     time.Time `json:"created_at"`     // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`     // 更新时间
}

// Log 日志
type Log struct {
	ID     string `json:"id"`
	Level  string `json:"level"`
	Time   int64  `json:"time"`
	Msg    string `json:"msg"`
	Source string `json:"source"`
}
