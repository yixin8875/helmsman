package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// strategies business-level http error codes.
// the strategiesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	strategiesNO       = 76
	strategiesName     = "strategies"
	strategiesBaseCode = errcode.HCode(strategiesNO)

	ErrCreateStrategies     = errcode.NewError(strategiesBaseCode+1, "failed to create "+strategiesName)
	ErrDeleteByIDStrategies = errcode.NewError(strategiesBaseCode+2, "failed to delete "+strategiesName)
	ErrUpdateByIDStrategies = errcode.NewError(strategiesBaseCode+3, "failed to update "+strategiesName)
	ErrGetByIDStrategies    = errcode.NewError(strategiesBaseCode+4, "failed to get "+strategiesName+" details")
	ErrListStrategies       = errcode.NewError(strategiesBaseCode+5, "failed to list of "+strategiesName)

	// error codes are globally unique, adding 1 to the previous error code
)
