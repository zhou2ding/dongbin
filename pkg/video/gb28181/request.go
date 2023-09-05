package gb28181

import (
	"blog/pkg/v"
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/ghettovoice/gosip/sip"
	"net/http"
)

func CreateRequest(Method sip.RequestMethod) (req sip.Request) {
	callId := sip.CallID(gofakeit.StreetNumber())
	userAgent := sip.UserAgentHeader("Monibuca")
	cseq := sip.CSeq{
		MethodName: Method,
	}
	port := sip.Port(v.GetViper().GetUint("sipServer.sipPort"))
	serverAddr := sip.Address{
		Uri: &sip.SipUri{
			FUser: sip.String{Str: "34020000002000000001"}, // gb28181 id
			FPort: &port,
		},
		Params: sip.NewParams().Add("tag", sip.String{Str: gofakeit.StreetNumber()}),
	}
	addr := sip.Address{}
	req = sip.NewRequest(
		"",
		Method,
		addr.Uri,
		"SIP/2.0",
		[]sip.Header{
			serverAddr.AsFromHeader(),
			&callId,
			&userAgent,
			&cseq,
			serverAddr.AsContactHeader(),
		},
		"",
		nil,
	)

	req.SetTransport("")
	req.SetDestination("")
	return
}

func Catalog() int {
	request := CreateRequest(sip.MESSAGE)
	expires := sip.Expires(3600)
	contentType := sip.ContentType("Application/MANSCDP+xml")

	request.AppendHeader(&contentType)
	request.AppendHeader(&expires)
	// 输出Sip请求设备通道信息信令
	resp, err := SipRequestForResponse(request)
	if err == nil && resp != nil {
		return int(resp.StatusCode())
	}
	return http.StatusRequestTimeout
}

func SipRequestForResponse(request sip.Request) (sip.Response, error) {
	return gServer.RequestWithContext(context.Background(), request)
}
