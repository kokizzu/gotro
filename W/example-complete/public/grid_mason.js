function GridMason( id, rec_type, columns, rows, actions, events, vue_options ) {
	rows = rows || [];
	events = events || {};
	vue_options = vue_options || {};
	if( !id || 'string'!= typeof id ) throw('GridMason: first argument should be an html (id)');
	if( !rec_type || 'string'!= typeof rec_type ) throw('GridMason: second argument should be a string (rec_type)');
	if( !columns || 'object'!= typeof columns ) throw('GridMason: third argument should be an array of (columns)');
	var z, y, ret;
	var local_storage_key = window.location + '|' + rec_type;
	// initialize callbacks
	var default_events = {
		OnGridCreated: function() {H.Trace( arguments )},
		OnGridMounted: function() {H.Trace( arguments )},
		OnGridUpdated: function() {H.Trace( arguments )},
		OnGridDestroyed: function() {H.Trace( arguments )},
		OnVueCreated: function() {H.Trace( arguments )},
		OnBeforeSearch: function() {H.Trace( arguments )},
		OnAfterSearch: function() {H.Trace( arguments )}
	};
	for( z in events ) default_events[ z ] = events[ z ];
	// initialize generate grid default html
	var ri = Icp( 'right chevron icon' );
	var li = Icp( 'left chevron icon' );
	var show_opts = {};
	var header = DIVc( 'ui icon input',
			'<input type="text" placeholder="Search..." v-model="term" v-focus @keydown="key_search_event"/>' +
			Icp( 'inverted circular search link icon', ZZ( '@click', 'search_event' ) )
		) + '<span> Page {{page}} of {{total_page}}, Total: {{count}} ' +
		DIVc( 'ui right labeled input',
			'<input type="number" v-model="limit" @blur="search_event" min="1" class="r"/>' +
			DIVc( 'ui label', 'Row(s) per page' )
		) + '</span> ' +
		BUTTONicoif( 'edit_event(0)', 'plus icon', 'add new ' + rec_type, 'can_add' ) +
		BUTTONicoif( 'filter_form_event', 'filter icon', 'filter ' + rec_type ) +
		'<span v-if="filter_count>0">' + ICOclick( 'clean_filter_event', 'clean filter', Icp( 'undo icon' ) ) + ' {{filter_count}} active filter{{ filter_count>1 ? "s" : "" }}</span> ' +
		DIVc( 'ui right floated pagination menu',
			'<a class="icon item" @click="page_event(1)">' + li + li + '</a>' +
			'<a class="icon item" @click="page_event(0,-1)">' + li + '</a>' +
			'<a v-for="idx in pager" class="item"' +
			'   :class="{ active: (idx==page), disabled: (idx>total_page || idx<1), invisible: (idx>total_page || idx<1) }"' +
			'   @click="page_event(idx)">{{idx}}</a>' +
			'<a class="icon item" @click="page_event(0,+1)">' + ri + '</a>' +
			'<a class="icon item" @click="page_event(0,+2)">' + ri + ri + '</a>'
		); // TODO: use global.js instead of raw html
	var thead = '<th>No.</th><th>A <i class="sort icon" @click="clear_order_event"></i></th>';
	var tbody = '<td class="r">{{offset+index+1}}</td><td class="u">' + actions + '</td>';
	var tfoot = '<th></th><th></th>';
	var show_first_n = 3;
	var wdw = $( window ).width();
	if( wdw>800 ) ++show_first_n;
	if( wdw>1024 ) ++show_first_n;
	if( wdw>1280 ) ++show_first_n;
	if( wdw>1440 ) ++show_first_n;
	if( wdw>1600 ) ++show_first_n;
	if( wdw>1920 ) ++show_first_n;
	if( wdw>2560 ) ++show_first_n;
	var first_n_columns = [];
	for( z = 0; z<columns.length; ++z ) {
		var column = columns[ z ];
		var col = '', cla = '';
		var key = column.key;
		if( z<show_first_n ) first_n_columns.push( key );
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
			case 'image_readonly':
				col += '<img style="height:100%;max-height:240px" :src="row.' + key + '" />';
				break;
			case 'video_readonly':
				col += '<video controls height="240"><source :src="row.' + key + '" type="video/mp4" /></video>';
				break;
			case 'image':
			case 'video':
				// file upload
				col = '<div class="file_upload not_initialized"' + ZZ( 'data-key', key ) + ZZ( ':data-id', 'row.id' ) + ZZ( 'data-type', typ ) + '>Upload</div>';
				col += '<span' + ZZ( 'v-if', 'row.' + key + ' && row.' + key + '_size' ) + '>';
				if( typ=='image' ) col += '<img style="height:100%;max-height:240px" :src="row.' + key + '" /><br/>';
				else if( typ=='video' ) col += '<video controls height="240"><source :src="row.' + key + '" type="video/mp4" /></video><br/>';
				col += '<a :href="row.' + key + '">{{row.' + key + '_size | filesize}} | {{row.' + key + '_time | datetime}}</a>';
				col += '</span>';
				break;
			case 'json':
				col = '<div' + ZZ( 'v-html', 'row.' + key + '' ) + '></div>';
				break;
			default:
				col = '{{row.' + key + '}}';
		}
		var shown_if = ZZ( 'v-if', 'col_shown_hash.' + key );
		show_opts[ key ] = column.label;
		var click = '@click="click_order_event(' + "'" + key + "'" + ',$event)"'
		var sorted = "{ascending:(order_hash." + key + "=='+'),descending:(order_hash." + key + "=='-'),sort:(!!order_hash." + key + "),icon:(!!order_hash." + key + ")}";
		thead += '<th' + shown_if + click + '>' + column.label + ' <i ' + ZZ( ':class', sorted ) + '>{{order_idx.' + key + '}}</i></th>';
		tbody += '<td' + shown_if + ZZ( 'class', cla ) + '>' + col + '</td>';
		// TODO: add sum, etc
		tfoot += '<th' + shown_if + '></th>';
	}
	thead = '<thead><tr>' + thead + '</tr></thead>';
	tbody = '<tbody><tr v-for="(row,index) in rows" :class="{ del:(!!row.is_deleted) }">' + tbody + '</tr></tbody>';
	tfoot = '<tfoot><tr>' + tfoot + '</tr></tfoot>';
	header += '<div class="ui floating labeled input button"><span class="yellow ui tiny label">Columns</span>' +
		'<select class="ui fluid search dropdown" multiple v-model="columns_shown">' + OPTIONs( show_opts, true ) + '</select>' +
		'</div>';
	var g_html = header + '<table class="ui celled very compact striped stackable table">' + thead + tbody + tfoot + '</table>';
	// initialize generate form default html
	var f_html = '';
	var mod = 3, shown = 0;
	columns.push( { type: 'bool', key: 'is_deleted', label: 'Is Deleted?' } );
	var filter_state = H.LocalLoad( local_storage_key + '|filter_state' ) || {};
	var filter_val = H.LocalLoad( local_storage_key + '|filter_val' ) || {};
	var order = H.LocalLoad( local_storage_key + '|order' ) || [];
	var column_by_key = {};
	for( z = 0; z<columns.length; ++z ) {
		var col = columns[ z ];
		// TODO: switch case per type, change to INPUT, TEXTAREA, DATE, TIME, etc
		var input = '';
		var typ = col.type;
		var key = col.key;
		column_by_key[ key ] = col;
		var props = ZZ( 'name', key ) + ZZ( 'v-model', 'filter_val.' + key );
		var def_val = filter_val[ key ] || '';
		switch( typ ) {
			case 'image':
			case 'video':
				continue;
			case 'textarea':
				input = '<textarea' + props + ' placeholder="eg. term a | term b"></textarea>';
				break;
			case 'bool':
			case 'filled':
				/** @namespace col.sub_type */
				var sub_typ = col.sub_type || '"Yes":"No"';
				input = '<input' + props + ' type="checkbox" tabindex="0" class="hidden"/>';
				input += '<label>{{ !!filter_val.' + key + ' ? ' + sub_typ + ' }}</label>';
				input = DIVc( 'ui slider checkbox', input );
				def_val = false;
				break;
			case 'integer':
				input = '<input' + props + ' placeholder="eg. >10 <=999 | 55" />';
				def_val = 0;
				break;
			case 'float2':
				input = '<input' + props + ' placeholder="eg. >=0.01 <9.99 | >15 <=30" />';
				def_val = 0;
				break;
			case 'multiselect':
			case 'select':
				props += ZZ( 'multiple', 'multiple' ); // fallthrough
				var sel_opts = DS[ col.sub_type ] || {};
				input = '<select class="dropdown"' + props + '>' + OPTIONs( sel_opts ) + '</select>';
				def_val = [];
				break;
			case 'datetime':
			case 'date':
				var tooltip = 'eg. >2016-01-01 <=2017-01-01T23:59 | 2017-03-06';
				input = '<input' + props + ZZ( 'placeholder', tooltip ) + ZZ( 'title', tooltip ) + '/>';
				break;
			case 'emails':
			case 'url':
			case 'phone':
				typ = 'text'; // fallthrough
			default:
				input = '<input ' + props + ZZ( 'type', typ ) + ' placeholder="eg. term a | term b" />';
		}
		if( !filter_val[ key ] ) filter_val[ key ] = def_val;
		var input2 = '<input v-model="filter_state.' + key + '" type="checkbox" tabindex="0" class="hidden"/>';
		var not_set = " <i v-if='undefined==filter_val." + key + "' class='red warning icon' title='filter value not set'></i>"
		input2 = DIVc( 'ui toggle checkbox', input2 );
		f_html += DIVfield( input2 + col.label + not_set, input );
		++shown;
		if( z<columns.length - 1 && (col.type=='textarea' || shown % mod==mod - 1) ) f_html += '</div><div class="fields">';
	}
	columns.pop(); // remove is_deleted
	f_html = DIVc( 'content', FORM( DIVc( 'fields', f_html ) ) );
	f_html = DIVc( 'header', "Filter for " + rec_type ) + f_html;
	var footer = DIVcp( 'ui yellow button',
		ZZ( '@click', "clear_filter_event" ),
		Icp( 'undo icon', '' ) + 'Clear'
	);
	footer += DIVcp( 'ui green button',
		ZZ( '@click', "close_filter_event" ),
		Icp( 'check icon', '' ) + 'Close'
	);
	f_html += DIVc( 'actions', footer );
	f_html = DIVidc( id + '__filter', f_html, 'ui modal' );
	// init
	var el_id = '#' + id;
	var el = $( el_id );
	if(!el.length) throw('Element '+el_id+ ' not found');
	el.html( g_html + f_html );
	var el_filter = $( el_id + '__filter' );
	var el_parent = el.parent();
	// load from local storage
	var saved_columns_shown = H.LocalLoad( local_storage_key );
	// default options for vue
	var default_opts = {
		el: el_parent[ 0 ],
		data: {
			can_add: true,
			columns: columns,
			rows: rows,
			term: '',
			rows_index: {}, // row.id to index position
			limit: 10,
			offset: 0,
			count: 0,
			loading: 0,
			column_by_key: column_by_key,
			columns_shown: (saved_columns_shown || first_n_columns),
			filter_val: filter_val, // {key:filter_val}
			filter_state: filter_state, // {key:bool}
			filter_count: 0, // cache
			order: order, // ["+key","-key"] + ascending, - descending
			order_hash: {}, // cache
			order_idx: {} // cache
		},
		computed: {
			col_shown_hash: function() {
				var res = {};
				var cs = this.columns_shown;
				for( var z in cs ) res[ cs[ z ] ] = true;
				H.LocalStore( local_storage_key, cs );
				return res;
			},
			total_page: function() {
				return Math.floor( (this.count - 1) / this.limit ) + 1
			},
			pager: function() {
				var span = 3;
				var max_space = this.page + span;
				var min_space = this.page - span;
				var arr = [];
				for( var z = min_space; z<=max_space; ++z ) arr.push( z );
				return arr;
			},
			page: function() {
				return Math.floor( this.offset / this.limit ) + 1;
			}
		},
		created: function() {
			// create reverse search index
			H.ExecIfFunction( default_events.OnGridCreated, ret );
		},
		mounted: function() {
			el = $( el_id );
			el_parent = el.parent();
			el_filter = $( el_id + '__filter' );
			H.ExecIfFunction( default_events.OnGridMounted, ret );
			H.InitFileUpload( el );
			el.find( '.ui.dropdown' ).dropdown();
			el_filter.modal( { closable: false, duration: 0 } );
			el_filter.find( '.ui.checkbox' ).checkbox();
			el_filter.find( 'select.dropdown' ).dropdown();
			if( Object.keys( filter_state ).length ) this.search_event();
			this.recompute_filter_count();
			this.recompute_order_hash();
		},
		updated: function() {
			H.ExecIfFunction( default_events.OnGridUpdated, ret );
			//H.InitFileUpload( el );
		},
		destroyed: function() {
			H.ExecIfFunction( default_events.OnGridDestroyed, ret );
		},
		methods: {
			recompute_filter_count: function() {
				var states = this.filter_state;
				var count = 0;
				for( var z in states ) if( states[ z ] ) ++count;
				this.filter_count = count;
			},
			recompute_order_hash: function() {
				var res = { '__ob__': true };
				for( var z = 0; z<this.order.length; ++z ) {
					var str = this.order[ z ], ch = str[ 0 ];
					if( ch=='+' || ch=='-' ) {
						var key = str.substr( 1 );
						res[ key ] = ch;
						Vue.set( this.order_hash, key, ch );
						Vue.set( this.order_idx, key, z + 1 );
					}
				}
				for( z in this.order_hash ) {
					if( !res[ z ] ) {
						Vue.delete( this.order_hash, z );
						Vue.delete( this.order_idx, z );
					}
				}
			},
			edit_event: function( row_id ) {
				MASTER_FORM.edit_event( row_id );
			},
			key_search_event: function( e ) {
				if( !H.IsEnter( e ) ) {
					if( !H.IsEscape( e ) ) return;
					this.term = '';
				}
				this.offset = 0;
				this.search_event();
			},
			search_event: function() {
				var values = {
					a: 'search',
					term: this.term,
					limit: this.limit,
					offset: this.offset,
					filter: this.collect_filter(),
					order: JSON.stringify( this.order )
				};
				this.loading += 1;
				var self = this;
				H.ExecIfFunction( default_events.OnBeforeSearch, values );
				H.Post( values, function( res ) {
					if( H.HasAjaxErrors( res ) ) return;
					var nrs = res.rows;
					var ors = self.rows;
					while( ors.length ) ors.pop();
					for( var z in nrs ) ors.push( nrs[ z ] );
					self.count = res.count;
					self.offset = res.offset;
					self.limit = res.limit;
					self.loading -= 1;
					H.ExecIfFunction( default_events.OnAfterSearch, res );
				} )
			},
			page_event: function( page, mod ) {
				var new_page = this.page;
				if( mod== -1 ) new_page -= 1;
				else if( mod== +1 ) new_page += 1;
				else if( mod== +2 ) new_page = this.total_page;
				else new_page = page;
				if( new_page>1 && new_page>this.total_page ) new_page = this.total_page;
				if( new_page<1 ) new_page = 1;
				this.offset = (new_page - 1) * this.limit;
				this.search_event();
			},
			filter_form_event: function() {
				el_filter.modal( 'show' );
			},
			collect_filter: function() {
				var states = this.filter_state;
				var vals = this.filter_val;
				var res = {};
				for( var z in states ) {
					if( states[ z ] ) {
						var val = vals[ z ];
						// convert date to integer (unix timestamp)
						//console.log( this.column_by_key[ z ].type );
						switch( this.column_by_key[ z ].type ) {
							case 'date':
							case 'datetime':
								var iso = '', prefix = '', nv = '';
								var append_nv = function() {
									if( !iso ) return;
									if( nv!='' && nv[ nv.length - 1 ]!=' ' ) nv += ' ';
									nv += prefix + H.Moment( iso ).unix();
									prefix = '';
									iso = '';
								};
								for( var y = 0; y<val.length; ++y ) {
									var ch = val[ y ];
									if( ch=='<' || ch=='>' ) {
										append_nv(); // if they forgot to add space
										prefix += ch;
									}
									else if( ch=='=' ) prefix += ch;
									else if( (ch=='-' || ch=='T' || ch==':') || (ch>='0' && ch<='9') ) iso += ch;
									else if( ch==' ' ) append_nv();
								}
								append_nv();
								res[ z ] = nv;
								break;
							default:
								res[ z ] = val;
						}
					}
				}
				return JSON.stringify( res );
			},
			close_filter_event: function() {
				el_filter.modal( 'hide' );
				this.search_event();
				this.save_filters();
			},
			clean_filter_event: function() {
				this.clear_filter_event();
				this.search_event();
				this.save_filters();
			},
			save_filters: function() {
				H.LocalStore( local_storage_key + '|filter_val', this.filter_val );
				H.LocalStore( local_storage_key + '|filter_state', this.filter_state );
				this.recompute_filter_count();
			},
			clear_filter_event: function() {
				this.filter_state = Object.assign( {}, {} );
				H.LocalStore( local_storage_key + '|filter_state', this.filter_state );
				this.recompute_filter_count();
			},
			clear_order_event: function() {
				while( this.order.length ) this.order.pop();
				this.save_order();
				this.search_event();
			},
			save_order: function() {
				this.recompute_order_hash();
				H.LocalStore( local_storage_key + '|order', this.order );
			},
			click_order_event: function( key, e ) {
				var order = this.order;
				var add = e.shiftKey || e.ctrlKey;
				var next = '';
				for( var z = 0; z<order.length; ++z ) {
					var str = order[ z ];
					if( str.substr( 1 )==key ) {
						order.splice( z, 1 );
						next = str[ 0 ];
						break;
					}
				}
				if( !add ) while( order.length ) order.pop(); // clear everything if shift or ctrl not pressed
				if( next=='' ) next = '+';
				else if( next=='+' ) next = '-';
				else next = '';
				if( next ) order.push( next + key );
				this.save_order();
				this.search_event();
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
