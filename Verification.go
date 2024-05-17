package hqtrivia

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Verification struct {
	verificationId string
}

func (v *Verification) GetID() string {
	return v.verificationId
}

func CreateVerification(phoneNumber string) (Verification, error) {
	body, err := json.Marshal(VerificationStart{Method: "sms", Phone: phoneNumber})
	if err != nil {
		return Verification{}, err
	}
	verificationResponse, _, err := request("POST", fmt.Sprintf("%s/verifications/verify-existing-phone", apiURL), body, "")
	if err != nil {
		return Verification{}, err
	}
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
	verificationResponse, _, err := request("POST", fmt.Sprintf("%s/verifications/%s", apiURL, verification.verificationId), body, "")
	if err != nil {
		return User{}, false, err
	}
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
	checkResponse, _, err := request("POST", fmt.Sprintf("%s%s", apiURL, endpoint), body, "")
	if err != nil {
		return false, err
	}
	if checkResponse["error"] != nil {
		return false, nil
	}
	return true, nil
}

func (verification *Verification) RegisterUser(username string, referringUsername string) (User, error) {
	var referralValid bool
	var err error
	usernameAvailable, err := checkUsername(username, false)
	if err != nil {
		return User{}, err
	}
	if referringUsername != "" {
		referralValid, err = checkUsername(referringUsername, true)
		if err != nil {
			return User{}, err
		}
	}
	if !usernameAvailable {
		return User{}, errors.New("The requested username is not available")
	}
	if !referralValid && referringUsername != "" {
		return User{}, errors.New("The provided referral code is not valid")
	}
	var bodyStruct *RegisterAccount
	if referringUsername == "" {
		bodyStruct = &RegisterAccount{VerificationID: verification.verificationId, Username: username, Country: "us", Locale: "en", TimeZone: "America/New_York"}
	} else {
		bodyStruct = &RegisterAccount{VerificationID: verification.verificationId, Username: username, ReferringUsername: referringUsername, Country: "us", Locale: "en", TimeZone: "America/New_York"}
	}
	registerBody, err := json.Marshal(bodyStruct)
	if err != nil {
		return User{}, err
	}
	registerRes, statusCode, err := request("POST", fmt.Sprintf("%s/users", apiURL), registerBody, "")
	if err != nil {
		return User{}, err
	}
	if statusCode != 200 {
		return User{}, errors.New("There was an error registering")
	}
	return User{loginToken: registerRes["loginToken"].(string), bearerToken: registerRes["authToken"].(string), id: int(registerRes["userId"].(float64)), username: registerRes["username"].(string)}, nil
}
