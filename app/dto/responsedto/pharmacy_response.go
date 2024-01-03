package responsedto

type PharmacyResponse struct {
	Id                    int64    `json:"id"`
	Name                  string   `json:"name"`
	Address               string   `json:"address"`
	SubDistrict           string   `json:"sub_district"`
	District              string   `json:"district"`
	City                  string   `json:"city"`
	Province              string   `json:"province"`
	PostalCode            string   `json:"postal_code"`
	Latitude              string   `json:"latitude"`
	Longitude             string   `json:"longitude"`
	PharmacistName        string   `json:"pharmacist_name"`
	PharmacistLicenseNo   string   `json:"pharmacist_license_no"`
	PharmacistPhoneNo     string   `json:"pharmacist_phone_no"`
	OperationalHoursOpen  int      `json:"operational_hours_open"`
	OperationalHoursClose int      `json:"operational_hours_close"`
	OperationalDays       []string `json:"operational_days"`
	PharmacyAdminId       int64    `json:"pharmacy_admin_id"`
}
