package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var _ Notifier = &PushBullet{}

type PushBullet struct {
	AccessToken string
}

func (n *PushBullet) Notify(title, message string) error {
	client := http.Client{}

	body, err := json.Marshal(struct {
		Type  string `json:"type"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		Type:  "note",
		Title: title,
		Body:  message,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.pushbullet.com/v2/pushes", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Access-Token", n.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("unable to read body")
	}

	if resp.StatusCode != 200 {
		var errJSON struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.Unmarshal(buf, &errJSON); err != nil {
			return fmt.Errorf("unknown error")
		}

		return fmt.Errorf(errJSON.Error.Message)
	}

	return nil
}
