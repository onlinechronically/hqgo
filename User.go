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

type Payout struct{}

type Opt struct {
	title       string `json:"title"`
	opt         string `json:"opt"`
	inMessage   string `json:"in"`
	outMessage  string `json:"out"`
	description string `json:"onboardingDescription"`
	opted       bool   `json:"opted"`
}

func (u *User) request(method string, endpoint string, body []byte) (map[string]interface{}, int, error) {
	if u.bearerToken == "" {
		return nil, -1, errors.New("You must be logged in to do this")
	}
	return request(method, endpoint, body, u.bearerToken)
}

func (u *User) UsersMe() (map[string]interface{}, error) {
	usersMe, _, err := u.request("GET", fmt.Sprintf("%s/users/me", apiURL), nil)
	if err != nil {
		return nil, err
	} else if usersMe["error"] != nil {
		return nil, errors.New(usersMe["error"].(string))
	}
	if usersMe["username"] != u.username {
		u.username = usersMe["username"].(string)
	}
	return usersMe, err
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
	usersMe, err := u.UsersMe()
	if err != nil {
		return -1, -1, -1, err
	}
	items := usersMe["items"].(map[string]interface{})
	if items == nil {
		return -1, -1, -1, nil
	}
	return int(items["lives"].(float64)), int(usersMe["erase1s"].(float64)), int(items["superSpins"].(float64)), nil
}

func (u *User) Coins() (int, error) {
	usersMe, err := u.UsersMe()
	if err != nil {
		return -1, err
	}
	return int(usersMe["coins"].(float64)), nil
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

func (u *User) Achievements() (map[string]interface{}, error) {
	resp, _, err := u.request("GET", fmt.Sprintf("%s/achievements/v2/me", apiURL), nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *User) Config(public bool) (map[string]interface{}, error) {
	var url string
	if public {
		url = fmt.Sprintf("%s/config/public", apiURL)
	} else {
		url = fmt.Sprintf("%s/config", apiURL)
	}
	resp, _, err := u.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *User) Opts() (opts []Opt, _ error) {
	resp, _, err := u.request("GET", fmt.Sprintf("%s/opt-in", apiURL), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println(resp)
	for _, opt := range resp["opts"].([]interface{}) {
		fmt.Println(opt)
		opt, isValid := opt.(map[string]interface{})
		if isValid {
			return opts, errors.New("Invalid Data")
		}
		newOpt := Opt{title: opt["title"].(string), opt: opt["opt"].(string), inMessage: opt["in"].(string), outMessage: opt["out"].(string), description: opt["onboardingDescription"].(string), opted: opt["opted"].(bool)}
		opts = append(opts, newOpt)
	}
	return opts, nil
}

func (u *User) OptIn() (map[string]interface{}, error) {
	resp, _, err := u.request("POST", fmt.Sprintf("%s/opt-in", apiURL), nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (u *User) Sample() (map[string]interface{}, error) {
	resp, _, err := u.request("GET", fmt.Sprintf("%s/opt-in", apiURL), nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
