package gb28181

import (
	"blog/pkg/l"
	"github.com/ghettovoice/gosip/sip"
	"net/http"
)

func onBye(req sip.Request, tx sip.ServerTransaction) {
	err := tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", ""))
	if err != nil {
		l.Logger().Errorf("onBye Respond error: %v", err)
	}
}
