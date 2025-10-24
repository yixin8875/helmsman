package model

type Users struct {
	ID           uint64 `gorm:"column:id;type:int(11);primary_key" json:"id"`
	Username     string `gorm:"column:username;type:text;not null" json:"username"`
	PasswordHash string `gorm:"column:password_hash;type:text;not null" json:"passwordHash"`
	CreatedAt    string `gorm:"column:created_at;type:varchar(100)" json:"createdAt"`
	UpdatedAt    string `gorm:"column:updated_at;type:varchar(100)" json:"updatedAt"`
}

// UsersColumnNames Whitelist for custom query fields to prevent sql injection attacks
var UsersColumnNames = map[string]bool{
	"id":            true,
	"username":      true,
	"password_hash": true,
	"created_at":    true,
	"updated_at":    true,
}
