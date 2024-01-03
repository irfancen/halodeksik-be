package responsedto

type LoginResponse struct {
	UserId     int64  `json:"user_id"`
	Email      string `json:"email"`
	UserRoleId int64  `json:"user_role_id"`
	Image      string `json:"image"`
	Token      string `json:"token"`
}
