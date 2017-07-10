// require last_range variable
// range_from element
// range_to element
function download_log( sta_id, typ ) {
	var to = '/guest/download_log/' + sta_id + '/' + last_range + '/' + typ;
	var left = moment( $( '#range_from' ).val() );
	var right = moment( $( '#range_to' ).val() );
	if( last_range=='range' ) to += '?from=' + left.unix() + '&to=' + right.unix();
	var a = document.createElement( 'A' );
	var time_format = 'YYYY-MM-DD_HHMMSS';
	a.href = to;
	a.download = sta_id + '__' + last_range + '__' + left.format( time_format ) + '__' + right.format( time_format ) + '.' + typ;
	document.body.appendChild( a );
	a.click();
	document.body.removeChild( a );
}
