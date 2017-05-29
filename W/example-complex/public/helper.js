var H = H || {};

// check if item is array
H.IsArray = function( arr ) {
	if( 'undefined'== typeof Array.isArray ) {
		Array.isArray = function( obj ) {
			return Object.prototype.toString.call( obj )==='[object Array]';
		}
	}
	return Array.isArray( arr );
};

// check if item is function
/**
 * @return {boolean}
 */
H.IsFunc = function( fun ) {
	return 'function'== typeof fun
};

// filter function, depend on Array.isArray
// arr = array or object
// fun = function(item, key, arr), returns bool
H.Filter = function( arr, fun ) {
	if( !H.IsFunc( fun ) ) throw('H.Filter 2nd paramter must be a function');
	var is_arr = H.IsArray( arr );
	var res;
	if( is_arr ) res = [];
	else res = {};
	for( var z in arr ) {
		if( !arr.hasOwnProperty( z ) ) continue;
		var el = arr[ z ];
		if( !fun( el, z, arr ) ) continue;
		if( is_arr ) res.push( el );
		else res[ z ] = el;
	}
	return res;
};

// stack trace
H.Trace = function() {
	if( !H.trace_function ) return;
	var stack = H.Filter( (new Error()).stack.split( "\n" ), function( txt ) {
		return txt.indexOf( 'jquery' )=== -1
	} );
	stack.shift();
	stack.shift();
	var str = '';
	var len = stack.length - 1;
	for( var z = 0; z<len; ++z ) str += ' '
	console.log( 'FUNCTION TRACE:', len, str, stack[ 0 ], arguments[ 0 ] )
};

// execute if function with single parameter
H.ExecIfFunction = function( func, params ) {
	//H.Trace(arguments);
	if( !func || !H.IsFunc( func ) ) return params;
	//console.log('warning:','cannon execute non-function', func, params);
	return func( params );
};

// replace all string
H.Replace = function( str, from, to ) {
	return str.replace( new RegExp( from, 'g' ), to )
};

// newline to break
H.NL2BR = function( str ) {
	return H.Replace( str, '\n', '<br/>' )
};

// break to newline
H.BR2NL = function( str ) {
	return H.Replace( str, '<br/>', '\n' )
};

// copy to clipboard, depend on jquery
/**
 * @return {boolean}
 */
H.Clipboard = function( txt ) {
	txt = H.BR2NL( txt );
	var oText = false;
	var bResult = false;
	try {
		oText = document.createElement( "textarea" );
		$( oText ).addClass( 'clipboardCopier' ).val( txt ).insertAfter( 'body' ).focus();
		oText.select();
		document.execCommand( "Copy" );
		bResult = true;
	} catch( e ) {
		H.GrowlInfo( 'failed to copy to clipboard' );
	}
	$( oText ).remove();
	return bResult;
};

// coalesce, returns empty string if no truthy value found
H.Coalesce = function() {
	for( var z = 0; z<arguments.length; ++z ) if( arguments[ z ] ) return arguments[ z ];
	return '';
};

// similar to coalesce but only one value, returns empty string if not truthy value
H.OrEmptyStr = function( a ) {
	if( !a ) return '';
	return a;
};

