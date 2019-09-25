package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// COGS consists of login credentials for the portal
type COGS struct {
	username, password string
}

func (c COGS) accessToken() (string, error) {
	body, err := json.Marshal(struct {
		Data interface{} `json:"data"`
	}{
		Data: struct {
			Type       string      `json:"type"`
			Attributes interface{} `json:"attributes"`
		}{
			Type: "auths",
			Attributes: struct {
				Username string `json:"userName"`
				Password string `json:"password"`
			}{
				Username: c.username,
				Password: c.password,
			},
		},
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post("https://cogs.10pearls.com/cogsapi/api/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		// too lazy to decode error response
		return "", fmt.Errorf("probably invalid username and/or password")
	}

	buf, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	var respJSON struct {
		Data struct {
			Attributes struct {
				AccessToken string `json:"access-token"`
			}
		}
	}
	if err := json.Unmarshal(buf, &respJSON); err != nil {
		return "", err
	}

	return respJSON.Data.Attributes.AccessToken, nil
}

// Lunch returns list of strings for lunch menu today
func (c COGS) Lunch() ([]string, error) {
	token, err := c.accessToken()
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://cogs.10pearls.com/cogsapi/api/Lunches/Weekly", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("unauthorized access")
	}

	buf, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var errJSON struct {
			Message string
		}

		if err := json.Unmarshal(buf, &errJSON); err != nil {
			return nil, fmt.Errorf("unknown error")
		}

		return nil, fmt.Errorf(errJSON.Message)
	}

	var respJSON struct {
		Data []struct {
			Attributes struct {
				Date      string `json:"lunch-date"`
				Menu      string `json:"menu-item"`
				LunchType string `json:"lunch-type"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(buf, &respJSON); err != nil {
		return nil, err
	}

	location, err := time.LoadLocation("Asia/Karachi")
	if err != nil {
		return nil, err
	}

	complete_menu := []string{}
	for _, v := range respJSON.Data {
		t, err := time.ParseInLocation("2006-01-02T15:04:05", v.Attributes.Date, location)
		if err != nil {
			return nil, fmt.Errorf("invalid date %s", v.Attributes.Date)
		}

		if time.Now().Year() == t.Year() && time.Now().YearDay() == t.YearDay() {
			menu := strings.Split(v.Attributes.Menu, ",")

			// clean up spaces and lines
			for k, v := range menu {
				menu[k] = strings.Trim(v, " \n")
			}

			if v.Attributes.LunchType == "R" {
				menu = append([]string{"*Regular Lunch*"}, menu...)
			} else if v.Attributes.LunchType == "L" {
				menu = append([]string{"*Low Calorie Lunch*"}, menu...)
			} else {
				menu = append([]string{"*SPECIAL LUNCH*"}, menu...)
			}

			complete_menu = append(complete_menu, menu...)
		}
	}

	if len(complete_menu) == 0 {
		return nil, fmt.Errorf("no lunch found for today")
	}

	return complete_menu, nil
}
