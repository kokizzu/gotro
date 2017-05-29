package main

/*
 * // documentation version
 * {`type`, `filename`},
 * {`js`, `jquery`}, // should be on public/lib/jquery.js
 * {`js`, `bootstrap`}, // should be on public/lib/bootstrap.js
 * {`css`, `bootstrap`}, // should be on public/lib/bootstrap.css
 * if started with slash, then it should be on public/ directory
 */
var ASSETS = [][2]string{
	//// http://api.jquery.com/ 1.11.1
	{`js`, `jquery`},
	{`js`, `vue`},
	////// http://momentjs.com/ 2.17.1
	//{`js`, `moment`},
	////// github.com/kokizzu/semantic-ui-daterangepicker
	//{`css`, `daterangepicker`},
	//{`js`, `daterangepicker`},
	////// http://semantic-ui.com 2.2 // should be below `js` and `css` items
	{`css`, `semantic/semantic`},
	{`js`, `semantic/semantic`},
	////// global, helpers, project specific
	//{`/css`, `global`},
	//{`/js`, `global`},
	//{`/js`, `data_sources`},
	{`/js`, `helper`},
	{`/js`, `global`},
	{`/js`, `vue-common`},
	{`/js`, `form_mason`},
	{`/js`, `grid_mason`},
	{`/js`, `list_mason`},
}
