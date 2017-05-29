// JS to HTML functions

// date formatter
var Const_DMY = 'D MMM YYYY';
var Const_MDY = 'MMMM D, YYYY';
var Const_DMYHM = 'D MMM YYYY HH:mm';
var Const_DMYHMS = 'D MMM YYYY HH:mm:ss';
var Const_HMS = 'HH:mm:ss';
var Const_YMD = 'YYYY-MM-DD';
var Const_YMDHMS = 'YYYY-MM-DD HH:mm:ss';
var Const_ISO = 'YYYYMMDD_HHmmss';
var Const_YM = 'YYYY-MM';
var Const_DM = 'D MMM';
var Const_DAY = 'dddd';
var Const_DAY_DMY = 'dddd, D MMM YYYY';
var Const_DAY_DMYHM = 'dddd, D MMM YYYY HH:mm';
var Const_MY = 'MMM YYYY';

// quote proprty with double quote
function ZZ( name, val ) {
	if( !val ) return '';
	return ' ' + name + '="' + val + '" ';
}

// <form>
function FORM( html ) {
	return '<form class="ui form">' + html + '</form>\n'
}

// id=
function DIVidc( id, html, classes ) {
	return '<div' + ZZ( 'id', id ) + ZZ( 'class', classes ) + '>' + html + '</div>\n'
}

// class=
function DIVc( classes, html ) {
	return '<div' + ZZ( 'class', classes ) + '>' + html + '</div>\n'
}

// <label>
function DIVfield( label, html_input ) {
	return '<div class="field"><label>' + label + '</label>\n' + html_input + '</div>\n'
}

// <div> class=
function DIVcp( classes, props, html ) {
	return '<div' + ZZ( 'class', classes ) + props + '>' + html + '</div>\n'
}

// <i> class=
function Icp( classes, props ) {
	return '<i' + ZZ( 'class', classes ) + ' ' + (props || '') + '></i>\n'
}

// <button if, click, icon
function BUTTONicoif( click, icon, tooltip, if_case ) {
	return '<button' + ZZ( 'v-if', if_case ) + ' class="micro ui icon green button" ' + ZZ( '@click', click ) + '>' + Icp( icon, ZZ( 'title', tooltip ) ) + '</button>'
}

// <a> @click html=
function ICOclick( func, tooltip, html ) {
	return '<a class="item" ' + ZZ( '@click', func ) + ZZ( 'title', tooltip ) + '>' + html + '</a>'
}

// <a> href= <i> class= href
function ICOlink( href, tooltip, ico_classes, target ) {
	var tgt = '', label = '', cls = 'micro ui icon blue button';
	if( target===true ) {
		cls = 'item';
		label = tooltip;
	}
	else tgt = ZZ( 'target', target );
	return '<a' + ZZ('class',cls) + ZZ( ':href', href ) + ZZ( 'title', tooltip ) + tgt + '>' + Icp( ico_classes + ' white' ) + label + '</a>';
}

// options for select
function OPTIONs( opts, hide_key ) {
	var html = '';
	for( var y in opts ) {
		var key_is_val = hide_key || (y==opts[ y ]);
		html += '<option' + ZZ( 'value', y ) + '>' + (key_is_val ? '' : y + ': ') + opts[ y ] + '</option>';
	}
	return html;
}

// dropdown menu
function DROPDOWN( items ) {
	var html = '<div class="ui dropdown"><i class="dropdown icon"></i><div class="menu">';
	for( var z in items ) html += '<div class="item">' + items[ z ] + '</div>';
	return html + '</div>'
}

// append breadcrumb and document title
function BREADCRUMB( desc, text, path ) {
	var el = $( '#breadcrumb' );
	if( el.childElementCount>0 ) {
		el.find( '.section' ).last().removeClass( 'active' );
		el.append( '<i class="right angle icon divider"></i>' );
	}
	var segment1 = path.split( '/' )[ 1 ] || '';
	if( !$( '#' + segment1 + '_menu' ).length ) path = '#';
	el.append( ' | <span class="active section">' + desc + '<a' + ZZ( 'href', path ) + '>' + text + '</a></span>' );
	document.title += ' | ' + desc + ' ' + text;
}