// semanticUiGrowl, depend on jquery
// modified from: tomitrescak/meteor-semantic-ui-growl
// msg is text, may contain html
// options is overrider, see $semanticUiGrowl.defaultOptions
$.semanticUiGrowl = function( msg, options ) {
	var $alert;
	var css;
	var offsetAmount;
	msg = (msg + '').replace( /\n\u003cbr\/\u003e/g, '\n' );
	msg = H.NL2BR( msg );
	
	options = $.extend( {}, $.semanticUiGrowl.defaultOptions, options );
	$alert = $( '<div>' );
	$alert.attr( 'class', 'ui message semanticui-growl ' + options.cls );
	if( options.type ) $alert.addClass( options.type );
	if( options.allow_dismiss ) $alert.append( '<i class="close icon button"></i>' ).find( 'i.close' ).on( 'click', function() { $alert.hide(); } );
	if( options.type=='red' ) $alert.append( '<i class="copy icon button"></i>' ).find( 'i.copy' ).on( 'click', function() {
		var txt = H.OrEmptyStr( options.header );
		if( txt ) txt += "\n";
		txt += msg;
		H.Clipboard( txt );
		H.GrowlInfo( 'Text copied successfully', 'Clipboard' );
	} );
	if( options.header ) $alert.append( '<div class="header">' + options.header + '</div>' );
	
	$alert.append( msg );
	/** @namespace options.top_offset */
	if( options.top_offset ) options.offset = { from: 'top', amount: options.top_offset };
	offsetAmount = options.offset.amount;
	
	$( '.semanticui-growl' ).each( function() {
		return offsetAmount = Math.max( offsetAmount, parseInt( $( this ).css( options.offset.from ) ) + $( this ).outerHeight() + options.stackup_spacing );
	} );
	css = {
		'position': (options.ele==='body' ? 'fixed' : 'absolute'),
		'margin': 0,
		'z-index': '9999',
		'display': 'none'
	};
	css[ options.offset.from ] = offsetAmount + 'px';
	$alert.css( css );
	if( options.width!=='auto' ) $alert.css( 'width', options.width + 'px' );
	$( options.ele ).append( $alert );
	switch( options.align ) {
		case 'center':
			$alert.css( {
				'left': '50%',
				'margin-left': '-' + ($alert.outerWidth() / 2) + 'px'
			} );
			break;
		case 'left':
			$alert.css( 'left', '20px' );
			break;
		default:
			$alert.css( 'right', '20px' );
	}
	$alert.fadeIn();
	if( options.delay>0 ) {
		$alert.delay( options.delay ).fadeOut( 500, function() { $alert.remove(); } );
	}
	return $alert;
};
$.semanticUiGrowl.defaultOptions = {
	ele: 'body',
	type: 'blue',
	header: '',
	offset: {
		from: 'top',
		amount: 20
	},
	align: 'right',
	width: 250,
	delay: 14000,
	allow_dismiss: true,
	stackup_spacing: 10,
	cls: ''
};

// growl error
H.GrowlError = function( msg, title ) {
	$.semanticUiGrowl( msg, { type: 'red', header: title || 'error' } )
};

// growl info
H.GrowlInfo = function( msg, title ) {
	var opt = (title ? { header: title } : {});
	if( 'object'== typeof msg ) msg = JSON.stringify( msg );
	$.semanticUiGrowl( msg, opt );
};

// ajax post, fill values._url if the target url differ from current url
H.Post = function( values, success_callback, error_callback ) {
	H.Trace( arguments );
	var url = '';
	if( values._url!=null ) {
		url = values._url;
		delete values._url
	}
	try {
		$.post( url, values, function( res ) {
			H.Trace( arguments );
			H.ExecIfFunction( success_callback, res )
		} ).fail( function( xhr, textStatus, errorThrown ) {
			H.Trace( arguments );
			xhr.textStatus = textStatus;
			xhr.errorThrown = errorThrown;
			H.ExecIfFunction( error_callback, xhr );
			if( !errorThrown ) errorThrown = 'unable to load resource, network connection or server is down?';
			H.GrowlError( textStatus + ' ' + errorThrown + '<br/>' + xhr.responseText );
		} );
	} catch( e ) {
		console.log( e );
		H.GrowlError( e );
	}
};

// is enter key, need event key
/**
 * @return {boolean}
 */
H.IsEnter = function( e ) {
	return e.keyCode==13 || e.which==13 || e.keyIdentifier=='Enter';
};
// is enter key, need event key
/**
 * @return {boolean}
 */
H.IsEscape = function( e ) {
	return e.keyCode==27 || e.which==27 || e.keyIdentifier=='Escape';
};

// add 0 on left
H.PadZero = function( s, len ) {
	s = '' + s;
	while( s.length<len ) s = '0' + s;
	return s;
};

/* required by common-vue.js */
/**
 * @return {boolean}
 */
H.InvalidDate = function( d ) {
	return !d || d==' ' || d=='-' || d=='0001-01-01T00:00:00Z';
};

/* required by common-vue.js */
/* require moment.js */
H.Moment = function( d ) {
	var typ = typeof(d);
	if( 'number'==typ || ('string'==typ && !isNaN( d )) ) d *= 1000;// moment use milisecond
	var md = moment( d );
	if( (d + '').indexOf( 'Z' )>0 ) md = md.zone( 0 );
	return md;
};

// check ajax response
/**
 * @return {boolean}
 */
H.HasAjaxErrors = function( res, title ) {
	var err = res ? res.errors : '';
	if( !err && 'string'== typeof res ) {
		err = [ res ];
	}
	if( err && err.length ) {
		err = err.join( '<br/>' );
		H.GrowlError( err, title );
		return true;
	}
	return false;
};

