// K misc formatting
var K = {};
// <input> v-focus
Vue.directive( 'focus', {
	inserted: function( el ) {
		el.focus()
	}
} );
// <input> v-required
Vue.directive( 'required', {
	componentUpdated: function( el ) {
		$( el ).css( 'border-color', (el.value ? 'blue' : 'red') );
	}
} );
/* requires: global.js Moment */
/* requires: data_source.js DS.WeekDay */
K.date = function( d ) {
	if( H.InvalidDate( d ) ) return '-';
	var md = H.Moment( d );
	return md.format( Const_DMY );
};
Vue.filter( 'date', K.date );
K.datetime = function( d ) {
	if( H.InvalidDate( d ) ) return '-';
	var md = H.Moment( d );
	return md.format( Const_DMYHM );
};
Vue.filter( 'datetime', K.datetime );
K.integer = function( n ) {
	n = +n;
	//if( !n ) return '<i>0</i>';
	if( !n ) return '0';
	var str = '';
	var odd = false;
	var prefix = '';
	if( n<0 ) {
		n = -n;
		prefix = '-';
	}
	while( n>0 ) {
		var rem = Math.floor( n % 1000 );
		n = Math.floor( n / 1000 );
		rem = ((n) ? (H.PadZero( rem, 3 )) : rem);
		// if( odd ) rem = '<span class="n">' + rem + '</span>';
		rem = ' ' + rem;
		str = rem + str;
		odd = !odd;
	}
	return prefix + str;
};
Vue.filter( 'integer', K.integer );
K.float2 = function( n ) {
	var two = Math.round( n * 100 % 100 );
	n = Math.floor( n );
	n = +n;
	if( !n && !two ) return '0.00';
	var str = '';
	var odd = false;
	var prefix = '';
	if( n<0 ) {
		n = -n;
		prefix = '-';
	}
	while( n>0 ) {
		var rem = Math.floor( n % 1000 );
		n = Math.floor( n / 1000 );
		rem = ((n) ? (H.PadZero( rem, 3 )) : rem);
		//if( odd ) rem = '<span class="n">' + rem + '</span>';
		rem = ' ' + rem;
		str = rem + str;
		odd = !odd;
	}
	two = H.PadZero( two, 2 );
	return prefix + str + '.' + two;
};
Vue.filter( 'float2', K.float2 );
K.bool = function( val ) {
	return !!val ? 'yes' : 'no';
};
Vue.filter( 'bool', K.bool );
Vue.filter( 'filled', K.bool );
K.weekdate = function( d ) {
	if( 'string'== typeof d ) d = H.Moment( d );
	if( !d || !d.isValid ) return 'Weekday, ' + Const_DMYHM;
	return DS.WeekDay[ d.isoWeekday() ] + ', ' + d.format( Const_DMYHM );
};
Vue.filter( 'weekdate', K.weekdate );
K.weekdatetime = function( d ) {
	if( 'string'== typeof d ) d = H.Moment( d );
	if( !d || !d.isValid ) return 'Weekday, ' + Const_DMY;
	return DS.WeekDay[ d.isoWeekday() ] + ', ' + d.format( Const_DMY );
};
Vue.filter( 'weekdatetime', K.weekdatetime );
K.filesize = function( b ) {
	var ones = { 0: 'Byte', 1: 'KB', 2: 'MB', 3: 'GB', 4: 'TB' };
	var level = 0;
	while( b>=0 ) {
		if( b<1024 ) return K.float2( b ) + ' ' + ones[ level ];
		b /= 1024;
		++level;
	}
};
Vue.filter( 'filesize', K.filesize );
K.json = function( json ) {
	if( typeof json!='string' ) json = JSON.stringify( json, null, 2 );
	json = json.replace( /</g, '&lt;' ).replace( />/g, '&gt;' ); // replace(/&/g, '&amp;')
	var pattern = /("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g;
	var html = json.replace( pattern, function( match ) {
		var cls = 'number';
		var suffix = '';
		if( /^"/.test( match ) ) {
			if( /:$/.test( match ) ) {
				cls = 'key';
				match = match.slice( 0, -1 );
				suffix = ':'
			} else {
				cls = 'string';
			}
		} else if( /true|false/.test( match ) ) {
			cls = 'boolean';
		} else if( /null/.test( match ) ) {
			cls = 'null';
		}
		return '<span class="' + cls + '">' + match + '</span>' + suffix;
	} );
	return html;
};
Vue.filter( 'json', K.json );