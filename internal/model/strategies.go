package model

type Strategies struct {
	ID          uint64 `gorm:"column:id;type:int(11);primary_key" json:"id"`
	UserID      int    `gorm:"column:user_id;type:int(11);not null" json:"userID"`
	Name        string `gorm:"column:name;type:text;not null" json:"name"`
	Description string `gorm:"column:description;type:text" json:"description"`
	CreatedAt   string `gorm:"column:created_at;type:varchar(100)" json:"createdAt"`
	UpdatedAt   string `gorm:"column:updated_at;type:varchar(100)" json:"updatedAt"`
}

// StrategiesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var StrategiesColumnNames = map[string]bool{
	"id":          true,
	"user_id":     true,
	"name":        true,
	"description": true,
	"created_at":  true,
	"updated_at":  true,
}
