package requestdto

import (
	"halodeksik-be/app/entity"
)

type RequestRegisterToken struct {
	Email string `json:"email" validate:"required,email"`
}

type RequestTokenUrl struct {
	Token string `form:"token" validate:"required"`
}

type RequestRegisterUser struct {
	Email      string `json:"email" form:"email" validate:"required"`
	Password   string `json:"password" form:"password" validate:"required,min=8,max=72"`
	UserRoleId int64  `json:"user_role_id" form:"user_role_id" validate:"required"`
}

func (u *RequestRegisterUser) ToUser() entity.User {
	return entity.User{
		Email:      u.Email,
		Password:   u.Password,
		UserRoleId: u.UserRoleId,
	}
}
