package helper

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type UserIDType string

const UserIDKey UserIDType = "userId"

type ResponseMessage struct {
	Message string `json:"messsage"`
}

func ConvertStringToInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		log.Printf("Error converting string to int64: %v", err)
		return 0, err
	}

	return i, nil
}

func WriteResponseMessage(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	responseMessage := ResponseMessage{Message: message}
	json.NewEncoder(w).Encode(responseMessage)
}

func WriteResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func DecodeRequestBody(r *http.Request, data interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		return err
	}

	return nil
}
