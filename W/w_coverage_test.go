package W

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/OneOfOne/cmap"
	"github.com/jordan-wright/email"
	"github.com/kokizzu/gotro/M"
	"github.com/valyala/fasthttp"
)

type globalsSnapshot struct {
	mailers   map[string]*SmtpConfig
	webmaster M.SS
	sessions  SessionConnector
	globals   SessionConnector
	routes    map[string]Action
	assets    [][2]string
	filters   []Action
	sessKey   string
	expireSec int64
	renewSec  int64
}

func snapshotGlobals() globalsSnapshot {
	return globalsSnapshot{
		mailers:   Mailers,
		webmaster: Webmasters,
		sessions:  Sessions,
		globals:   Globals,
		routes:    Routes,
		assets:    Assets,
		filters:   Filters,
		sessKey:   SESS_KEY,
		expireSec: EXPIRE_SEC,
		renewSec:  RENEW_SEC,
	}
}

func restoreGlobals(s globalsSnapshot) {
	Mailers = s.mailers
	Webmasters = s.webmaster
	Sessions = s.sessions
	Globals = s.globals
	Routes = s.routes
	Assets = s.assets
	Filters = s.filters
	SESS_KEY = s.sessKey
	EXPIRE_SEC = s.expireSec
	RENEW_SEC = s.renewSec
}

type mockSessionConnector struct {
	str    map[string]string
	ints   map[string]int64
	msx    map[string]M.SX
	expiry map[string]int64
	lists  map[string][]string
}

func newMockSessionConnector() *mockSessionConnector {
	return &mockSessionConnector{
		str:    map[string]string{},
		ints:   map[string]int64{},
		msx:    map[string]M.SX{},
		expiry: map[string]int64{},
		lists:  map[string][]string{},
	}
}

func cloneSX(v M.SX) M.SX {
	if v == nil {
		return M.SX{}
	}
	res := M.SX{}
	for k, vv := range v {
		res[k] = vv
	}
	return res
}

func (m *mockSessionConnector) Del(key string) {
	delete(m.str, key)
	delete(m.ints, key)
	delete(m.msx, key)
	delete(m.expiry, key)
}
func (m *mockSessionConnector) Expiry(key string) int64 {
	if v, ok := m.expiry[key]; ok {
		return v
	}
	return 0
}
func (m *mockSessionConnector) FadeStr(key, val string, ttl int64) {
	m.str[key] = val
	m.expiry[key] = ttl
}
func (m *mockSessionConnector) FadeInt(key string, val int64, ttl int64) {
	m.ints[key] = val
	m.expiry[key] = ttl
}
func (m *mockSessionConnector) FadeMSX(key string, val M.SX, ttl int64) {
	m.msx[key] = cloneSX(val)
	m.expiry[key] = ttl
}
func (m *mockSessionConnector) GetStr(key string) string { return m.str[key] }
func (m *mockSessionConnector) GetInt(key string) int64  { return m.ints[key] }
func (m *mockSessionConnector) GetMSX(key string) M.SX   { return cloneSX(m.msx[key]) }
func (m *mockSessionConnector) Inc(key string) int64     { m.ints[key]++; return m.ints[key] }
func (m *mockSessionConnector) Dec(key string) int64     { m.ints[key]--; return m.ints[key] }
func (m *mockSessionConnector) SetStr(key, val string)   { m.str[key] = val }
func (m *mockSessionConnector) SetInt(key string, val int64) {
	m.ints[key] = val
}
func (m *mockSessionConnector) SetMSX(key string, val M.SX) {
	m.msx[key] = cloneSX(val)
}
func (m *mockSessionConnector) SetMSS(key string, val M.SS) {
	res := M.SX{}
	for k, v := range val {
		res[k] = v
	}
	m.msx[key] = res
}
func (m *mockSessionConnector) Product() string { return "mock" }
func (m *mockSessionConnector) Lpush(key string, val string) {
	m.lists[key] = append([]string{val}, m.lists[key]...)
}
func (m *mockSessionConnector) Rpush(key string, val string) {
	m.lists[key] = append(m.lists[key], val)
}
func (m *mockSessionConnector) Lrange(key string, first, last int64) []string {
	src := m.lists[key]
	if len(src) == 0 {
		return []string{}
	}
	if first < 0 {
		first = 0
	}
	if last < first {
		return []string{}
	}
	if int(last) >= len(src) {
		last = int64(len(src) - 1)
	}
	res := make([]string, 0, last-first+1)
	for i := first; i <= last; i++ {
		res = append(res, src[i])
	}
	return res
}

