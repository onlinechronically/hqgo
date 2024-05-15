package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Verification struct {
	verificationId string
}

func CreateVerification(phoneNumber string) (Verification, error) {
	body, err := json.Marshal(VerificationStart{Method: "sms", Phone: phoneNumber})
	if err != nil {
		return Verification{}, err
	}
	verificationResponse := request("POST", fmt.Sprintf("%s/verifications/verify-existing-phone", apiURL), body, "")
	if verificationResponse["error"] != nil {
		return Verification{}, errors.New(verificationResponse["error"].(string))
	}
	return Verification{verificationId: verificationResponse["verificationId"].(string)}, nil
}

func (verification *Verification) Verify(code string) (User, bool, error) {
	body, err := json.Marshal(VerificationFinish{Code: code})
	if err != nil {
		return User{}, false, err
	}
	verificationResponse := request("POST", fmt.Sprintf("%s/verifications/%s", apiURL, verification.verificationId), body, "")
	if verificationResponse["error"] != nil {
		return User{}, false, errors.New(verificationResponse["error"].(string))
	}
	authData := verificationResponse["auth"]
	if authData == nil {
		return User{}, false, nil
	} else {
		authData := authData.(map[string]interface{})
		return User{loginToken: authData["loginToken"].(string), bearerToken: authData["authToken"].(string), id: int(authData["userId"].(float64)), username: authData["username"].(string)}, true, nil
	}
}

func checkUsername(username string, isReferral bool) (bool, error) {
	var body []byte
	var err error
	var endpoint string
	if isReferral {
		body, err = json.Marshal(CheckReferral{ReferralCode: username})
		endpoint = "/referral-code/valid"
	} else {
		body, err = json.Marshal(CheckUsername{Username: username})
		endpoint = "/usernames/available"
	}
	if err != nil {
		return false, err
	}
	checkResponse := request("POST", fmt.Sprintf("%s%s", apiURL, endpoint), body, "")
	if checkResponse["error"] != nil {
		return false, nil
	}
	return true, nil
}
