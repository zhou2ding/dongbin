package captcha

import (
	"blog/pkg/v"
	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	ID      string
	Content string
}

func GetCaptcha() (*Captcha, error) {
	size := v.GetViper().GetIntSlice("captcha.size")
	dr := base64Captcha.NewDriverDigit(size[0], size[1], size[2], 0.7, size[3])
	capt := base64Captcha.NewCaptcha(dr, base64Captcha.DefaultMemStore)
	id, b64, err := capt.Generate()
	if err != nil {
		return nil, err
	}
	return &Captcha{id, b64}, nil
}

func (c *Captcha) Ok() bool {
	return base64Captcha.DefaultMemStore.Verify(c.ID, c.Content, true)
}