func newTestRequestCtx(method, uri, body, contentType string) *fasthttp.RequestCtx {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.Header.SetUserAgent("gotro-test-agent")
	req.SetRequestURI(uri)
	req.SetBodyString(body)
	if contentType != "" {
		req.Header.SetContentType(contentType)
	}
	ctx := &fasthttp.RequestCtx{}
	ctx.Init(req, &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 6789}, nil)
	return ctx
}

func createBaseDir(t *testing.T) string {
	t.Helper()
	baseDir := filepath.ToSlash(t.TempDir()) + `/`
	if err := os.MkdirAll(baseDir+VIEWS_SUBDIR, DEFAULT_FILEDIR_PERM); err != nil {
		t.Fatalf("mkdir views: %v", err)
	}
	if err := os.MkdirAll(baseDir+PUBLIC_SUBDIR+`lib/`, DEFAULT_FILEDIR_PERM); err != nil {
		t.Fatalf("mkdir public/lib: %v", err)
	}
	layout := `layout|#{title}|#{project_name}|#{assets}|#{is_superadmin}|#{debug_mode}|#{contents}`
	errPage := `error|#{error_code}|#{error_title}|#{error_detail}|#{requested_path}|#{project_name}|#{webmaster}`
	partial := `hello #{name}`
	if err := os.WriteFile(baseDir+VIEWS_SUBDIR+`layout.html`, []byte(layout), DEFAULT_FILEDIR_PERM); err != nil {
		t.Fatalf("write layout: %v", err)
	}
	if err := os.WriteFile(baseDir+VIEWS_SUBDIR+`error.html`, []byte(errPage), DEFAULT_FILEDIR_PERM); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if err := os.WriteFile(baseDir+VIEWS_SUBDIR+`partial.html`, []byte(partial), DEFAULT_FILEDIR_PERM); err != nil {
		t.Fatalf("write partial: %v", err)
	}
	return baseDir
}

func newTestEngine(t *testing.T, baseDir string, debug bool) *Engine {
	t.Helper()
	logPath := filepath.ToSlash(t.TempDir()) + `/w.log`
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, DEFAULT_FILEDIR_PERM)
	if err != nil {
		t.Fatalf("open logger: %v", err)
	}
	t.Cleanup(func() { _ = f.Close() })
	return &Engine{
		DebugMode:       debug,
		Name:            `WTest`,
		BaseDir:         baseDir,
		PublicDir:       baseDir + PUBLIC_SUBDIR,
		ViewCache:       cmap.New(),
		Router:          nil,
		Logger:          f,
		LogPath:         logPath,
		Assets:          `<script src="/x.js"></script>`,
		WebMasterAnchor: `wm`,
	}
}

func TestAjaxAndGlobalsCoverage(t *testing.T) {
	snap := snapshotGlobals()
	t.Cleanup(func() { restoreGlobals(snap) })

	if MIME2EXT[`video/mp4`] != `mp4` {
		t.Fatalf("globals init should merge youtube mime")
	}
	if len(YOUTUBE_MIME_LIST) == 0 {
		t.Fatalf("YOUTUBE_MIME_LIST should not be empty")
	}

	j := NewAjax()
	if j.HasError() {
		t.Fatalf("new ajax should have no errors")
	}
	j.Info("a")
	j.Info("b")
	if !strings.Contains(j.SX[`info`].(string), "a") || !strings.Contains(j.SX[`info`].(string), "b") {
		t.Fatalf("Info should append")
	}
	if j.Error("") != "" {
		t.Fatalf("Error empty should return empty")
	}
	j.Error("e1")
	if !j.HasError() || j.LastError() != "e1" {
		t.Fatalf("Error/LastError mismatch: %#v", j.SX)
	}
	if j.ErrorIf(nil, "noop") {
		t.Fatalf("ErrorIf(nil) should be false")
	}
	if !j.ErrorIf(errors.New("x"), "e2") || j.LastError() != "e2" {
		t.Fatalf("ErrorIf(err) should append")
	}
	j.OverwriteInfo("reset")
	if j.SX[`info`] != "reset" {
		t.Fatalf("OverwriteInfo mismatch")
	}
	j.ClearErrors()
	if j.HasError() || j.LastError() != "" {
		t.Fatalf("ClearErrors mismatch: %#v", j.SX[`errors`])
	}
	if j.TestError(nil, "ignored") {
		t.Fatalf("TestError(nil) should be false")
	}
	if !j.TestError(errors.New("boom"), "boom") {
		t.Fatalf("TestError(err) should be true")
	}
}

