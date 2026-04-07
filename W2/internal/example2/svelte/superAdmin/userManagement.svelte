<script>
    // @ts-nocheck
    import Menu from '../_components/Menu.svelte';
    import AdminSubMenu from '../_components/AdminSubMenu.svelte';
    import ProfileHeader from '../_components/ProfileHeader.svelte';
    import Footer from '../_components/Footer.svelte';
    import TableView from '../_components/TableView.svelte';
    import ModalForm from '../_components/ModalForm.svelte';
    import Growl from '../_components/Growl.svelte';

    import { Icon } from 'svelte-icons-pack';
    import { FaSolidCirclePlus } from 'svelte-icons-pack/fa';

    let segments = {/* segments */};
    let fields = [/* fields */];
    let users = $state([/* users */]);
    let pager = $state({/* pager */});

    // $: console.log( users, fields, pager );
    let growl = Growl;
    let apiModule;

    async function loadApi() {
        apiModule ??= await import('../jsApi.GEN.cjs');
        return apiModule;
    }

    // return true if got error
    function handleResponse( res ) {
        console.log( res );
        if( res.error ) {
            growl.showError(res.error );
            return true;
        }
        if( res.users && res.users.length ) users = res.users;
        if( res.pager && res.pager.page ) pager = res.pager;
    }

    async function refreshTableView( pagerIn ) {
        // console.log( 'pagerIn=',pagerIn );
        const { SuperAdminUserManagement } = await loadApi();
        await SuperAdminUserManagement( {
            pager: pagerIn,
            cmd: 'list',
        }, function( res ) {
            handleResponse( res );
        } );
    }

    let form = ModalForm; // for lookup

    async function editRow( id, row ) {
        const { SuperAdminUserManagement } = await loadApi();
        await SuperAdminUserManagement( {
            user: {id},
            cmd: 'form',
        }, function( res ) {
            if( !handleResponse( res ) )
                form.showModal( res.user );
        } );
    }

    function addRow() {
        form.showModal( {id: ''} );
    }

    async function saveRow( action, row ) {
        let user = {...row};
        if( !user.id ) user.id = '0';
        const { SuperAdminUserManagement } = await loadApi();
        await SuperAdminUserManagement( {
            user: user,
            cmd: action,
            pager: pager, // force refresh page, will be slow
        }, function( res ) {
            if( handleResponse( res ) ) {
                return form.setLoading( false ); // has error
            }
            form.hideModal(); // success
        } );
    }
</script>


<section class='dashboard'>
    <Menu access={segments} />
    <div class='dashboard_main_content'>
            <ProfileHeader />
            <AdminSubMenu />
        <div class='content'>
            <ModalForm {fields}
                       rowType='User'
                       bind:this={form}
                       onConfirm={saveRow}
            />
            <section class='tableview_container'>
                <TableView {fields}
                           bind:pager={pager}
                           rows={users}
                           onRefreshTableView={refreshTableView}
                           onEditRow={editRow}
                >
                    <button onclick={addRow} class='action_btn'>
                        <Icon size={17} color='#FFF' src={FaSolidCirclePlus} />
                        <span>Add</span>
                    </button>
                </TableView>
            </section>
        </div>
        <Footer />
    </div>
</section>

<style>
</style>
