package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateTagsRequest request params
type CreateTagsRequest struct {
	UserID int    `json:"userID" binding:""`
	Name   string `json:"name" binding:""`
	Color  string `json:"color" binding:""`
}

// UpdateTagsByIDRequest request params
type UpdateTagsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	UserID int    `json:"userID" binding:""`
	Name   string `json:"name" binding:""`
	Color  string `json:"color" binding:""`
}

// TagsObjDetail detail
type TagsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	UserID    int    `json:"userID"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// CreateTagsReply only for api docs
type CreateTagsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteTagsByIDReply only for api docs
type DeleteTagsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateTagsByIDReply only for api docs
type UpdateTagsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetTagsByIDReply only for api docs
type GetTagsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Tags TagsObjDetail `json:"tags"`
	} `json:"data"` // return data
}

// ListTagssRequest request params
type ListTagssRequest struct {
	query.Params
}

// ListTagssReply only for api docs
type ListTagssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Tagss []TagsObjDetail `json:"tagss"`
	} `json:"data"` // return data
}
