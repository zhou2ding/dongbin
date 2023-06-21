package rpcserver

import (
	"blog/pkg/internal/rpcencrypt"
	"blog/pkg/l"
	"blog/pkg/rand"
	"blog/pkg/rpcservice"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

const (
	NotLogin = iota
	FirstLogin
	SecondLogin
	ClientLogout
)

type rpcUser struct {
	user         string
	rand         string
	accessStatus atomicStatus
	checker      UserChecker
}

type loginValue struct {
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Random   string `json:"random,omitempty"`
}

func newRpcUser(checker UserChecker) *rpcUser {
	return &rpcUser{
		checker: checker,
	}
}

func (u *rpcUser) getStatus() int {
	return u.accessStatus.GetValue()
}

func (u *rpcUser) setStatus(status int) {
	u.accessStatus.SetValue(status)
}

func (u *rpcUser) login(jsonData []byte) ([]byte, int, error) {
	var req jsonRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		l.Logger().Error("login Unmarshal error", zap.Error(err))
		return nil, 0, err
	}

	var (
		id       int
		value    []byte
		loginVal loginValue
	)
	resp := jsonResponse{
		Domain: req.Domain,
		Id:     req.Id,
		Key:    req.Key,
	}
	if req.Key != "System.login" {
		resp.RetVal = rpcservice.ErrInvalidLoginRequest
		resp.ErrMsg = "Invalid Login Request"
	} else {
		if u.accessStatus.GetValue() == NotLogin {
			if len(req.Value) == 0 || json.Unmarshal(req.Value, &loginVal) != nil || len(loginVal.User) == 0 {
				l.Logger().Error("rpc request value is wrong")
				resp.RetVal = rpcservice.ErrInvalidLoginRequest
				resp.ErrMsg = rpcservice.StatusText(rpcservice.ErrInvalidLoginRequest)
			} else {
				var err error
				l.Logger().Info("first login", zap.String("user", loginVal.User))
				id, err = u.checker.Identify(loginVal.User)
				if err != nil {
					l.Logger().Error("first login failed", zap.Error(err))
					resp.RetVal = rpcservice.ErrInvalidUserName
					resp.ErrMsg = rpcservice.StatusText(rpcservice.ErrInvalidUserName)
				} else {
					random := rand.GetRandGeneratorInstance().GetRandomString(8)
					value = []byte(fmt.Sprintf("{\"user\":\"%s\",\"random\":\"%s\"}", loginVal.User, random))
					resp.RetVal = 0
					u.user = loginVal.User
					u.rand = random
					u.accessStatus.SetValue(FirstLogin)
				}
			}
		} else if u.accessStatus.GetValue() == FirstLogin {
			if len(req.Value) == 0 || json.Unmarshal(req.Value, &loginVal) != nil || len(loginVal.User) == 0 {
				l.Logger().Error("rpc request value is wrong")
				resp.RetVal = rpcservice.ErrInvalidLoginRequest
				resp.ErrMsg = rpcservice.StatusText(rpcservice.ErrInvalidLoginRequest)
				u.accessStatus.SetValue(NotLogin)
			} else {
				l.Logger().Info("second login", zap.String("user", loginVal.User))
				loginOk := false
				if loginVal.User == u.user && loginVal.Random == u.rand {
					pwd, err := u.checker.GetToken(loginVal.User)
					if err == nil {
						encPwd := rpcencrypt.GetEncryptHelperInstance().Encrypt(loginVal.User, pwd, loginVal.Random)
						if encPwd == loginVal.Password {
							loginOk = true
						}
					}
				}
				if loginOk {
					value = []byte(fmt.Sprintf("{\"user\":\"%s\"}", u.user))
					resp.RetVal = 0
					u.accessStatus.SetValue(SecondLogin)
				} else {
					resp.RetVal = rpcservice.ErrInvalidLoginRequest
					resp.ErrMsg = rpcservice.StatusText(rpcservice.ErrInvalidLoginRequest)
					u.accessStatus.SetValue(NotLogin)
				}
			}
		}
	}
	if len(value) != 0 {
		resp.Value = value
	}
	l.Logger().Debug("doLogin done", zap.Any("login response", resp))
	jsonResp, _ := json.Marshal(resp)
	return jsonResp, id, nil
}
