package services

import (
	"sync"

	"htst/internal/app/models"
	"htst/internal/pkg/code"
)

var initOnce sync.Once
var (
	IUser IUserService
	IInfo IInfoService
)

func Init() {
	initOnce.Do(func() {
		IUser = NewUserService()
		IInfo = NewInfoService()
	})
}

var successInfo = models.RespInfo{
	Code: code.Success,
	Msg:  code.MsgSuccess,
}
