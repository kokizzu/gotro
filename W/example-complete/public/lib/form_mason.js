function FormMason( id, rec_type, fields, events, vue_options ) {
	events = events || {};
	vue_options = vue_options || {};
	if( !id || 'string'!= typeof id ) throw('FormMason: first argument should be an html (id)');
	if( !rec_type || 'string'!= typeof rec_type ) throw('FormMason: second argument should be a string (rec_type)');
	if( !fields || 'object'!= typeof fields ) throw('FormMason: third argument should be an array of (fields)');
	var z, y, ret;
	// initialize callbacks
	var default_events = {
		OnFormCreated: function() {H.Trace( arguments )},
		OnFormMounted: function() {H.Trace( arguments )},
		OnFormUpdated: function() {H.Trace( arguments )},
		OnFormDestroyed: function() {H.Trace( arguments )},
		OnVueCreated: function() {H.Trace( arguments )},
		OnBeforeEdit: function() {H.Trace( arguments )},
		OnAfterEdit: function() {H.Trace( arguments )},
		OnBeforeSave: function() {H.Trace( arguments )},
		OnAfterSave: function() {H.Trace( arguments )},
		OnBeforeReset: function() {H.Trace( arguments )},
		OnAfterReset: function() {H.Trace( arguments )},
		OnBeforeCancel: function() {H.Trace( arguments )},
		OnAfterCancel: function() {H.Trace( arguments )}
	};
	for( z in events ) if( events.hasOwnProperty( z ) ) default_events[ z ] = events[ z ];
	// initialize generate form default html
	var html = '';
	var mod = 4;
	var default_row = {};
	var dropdowns = {};
	var editors = {};
	for( z = 0; z<fields.length; ++z ) {
		var field = fields[ z ];
		// TODO: switch case per type, change to INPUT, TEXTAREA, DATE, TIME, etc
		var input = '';
		var typ = field.type;
		var key = field.key;
		var props = ZZ( 'name', key ) + ZZ( 'v-model', 'row.' + key ) + ZZ( 'placeholder', field.placeholder || field.tooltip ) + ZZ( 'title', field.placeholder || field.tooltip );
		var def_val = '';
		switch( typ ) {
			case 'textarea':
				input = '<textarea' + props + ' rows="3"></textarea>';
				break;
			case 'bool':
				/** @namespace field.sub_type */
				var sub_typ = field.sub_type || '"Yes":"No"';
				input = '<input' + props + ' type="checkbox" tabindex="0" class="hidden">';
				input += '<label>{{ !!row.' + key + ' ? ' + sub_typ + ' }}</label>';
				input = DIVc( 'ui toggle checkbox', input );
				break;
			case 'integer':
				input = '<input' + props + 'type="number" />';
				def_val = 0;
				break;
			case 'float2':
				input = '<input' + props + 'type="number" min="0.00" step="0.01" />';
				def_val = 0;
				break;
			case 'float': // latitude longitude
				input = '<input' + props + 'type="number" step="0.000001" />';
				def_val = 0;
				break;
			case 'multiselect':
				props += ZZ( 'multiple', 'multiple' ); // fallthrough
			case 'select':
				input = '<select class="dropdown"' + props + '><option v-for="(label,idx) in DS.' + field.sub_type + '" :value="idx" :selected="row.' + key + ' == idx ? \'selected\' : \'\'">{{label}}</option></select>';
				dropdowns[ key ] = field;
				break;
			case 'datetime':
			case 'date':
				input = '<div class="ui right labeled input">' +
					'<input' + props + ' style="width: 1em" class="date-range-picker not_initialized" ' + ZZ( 'data-type', typ ) + ' type="text" />' +
					'<div class="ui label">{{ row.' + key + ' | ' + typ + '}}</div>' +
					'</div>';
				break;
			case 'json':
				input = '<textarea' + props + ' rows="2" cols="32" style="width: 240px" class="code json not_initialized"></textarea>';
				editors[ key ] = field;
				break;
			case 'emails':
			case 'url':
			case 'phone':
				typ = 'text'; // fallthrough
			default:
				input = '<input ' + props + ZZ( 'type', typ ) + ' />';
		}
		default_row[ key ] = def_val;
		html += DIVfield( field.label, input );
		if( z<fields.length - 1 && (z % mod==mod - 1) ) html += '</div><div class="fields">'
	}
	html = DIVc( 'content', FORM( DIVc( 'fields', html ) ) );
	html = DIVc( 'header', "{{row_id ? 'Edit '+rec_type+' Record #' + row_id : 'New '+rec_type+' Record' }}" ) + html;
	var footer = DIVcp( 'ui blue button',
		'@click="cancel_event"',
		Icp( 'arrow left icon', '' ) + 'Cancel'
	);
	footer += DIVcp( 'ui button',
		ZZ( 'v-if', 'can_restore_delete && row_id > 0' ) + ZZ( ':class', '{ red:(!row.is_deleted), orange:(!!row.is_deleted) }' ) + ZZ( '@click', "delete_restore_event" ),
		Icp( 'icon', ZZ( ':class', '{ remove:(!row.is_deleted), checkmark:(!!row.is_deleted) }' ) ) + "{{row.is_deleted ? 'Restore' : 'Delete'}}"
	);
	footer += DIVcp( 'ui yellow button',
		ZZ( '@click', "reset_event" ),
		Icp( 'undo icon', '' ) + 'Reset'
	);
	footer += DIVcp( 'ui green button',
		ZZ( '@click', "save_event" ),
		Icp( 'undo icon', '' ) + '{{row_id ? update_label : create_label}}'
	);
	html += DIVc( 'actions', footer );
	var el_id = '#' + id;
	var el = $( el_id );
	if( !el.length ) throw('FormMason: element not found: ' + el_id);
	el.prop( 'class', 'ui modal' ).html( html );
	var el_parent = el.parent();
	// init function
	var init_daterangepicker = function( el ) {
		el.find( 'input.date-range-picker.not_initialized' ).each( function() { // initialize file_upload
			var el = $( this );
			el.removeClass( 'not_initialized' );
			var opt = {
				showDropdowns: true,
				showWeekNumbers: true,
				timePickerIncrement: 1,
				timePicker12Hour: false,
				timePickerSeconds: true,
				singleDatePicker: true,
				format: 'X'
			};
			if( el.data( 'type' )=='datetime' ) opt.timePicker = true;
			el.daterangepicker( opt );
			el.on( 'apply.daterangepicker', function( ev, picker ) { // update vue from datetimepicker
				var epoch = picker.startDate.format( 'X' );
				H.TriggerVue( picker.element, epoch );
			} );
			el.on( 'show.daterangepicker', function( ev, picker ) { // #way1 reinit datetimepicker
				var date = H.Moment( picker.element.val() || moment().unix() );
				picker.setStartDate( date );
				picker.setEndDate( date )
			} )
		} );
	};
	var init_jsoneditor = function( el, vue ) {
		var aces = el.find( 'textarea.code.json' );
		var aceInit = function() {
			aces.ace( { theme: 'eclipse', lang: 'json', blockScrolling: Infinity } ).each( function( idx, editor ) {
				var el = $( editor );
				var key = el.prop( 'name' );
				var editor_prop = el.data( 'ace' ).editor;
				var ace = editor_prop.ace;
				ace.$blockScrolling = Infinity;
				ace.setReadOnly( el.prop( 'disabled' ) );
				ace.setOption( "maxLines", 10 );
				ace.setOption( "minLines", 2 );
				ace.on( 'change', function() {
					var val = ace.getValue();
					vue.row[ key ] = val;
					H.TriggerVue( el, val )
				} )
			} );
			aces.removeClass( 'not_initialized' );
		};
		if( $.ace ) aceInit();
		else H.LoadScripts( [ '/js/ace/ace.js', '/js/ace/mode-json.js', '/js/jquery-ace.js' ], aceInit );
	};
	// default options for vue
	var default_opts = {
		el: el_parent[ 0 ],
		data: {
			rec_type: (rec_type || id),
			fields: fields,
			dropdowns: dropdowns,
			editors: editors,
			row: {},
			orig_row: {},
			override_row: {}, // to replace selected value when form shown
			row_id: 0,
			can_restore_delete: true,
			action_suffix: '',
			create_label: 'Create',
			update_label: 'Update',
			save_action: 'save',
			save_action_suffix: '',
			form_action: 'form',
			form_action_suffix: '',
			restore_delete_action_suffix: '',
			DS: DS
		},
		created: function() {
			H.ExecIfFunction( default_events.OnFormCreated, ret );
		},
		mounted: function() {
			el = $( el_id );
			el_parent = el.parent();
			H.ExecIfFunction( default_events.OnFormMounted, ret );
			el.modal( { closable: false, duration: 0 } );
			el.find( '.ui.checkbox' ).checkbox();
			el.find( 'select.dropdown' ).dropdown( {
				onChange: function( value, text, $selectedItem ) { // update vue from dropdown
					var el = $selectedItem.parent().parent().find( 'select' );
					H.TriggerVue( el, value );
				}
			} );
			init_daterangepicker( el );
			if( Object.keys( this.editors ).length ) init_jsoneditor( el, this );
		},
		updated: function() {
			H.ExecIfFunction( default_events.OnFormUpdated, ret );
		},
		destroyed: function() {
			H.ExecIfFunction( default_events.OnFormDestroyed, ret );
		},
		methods: {
			// when edit triggered
			edit_event: function( row_id ) {
				this.row_id = row_id;
				var action = this.form_action + (this.form_action_suffix || this.action_suffix );
				var values = {
					id: row_id,
					a: action
				};
				H.ExecIfFunction( default_events.OnBeforeEdit, this );
				if( !row_id ) {
					this.row = Object.assign( {}, this.row, default_row, this.override_row );
					this.orig_row = Object.assign( {}, this.orig_row, default_row );
					el.modal( 'show' );
					H.ExecIfFunction( default_events.OnAfterEdit, this );
					console.log( this.row.parent_id, this.override_row )
					return
				}
				var self = this;
				H.Post( values, function( res ) {
					if( H.HasAjaxErrors( res ) ) return;
					res.password = '';
					// TODO: set and unset last_mod, diff the changes only
					self.row = Object.assign( {}, self.row, res );
					self.orig_row = Object.assign( {}, self.orig_row, res );
					// #way2 reinit select
					for( var z in self.dropdowns ) el.find( 'select[name="' + z + '"]' ).dropdown( 'set exactly', self.row[ z ] );
					// reinit ace editor
					for( z in self.editors ) {
						var val = JSON.parse( self.row[ z ] || '[]' );
						val = JSON.stringify( val, null, 2 );
						el.find( 'textarea[name="' + z + '"]' ).data( 'ace' ).editor.ace.setValue( val );
					}
					el.modal( 'show' );
					el.css( 'margin-top', 0 );
					el.css( 'top', '1em' );
					H.ExecIfFunction( default_events.OnAfterEdit, self );
				} );
			},
			// collecting changes
			collect_changes: function( a ) {
				var id = this.row_id;
				var changes = { a: a, id: id };
				var changed = 0;
				var row = this.row, orig = this.orig_row;
				for( var z in orig ) {
					if( !orig.hasOwnProperty( z ) ) continue;
					var rz = row[ z ];
					var oz = orig[ z ];
					//console.log( z, rz, oz );
					if( rz==oz ) continue;
					++changed;
					changes[ z ] = rz;
				}
				// TODO: find reactive way to show changes, instead of this
				// TODO: make this more complete just like FormBuilder
				changes._changed = changed;
				console.log( changes );
				return changes;
			},
			// when save button pressed
			save_event: function() {
				// TODO: confirm changes
				var action = this.save_action + (this.save_action_suffix || this.action_suffix );
				var values = this.collect_changes( action );
				if( !values._changed ) return H.GrowlInfo( 'Nothing changed? Use "Cancel" button to close dialog' );
				this.update_row_event( values )
			},
			// when reset button pressed
			reset_event: function() {
				H.ExecIfFunction( default_events.OnBeforeReset, this );
				// TODO: check reset
				this.row = Object.assign( {}, this.orig_row );
				H.ExecIfFunction( default_events.OnAfterReset, this );
			},
			// when delete or restore button pressed
			delete_restore_event: function() {
				if( !this.can_restore_delete ) return H.GrowlError( 'Unable to restore/delete' );
				// TODO: confirm changes and confirm delete
				// TODO: refactor, most parts are similar to save_event
				var action = (this.row.is_deleted ? 'restore' : 'delete') + (this.restore_delete_action_suffix || this.action_suffix );
				var values = this.collect_changes( action );
				this.update_row_event( values )
			},
			update_row_event: function( values ) {
				var cancel = H.ExecIfFunction( default_events.OnBeforeSave, values );
				if( cancel ) return;
				H.Post( values, function( res ) {
					if( H.HasAjaxErrors( res ) ) return;
					H.GrowlInfo( res.info + values.a + ' completed!' );
					el.modal( 'hide' );
					H.ExecIfFunction( default_events.OnAfterSave, res );
				} );
			},
			cancel_event: function() {
				H.ExecIfFunction( default_events.OnBeforeCancel, this );
				// TODO: confirm close if changed
				el.modal( 'hide' );
				H.ExecIfFunction( default_events.OnAfterCancel, this );
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
	var rev = {};
	for( z in ret.fields ) {
		var rdf = ret.fields[ z ];
		rev[ rdf.key ] = rdf;
	}
	ret.column_by_key = rev;
	H.ExecIfFunction( default_events.OnVueCreated );
	return ret;
}