package responsedto

type DoctorProfileResponse struct {
	Id                   int64  `json:"id"`
	Email                string `json:"email"`
	UserRoleID           int64  `json:"user_role_id"`
	IsVerified           bool   `json:"is_verified"`
	Name                 string `json:"name"`
	ProfilePhoto         string `json:"profile_photo"`
	StartingYear         int32  `json:"starting_year"`
	DoctorCertificate    string `json:"doctor_certificate"`
	DoctorSpecialization string `json:"doctor_specialization"`
	ConsultationFee      string `json:"consultation_fee"`
	IsOnline             bool   `json:"is_online"`
}
