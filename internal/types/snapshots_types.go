package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateSnapshotsRequest request params
type CreateSnapshotsRequest struct {
	TradeID  int    `json:"tradeID" binding:""`
	Type     string `json:"type" binding:""`
	ImageURL string `json:"imageURL" binding:""`
}

// UpdateSnapshotsByIDRequest request params
type UpdateSnapshotsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	TradeID  int    `json:"tradeID" binding:""`
	Type     string `json:"type" binding:""`
	ImageURL string `json:"imageURL" binding:""`
}

// SnapshotsObjDetail detail
type SnapshotsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	TradeID   int    `json:"tradeID"`
	Type      string `json:"type"`
	ImageURL  string `json:"imageURL"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// CreateSnapshotsReply only for api docs
type CreateSnapshotsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteSnapshotsByIDReply only for api docs
type DeleteSnapshotsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateSnapshotsByIDReply only for api docs
type UpdateSnapshotsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetSnapshotsByIDReply only for api docs
type GetSnapshotsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Snapshots SnapshotsObjDetail `json:"snapshots"`
	} `json:"data"` // return data
}

// ListSnapshotssRequest request params
type ListSnapshotssRequest struct {
	query.Params
}

// ListSnapshotssReply only for api docs
type ListSnapshotssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Snapshotss []SnapshotsObjDetail `json:"snapshotss"`
	} `json:"data"` // return data
}