func TestQueryPostsRequestModelCoverage(t *testing.T) {
	args := &fasthttp.Args{}
	args.Set("i", "12")
	args.Set("s", "abc")
	args.Set("f", "3.5")
	q := &QueryParams{args}
	if q.GetInt("i") != 12 || q.GetStr("s") != "abc" || q.GetFloat("f") != 3.5 {
		t.Fatalf("query params mismatch")
	}

	p := &Posts{SS: M.SS{
		`ok`:       `1`,
		`zero`:     `0`,
		`json`:     `{"x":1}`,
		`arrs`:     `["a","b"]`,
		`arri`:     `[1,2]`,
		`arro`:     `[{"a":1}]`,
		`password`: `secret`,
		`long`:     strings.Repeat("z", 80),
	}}
	if !p.GetBool("ok") || p.GetBool("zero") {
		t.Fatalf("posts bool mismatch")
	}
	if !p.IsSet("ok") || p.IsSet("none") {
		t.Fatalf("posts isset mismatch")
	}
	if p.GetJsonMap("json").GetInt("x") != 1 || len(p.GetJsonStrArr("arrs")) != 2 ||
		len(p.GetJsonIntArr("arri")) != 2 || len(p.GetJsonObjArr("arro")) != 1 {
		t.Fatalf("json conversion mismatch")
	}
	if s := p.String(); !strings.Contains(s, "password") || !strings.Contains(s, "***") {
		t.Fatalf("posts String should mask password: %q", s)
	}
	if s := p.NewlineString(); !strings.Contains(s, "\n\t") {
		t.Fatalf("posts NewlineString mismatch: %q", s)
	}

	ctxForm := &Context{RequestCtx: newTestRequestCtx("POST", "http://x/path", "a=1&password=secret", "application/x-www-form-urlencoded")}
	p2 := &Posts{}
	p2.FromContext(ctxForm)
	if p2.GetStr("a") != "1" || p2.GetStr("password") != "secret" {
		t.Fatalf("FromContext urlencoded mismatch: %#v", p2.SS)
	}

	body := bytes.Buffer{}
	w := multipart.NewWriter(&body)
	if err := w.WriteField("mf", "v1"); err != nil {
		t.Fatalf("multipart field: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("multipart close: %v", err)
	}
	ctxMP := &Context{RequestCtx: newTestRequestCtx("POST", "http://x/upload", body.String(), w.FormDataContentType())}
	p3 := &Posts{}
	p3.FromContext(ctxMP)
	if p3.GetStr("mf") != "v1" {
		t.Fatalf("FromContext multipart mismatch: %#v", p3.SS)
	}

	rm1 := NewRequestModel_ById_ByDbActor_ByAjax("123", "actor1", Ajax{})
	if rm1.Id != "123" || rm1.DbActor != "actor1" || rm1.Ajax.SX == nil || rm1.IdInt() != 123 {
		t.Fatalf("request model by id mismatch: %#v", rm1)
	}
	rm2 := NewRequestModel_ByUniq_ByDbActor_ByAjax("uniq-1", "actor2", Ajax{})
	if rm2.Uniq != "uniq-1" || rm2.DbActor != "actor2" || rm2.Ajax.SX == nil {
		t.Fatalf("request model by uniq mismatch: %#v", rm2)
	}
	rm1.Ctx = &Context{RequestCtx: newTestRequestCtx("POST", "http://x/ajax", "", "")}
	if !rm1.IsAjax() {
		t.Fatalf("request model IsAjax should be true")
	}
}

func TestMailerCoverage(t *testing.T) {
	m := NewMailer("u@example.com", "p", "127.0.0.1", 1)
	m.Name = "Tester"
	if m.Address() != "127.0.0.1:1" || m.From() != "Tester <u@example.com>" || m.Auth() == nil {
		t.Fatalf("mailer basic methods mismatch")
	}
	if errStr := m.SendSyncBCC([]string{"x@y.z"}, "subj", "msg"); errStr == "" {
		t.Fatalf("SendSyncBCC should fail on unreachable smtp server")
	}
	e := email.NewEmail()
	if err := m.SendRaw(e); err == nil {
		t.Fatalf("SendRaw(non-gmail) should fail on unreachable smtp server")
	}
	gm := &SmtpConfig{Name: "G", Username: "u@gmail.com", Password: "p", Hostname: "xgmail.com", Port: 1}
	if err := gm.SendRaw(email.NewEmail()); err == nil {
		t.Fatalf("SendRaw(gmail branch) should fail on unreachable smtp server")
	}
}

func TestSessionCoverage(t *testing.T) {
	snap := snapshotGlobals()
	t.Cleanup(func() { restoreGlobals(snap) })

	conn := newMockSessionConnector()
	InitSession("CK", 20*time.Second, 5*time.Second, conn, conn)

	ctx := &Context{RequestCtx: newTestRequestCtx("GET", "http://x/s", "", "")}
	s := &Session{}
	s.Load(ctx)
	if s.Key == "" || s.UserAgent == "" || s.IpAddr == "" || !s.Changed {
		t.Fatalf("load without cookie should create random key: %#v", s)
	}
	s.Save(ctx)
	ck := fasthttp.Cookie{}
	ck.SetKey(SESS_KEY)
	if !ctx.Response.Header.Cookie(&ck) || len(ck.Value()) == 0 {
		t.Fatalf("Save should set cookie")
	}
	if !strings.Contains(s.StateCSRF(), "|") {
		t.Fatalf("StateCSRF format mismatch")
	}

	s.Login(M.SX{"id": "7", "email": "wm@example.com"})
	if s.GetStr("id") != "7" || s.String() == "" || s.NewlineString() == "" || s.HeaderString() == "" {
		t.Fatalf("login/string mismatch: %#v", s.SX)
	}

	conn.expiry[s.Key] = 1
	s.Changed = false
	s.Touch()
	if !s.Changed || s.SX["renew_at"] == nil {
		t.Fatalf("Touch should renew when expiring soon: %#v", s.SX)
	}

	oldKey := s.Key
	s.Logout()
	if s.Key == "" || s.Key == oldKey {
		t.Fatalf("Logout should reset key")
	}

	badCtx := &Context{RequestCtx: newTestRequestCtx("GET", "http://x/s2", "", "")}
	badCtx.Request.Header.SetCookie(SESS_KEY, "bad-cookie")
	s2 := &Session{}
	s2.Load(badCtx)
	if s2.Key == "bad-cookie" || s2.Key == "" {
		t.Fatalf("invalid cookie should trigger logout/random key: %q", s2.Key)
	}

	goodCtx := &Context{RequestCtx: newTestRequestCtx("GET", "http://x/s3", "", "")}
	key := s2.StateCSRF() + "fixed"
	conn.msx[key] = M.SX{"id": "99", "email": "a@b.c"}
	goodCtx.Request.Header.SetCookie(SESS_KEY, key)
	s3 := &Session{}
	s3.Load(goodCtx)
	if s3.Key != key || s3.GetStr("id") != "99" {
		t.Fatalf("valid cookie should load existing session: %#v", s3.SX)
	}

	emptyCtx := &Context{RequestCtx: newTestRequestCtx("GET", "http://x/s4", "", "")}
	key2 := s3.StateCSRF() + "expired"
	emptyCtx.Request.Header.SetCookie(SESS_KEY, key2)
	s4 := &Session{}
	s4.Load(emptyCtx)
	if s4.Key == key2 || s4.Key == "" {
		t.Fatalf("expired cookie should trigger logout/random key")
	}
}

func TestContextAndFilterCoverage(t *testing.T) {
	snap := snapshotGlobals()
	t.Cleanup(func() { restoreGlobals(snap) })

	baseDir := createBaseDir(t)
	engine := newTestEngine(t, baseDir, true)
	engine.LoadLayout()
	Webmasters = M.SS{"wm@example.com": "wm"}
	conn := newMockSessionConnector()
	InitSession("CK", 20*time.Second, 5*time.Second, conn, conn)

	rctx := newTestRequestCtx("POST", "http://example.com/first/next?q=7", "name=John&password=secret", "application/x-www-form-urlencoded")
	rctx.Request.Header.Set("x-forwarded-proto", "https")
	rctx.Request.Header.Set("Referer", "http://ref.example")
	rctx.SetUserValue("id", "21")
	rctx.SetUserValue("flag", "true")
	rctx.SetUserValue("jmap", `{"a":1}`)
	rctx.SetUserValue("jarr", `["x","y"]`)
	ctx := &Context{
		RequestCtx:  rctx,
		Session:     &Session{UserAgent: "ua", IpAddr: "1.2.3.4", Key: "k1", SX: M.SX{"email": "wm@example.com", "id": "7"}},
		Title:       "MyTitle",
		Engine:      engine,
		ContentType: "text/plain",
		NoLayout:    true,
		Actions:     []Action{func(c *Context) { c.AppendString("next-called") }},
	}

	if ctx.Proto() != "https://" || ctx.Host() != "https://example.com" || !ctx.IsAjax() {
		t.Fatalf("proto/host/ajax mismatch")
	}
	if ctx.Headers().GetStr("x-forwarded-proto") != "https" {
		t.Fatalf("headers parse mismatch")
	}
	if ctx.ParamStr("id") != "21" || ctx.ParamInt("id") != 21 || !ctx.ParamBool("flag") {
		t.Fatalf("param primitive mismatch")
	}
	if ctx.ParamJsonMap("jmap").GetInt("a") != 1 || len(ctx.ParamJsonStrArr("jarr")) != 2 {
		t.Fatalf("param json mismatch")
	}
	if ctx.FirstPath() != "first" {
		t.Fatalf("FirstPath mismatch: %q", ctx.FirstPath())
	}
	if !ctx.IsWebMaster() {
		t.Fatalf("IsWebMaster mismatch")
	}
	if ctx.QueryParams().GetInt("q") != 7 {
		t.Fatalf("QueryParams mismatch")
	}
	if ctx.RequestURL() != "/first/next?q=7" {
		t.Fatalf("RequestURL mismatch: %q", ctx.RequestURL())
	}

	ctx.AppendBytes([]byte("A"))
	b := bytes.Buffer{}
	b.WriteString("B")
	ctx.AppendBuffer(b)
	ctx.AppendString("C")
	ctx.AppendMap(M.SX{"n": 1})
	aj := NewAjax()
	aj.Error("bad")
	ctx.AppendAjax(aj)
	ctx.Render("partial", M.SX{"name": "Neo"})
	if !strings.Contains(ctx.PartialNoDebug("partial", M.SX{"name": "Morpheus"}), "Morpheus") {
		t.Fatalf("PartialNoDebug mismatch")
	}

	if got := ctx.Posts(); got == nil || got.GetStr("name") != "John" {
		t.Fatalf("Posts cache load mismatch: %#v", got)
	}
	if ctx.Posts() != ctx.PostCache {
		t.Fatalf("Posts should use cache")
	}

	if next := ctx.Next(); next == nil {
		t.Fatalf("Next should return action")
	} else {
		next(ctx)
	}
	if !strings.Contains(ctx.Buffer.String(), "next-called") {
		t.Fatalf("Next action not executed")
	}

	ctxErr := &Context{
		RequestCtx:  newTestRequestCtx("GET", "http://example.com/path", "", ""),
		Session:     &Session{UserAgent: "ua", IpAddr: "1.2.3.4", Key: "k2", SX: M.SX{"id": "8"}},
		Engine:      engine,
		Title:       "ErrorTitle",
		ContentType: "text/html",
		NoLayout:    true,
	}
	ctxErr.Error(404, "missing")
	if ctxErr.Response.StatusCode() != 404 || !strings.Contains(ctxErr.Buffer.String(), "error|404") {
		t.Fatalf("Error render mismatch: status=%d body=%q", ctxErr.Response.StatusCode(), ctxErr.Buffer.String())
	}
	if _, _, _, reader := ctxErr.UploadedFile("missing"); reader != nil {
		t.Fatalf("UploadedFile missing should return nil reader")
	}
	if !strings.Contains(ctxErr.RequestLogStr(), "GET /path") || !strings.Contains(ctxErr.RequestDebugStr(), "Session:") ||
		!strings.Contains(ctxErr.RequestHtmlStr(), "Request Path:") {
		t.Fatalf("request debug/log/html strings mismatch")
	}

	ctxFinish := &Context{
		RequestCtx:  newTestRequestCtx("GET", "http://example.com/f", "", ""),
		Session:     &Session{SX: M.SX{"email": "wm@example.com"}},
		Engine:      engine,
		Title:       "T",
		ContentType: "text/plain",
		NoLayout:    true,
	}
	ctxFinish.AppendString("plain")
	ctxFinish.Finish()
	if string(ctxFinish.Response.Body()) != "plain" {
		t.Fatalf("Finish no-layout mismatch: %q", string(ctxFinish.Response.Body()))
	}

	ctxLayout := &Context{
		RequestCtx:  newTestRequestCtx("GET", "http://example.com/l", "", ""),
		Session:     &Session{SX: M.SX{"email": "wm@example.com"}},
		Engine:      engine,
		Title:       "L",
		ContentType: "text/html",
		NoLayout:    false,
	}
	ctxLayout.AppendString("inside")
	ctxLayout.Finish()
	if !strings.Contains(string(ctxLayout.Response.Body()), "layout|L|WTest") {
		t.Fatalf("Finish layout mismatch: %q", string(ctxLayout.Response.Body()))
	}

	ctxPanic := &Context{
		RequestCtx: newTestRequestCtx("GET", "http://example.com/panic", "", ""),
		Session:    &Session{UserAgent: "ua", IpAddr: "2.2.2.2", Key: "k3", SX: M.SX{"id": "9"}},
		Engine:     engine,
		Title:      "panic",
		Actions: []Action{
			func(c *Context) { panic(errors.New("boom")) },
		},
	}
	PanicFilter(ctxPanic)
	if ctxPanic.Response.StatusCode() != 500 || !strings.Contains(ctxPanic.Buffer.String(), "error|500") {
		t.Fatalf("PanicFilter should render 500, got=%d body=%q", ctxPanic.Response.StatusCode(), ctxPanic.Buffer.String())
	}

	ctxNoPanic := &Context{
		RequestCtx: newTestRequestCtx("GET", "http://example.com/ok", "", ""),
		Session:    &Session{UserAgent: "ua", IpAddr: "2.2.2.3", Key: "k4", SX: M.SX{"id": "10"}},
		Engine:     engine,
		Title:      "ok",
		Actions: []Action{
			func(c *Context) { c.AppendString("ok") },
		},
	}
	PanicFilter(ctxNoPanic)
	if !strings.Contains(ctxNoPanic.Buffer.String(), "ok") {
		t.Fatalf("PanicFilter normal path should call action")
	}

	ctxLog := &Context{
		RequestCtx: newTestRequestCtx("POST", "http://example.com/log", "x=1", "application/x-www-form-urlencoded"),
		Session:    &Session{UserAgent: "ua", IpAddr: "3.3.3.3", Key: "k5", SX: M.SX{"id": "11"}},
		Engine:     engine,
		Actions: []Action{
			func(c *Context) { c.SetStatusCode(201); c.AppendString("logged") },
		},
	}
	LogFilter(ctxLog)
	if ctxLog.Response.StatusCode() != 201 {
		t.Fatalf("LogFilter should preserve status code")
	}

	ctxSess := &Context{
		RequestCtx: newTestRequestCtx("GET", "http://example.com/sess", "", ""),
		Engine:     engine,
		Actions: []Action{
			func(c *Context) {
				c.Session.SX["id"] = "77"
				c.Session.Changed = true
			},
		},
	}
	SessionFilter(ctxSess)
	ck2 := fasthttp.Cookie{}
	ck2.SetKey(SESS_KEY)
	if !ctxSess.Response.Header.Cookie(&ck2) || len(ck2.Value()) == 0 {
		t.Fatalf("SessionFilter should save cookie when changed")
	}
}

func TestEngineCoverage(t *testing.T) {
	snap := snapshotGlobals()
	t.Cleanup(func() { restoreGlobals(snap) })

	errs := checkMailers(nil)
	if len(errs) == 0 {
		t.Fatalf("checkMailers should fail when nil")
	}
	Mailers = map[string]*SmtpConfig{
		`debug`: {Name: "N", Username: "debug@example.com", Hostname: "127.0.0.1", Port: 1},
	}
	errs = checkMailers(nil)
	if Mailers[``] == nil || Mailers[`debug`] == nil || len(errs) != 0 {
		t.Fatalf("checkMailers should fill default/debug pair, errs=%v", errs)
	}

	errs = checkSessions(nil)
	if len(errs) == 0 {
		t.Fatalf("checkSessions should fail when not initialized")
	}
	conn := newMockSessionConnector()
	Sessions = conn
	Globals = conn
	errs = checkSessions(nil)
	if len(errs) != 0 {
		t.Fatalf("checkSessions should pass after init: %v", errs)
	}

	Webmasters = nil
	Mailers = nil
	errs = checkWebmasters(nil)
	if len(errs) == 0 {
		t.Fatalf("checkWebmasters should fail when webmasters+mailers nil")
	}
	Mailers = map[string]*SmtpConfig{
		`debug`: {Name: "N", Username: "wm@example.com", Hostname: "127.0.0.1", Port: 1},
		``:      {Name: "N", Username: "wm@example.com", Hostname: "127.0.0.1", Port: 1},
	}
	Webmasters = nil
	errs = checkWebmasters(nil)
	if len(errs) != 0 || Webmasters["wm@example.com"] == "" {
		t.Fatalf("checkWebmasters should auto-set from mailer: errs=%v webmasters=%v", errs, Webmasters)
	}

	if len(checkRoutes(nil)) == 0 {
		t.Fatalf("checkRoutes should fail on nil routes")
	}
	Routes = map[string]Action{}
	if len(checkRoutes(nil)) != 0 {
		t.Fatalf("checkRoutes should pass on non-nil map")
	}
	if len(checkAssets(nil)) == 0 {
		t.Fatalf("checkAssets should fail on nil assets")
	}
	Assets = [][2]string{}
	if len(checkAssets(nil)) != 0 {
		t.Fatalf("checkAssets should pass on non-nil assets")
	}

	tmpBase := filepath.ToSlash(t.TempDir()) + `/`
	errs = checkRequiredDir(nil, tmpBase+`need/public/`, `PUBLIC_SUBDIR`, false)
	if len(errs) != 0 {
		t.Fatalf("checkRequiredDir should create dirs: %v", errs)
	}
	if st, err := os.Stat(tmpBase + `need/public/lib/`); err != nil || !st.IsDir() {
		t.Fatalf("checkRequiredDir should create public/lib")
	}
	errs = checkRequiredDirs(nil, tmpBase+`need2/`)
	if len(errs) != 0 {
		t.Fatalf("checkRequiredDirs should create both dirs: %v", errs)
	}
	filePath := tmpBase + `need2/` + VIEWS_SUBDIR + `x.html`
	errs = checkRequiredFile(nil, filePath, `label`, `default-x`)
	if len(errs) != 0 {
		t.Fatalf("checkRequiredFile create mismatch: %v", errs)
	}
	if b, err := os.ReadFile(filePath); err != nil || !strings.Contains(string(b), "default-x") {
		t.Fatalf("checkRequiredFile should write defaults")
	}

	baseDir := createBaseDir(t)
	if err := os.WriteFile(baseDir+VIEWS_SUBDIR+`page.html`, []byte(`page #{x}`), DEFAULT_FILEDIR_PERM); err != nil {
		t.Fatalf("write page: %v", err)
	}

	engine := newTestEngine(t, baseDir, true)
	engine.LoadLayout()
	tpl1 := engine.Template("page")
	tpl2 := engine.Template("page")
	if tpl1 == nil || tpl1 != tpl2 {
		t.Fatalf("Template should load/cache page template")
	}
	engine.Log("line1")
	if st, err := os.Stat(engine.LogPath); err != nil || st.Size() == 0 {
		t.Fatalf("Engine.Log should write file")
	}

	if engine.SendMailSync("nope", nil, "s", "m") == "" {
		t.Fatalf("SendMailSync invalid id should return error")
	}
	Mailers = nil
	engine.SendDebugMail("dbg")
	Mailers = map[string]*SmtpConfig{
		`debug`: {Name: "N", Username: "wm@example.com", Hostname: "127.0.0.1", Port: 1},
	}
	engine.SendMail("debug", nil, "subject", "msg")

	Assets = [][2]string{
		{`js`, `lib1`},
		{`css`, `lib1`},
		{`/js`, `mod1`},
		{`/css`, `mod1`},
		{`weird`, `x`},
	}
	if err := os.WriteFile(baseDir+PUBLIC_SUBDIR+`lib/lib1.js`, []byte(`var  a = 1;`), DEFAULT_FILEDIR_PERM); err != nil {
		t.Fatalf("write lib1.js: %v", err)
	}
	if err := os.WriteFile(baseDir+PUBLIC_SUBDIR+`lib/lib1.css`, []byte(`body { color: red; }`), DEFAULT_FILEDIR_PERM); err != nil {
		t.Fatalf("write lib1.css: %v", err)
	}
	if err := os.WriteFile(baseDir+PUBLIC_SUBDIR+`mod1.js`, []byte(`var  b = 2;`), DEFAULT_FILEDIR_PERM); err != nil {
		t.Fatalf("write mod1.js: %v", err)
	}
	if err := os.WriteFile(baseDir+PUBLIC_SUBDIR+`mod1.css`, []byte(`h1 { color: blue; }`), DEFAULT_FILEDIR_PERM); err != nil {
		t.Fatalf("write mod1.css: %v", err)
	}
	engine.DebugMode = true
	engine.Assets = ""
	engine.MinifyAssets()
	if !strings.Contains(engine.Assets, "lib/lib1.js") || !strings.Contains(engine.Assets, "mod1.css") {
		t.Fatalf("MinifyAssets debug should output tag list: %q", engine.Assets)
	}

	engine.DebugMode = false
	engine.Assets = ""
	engine.MinifyAssets()
	for _, fname := range []string{`lib.css`, `mod.css`, `lib.js`, `mod.js`} {
		if st, err := os.Stat(baseDir + PUBLIC_SUBDIR + fname); err != nil || st.Size() == 0 {
			t.Fatalf("MinifyAssets production should write %s", fname)
		}
	}
	if !strings.Contains(engine.Assets, "/lib.js?WTest") || !strings.Contains(engine.Assets, "/mod.css?WTest") {
		t.Fatalf("MinifyAssets production assets mismatch: %q", engine.Assets)
	}

	Mailers = map[string]*SmtpConfig{
		`debug`: {Name: "N", Username: "wm@example.com", Hostname: "127.0.0.1", Port: 1},
		``:      {Name: "N", Username: "wm@example.com", Hostname: "127.0.0.1", Port: 1},
	}
	Webmasters = M.SS{"wm@example.com": "wm"}
	Sessions = conn
	Globals = conn
	Assets = [][2]string{}
	Filters = []Action{}
	Routes = map[string]Action{
		`home`: func(c *Context) {
			c.Title = "Home"
			c.AppendString("HOME")
		},
	}

	baseDir2 := createBaseDir(t)
	eng2 := NewEngine(true, false, "ProjectX", baseDir2)
	defer eng2.Logger.Close()
	if eng2.Router == nil || eng2.ViewCache == nil || eng2.Name != "ProjectX" {
		t.Fatalf("NewEngine basic fields mismatch: %#v", eng2)
	}
	reqCtx := newTestRequestCtx("GET", "http://localhost/home", "", "")
	eng2.Router.Handler(reqCtx)
	if reqCtx.Response.StatusCode() != 200 || !strings.Contains(string(reqCtx.Response.Body()), "HOME") {
		t.Fatalf("route handler response mismatch: code=%d body=%q", reqCtx.Response.StatusCode(), string(reqCtx.Response.Body()))
	}
	reqCtx2 := newTestRequestCtx("POST", "http://localhost/home", "x=1", "application/x-www-form-urlencoded")
	eng2.Router.Handler(reqCtx2)
	if !strings.Contains(string(reqCtx2.Response.Header.ContentType()), "application/json") {
		t.Fatalf("POST route should set ajax content-type, got %q", string(reqCtx2.Response.Header.ContentType()))
	}

	panicCtx := newTestRequestCtx("GET", "http://localhost/unknown", "", "")
	eng2.Router.PanicHandler(panicCtx, "boom")
	if panicCtx.Response.StatusCode() != 504 || !strings.Contains(string(panicCtx.Response.Body()), `"is_success":false`) {
		t.Fatalf("panic handler mismatch: code=%d body=%q", panicCtx.Response.StatusCode(), string(panicCtx.Response.Body()))
	}
}
