package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateAccountsRequest request params
type CreateAccountsRequest struct {
	UserID         int     `json:"userID" binding:""`
	Name           string  `json:"name" binding:""`
	InitialBalance float64 `json:"initialBalance" binding:""`
	Currency       string  `json:"currency" binding:""`
}

// UpdateAccountsByIDRequest request params
type UpdateAccountsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	UserID         int     `json:"userID" binding:""`
	Name           string  `json:"name" binding:""`
	InitialBalance float64 `json:"initialBalance" binding:""`
	Currency       string  `json:"currency" binding:""`
}

// AccountsObjDetail detail
type AccountsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	UserID         int     `json:"userID"`
	Name           string  `json:"name"`
	InitialBalance float64 `json:"initialBalance"`
	Currency       string  `json:"currency"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

// CreateAccountsReply only for api docs
type CreateAccountsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteAccountsByIDReply only for api docs
type DeleteAccountsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateAccountsByIDReply only for api docs
type UpdateAccountsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetAccountsByIDReply only for api docs
type GetAccountsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Accounts AccountsObjDetail `json:"accounts"`
	} `json:"data"` // return data
}

// ListAccountssRequest request params
type ListAccountssRequest struct {
	query.Params
}

// ListAccountssReply only for api docs
type ListAccountssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Accountss []AccountsObjDetail `json:"accountss"`
	} `json:"data"` // return data
}
