package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateUsersRequest request params
type CreateUsersRequest struct {
	Username     string `json:"username" binding:""`
	PasswordHash string `json:"passwordHash" binding:""`
}

// UpdateUsersByIDRequest request params
type UpdateUsersByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Username     string `json:"username" binding:""`
	PasswordHash string `json:"passwordHash" binding:""`
}

// UsersObjDetail detail
type UsersObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// CreateUsersReply only for api docs
type CreateUsersReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateUsersByIDReply only for api docs
type UpdateUsersByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetUsersByIDReply only for api docs
type GetUsersByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Users UsersObjDetail `json:"users"`
	} `json:"data"` // return data
}

// DeleteUsersByIDReply only for api docs
type DeleteUsersByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteUserssByIDsReply only for api docs
type DeleteUserssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListUserssRequest request params
type ListUserssRequest struct {
	query.Params
}

// GetUsersByConditionRequest request params
type GetUsersByConditionRequest struct {
	query.Conditions
}

// GetUsersByConditionReply only for api docs
type GetUsersByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Users UsersObjDetail `json:"users"`
	} `json:"data"` // return data
}
