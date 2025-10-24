package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// snapshots business-level http error codes.
// the snapshotsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	snapshotsNO       = 73
	snapshotsName     = "snapshots"
	snapshotsBaseCode = errcode.HCode(snapshotsNO)

	ErrCreateSnapshots     = errcode.NewError(snapshotsBaseCode+1, "failed to create "+snapshotsName)
	ErrDeleteByIDSnapshots = errcode.NewError(snapshotsBaseCode+2, "failed to delete "+snapshotsName)
	ErrUpdateByIDSnapshots = errcode.NewError(snapshotsBaseCode+3, "failed to update "+snapshotsName)
	ErrGetByIDSnapshots    = errcode.NewError(snapshotsBaseCode+4, "failed to get "+snapshotsName+" details")
	ErrListSnapshots       = errcode.NewError(snapshotsBaseCode+5, "failed to list of "+snapshotsName)

	// error codes are globally unique, adding 1 to the previous error code
)
