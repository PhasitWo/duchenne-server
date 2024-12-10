package notification

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/PhasitWo/duchenne-server/model"
	expo "github.com/PhasitWo/duchenne-server/notification/expo/exponent-server-sdk-golang-master/sdk"
	"github.com/gin-gonic/gin"
)

func TestPushNotification(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		MockupScheduleNotifications(db, mockSendRequest)
	}
}

func MockupScheduleNotifications(db *sql.DB, sendRequestFunc func([]expo.PushMessage)) {
	// query
	res, err := queryDB(db)
	if err != nil {
		fmt.Println("Notification: Can't query database")
		return
	}
	if res == nil {
		fmt.Println("Notification: No appointment..")
		return
	}
	// prepare messages
	// 1 appointmemnt -> 1 message -- to --> multiple receivers
	messagesPool := []expo.PushMessage{}
	var newMessage expo.PushMessage
	prior := -1
	for _, elem := range res {
		if elem.AppointmentId != prior {
			if prior != -1 {
				messagesPool = append(messagesPool, newMessage)
			}
			newMessage = expo.PushMessage{
				To:       []expo.ExponentPushToken{expo.ExponentPushToken(elem.ExpoToken)},
				Body:     formatTimeOutput(elem.Date, int(time.Now().Unix())),
				Sound:    "default",
				Title:    "Test Notification",
				Priority: expo.HighPriority,
			}
		} else {
			newMessage.To = append(newMessage.To, expo.ExponentPushToken(elem.ExpoToken))
		}
		prior = elem.AppointmentId
	}
	// log result
	fmt.Printf("Preparing Messages\n.\n.\n.\n")
	for _, m := range messagesPool {
		fmt.Printf("Message: %v\n", m.Body)
		fmt.Println("To:")
		for _, t := range m.To {
			fmt.Printf("\t%v\n", t)
		}
	}
	// 1 request can contain up to 100 messages, for safety purpose -> 1 request should contain only up to 80 messages
	// divide len([]message) with 80 -> split up to multiple request
	fmt.Printf("Splitting up messages to multiple request\n.\n.\n.\n")
	const MAX_MESSAGES_PER_REQUEST = 80
	var messageCnt = float64(len(messagesPool))
	var cnt float64 = math.Ceil(float64(messageCnt) / MAX_MESSAGES_PER_REQUEST)
	for i := 0; i < int(cnt); i++ {
		// calculate base and limit for slicing slice
		base := float64(i * MAX_MESSAGES_PER_REQUEST)
		limit := base + math.Min(messageCnt-base, MAX_MESSAGES_PER_REQUEST)
		fmt.Printf("request %v => messagesPool[%v:%v]\n", i, base, limit)

		// send request
		fmt.Printf("sending request %v\n", i)
		sendRequestFunc(messagesPool[int(base):int(limit)])
	}

}

func mockSendRequest(messages []expo.PushMessage) {

}

func sendRequest(messages []expo.PushMessage) {
	client := expo.NewPushClient(nil)
	// Publish message
	responses, err := client.PublishMultiple(messages)

	// Check errors
	if err != nil && !strings.Contains(err.Error(), "Mismatched response length"){
		panic(err)
	}
	fmt.Println("validate")
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
select appointment.id ,date, device.id, device.device_name , expo_token, appointment.patient_id from appointment 
inner join device on appointment.patient_id = device.patient_id
where device.expo_token != "" AND appointment.date > ?
order by appointment.id asc
`

func queryDB(db *sql.DB) ([]model.AppointmentDevice, error) {
	rows, err := db.Query(apmtQuery, time.Now().Unix())
	if err != nil {
		fmt.Println("queryDB : Can't query database")
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
	return res, nil
}

func formatTimeOutput(dueTimestamp int, nowTimestamp int) string {
	sec := (dueTimestamp - nowTimestamp)
	minute := sec / 60
	hour := minute / 60
	day := hour / 24
	baseStr := "You've got an appointment coming up in "
	var output string
	if minute == 0 {
		output = "several minutes"
	} else if hour == 0 {
		output = fmt.Sprintf("%d minutes", minute)
	} else if day == 0 {
		output = fmt.Sprintf("%d hour(s) %d minute(s)", hour, minute%60)
	} else {
		output = fmt.Sprintf("%d day(s) %d hour(s)", day, hour%24)
	}
	return baseStr + output
}