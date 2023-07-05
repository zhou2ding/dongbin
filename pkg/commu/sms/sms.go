package sms

import (
	"blog/pkg/encrypt"
	"blog/pkg/l"
	"blog/pkg/v"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func SendSms(mobile, content string) {
	reqUri := url.Values{}
	now := strconv.FormatInt(time.Now().Unix(), 10)
	account := v.GetViper().GetString("sms.account")
	password := v.GetViper().GetString("sms.password")
	reqUri.Set("account", account)
	reqUri.Set("password", encrypt.GetMd5String(account+password+mobile+content+now))
	reqUri.Set("mobile", mobile)
	reqUri.Set("content", content)
	reqUri.Set("time", now)
	body := strings.NewReader(reqUri.Encode())
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://106.ihuyi.com/webservice/sms.php?method=Submit&format=json", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, err := client.Do(req)
	if err != nil {
		l.Logger().Errorf("send sms error: %v", err)
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	l.Logger().Infof("send sms done, resp: %s", string(data))
}
