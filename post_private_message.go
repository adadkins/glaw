package glaw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func (lc *LemmyClient) SendPrivateMessage(messageBody string, targetUserID int) error {
	auth := ""
	// Set the API token for authentication (if required)
	if lc.APIToken != "" {
		auth = lc.APIToken
	}
	if lc.jwtCookie != "" {
		test := strings.Split(lc.jwtCookie, "jwt=")
		if len(test) != 2 {
			return fmt.Errorf("JWT is not correctly formatted")
		}
		auth = test[1]
	}
	pmStruct := struct {
		Recipient_id int    `json:"recipient_id"`
		Content      string `json:"content"`
		Auth         string `json:"auth"`
	}{
		Recipient_id: targetUserID,
		Content:      messageBody,
		Auth:         auth,
	}

	jsonData, err := json.Marshal(pmStruct)
	if err != nil {
		lc.logger.Info(err.Error())
		return err
	}

	lc.logger.Sugar().Infoln(string(jsonData))
	_, err = lc.callLemmyAPI("POST", "private_message", bytes.NewBuffer(jsonData))
	if err != nil {
		lc.logger.Info(err.Error())
		return err
	}

	return nil
}
