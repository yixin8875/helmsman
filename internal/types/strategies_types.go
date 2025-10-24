package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateStrategiesRequest request params
type CreateStrategiesRequest struct {
	UserID      int    `json:"userID" binding:""`
	Name        string `json:"name" binding:""`
	Description string `json:"description" binding:""`
}

// UpdateStrategiesByIDRequest request params
type UpdateStrategiesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	UserID      int    `json:"userID" binding:""`
	Name        string `json:"name" binding:""`
	Description string `json:"description" binding:""`
}

// StrategiesObjDetail detail
type StrategiesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	UserID      int    `json:"userID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateStrategiesReply only for api docs
type CreateStrategiesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteStrategiesByIDReply only for api docs
type DeleteStrategiesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateStrategiesByIDReply only for api docs
type UpdateStrategiesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetStrategiesByIDReply only for api docs
type GetStrategiesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Strategies StrategiesObjDetail `json:"strategies"`
	} `json:"data"` // return data
}

// ListStrategiessRequest request params
type ListStrategiessRequest struct {
	query.Params
}

// ListStrategiessReply only for api docs
type ListStrategiessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Strategiess []StrategiesObjDetail `json:"strategiess"`
	} `json:"data"` // return data
}
