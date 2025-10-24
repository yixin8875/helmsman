package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// accounts business-level http error codes.
// the accountsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	accountsNO       = 79
	accountsName     = "accounts"
	accountsBaseCode = errcode.HCode(accountsNO)

	ErrCreateAccounts     = errcode.NewError(accountsBaseCode+1, "failed to create "+accountsName)
	ErrDeleteByIDAccounts = errcode.NewError(accountsBaseCode+2, "failed to delete "+accountsName)
	ErrUpdateByIDAccounts = errcode.NewError(accountsBaseCode+3, "failed to update "+accountsName)
	ErrGetByIDAccounts    = errcode.NewError(accountsBaseCode+4, "failed to get "+accountsName+" details")
	ErrListAccounts       = errcode.NewError(accountsBaseCode+5, "failed to list of "+accountsName)

	// error codes are globally unique, adding 1 to the previous error code
)
