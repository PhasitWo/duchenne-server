package model

type Device struct {
	ID         int    `json:"id"`
	LoginAt    int    `json:"loginAt"`
	DeviceName string `json:"deviceName"`
	ExpoToken  string `json:"expoToken"`
	PatientId  int    `json:"patientId"`
}

type AppointmentDevice struct {
	AppointmentId int    `json:"appointment_id"`
	Date          int    `json:"date"`
	DeviceId      int    `json:"device_id"`
	DeviceName    string `json:"device_name"`
	ExpoToken     string `json:"expoToken"`
	PatientId     int    `json:"patient_id"`
}
