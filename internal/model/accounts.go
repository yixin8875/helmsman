package model

type Accounts struct {
	ID             uint64  `gorm:"column:id;type:int(11);primary_key" json:"id"`
	UserID         int     `gorm:"column:user_id;type:int(11);not null" json:"userID"`
	Name           string  `gorm:"column:name;type:text;not null" json:"name"`
	InitialBalance float64 `gorm:"column:initial_balance;type:float;not null" json:"initialBalance"`
	Currency       string  `gorm:"column:currency;type:text" json:"currency"`
	CreatedAt      string  `gorm:"column:created_at;type:varchar(100)" json:"createdAt"`
	UpdatedAt      string  `gorm:"column:updated_at;type:varchar(100)" json:"updatedAt"`
}

// AccountsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var AccountsColumnNames = map[string]bool{
	"id":              true,
	"user_id":         true,
	"name":            true,
	"initial_balance": true,
	"currency":        true,
	"created_at":      true,
	"updated_at":      true,
}
