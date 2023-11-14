<script>
    import { onMount } from 'svelte';
    let segment2;
    onMount(() => {
        console.log('onMount.AdminSubMenu')
        segment2 = window.location.pathname.split('/').pop() || '';
    });
</script>

<div class="admin_submenu">
    <nav class='menu'>
        <a href="/superAdmin/userManagement" class:active={segment2==='userManagement'}>Users</a>
<!--        <a href="/admin/accessLogs" class:active={segment2==='accessLogs'}>Access Logs</a>-->
<!--        <a href="/admin/files" class:active={segment2==='files'}>Files</a>-->
    </nav>
</div>
<style>
    .admin_submenu {
        position: relative;
        margin-top: -40px;
        margin-left: auto;
        margin-right: auto;
        display: flex;
        background-color: #EF4444;
        padding-bottom: 70px;
    }
    .admin_submenu .menu {
        width: 88%;
        margin-left: auto;
        margin-right: auto;
        display: flex;
        flex-direction: row;
    }
    .admin_submenu .menu a {
        flex-grow: 1;
        color: #475569;
        background-color: white;
        border-radius: 6px;
        filter: drop-shadow(0 10px 8px rgb(0 0 0 / 0.04)) drop-shadow(0 4px 3px rgb(0 0 0 / 0.1));
        text-transform: uppercase;
        font-size: 13px;
        text-decoration: none;
        font-weight: 700;
        padding-top: 11px;
        padding-bottom: 11px;
        text-align: center;
        margin-left: 20px;
    }
    .admin_submenu .menu a:hover {
        color: #EF4444;
    }
    .admin_submenu .menu a:nth-child(1) {
        margin-left: 0 !important;
    }
    .admin_submenu .menu .active {
        background-color : #334155 !important;
        color: white !important;
    }
</style>
