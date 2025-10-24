package model

type TradeTags struct {
	TradeID   int    `gorm:"column:trade_id;type:int(11);primary_key" json:"tradeID"`
	TagID     int    `gorm:"column:tag_id;type:int(11);not null" json:"tagID"`
	CreatedAt string `gorm:"column:created_at;type:varchar(100)" json:"createdAt"`
	UpdatedAt string `gorm:"column:updated_at;type:varchar(100)" json:"updatedAt"`
}

// TradeTagsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var TradeTagsColumnNames = map[string]bool{
	"trade_id":   true,
	"tag_id":     true,
	"created_at": true,
	"updated_at": true,
}
