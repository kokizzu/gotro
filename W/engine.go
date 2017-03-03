package W

import "github.com/kokizzu/gotro/M"

type Engine struct {
	DebugMode       bool
	WebMasterAnchor M.SX
	GlobalInt       M.SI
	GlobalStr       M.SS
	GlobalAny       M.SX
}

func (engine *Engine) SyncSendMail(mail_id string, bcc []string, subject string, message string) string {
	// TODO: continue this
	return ``
}
func (engine *Engine) StartServer(addressPort string) {
	// TODO: continue this
	// engine.MinifyAssets()
	msg := `[DEVELOPMENT]`
	if !engine.DebugMode {
		msg = `[PRODUCTION]`
	}
	_ = msg
}

// attach a middleware on non-static files
func (engine *Engine) Use(m Action) {
	//engine.Filters = append(engine.Filters, m)
	// TODO: continue this
}

func NewEngine(debugMode bool, multiApp bool, projectName string, mailAccounts map[string]*SmtpConfig, webMaster M.SS, baseDir string, assets [][2]string, baseUrls M.SS, staticSubdir ...string) *Engine {
	// TODO: continue this
	return &Engine{}
}

// register post and get
func (engine *Engine) REGISTER(url string, action Action) {
	// url += `/*all`
	// TODO: continue this
	//engine.GET(url, get_action)
	//engine.POST(url, post_action)
}
