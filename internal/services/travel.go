package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

func SetDayOfWeek(date *time.Time) int {
	var targetDate time.Time

	if date == nil || date.IsZero() {
		targetDate = time.Now()
	} else {
		targetDate = *date
	}

	dayOfWeek := int(targetDate.Weekday())
	dayOfWeek = (dayOfWeek + 1) % 7

	return dayOfWeek
}

func EncodeSearchID(scheduleID uint, busID uint) string {
	randomBytes := make([]byte, 6)
	timeStamp := time.Now()
	encodedIDLength := 20
	_, _ = rand.Read(randomBytes)

	raw := fmt.Sprintf("%d:%d:%x:%s", scheduleID, busID, randomBytes, timeStamp)
	encoded := base64.RawURLEncoding.EncodeToString([]byte(raw))

	if len(encoded) > encodedIDLength {
		return encoded[:encodedIDLength]
	}
	for len(encoded) < encodedIDLength {
		encoded += "X" // padding
	}
	return encoded
}

func DecodeSearchID(searchID string) (scheduleID uint, busID uint, err error) {
	decodedBytes, err := base64.RawURLEncoding.DecodeString(searchID)
	if err != nil {
		return 0, 0, err
	}

	var randomHex string
	_, err = fmt.Sscanf(string(decodedBytes), "%d:%d:%s", &scheduleID, &busID, &randomHex)
	if err != nil {
		return 0, 0, err
	}

	return scheduleID, busID, nil
}
