package notification

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	expo "github.com/PhasitWo/duchenne-server/services/notification/expo/exponent-server-sdk-golang-master/sdk"
	"gorm.io/gorm"
)

type INotificationService interface {
	SendDailyNotifications(dayRange *int) error
	SendNotiByPatientId(id int, title string, body string) error
}

type service struct {
	Repo  repository.IRepo
	sqldb *sql.DB
}

var NotiLogger = log.New(os.Stdout, "[NOTI] ", log.LstdFlags)
var ErrDevicesNotFound = errors.New("error not found any devices")

func NewService(db *gorm.DB) *service {
	sqldb, err := db.DB()
	if err != nil {
		panic("can't get *sql.DB from gorm")
	}
	return &service{
		Repo:  repository.New(db),
		sqldb: sqldb,
	}
}

func (n *service) SendNotiByPatientId(id int, title string, body string) error {
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

/*
send daily notification about upcoming appointments,
day_range = nil will use default value from config
*/
func (n *service) SendDailyNotifications(dayRange *int) error {
	if dayRange == nil {
		dayRange = &config.AppConfig.NOTIFY_IN_RANGE
		NotiLogger.Println("using default day range from config")
	}
	res, err := queryDB(n.sqldb, *dayRange)
	if err != nil {
		msg := "can't query database"
		NotiLogger.Println(msg)
		return errors.New(msg)
	}
	if len(res) == 0 {
		NotiLogger.Printf("no upcoming appointments in the next %v days\n", *dayRange)
		return nil
	}
	NotiLogger.Printf("preparing messages..\n")
	// prepare messages
	// 1 appointmemnt -> 1 message -- to --> multiple receivers
	messagesPool := []expo.PushMessage{}
	var newMessage expo.PushMessage
	prior := -1
	for index, elem := range res {
		if elem.AppointmentId != prior {
			// add prior new message to pool
			if prior != -1 {
				messagesPool = append(messagesPool, newMessage)
			}
			// create new message
			newMessage = expo.PushMessage{
				To:       []expo.ExponentPushToken{expo.ExponentPushToken(elem.ExpoToken)},
				Body:     formatTimeOutput(elem.Date, int(time.Now().Unix())),
				Sound:    "default",
				Title:    "อย่าลืมนัดหมายของคุณ!",
				Priority: expo.HighPriority,
			}
		} else {
			newMessage.To = append(newMessage.To, expo.ExponentPushToken(elem.ExpoToken))
		}
		// special case -> if this new message is the last message
		if index == len(res)-1 {
			messagesPool = append(messagesPool, newMessage)
		}
		prior = elem.AppointmentId
	}

	// 1 request can contain up to 100 messages, for safety purpose -> 1 request should contain only up to 80 messages
	// divide len([]message) with 80 -> split up to multiple request
	NotiLogger.Printf("splitting up messages to multiple request\n")
	const MAX_MESSAGES_PER_REQUEST = 80
	var messageCnt = float64(len(messagesPool))
	var cnt float64 = math.Ceil(float64(messageCnt) / MAX_MESSAGES_PER_REQUEST)
	for i := 0; i < int(cnt); i++ {
		// calculate base and limit for slicing slice
		base := float64(i * MAX_MESSAGES_PER_REQUEST)
		limit := base + math.Min(messageCnt-base, MAX_MESSAGES_PER_REQUEST)

		// send request
		NotiLogger.Printf("sending request %v => messagesPool[%v:%v]\n", i, base, limit)
		SendRequest(messagesPool[int(base):int(limit)])
	}
	return nil
}

func SendRequest(messages []expo.PushMessage) {
	client := expo.NewPushClient(nil)
	// Publish message
	responses, err := client.PublishMultiple(messages)

	// Check errors
	if err != nil && !strings.Contains(err.Error(), "mismatched response length") {
		NotiLogger.Panic(err)
	}
	NotiLogger.Println("validating..")
	// Validate responses
	for index, response := range responses {
		fmt.Printf("push ticket %v =>", index)
		if response.ValidateResponse() != nil {
			fmt.Println(response.PushMessage.To, "failed")
		} else {
			fmt.Println(response.PushMessage.To, "succeed")
		}
	}
}

var apmtQuery = `
select appointments.id ,date, devices.id, devices.device_name , devices.expo_token, appointments.patient_id from appointments
inner join devices on appointments.patient_id = devices.patient_id
where devices.expo_token != "" AND appointments.approve_at IS NOT NULL AND appointments.date > ? AND appointments.date < ?
order by appointments.id asc
`

func queryDB(db *sql.DB, dayRange int) ([]model.AppointmentDevice, error) {
	now := int(time.Now().Unix())
	limit := now + dayRange*24*60*60
	rows, err := db.Query(apmtQuery, now, limit)
	if err != nil {
		NotiLogger.Println("queryDB : Can't query database")
		NotiLogger.Println(err.Error())
		return nil, err
	}
	defer rows.Close()
	res := []model.AppointmentDevice{}
	for rows.Next() {
		var ad model.AppointmentDevice
		if err := rows.Scan(
			&ad.AppointmentId,
			&ad.Date,
			&ad.DeviceId,
			&ad.DeviceName,
			&ad.ExpoToken,
			&ad.PatientId,
		); err != nil {
			fmt.Printf("queryDB : %v", err.Error())
			return nil, err
		}
		res = append(res, ad)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("queryDB : %v", err.Error())
		return nil, err
	}
	NotiLogger.Printf("now: %v, limit: %v (day range: %v)\n", now, limit, dayRange)
	return res, nil
}

func formatTimeOutput(dueTimestamp int, nowTimestamp int) string {
	sec := (dueTimestamp - nowTimestamp)
	minute := sec / 60
	hour := minute / 60
	day := hour / 24
	baseStr := "คุณมีนัดหมายในอีก "
	var output string
	if minute == 0 {
		output = "ไม่กี่นาที"
	} else if hour == 0 {
		output = fmt.Sprintf("%d นาที", minute)
	} else if day == 0 {
		output = fmt.Sprintf("%d ชั่วโมง %d นาที", hour, minute%60)
	} else {
		output = fmt.Sprintf("%d วัน %d ชั่วโมง", day, hour%24)
	}
	return baseStr + output
}
