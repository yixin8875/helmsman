package model

type Trades struct {
	ID                uint64  `gorm:"column:id;type:int(11);primary_key" json:"id"`
	AccountID         int     `gorm:"column:account_id;type:int(11);not null" json:"accountID"`
	StrategyID        int     `gorm:"column:strategy_id;type:int(11)" json:"strategyID"`
	Status            string  `gorm:"column:status;type:text;not null" json:"status"`
	Symbol            string  `gorm:"column:symbol;type:text;not null" json:"symbol"`
	Direction         string  `gorm:"column:direction;type:text;not null" json:"direction"`
	PlannedEntryPrice float64 `gorm:"column:planned_entry_price;type:float" json:"plannedEntryPrice"`
	PlannedStopLoss   float64 `gorm:"column:planned_stop_loss;type:float" json:"plannedStopLoss"`
	PlannedTakeProfit float64 `gorm:"column:planned_take_profit;type:float" json:"plannedTakeProfit"`
	PositionSize      float64 `gorm:"column:position_size;type:float" json:"positionSize"`
	PlannedRiskAmount float64 `gorm:"column:planned_risk_amount;type:float" json:"plannedRiskAmount"`
	PlanNotes         string  `gorm:"column:plan_notes;type:text" json:"planNotes"`
	ActualEntryTime   string  `gorm:"column:actual_entry_time;type:varchar(100)" json:"actualEntryTime"`
	ActualEntryPrice  float64 `gorm:"column:actual_entry_price;type:float" json:"actualEntryPrice"`
	ActualExitTime    string  `gorm:"column:actual_exit_time;type:varchar(100)" json:"actualExitTime"`
	ActualExitPrice   float64 `gorm:"column:actual_exit_price;type:float" json:"actualExitPrice"`
	Commission        float64 `gorm:"column:commission;type:float" json:"commission"`
	Pnl               float64 `gorm:"column:pnl;type:float" json:"pnl"`
	RMultiple         float64 `gorm:"column:r_multiple;type:float" json:"rMultiple"`
	ExitReason        string  `gorm:"column:exit_reason;type:text" json:"exitReason"`
	ExecutionScore    int     `gorm:"column:execution_score;type:int(11)" json:"executionScore"`
	ReflectionNotes   string  `gorm:"column:reflection_notes;type:text" json:"reflectionNotes"`
	CreatedAt         string  `gorm:"column:created_at;type:varchar(100)" json:"createdAt"`
	UpdatedAt         string  `gorm:"column:updated_at;type:varchar(100)" json:"updatedAt"`
}

// TradesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var TradesColumnNames = map[string]bool{
	"id":                  true,
	"account_id":          true,
	"strategy_id":         true,
	"status":              true,
	"symbol":              true,
	"direction":           true,
	"planned_entry_price": true,
	"planned_stop_loss":   true,
	"planned_take_profit": true,
	"position_size":       true,
	"planned_risk_amount": true,
	"plan_notes":          true,
	"actual_entry_time":   true,
	"actual_entry_price":  true,
	"actual_exit_time":    true,
	"actual_exit_price":   true,
	"commission":          true,
	"pnl":                 true,
	"r_multiple":          true,
	"exit_reason":         true,
	"execution_score":     true,
	"reflection_notes":    true,
	"created_at":          true,
	"updated_at":          true,
}