// clone an object, non-recursive copy
H.Clone = function( obj ) {
	//return Object.create(obj);
	//return JSON.parse(JSON.stringify(obj));
	if( H.IsArray( obj ) ) return obj.slice( 0 );
	var target = {};
	for( var i in obj ) {
		if( obj.hasOwnProperty( i ) ) target[ i ] = obj[ i ];
	}
	return target;
};

// check if item is object
/**
 * @return {boolean}
 */
H.IsArray = function( obj ) {
	return !!obj && Array===obj.constructor;
};

// deep clone
H.DeepClone = function( obj ) {
	return JSON.parse( JSON.stringify( obj ) );
};

// TODO: diff 2 object

// get cookie
H.Cookie = function( name ) {
	var value = '; ' + document.cookie;
	var parts = value.split( '; ' + name + '=' );
	if( parts.length>=2 ) return parts.pop().split( ';' ).shift();
};

// get localStorage
H.LocalLoad = function( key, not_exists ) {
	H.Trace( arguments );
	if( !localStorage || !localStorage.getItem ) return not_exists;
	var res = localStorage.getItem( key );
	if( res==null || res==undefined ) return not_exists;
	return JSON.parse( res ).val;
};

// set localStorage
H.LocalStore = function( key, value ) {
	H.Trace( arguments );
	if( !localStorage || !localStorage.getItem ) return;
	var val = JSON.stringify( { val: value } );
	if( val==localStorage.getItem( key ) ) return;
	localStorage.setItem( key, val )
};

// load script async
H.LoadScript = function( src, callback ) {
	H.Trace( arguments );
	var script = document.createElement( 'script' );
	script.type = 'text/javascript';
	script.src = src;
	script.addEventListener( 'load', function( e ) {
		H.ExecIfFunction( callback, e );
	}, false );
	var head = document.getElementsByTagName( 'head' )[ 0 ];
	head.appendChild( script );
};

// load scripts async
H.LoadScripts = function( sources, callback ) {
	H.Trace( arguments );
	var len = sources.length;
	var now = 0;
	var loadOne = function() {
		if( now<len ) return H.LoadScript( sources[ now++ ], loadOne );
		if( now>=len ) H.ExecIfFunction( callback );
	};
	loadOne();
};

// trigger element
H.TriggerVue = function( el, val ) {
	el.val( val );
	var e = document.createEvent( 'HTMLEvents' );
	e.initEvent( 'input', true, true );
	el[ 0 ].dispatchEvent( e );
};

// get file upload options
H.UploadOptions = function( typ ) {
	var allowed = { image: 'jpg,png', video: 'mp4,mov' };
	var sizes = { image: 8, video: 40 };
	return {
		url: '',
		dragDropStr: '',
		uploadButtonClass: 'ui mini button file-upload',
		statusBarWidth: 320,
		dragdropWidth: 240,
		allowedTypes: (allowed[ typ ] || 'jpg,png,pdf,docx,odt'),
		multiple: false,
		grids: [],
		maxFileSize: (sizes[ typ ] || 12) * 1024 * 1024, // 8MB
		onSuccess: function( files, data, xhr ) {
			H.Trace( files, data, xhr );
			/** @namespace xhr.responseJSON */
			var res = xhr.responseJSON || xhr.responseText;
			if( 'string'== typeof res ) res = JSON.parse( res );
			if( H.HasAjaxErrors( res ) ) return;
			H.GrowlInfo( 'File(s) Uploaded Successfully\n' + res.info );
		},
		onError: function( files, status, errMsg ) {
			H.Trace( files, status, errMsg );
			H.GrowlError( status + ' ' + errMsg + ' ' + files );
		},
		onSubmit: function( files ) {
			H.Trace( files );
		}
	};
};

// init file upload
H.InitFileUpload = function( el_parent ) {
	el_parent.find( 'div.file_upload.not_initialized' ).each( function() { // initialize file_upload
		var el = $( this );
		el.removeClass( 'not_initialized' );
		var opt = H.UploadOptions( el.data( 'type' ) );
		opt = $.extend( opt, {
			formData: {
				a: 'file_upload',
				key: el.data( 'key' ),
				id: el.data( 'id' )
			}
		} );
		el.uploadFile( opt );
	} );
};