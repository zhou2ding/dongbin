package captcha

import (
	"blog/pkg/v"
	"github.com/mojocn/base64Captcha"
)

func GetCaptcha() (string, string, error) {
	size := v.GetViper().GetIntSlice("captcha.size")
	dr := base64Captcha.NewDriverDigit(size[0], size[1], size[2], 0.7, size[3])
	capt := base64Captcha.NewCaptcha(dr, base64Captcha.DefaultMemStore)
	id, b64, err := capt.Generate()
	if err != nil {
		return "", "", err
	}
	return id, b64, nil
}
