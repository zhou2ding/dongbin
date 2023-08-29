package gb28181

import (
	"blog/pkg/l"
	"github.com/ghettovoice/gosip/sip"
	"github.com/gogf/gf/util/grand"
	"net/http"
	"time"
)

func onRegister(req sip.Request, tx sip.ServerTransaction) {
	from, ok := req.From()
	if from == nil || !ok || from.Address == nil || from.Address.User() == nil {
		l.Logger().Warningf("Server <- request nil error %s", req.String())
		return
	}

	for _, h := range req.Headers() {
		if h.Name() == "Expires" && h.Value() == "0" {
			resp := sip.NewResponseFromRequest("", req, http.StatusOK, "OK", "")
			to, _ := resp.To()
			resp.ReplaceHeaders("To", []sip.Header{
				&sip.ToHeader{
					Address: to.Address,
					Params:  sip.NewParams().Add("tag", sip.String{Str: grand.Str("1234567890", 9)}),
				},
			})
			resp.RemoveHeader("Allow")
			exp := sip.Expires(0)
			resp.AppendHeader(&exp)
			resp.AppendHeader(&sip.GenericHeader{
				HeaderName: "Date",
				Contents:   time.Now().Format("2006-01-02 15:04:05"),
			})
			if err := tx.Respond(resp); err != nil {
				l.Logger().Errorf("onRegister Respond error: %v", err)
			}
		}
	}
}

func onMessage(req sip.Request, tx sip.ServerTransaction) {

}

func onBye(req sip.Request, tx sip.ServerTransaction) {
	err := tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", ""))
	if err != nil {
		l.Logger().Errorf("onBye Respond error: %v", err)
	}
}
