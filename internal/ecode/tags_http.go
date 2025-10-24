package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// tags business-level http error codes.
// the tagsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	tagsNO       = 36
	tagsName     = "tags"
	tagsBaseCode = errcode.HCode(tagsNO)

	ErrCreateTags     = errcode.NewError(tagsBaseCode+1, "failed to create "+tagsName)
	ErrDeleteByIDTags = errcode.NewError(tagsBaseCode+2, "failed to delete "+tagsName)
	ErrUpdateByIDTags = errcode.NewError(tagsBaseCode+3, "failed to update "+tagsName)
	ErrGetByIDTags    = errcode.NewError(tagsBaseCode+4, "failed to get "+tagsName+" details")
	ErrListTags       = errcode.NewError(tagsBaseCode+5, "failed to list of "+tagsName)

	// error codes are globally unique, adding 1 to the previous error code
)
