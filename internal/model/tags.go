package model

type Tags struct {
	ID        uint64 `gorm:"column:id;type:int(11);primary_key" json:"id"`
	UserID    int    `gorm:"column:user_id;type:int(11);not null" json:"userID"`
	Name      string `gorm:"column:name;type:text;not null" json:"name"`
	Color     string `gorm:"column:color;type:text" json:"color"`
	CreatedAt string `gorm:"column:created_at;type:varchar(100)" json:"createdAt"`
	UpdatedAt string `gorm:"column:updated_at;type:varchar(100)" json:"updatedAt"`
}

// TagsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var TagsColumnNames = map[string]bool{
	"id":         true,
	"user_id":    true,
	"name":       true,
	"color":      true,
	"created_at": true,
	"updated_at": true,
}
