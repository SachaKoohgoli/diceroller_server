package helpers

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

type TimeBasedAccessToken struct {
	Timestamp  time.Time
	EncodedVal string
}

// GenerateToken generates a new TimeBasedAccessToken anchored on the passed in time
// If I was writing tests, I would pass in times to make sure things were encoded correctly
func GenerateToken(refTime time.Time) TimeBasedAccessToken {
	timestamp := refTime.Unix()

	// Encode the timestamp using base64 encoding
	encodedToken := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", timestamp)))

	return TimeBasedAccessToken{
		Timestamp:  refTime,
		EncodedVal: encodedToken,
	}
}

// DecodeToken decodes a base64 string into a TimeBasedAccessToken
// If I was writing tests, I would pass in times to make sure things were encoded correctly
func DecodeToken(encodedToken string) (TimeBasedAccessToken, error) {
	// Decode the base64 encoded string
	decodedBytes, err := base64.URLEncoding.DecodeString(encodedToken)
	if err != nil {
		fmt.Println("Error decoding token:", err)
		return TimeBasedAccessToken{}, err
	}

	// Convert the string back to an integer (Unix timestamp)
	timestamp, err := strconv.ParseInt(string(decodedBytes), 10, 64)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return TimeBasedAccessToken{}, err
	}

	// Return the time.Time object representing the timestamp
	return TimeBasedAccessToken{Timestamp: time.Unix(timestamp, 0), EncodedVal: encodedToken}, nil
}

// IsValid determines if the current token is within the time window of the reference time
func (t TimeBasedAccessToken) IsValid(referenceTime time.Time, timeWindowMinutes time.Duration) bool {
	return referenceTime.Add(timeWindowMinutes).After(t.Timestamp)
}
