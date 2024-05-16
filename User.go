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

func (u *User) request(method string, endpoint string, body []byte) (map[string]interface{}, int) {
	if u.bearerToken == "" {
		fatalError(errors.New("You must be logged in to do this"), "")
	}
	return request(method, endpoint, body, u.bearerToken)
}

func (u *User) RefreshTokens() error {
	body, err := json.Marshal(RefreshTokens{Token: u.loginToken})
	if err != nil {
		return err
	}
	tokenReq, _ := request("POST", fmt.Sprintf("%s/tokens", apiURL), body, "")
	if tokenReq["error"] != nil || tokenReq == nil {
		return errors.New(tokenReq["error"].(string))
	}
	u.bearerToken = tokenReq["authToken"].(string)
	u.username = tokenReq["username"].(string)
	return nil
}

func (u *User) Powerups() (int, int, int, error) {
	usersMe, _ := u.request("GET", fmt.Sprintf("%s/users/me", apiURL), nil)
	if usersMe["error"] != nil {
		return -1, -1, -1, errors.New(usersMe["error"].(string))
	}
	items := usersMe["items"].(map[string]interface{})
	if items == nil {
		return -1, -1, -1, nil
	}
	return int(items["lives"].(float64)), int(usersMe["erase1s"].(float64)), int(items["superSpins"].(float64)), nil
}
