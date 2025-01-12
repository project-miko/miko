package strutils

import uuid "github.com/satori/go.uuid"

func GetUUID() string {
	return uuid.NewV4().String()
}

func UUIDToInt(u string) (int, error) {
	_, err := uuid.FromString(u) // check
	if err != nil {
		return 0, err
	}

	var sum int = 0
	for i := 0; i < len(u); i++ {
		sum = sum + int(u[i])
	}
	return sum, nil
}
