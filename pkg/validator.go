package api

import (
	"strconv"
	"strings"
)

type HashRequest struct {
	Password string
	Errors   map[string]string
}

type HashGetRequest struct {
	Errors map[string]string
}

func ValidateHashRequest(password string) (*HashRequest, bool) {

	msg := &HashRequest{
		Password: password,
	}

	msg.Errors = make(map[string]string)

	if strings.TrimSpace(msg.Password) == "" {
		msg.Errors["Password"] = "Please enter a valid password"
		return msg, false
	}

	return msg, true
}

func ValidateHashGetRequest(hashId string) (int, *HashGetRequest, bool) {

	msg := &HashGetRequest{}

	msg.Errors = make(map[string]string)

	id, err := strconv.Atoi(hashId)

	if err != nil {
		msg.Errors["Id"] = "Please enter a valid password id"
		return 0, msg, false
	}

	return id, msg, true
}
