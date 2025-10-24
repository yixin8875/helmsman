package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// tradeTags business-level http error codes.
// the tradeTagsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	tradeTagsNO       = 1
	tradeTagsName     = "tradeTags"
	tradeTagsBaseCode = errcode.HCode(tradeTagsNO)

	ErrCreateTradeTags          = errcode.NewError(tradeTagsBaseCode+1, "failed to create "+tradeTagsName)
	ErrDeleteByTradeIDTradeTags = errcode.NewError(tradeTagsBaseCode+2, "failed to delete "+tradeTagsName)
	ErrUpdateByTradeIDTradeTags = errcode.NewError(tradeTagsBaseCode+3, "failed to update "+tradeTagsName)
	ErrGetByTradeIDTradeTags    = errcode.NewError(tradeTagsBaseCode+4, "failed to get "+tradeTagsName+" details")
	ErrListTradeTags            = errcode.NewError(tradeTagsBaseCode+5, "failed to list of "+tradeTagsName)

	// error codes are globally unique, adding 1 to the previous error code
)
