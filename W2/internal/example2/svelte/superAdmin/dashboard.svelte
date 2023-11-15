<script>
    // @ts-nocheck
    import Menu from '../_components/Menu.svelte';
    import AdminSubMenu from '../_components/AdminSubMenu.svelte';
    import ProfileHeader from '../_components/ProfileHeader.svelte';
    import Footer from '../_components/Footer.svelte';

    let segments = {/* segments */};

</script>


<section class='dashboard'>
    <Menu access={segments} />
    <div class='dashboard_main_content'>
        <ProfileHeader></ProfileHeader>
        <AdminSubMenu></AdminSubMenu>
        <div class='content'>
            <section class='tableview_container'>
                dashboard admin
            </section>
        </div>
        <Footer></Footer>
    </div>
</section>

<style>
</style>