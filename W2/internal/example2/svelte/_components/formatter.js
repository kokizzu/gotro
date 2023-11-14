function datetime( unixSec, humanize ) {
    if( !unixSec ) return '';
    if( typeof unixSec==='string' ) return unixSec; // might not be unix time
    let dt = new Date( unixSec * 1000 );
    if( !humanize ) {
        dt = dt.toISOString();
        return dt.substring( 0, 10 ) + ' ' + dt.substring( 11, 16 );
    }
    const options = {day: '2-digit', month: 'long', year: 'numeric'};
    const formattedDate = dt.toLocaleDateString( undefined, options );
    return formattedDate;
}


function localeDatetime( unixSec ) {
    if( !unixSec ) return '';
    const dt = new Date( unixSec * 1000 );
    const day = dt.toLocaleDateString( 'default', {weekday: 'long'} );
    const date = dt.getDate();
    const month = dt.toLocaleDateString( 'default', {month: 'long'} );
    const year = dt.getFullYear();
    let hh = dt.getHours();
    if( hh<10 ) hh = '0' + hh;
    let mm = dt.getMinutes();
    if( mm<10 ) mm = '0' + mm;
    const formattedDate = `${day}, ${date} ${month} ${year} ${hh}:${mm}`;
    return formattedDate;
}

module.exports = {
    datetime: datetime,
    localeDatetime: localeDatetime,
};