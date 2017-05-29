function ListMason( id, columns, row, extra_action, events, vue_options ) {
	row = row || [];
	if( !id || 'string'!= typeof id ) throw('ListMason: first argument should be an html (id)');
	if( !columns || 'object'!= typeof columns ) throw('ListMason: second argument should be an array of (columns)');
	if( 'string'!= typeof extra_action ) throw('ListMason: fourth argument should be a string');
	var z, y, ret;
	// initialize callbacks
	var default_events = {
		OnListCreated: function() {H.Trace( arguments )},
		OnListMounted: function() {H.Trace( arguments )},
		OnListUpdated: function() {H.Trace( arguments )},
		OnListdDestroyed: function() {H.Trace( arguments )},
		OnBeforeRefresh: function() {H.Trace( arguments )},
		OnAfterRefresh: function() {H.Trace( arguments )},
		OnVueCreated: function() {H.Trace( arguments )}
	};
	for( z in events ) default_events[ z ] = events[ z ];
	// initialize generate grid default html
	var column_by_key = {};
	var html = '<div class="ui center aligned divided selection list"><div class="item" v-if="can_edit || extra_action">' + BUTTONicoif( 'edit_event', 'edit icon', 'edit', 'can_edit' ) + extra_action + '</div>';
	for( z = 0; z<columns.length; ++z ) {
		var column = columns[ z ];
		var col = '', cla = '';
		var key = column.key;
		column_by_key[ key ] = col;
		var typ = column.type;
		switch( typ ) {
			case 'float2':
			case 'integer':
				cla = 'r u';
				col = '{{row.' + key + ' | ' + typ + '}}';
				break;
			case 'bool':
			case 'filled': // not empty
				cla = 'c';
				col = '{{row.' + key + ' | ' + typ + '}}';
				break;
			case 'date':
			case 'datetime':
				cla = 'r';
				col = '{{row.' + key + ' | ' + typ + '}}';
				break;
			case 'image':
			case 'video':
				if( vue_options && vue_options.data && vue_options.data.can_edit ) {
					col = '<div class="file_upload not_initialized"' + ZZ( 'data-key', key ) + ZZ( ':data-id', 'row.id' ) + ZZ( 'data-type', typ ) + '>Upload</div>';
					col += '<span' + ZZ( 'v-if', 'row.' + key + ' && row.' + key + '_size' ) + '>';
					if( typ=='image' ) col += '<img style="height:100%;max-height:240px" :src="row.' + key + '" /><br/>';
					else if( typ=='video' ) col += '<video controls height="240"><source :src="row.' + key + '" type="video/mp4" /></video><br/>';
					col += '<a :href="row.' + key + '">{{row.' + key + '_size | filesize}} | {{row.' + key + '_time | datetime}}</a>';
					col += '</span>';
				} else {
					if( typ=='image' ) col += '<br/><img style="height:100%;max-height:240px" :src="row.' + key + '" />';
					else if( typ=='video' ) col += '<br/><video controls height="240"><source :src="row.' + key + '" type="video/mp4" /></video>';
				}
				break;
			case 'json':
				col = '<div' + ZZ( 'v-html', 'row.' + key + '' ) + '></div>';
				break;
			default:
				col = '{{row.' + key + '}}';
		}
		html += '<div class="item"><div class="ui black horizontal label">' + column.label + '</div>' + col + '</div>';
	}
	html += '</div>';
	var el_id = '#' + id;
	var el = $( el_id );
	el.html( html );
	var el_parent = el.parent();
	var default_opts = {
		el: el_parent[ 0 ],
		data: {
			can_edit: false,
			columns: columns,
			loading: 0,
			row: row,
			term: '',
			rows_index: {}, // row.id to index position
			limit: 10,
			column_by_key: column_by_key,
			extra_action: extra_action
		},
		computed: {
			col_shown_hash: function() {
				var res = {};
				var cs = this.columns_shown;
				for( var z in cs ) res[ cs[ z ] ] = true;
				return res;
			}
		},
		created: function() {
			H.ExecIfFunction( default_events.OnGridCreated, ret );
		},
		mounted: function() {
			el = $( el_id );
			el_parent = el.parent();
			H.ExecIfFunction( default_events.OnGridMounted, ret );
			H.InitFileUpload( el );
			el.find( 'select.dropdown' ).dropdown();
		},
		updated: function() {
			H.ExecIfFunction( default_events.OnGridUpdated, ret );
		},
		destroyed: function() {
			H.ExecIfFunction( default_events.OnGridDestroyed, ret );
		},
		methods: {
			edit_event: function() {
				MASTER_FORM.edit_event( row.id );
			},
			refresh_event: function() {
				H.Trace( arguments );
				var values = {
					a: 'form',
					id: row.id
				};
				this.loading += 1;
				var self = this;
				H.ExecIfFunction( default_events.OnBeforeRefresh, values );
				H.Post( values, function( res ) {
					if( H.HasAjaxErrors( res ) ) return;
					self.row = Object.assign( {}, self.row, res );
					H.ExecIfFunction( default_events.OnAfterRefresh, res );
					self.loading -= 1;
				} )
			}
		}
	};
	// merge options with vue
	for( z in vue_options ) {
		var oz = vue_options[ z ];
		var nz = default_opts[ z ];
		if( nz && oz && 'object'== typeof nz && 'object'== typeof oz ) {
			for( y in oz ) default_opts[ z ][ y ] = oz[ y ];
			continue;
		}
		if( oz ) default_opts[ z ] = oz;
	}
	// create object and generate reverse index
	ret = new Vue( default_opts );
	console.log( 'Vue', ret, default_opts );
	H.ExecIfFunction( default_events.OnVueCreated );
	return ret;
}
