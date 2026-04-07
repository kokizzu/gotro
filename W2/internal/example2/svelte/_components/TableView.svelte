<script>
    // @ts-nocheck
    import { onMount } from 'svelte';
    import { datetime } from './formatter.js';

    import { Icon } from 'svelte-icons-pack';
    import { HiOutlinePencil } from 'svelte-icons-pack/hi';
    import {
        FaSolidAnglesLeft,
        FaSolidAnglesRight,
        FaSolidAngleLeft,
        FaSolidAngleRight,
        FaSolidFilter,
        FaSolidRotate
    } from 'svelte-icons-pack/fa';

    let {
        arrayOfArray = true,
        fields = [],
        rows = [],
        pager = $bindable({}),
        extraActions = [],
        onRefreshTableView = function( pager ) {
        console.log( 'TableView.onRefreshTableView', pager );
        },
        onEditRow = function( id, row ) {
        console.log( 'TableView.onEditRow', id, row );
        },
        widths = {},
        children
    } = $props();

    let deletedAtIdx = $state(-1);
    onMount( () => {
        console.log( 'onMount.TableView' );
        console.log( 'fields=', fields );
        console.log( 'rows=', rows );
        console.log( 'pager=', pager );
        console.log( 'extraActions=', extraActions );

        for( let z = 0; z<fields.length; z++ ) {
            let field = fields[ z ];
            // find deletedAt index
            if( field.name==='deletedAt' ) {
                console.log( 'deletedAtIdx', z );
                deletedAtIdx = z;
            }
            // empty all filter at beginning
            filtersMap[ field.name ] = '';
        }
        oldFilterStr = JSON.stringify( filtersMap );
    } );

    let oldFilterStr = $state('{}');
    let filtersMap = $state({});
    let newFilterStr = $derived(JSON.stringify(filtersMap));

    function filterKeyDown( event ) {
        if( event.key==='Enter' ) applyFilter();
    }

    function applyFilter() {
        console.log('masuk')
        let filters = {};
        for( let key in filtersMap ) {
            let value = filtersMap[ key ];
            if( value ) {
                filters[ key ] = value.split( '|' );
            }
        }
        console.log( 'filters=', filters );
        onRefreshTableView( {
            ...pager,
            filters: filters,
        } );
        oldFilterStr = newFilterStr;
    }

    function gotoPage( page ) {
        onRefreshTableView( {
            ...pager,
            page,
        } );
    }

    function changePerPage( perPage ) {
        onRefreshTableView( {
            ...pager,
            perPage,
        } );
    }

    function cell( row, i, field ) {
        if( arrayOfArray ) return row[ i ] || '';
        return row[ field.name ] || '';
    }

    let allowPrevPage = $derived((pager.page || 0) > 1);
    let allowNextPage = $derived((pager.page || 0) < (pager.pages || 0));

</script>

