package entity

import (
	"database/sql"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type User struct {
	Id         int64        `json:"id"`
	Email      string       `json:"email"`
	Password   string       `json:"password"`
	UserRoleId int64        `json:"user_role_id"`
	IsVerified bool         `json:"is_verified"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	DeletedAt  sql.NullTime `json:"deleted_at"`
}

func (u *User) GetEntityName() string {
	return "users"
}

func (u *User) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(u).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (u *User) ToUserResponse() *responsedto.UserResponse {
	return &responsedto.UserResponse{
		Id:         u.Id,
		Email:      u.Email,
		UserRoleId: u.UserRoleId,
		IsVerified: u.IsVerified,
	}
}
