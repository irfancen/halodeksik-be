package api

import (
	"halodeksik-be/app/util"
)

type AllUtil struct {
	AuthUtil util.AuthUtil
	MailUtil util.EmailUtil
	LocUtil  util.LocationUtil
}

func InitializeUtil() *AllUtil {
	return &AllUtil{
		AuthUtil: util.NewAuthUtil(),
		MailUtil: util.NewEmailUtil(),
		LocUtil:  util.NewLocationUtil("id"),
	}
}