<section>
    <div class='action_options_container'>
        <div class='left'>
            {@render children?.()}
            <button class='action_btn' disabled={oldFilterStr===newFilterStr} onclick={applyFilter}>
                <Icon color={oldFilterStr === newFilterStr ? '#5C646F' : '#FFF'} size={17} src={FaSolidFilter} />
                <span>Apply Filter</span>
            </button>
            <button class='action_btn' onclick={() => gotoPage(pager.page)}>
                <Icon color='#FFF' size={17} src={FaSolidRotate} />
                <span>Refresh</span>
            </button>
        </div>
    </div>
    <div class='table_container'>
        <table class='table_users'>
            <thead>
            <tr>
                {#each fields as field}
                    {#if field.name==='id'}
                        <th class='col_action'>Action</th>
                    {:else}
                        <th class='table_header'
                            style='{widths[field.name] ? "min-width: "+widths[field.name] : ""}'>
                            <label for='th_{field.name}'>{field.label}</label><br />
                            <input id='th_{field.name}'
                                   title='separate with pipe for multiple values, for example:
  >=100|<50|61|72 will show values greater equal to 100, OR less than 50, OR exactly 61 OR 72
    filtering with greater or less than will only show correct result if data is saved as number
    currently price and size NOT stored as number
  <>abc* will show values NOT started with abc*
  abc*def|ghi will show values started with abc ends with def OR exactly ghi
  *jkl* will show values containing jkl substring
multiple filter from other fields will do AND operation'
                                   type='text'
                                   class='input_filter'
                                   bind:value={filtersMap[field.name]}
                                   onkeydown={filterKeyDown}
                            />
                        </th>
                    {/if}
                {/each}
            </tr>
            </thead>
            <tbody>
            {#each rows as row, no}
                <tr class:deleted={row[deletedAtIdx] > 0}>
                    {#each fields as field, i}
                        {#if field.name==='id'}
                            <td class='col_action'>
                                <div>
                                    <button class='action' title='Edit' onclick={() => onEditRow(cell(row,i,field), row)}>
                                        <Icon src={HiOutlinePencil} />
                                    </button>
                                    {#each extraActions as action}
                                        {#if action.link}
                                            <a href='{action.link(row)}' class='action' target='_blank' title='{action.label || ""}'>
                                                <Icon src={action.icon} />
                                            </a>
                                        {:else}
                                            <button class='action' title='{action.label || ""}' onclick={() => action.onClick(row)}>
                                                <Icon src={action.icon} />
                                            </button>
                                        {/if}
                                    {/each}
                                </div>
                            </td>
                        {:else if field.inputType==='checkbox'}
                            <td class='table_data'>{!!cell( row, i, field )}</td>
                        {:else if field.inputType==='datetime' || field.name==='deletedAt'}
                            <td class='table_data'>{datetime( cell( row, i, field ) )}</td>
                        {:else if field.inputType==='number'}
                            <td>{(cell( row, i, field ) || 0).toLocaleString()}</td>
                        {:else}
                            <td class='table_data'>{cell( row, i, field )}</td>
                        {/if}
                    {/each}
                </tr>
            {/each}
            </tbody>
        </table>
    </div>
    <div class='pages_set'>
        <div class='page_and_rows_count'>
            <span>Page {pager.page} of {pager.pages},</span>
            <input
                    bind:value={pager.perPage}
                    class='perPage'
                    id='perPage'
                    min='0'
                    onchange={() => changePerPage(pager.perPage)}
                    type='number'
            />
            <span>rows per page.</span>
        </div>

        <p>Total: {pager.countResult | 0}</p>

        <div class='pagination'>
            <button disabled={!allowPrevPage} onclick={() => gotoPage(1)} title='Go to first page'>
                <Icon color={!allowPrevPage ? '#5C646F' : '#FFF'} size={18} src={FaSolidAnglesLeft} />
            </button>
            <button disabled={!allowPrevPage} onclick={() => gotoPage(pager.page - 1)} title='Go to previous page'>
                <Icon color={!allowPrevPage ? '#5C646F' : '#FFF'} size={18} src={FaSolidAngleLeft} />
            </button>
            <button disabled={!allowNextPage} onclick={() => gotoPage(pager.page + 1)} title='Go to next page'>
                <Icon color={!allowNextPage ? '#5C646F' : '#FFF'} size={18} src={FaSolidAngleRight} />
            </button>
            <button disabled={!allowNextPage} onclick={() => gotoPage(pager.pages)} title='Go to last page'>
                <Icon color={!allowNextPage ? '#5C646F' : '#FFF'} size={18} src={FaSolidAnglesRight} />
            </button>
        </div>
    </div>
</section>

<style>
    .action_options_container {
        display         : flex;
        flex-direction  : row;
        justify-content : space-between;
        align-items     : center;
    }

    .action_options_container .left {
        display     : flex;
        gap         : 8px;
        align-items : center;
    }

    .table_container {
        overflow-x : auto;
    }

    .table_container::-webkit-scrollbar-thumb {
        background-color : #EF4444;
        border-radius    : 4px;
    }

    .table_container::-webkit-scrollbar-thumb:hover {
        background-color : #EC6262;
    }

    .table_container::-webkit-scrollbar {
        height : 10px;
    }

    .table_container::-webkit-scrollbar-track {
        background-color : transparent;
    }

    .table_container .table_users {
        margin           : 10px 1px 10px 0;
        background-color : white;
        border-collapse  : collapse;
        font-size        : 14px;
        width            : 100%;
        color            : #475569;
    }

    .table_container .table_users th {
        color   : #6366F1;
        border  : 1px solid #CBD5E1;
        padding : 10px 7px;
    }

    .table_container .table_users td {
        border  : 1px solid #CBD5E1;
        padding : 2px 3px 2px 3px;
    }

    .table_users .table_header {
        text-align  : left;
        white-space : nowrap;
    }

    .table_users .col_action {
        text-align : center;
        padding    : 0 10px;
        width      : fit-content;
    }

    .table_users td.col_action {
        text-align : center;
        padding    : 0 10px;
        width      : fit-content;
    }

    .table_users td.col_action div {
        display        : flex;
        flex-direction : row;
        align-items    : center;
    }

    .table_users tr, .table_users td {
        height : 2em;
    }

    .table_users .col_action .action {
        background    : none;
        color         : #475569;
        text-align    : center;
        cursor        : pointer;
        padding       : 8px 10px;
        height        : fit-content;
        width         : fit-content;
        border        : none;
        display       : inline;
        border-radius : 5px;
    }

    .table_users .col_action .action:hover {
        color            : #EF4444;
        background-color : #F0F0F0;
    }

    .table_users .deleted td {
        text-decoration           : line-through;
        text-decoration-color     : rgba(239, 68, 68, 0.5);
        text-decoration-thickness : 5px;
    }

    .pages_set {
        padding-top     : 10px;
        display         : flex;
        flex-direction  : row;
        align-items     : center;
        justify-content : space-between;
        align-content   : center;
        font-size       : 15px;
        color           : #161616;
    }

    .pages_set .page_and_rows_count {
        display        : flex;
        flex-direction : row;
        align-content  : center;
        align-items    : center;
    }

    .pages_set .page_and_rows_count .perPage {
        margin        : auto 5px;
        width         : 4em;
        border        : 1px solid #CBD5E1;
        padding       : 5px;
        font-size     : 14pt;
        font-weight   : bold;
        text-align    : center;
        color         : #161616;
        outline-color : #6366F1;
    }

    .pages_set .page_and_rows_count .perPage::-webkit-inner-spin-button,
    .pages_set .page_and_rows_count .perPage::-webkit-outer-spin-button {
        opacity : 1;
    }

    .pagination button {
        color            : white;
        background-color : #6366F1;
        padding          : 8px;
        border-radius    : 5px;
        filter           : drop-shadow(0 10px 8px rgb(0 0 0 / 0.04)) drop-shadow(0 4px 3px rgb(0 0 0 / 0.1));
        margin-left      : 4px;
        cursor           : pointer;
        border           : none;
    }

    .pagination button:hover {
        background-color : #7E80F1;
    }

    .pagination button:disabled {
        cursor     : not-allowed;
        border     : 1px solid #CBD5E1 !important;
        background : none !important;
        color      : #5C646F;
        filter     : drop-shadow(0 10px 8px rgb(0 0 0 / 0.04)) drop-shadow(0 4px 3px rgb(0 0 0 / 0.1));
    }

    .pagination button:disabled:hover {
        background : none;
    }

    input.input_filter {
        width      : 0;
        min-width  : 100%;
        box-sizing : border-box;
    }
</style>
