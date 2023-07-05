package dingding

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type DingTalkMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func encrypt() (int64, string) {
	timestamp := time.Now().UnixMilli()
	secret := "SEC03522b1c0c17e44cb0647ba3bc4ca61af47101408a0f0d97b44c4dd95094427a"

	stringToSign := gconv.String(timestamp) + "\n" + secret
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(stringToSign))
	sign := url.QueryEscape(base64.StdEncoding.EncodeToString(hash.Sum(nil)))
	return timestamp, sign
}

func sendDingTalkMessage(message string) error {
	t, s := encrypt()
	webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=ed3b5bc1200a7d67e0217146aac0f0a17eb90c9a0f5084237826b07918d16d00" +
		fmt.Sprintf("&timestamp=%d&sign=%s", t, s)
	dingTalkMessage := DingTalkMessage{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: fmt.Sprintf("%s", message),
		},
	}

	jsonBody, err := json.Marshal(dingTalkMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal DingTalk message: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to send DingTalk message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send DingTalk message: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Println(string(body))

	return nil
}
