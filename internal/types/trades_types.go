package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateTradesRequest request params
type CreateTradesRequest struct {
	AccountID         int     `json:"accountID" binding:""`
	StrategyID        int     `json:"strategyID" binding:""`
	Status            string  `json:"status" binding:""`
	Symbol            string  `json:"symbol" binding:""`
	Direction         string  `json:"direction" binding:""`
	PlannedEntryPrice float64 `json:"plannedEntryPrice" binding:""`
	PlannedStopLoss   float64 `json:"plannedStopLoss" binding:""`
	PlannedTakeProfit float64 `json:"plannedTakeProfit" binding:""`
	PositionSize      float64 `json:"positionSize" binding:""`
	PlannedRiskAmount float64 `json:"plannedRiskAmount" binding:""`
	PlanNotes         string  `json:"planNotes" binding:""`
	ActualEntryTime   string  `json:"actualEntryTime" binding:""`
	ActualEntryPrice  float64 `json:"actualEntryPrice" binding:""`
	ActualExitTime    string  `json:"actualExitTime" binding:""`
	ActualExitPrice   float64 `json:"actualExitPrice" binding:""`
	Commission        float64 `json:"commission" binding:""`
	Pnl               float64 `json:"pnl" binding:""`
	RMultiple         float64 `json:"rMultiple" binding:""`
	ExitReason        string  `json:"exitReason" binding:""`
	ExecutionScore    int     `json:"executionScore" binding:""`
	ReflectionNotes   string  `json:"reflectionNotes" binding:""`
}

// UpdateTradesByIDRequest request params
type UpdateTradesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	AccountID         int     `json:"accountID" binding:""`
	StrategyID        int     `json:"strategyID" binding:""`
	Status            string  `json:"status" binding:""`
	Symbol            string  `json:"symbol" binding:""`
	Direction         string  `json:"direction" binding:""`
	PlannedEntryPrice float64 `json:"plannedEntryPrice" binding:""`
	PlannedStopLoss   float64 `json:"plannedStopLoss" binding:""`
	PlannedTakeProfit float64 `json:"plannedTakeProfit" binding:""`
	PositionSize      float64 `json:"positionSize" binding:""`
	PlannedRiskAmount float64 `json:"plannedRiskAmount" binding:""`
	PlanNotes         string  `json:"planNotes" binding:""`
	ActualEntryTime   string  `json:"actualEntryTime" binding:""`
	ActualEntryPrice  float64 `json:"actualEntryPrice" binding:""`
	ActualExitTime    string  `json:"actualExitTime" binding:""`
	ActualExitPrice   float64 `json:"actualExitPrice" binding:""`
	Commission        float64 `json:"commission" binding:""`
	Pnl               float64 `json:"pnl" binding:""`
	RMultiple         float64 `json:"rMultiple" binding:""`
	ExitReason        string  `json:"exitReason" binding:""`
	ExecutionScore    int     `json:"executionScore" binding:""`
	ReflectionNotes   string  `json:"reflectionNotes" binding:""`
}

// TradesObjDetail detail
type TradesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	AccountID         int     `json:"accountID"`
	StrategyID        int     `json:"strategyID"`
	Status            string  `json:"status"`
	Symbol            string  `json:"symbol"`
	Direction         string  `json:"direction"`
	PlannedEntryPrice float64 `json:"plannedEntryPrice"`
	PlannedStopLoss   float64 `json:"plannedStopLoss"`
	PlannedTakeProfit float64 `json:"plannedTakeProfit"`
	PositionSize      float64 `json:"positionSize"`
	PlannedRiskAmount float64 `json:"plannedRiskAmount"`
	PlanNotes         string  `json:"planNotes"`
	ActualEntryTime   string  `json:"actualEntryTime"`
	ActualEntryPrice  float64 `json:"actualEntryPrice"`
	ActualExitTime    string  `json:"actualExitTime"`
	ActualExitPrice   float64 `json:"actualExitPrice"`
	Commission        float64 `json:"commission"`
	Pnl               float64 `json:"pnl"`
	RMultiple         float64 `json:"rMultiple"`
	ExitReason        string  `json:"exitReason"`
	ExecutionScore    int     `json:"executionScore"`
	ReflectionNotes   string  `json:"reflectionNotes"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

// CreateTradesReply only for api docs
type CreateTradesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteTradesByIDReply only for api docs
type DeleteTradesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateTradesByIDReply only for api docs
type UpdateTradesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetTradesByIDReply only for api docs
type GetTradesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Trades TradesObjDetail `json:"trades"`
	} `json:"data"` // return data
}

// ListTradessRequest request params
type ListTradessRequest struct {
	query.Params
}

// ListTradessReply only for api docs
type ListTradessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Tradess []TradesObjDetail `json:"tradess"`
	} `json:"data"` // return data
}
