package notification

import (
	expo "github.com/PhasitWo/duchenne-server/notification/expo/exponent-server-sdk-golang-master/sdk"
	"github.com/PhasitWo/duchenne-server/repository"
	"gorm.io/gorm"
)

type INotificationService interface {
	SendNotiByPatientId(id int, title string, body string) error
}

type Service struct {
	Repo repository.IRepo
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		Repo: repository.New(db),
	}
}

func (n *Service) SendNotiByPatientId(id int, title string, body string) error {
	devices, err := n.Repo.GetAllDevice(repository.Criteria{QueryCriteria: repository.PATIENTID, Value: id})
	if err != nil {
		NotiLogger.Println("Error can't get devices to push notifications")
		return err
	}
	if len(devices) == 0 {
		NotiLogger.Println("Error no devices to push notifications")
		return ErrDevicesNotFound
	}
	msg := expo.PushMessage{To: []expo.ExponentPushToken{}, Title: title, Body: body, Sound: "default", Priority: expo.HighPriority}
	for _, d := range devices {
		msg.To = append(msg.To, expo.ExponentPushToken(d.ExpoToken))
	}
	SendRequest([]expo.PushMessage{msg})
	return nil
}