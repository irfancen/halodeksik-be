package requestdto

import (
	"github.com/shopspring/decimal"
	"halodeksik-be/app/entity"
	"mime/multipart"
)

type RequestDoctorProfile struct {
	Name                   string `json:"name" form:"name" validate:"required"`
	StartingYear           int32  `json:"starting_year" form:"starting_year"  validate:"required,min=1800"`
	DoctorSpecializationId int64  `json:"doctor_specialization_id" form:"doctor_specialization_id"  validate:"required"`
	ConsultationFee        string `json:"consultation_fee" form:"starting_year"  validate:"required,numeric,numericgt=0"`
}

type RequestDoctorCertificate struct {
	Certificate *multipart.FileHeader `json:"certificate" form:"certificate" validate:"omitempty,filetype=pdf png jpg jpeg,filesize=500"`
}

type RequestProfilePhoto struct {
	ProfilePhoto *multipart.FileHeader `json:"profile_photo" form:"profile_photo" validate:"omitempty,filetype=png jpg jpeg,filesize=500"`
}

func (p RequestDoctorProfile) ToDoctorProfile() entity.DoctorProfile {
	fee, _ := decimal.NewFromString(p.ConsultationFee)
	return entity.DoctorProfile{
		Name:                   p.Name,
		StartingYear:           p.StartingYear,
		DoctorSpecializationId: p.DoctorSpecializationId,
		ConsultationFee:        fee,
	}
}
