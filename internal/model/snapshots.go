package model

type Snapshots struct {
	ID        uint64 `gorm:"column:id;type:int(11);primary_key" json:"id"`
	TradeID   int    `gorm:"column:trade_id;type:int(11);not null" json:"tradeID"`
	Type      string `gorm:"column:type;type:text;not null" json:"type"`
	ImageURL  string `gorm:"column:image_url;type:text;not null" json:"imageURL"`
	CreatedAt string `gorm:"column:created_at;type:varchar(100)" json:"createdAt"`
	UpdatedAt string `gorm:"column:updated_at;type:varchar(100)" json:"updatedAt"`
}

// SnapshotsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var SnapshotsColumnNames = map[string]bool{
	"id":         true,
	"trade_id":   true,
	"type":       true,
	"image_url":  true,
	"created_at": true,
	"updated_at": true,
}
