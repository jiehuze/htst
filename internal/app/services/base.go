package services

import (
	"sync"

	"htst/internal/app/models"
	"htst/internal/pkg/code"
)

var initOnce sync.Once
var (
	IUser    IUserService
	IInfo    IInfoService
	IDoctor  IDoctorInfoService
	IContact IContactService
)

func Init() {
	initOnce.Do(func() {
		IUser = NewUserService()
		IInfo = NewInfoService()
		IDoctor = NewDoctorInfoService()
		IContact = NewContactService()
	})
}

var successInfo = models.RespInfo{
	Code: code.Success,
	Msg:  code.MsgSuccess,
}
