package plugin

import (
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"strconv"
	"syscall"

	"github.com/qiniu/httptest"
)

// ---------------------------------------------------------------------------

type authstubArgs struct {
	Uid     uint   `flag:"uid"`
	Utype   uint   `flag:"utype"`
	Sudoer  uint   `flag:"suid"`
	UtypeSu uint   `flag:"sut"`
	App     string `arg:"app,opt"`
}

func appendUint(form []byte, k string, v uint64) []byte {

	form = append(form, k...)
	return strconv.AppendUint(form, v, 10)
}

func formatAuthstub(user *authstubArgs, appid uint64) string {

	return "QiniuStub " + formatAuthstubToken(user, appid)
}

func formatAuthstubToken(user *authstubArgs, appid uint64) string {

	form := make([]byte, 0, 128)
	form = appendUint(form, "uid=", uint64(user.Uid))
	form = appendUint(form, "&ut=", uint64(user.Utype))
	if appid != 0 {
		form = appendUint(form, "&app=", uint64(appid))
	}
	if user.Sudoer != 0 {
		form = appendUint(form, "&suid=", uint64(user.Sudoer))
		if user.UtypeSu != 0 {
			form = appendUint(form, "&sut=", uint64(user.UtypeSu))
		}
	}
	return string(form)
}

// ---------------------------------------------------------------------------

type authstubTransport struct {
	auth      string
	Transport http.RoundTripper
}

func (t *authstubTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	req.Header.Set("Authorization", t.auth)
	return t.Transport.RoundTrip(req)
}

func (t *authstubTransport) NestedObject() interface{} {

	return t.Transport
}

func authstubNewTransport(
	user *authstubArgs, appid uint64, transport http.RoundTripper) *authstubTransport {

	if transport == nil {
		transport = http.DefaultTransport
	}
	return &authstubTransport{formatAuthstub(user, appid), transport}
}

// ---------------------------------------------------------------------------

type authstubTransportComposer struct {
	args *authstubArgs
	ctx  *httptest.Context
}

func toAppId(app string) (appId uint64, err error) {

	b, err := base64.URLEncoding.DecodeString(app)
	if err != nil {
		return
	}
	if len(b) != 12 {
		return 0, syscall.EINVAL
	}
	return binary.LittleEndian.Uint64(b[4:]), nil
}

func (p *authstubTransportComposer) Compose(base http.RoundTripper) http.RoundTripper {

	var appid uint64
	var err error

	if p.args.App != "" {
		appid, err = strconv.ParseUint(p.args.App, 10, 64)
		if err != nil {
			appid, err = toAppId(p.args.App)
			if err != nil {
				p.ctx.Fatal("Parse arg `app` failed:", err)
			}
		}
	}
	return authstubNewTransport(p.args, appid, base)
}

func (p *subContext) Eval_authstub(
	ctx *httptest.Context, args *authstubArgs) (httptest.TransportComposer, error) {

	return &authstubTransportComposer{args, ctx}, nil
}

// ---------------------------------------------------------------------------
