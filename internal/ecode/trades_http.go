package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// trades business-level http error codes.
// the tradesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	tradesNO       = 80
	tradesName     = "trades"
	tradesBaseCode = errcode.HCode(tradesNO)

	ErrCreateTrades     = errcode.NewError(tradesBaseCode+1, "failed to create "+tradesName)
	ErrDeleteByIDTrades = errcode.NewError(tradesBaseCode+2, "failed to delete "+tradesName)
	ErrUpdateByIDTrades = errcode.NewError(tradesBaseCode+3, "failed to update "+tradesName)
	ErrGetByIDTrades    = errcode.NewError(tradesBaseCode+4, "failed to get "+tradesName+" details")
	ErrListTrades       = errcode.NewError(tradesBaseCode+5, "failed to list of "+tradesName)

	// error codes are globally unique, adding 1 to the previous error code
)
