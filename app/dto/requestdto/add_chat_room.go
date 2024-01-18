package requestdto

type AddChatRoom struct {
	DoctorId  int64 `json:"doctor_id" validate:"required"`
	PatientId int64 `json:"patient_id" validate:"required"`
}
