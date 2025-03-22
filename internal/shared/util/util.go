package util

import (
	"crypto/rand"
	"encoding/json"
)

func GenerateRandomCode(length int) string {
	const charset = "0123456789"
	if length <= 0 {
		length = 6
	}

	// Create a byte slice to store the result
	result := make([]byte, length)

	// Use crypto/rand for better randomness
	rand.Read(result)

	// Map the random bytes to characters in the charset
	for i := range result {
		result[i] = charset[int(result[i])%len(charset)]
	}

	return string(result)
}

func StructToMap(obj any) (map[string]any, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	return result, err
}
