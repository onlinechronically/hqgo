package hqtrivia

import (
	"encoding/json"
	"errors"
	"fmt"
)

type User struct {
	loginToken  string
	bearerToken string
	id          int
	username    string
}

func (u *User) request(method string, endpoint string, body []byte) (map[string]interface{}, int, error) {
	if u.bearerToken == "" {
		return nil, -1, errors.New("You must be logged in to do this")
	}
	return request(method, endpoint, body, u.bearerToken)
}

func (u *User) RefreshTokens() error {
	body, err := json.Marshal(RefreshTokens{Token: u.loginToken})
	if err != nil {
		return err
	}
	tokenReq, _, err := request("POST", fmt.Sprintf("%s/tokens", apiURL), body, "")
	if err != nil {
		return err
	}
	if tokenReq["error"] != nil || tokenReq == nil {
		return errors.New(tokenReq["error"].(string))
	}
	u.bearerToken = tokenReq["authToken"].(string)
	u.username = tokenReq["username"].(string)
	return nil
}

func (u *User) Powerups() (int, int, int, error) {
	usersMe, _, err := u.request("GET", fmt.Sprintf("%s/users/me", apiURL), nil)
	if err != nil {
		return -1, -1, -1, err
	}
	if usersMe["error"] != nil {
		return -1, -1, -1, errors.New(usersMe["error"].(string))
	}
	items := usersMe["items"].(map[string]interface{})
	if items == nil {
		return -1, -1, -1, nil
	}
	return int(items["lives"].(float64)), int(usersMe["erase1s"].(float64)), int(items["superSpins"].(float64)), nil
}

func (u *User) ChangeUsername(username string) error {
	isValid, err := checkUsername(username, false)
	if err != nil {
		return err
	} else if !isValid {
		return errors.New("Invalid Username")
	}
	body, err := json.Marshal(ChangeUsername{Username: username})
	if err != nil {
		return err
	}
	resp, _, err := u.request("PATCH", fmt.Sprintf("%s/users/me", apiURL), body)
	if err != nil {
		return err
	}
	if resp["error"] != nil {
		return errors.New(resp["error"].(string))
	}
	return nil
}

func (u *User) EasterEgg() error {
	resp, statusCode, err := u.request("POST", fmt.Sprintf("%s/easter-eggs/makeItRain", apiURL), nil)
	if err != nil {
		return err
	} else if statusCode != 200 {
		return errors.New("Too Recent")
	} else if resp["error"] != nil {
		return errors.New(resp["error"].(string))
	} else {
		return nil
	}
}
