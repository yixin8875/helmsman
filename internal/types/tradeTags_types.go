package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateTradeTagsRequest request params
type CreateTradeTagsRequest struct {
	TradeID int `json:"tradeID" binding:""`
	TagID   int `json:"tagID" binding:""`
}

// UpdateTradeTagsByTradeIDRequest request params
type UpdateTradeTagsByTradeIDRequest struct {
	TradeID int `json:"tradeID" binding:""`
	TagID   int `json:"tagID" binding:""`
}

// TradeTagsObjDetail detail
type TradeTagsObjDetail struct {
	TradeID   int    `json:"tradeID"`
	TagID     int    `json:"tagID"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// CreateTradeTagsReply only for api docs
type CreateTradeTagsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		TradeID int `json:"tradeID"`
	} `json:"data"` // return data
}

// DeleteTradeTagsByTradeIDReply only for api docs
type DeleteTradeTagsByTradeIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateTradeTagsByTradeIDReply only for api docs
type UpdateTradeTagsByTradeIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetTradeTagsByTradeIDReply only for api docs
type GetTradeTagsByTradeIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		TradeTags TradeTagsObjDetail `json:"tradeTags"`
	} `json:"data"` // return data
}

// ListTradeTagsRequest request params
type ListTradeTagsRequest struct {
	query.Params
}

// ListTradeTagsReply only for api docs
type ListTradeTagsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		TradeTags []TradeTagsObjDetail `json:"tradeTags"`
	} `json:"data"` // return data
}
