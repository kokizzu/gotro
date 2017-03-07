package W

import (
	"github.com/kokizzu/gotro/M"
	"os"
)

var Mailers map[string]*SmtpConfig // used in mailer.go
var Webmasters M.SS                // used in engine.go
var Sessions SessionConnector      // from session.go
var Routes map[string]Action
var Assets [][2]string // []{{`js css /js /css`,`filename`}, ...}
var Filters []Action
var PUBLIC_SUBDIR = `public/`
var VIEWS_SUBDIR = `views/`

var Errors = map[int]string{
	0:   `Unknown Error`,
	500: `Internal Server Error`,
	404: `Page Not Found`,
	403: `Forbidden`,
	503: `Server is Overloaded`,
}

var DEFAULT_FILEDIR_PERM = os.FileMode(0755)

// locals: error_code, error_title, error_detail, project_name, requested_path, webmaster
const ERROR_DEFAULT_CONTENT = `
<div style="text-align: center">
	<div style="display: inline-block;">
		<div class="ui message">
			<div class="header">
				Error #{error_code}: #{error_title}
			</div>
		</div>
		<div class="ui grid" id="vm">
			<div class="ui sixteen wide centered column active" style="min-height: 200px">
				<div class="eight wide centered column">
				</div>
				<div class="ui success message">
					<div class="header">
						#{requested_path}
					</div>
					<div class="detail">
						#{error_detail}
					</div>
				</div>
				<div class="ui negative message">
					<div class="header">
						<p id="err403" style="display:none">It's forbidden for you to access this page (insufficient privilege), either contact the #{webmaster} or
							try again. Use your browser's <b>Back</b> button to navigate to the page you have previously come from</p>
						<p id="err404" style="display:none">The page you requested could not be found, either contact the #{webmaster} or try again. Use your
							browser's <b>Back</b> button to navigate to the page you have previously come from</p>
						<button class="ui blue button" onclick="window.history.back()">
						<p id="err" style="display:none">There is a possible programming mistake on this page, either contact the #{webmaster} or try again. Use
							your browser's <b>Back</b> button to navigate to the page you have previously come from</p>
						<p id="err503">Server is overloaded, please wait 2 minutes then try again by clicking your browser's <b>Refresh</b>
							button (F5/Ctrl+R)</p>
						<button>
							<i class="glyphicon glyphicon-arrow-left"></i> Take Me Back
						</button>
						<button class="ui orange button" onclick="window.location='/'">
							<i class="glyphicon glyphicon-home"></i> Take Me Home
						</button>
					</div>
				</div>
				<script>
					var code = '#{err_code}';
					if( code != '404' && code != '403' && code != '503' ) code = '';
					document.getElementById('err'+code).style.display = 'block'
				</script>
			</div>
		</div>
		<hr/>
		<div class="ui message">
			<div class="header">
				&copy; <script>document.write(new Date().getYear())</script> #{project_name}
			</div>
		</div>
	</div>
</div>
`

// locals: title, project_name, assets, contents, is_superadmin, debug_mode
const LAYOUT_DEFAULT_CONTENT = `
<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>#{title} | #{project_name}</title>
	<meta name="apple-mobile-web-app-capable" content="yes">
	<meta name="apple-mobile-web-app-status-bar-style" content="black">
	<link rel="apple-touch-icon" sizes="57x57" href="/apple-icon-57x57.png">
	<link rel="apple-touch-icon" sizes="60x60" href="/apple-icon-60x60.png">
	<link rel="apple-touch-icon" sizes="72x72" href="/apple-icon-72x72.png">
	<link rel="apple-touch-icon" sizes="76x76" href="/apple-icon-76x76.png">
	<link rel="apple-touch-icon" sizes="114x114" href="/apple-icon-114x114.png">
	<link rel="apple-touch-icon" sizes="120x120" href="/apple-icon-120x120.png">
	<link rel="apple-touch-icon" sizes="144x144" href="/apple-icon-144x144.png">
	<link rel="apple-touch-icon" sizes="152x152" href="/apple-icon-152x152.png">
	<link rel="apple-touch-icon" sizes="180x180" href="/apple-icon-180x180.png">
	<link rel="icon" type="image/png" sizes="192x192" href="/android-icon-192x192.png">
	<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
	<link rel="icon" type="image/png" sizes="96x96" href="/favicon-96x96.png">
	<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
	<link rel="manifest" href="/manifest.json">
	<meta name="msapplication-TileColor" content="#ffffff">
	<meta name="msapplication-TileImage" content="/ms-icon-144x144.png">
	<meta name="theme-color" content="#ffffff">
</head>
<body onload="_loaded()" style="width:100%; height: 100%; overflow-x: visible">
	<div id='_backdrop'
	     style="position:absolute;top:0;left:0;width:100%;height:100%;z-index:10;text-align:center;vertical-align:middle;background-color:rgba(0, 0, 0, 0.5);">
		<img src='/img/horizontal-loading.gif'/><br/>Loading scripts, please wait..
	</div>
	#{assets}
	<script>
		var is_superadmin = '#{is_superadmin}'=='true';
		var debug_mode = '#{debug_mode}'=='true';
		function _loaded() {$( '#_backdrop' ).remove();}
	</script>
	<noscript>Please enable Javascript or use Javascript-enabled browser.</noscript>
	#{contents}
	<!-- </body></html> anti telkom ad-inject -->
</body>
</html>
`
