<script>
    // @ts-nocheck
    import {isSideMenuOpen, langOptions} from './uiState.js';
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import FaSolidBars from 'svelte-icons-pack/fa/FaSolidBars';
    import {onMount} from 'svelte';
    import {UserUpdateProfile} from '../jsApi.GEN.js';

    export let user = null;

    function openSideMenu() {
        isSideMenuOpen.set(!$isSideMenuOpen);
    }

    let selectedLanguage = '';
    onMount(() => {
        console.log('onMount.ProfileHeader')
        selectedLanguage = 'EN';
    });


</script>

<header class='profile_header'>
    <nav class='navbar'>
        <div class='label_menu'>
            <button on:click|preventDefault={openSideMenu}>
                <Icon color='#FFF' size={20} src={FaSolidBars}/>
            </button>
            <p>DASHBOARD</p>
        </div>
        <div class='right_nav'>
            <select bind:value={selectedLanguage} id='lang' name='lang'>
                {#each Object.values(langOptions) as lang}
                    <option value={lang}>{lang}</option>
                {/each}
            </select>
            <button class='profile_button'>
                <img alt='profile' src='/assets/img/team-1-200x200.jpg'/>
            </button>
        </div>
    </nav>
</header>

<style>
    .profile_header {
        background-color: #EF4444;
        height: fit-content;
        position: relative;
        padding: 18px 50px 80px 50px;
        display: flex;
        flex-direction: column;
    }

    .profile_header .navbar {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        align-items: center;
    }

    .profile_header .navbar .label_menu {
        display: flex;
        flex-direction: row;
        align-items: center;
        color: white;
    }

    .profile_header .navbar .label_menu button {
        padding: 10px;
        border: none;
        background: none;
        border-radius: 5px;
        font-size: 14px;
        color: white;
        cursor: pointer;
    }

    .profile_header .navbar .label_menu button:hover {
        background-color: rgba(255, 255, 255, 0.3);
    }

    .profile_header .navbar .label_menu p {
        font-size: 17px;
        font-weight: 600;
        line-height: 1.5rem;
        padding: 0;
        margin: 0 0 0 15px;
    }

    .profile_header .navbar .right_nav {
        display: flex;
        flex-direction: row;
        gap: 30px;
        align-items: center;
    }

    .profile_header .navbar #lang {
        padding: 6px 8px;
        width: 60px;
        border: none;
        border-radius: 5px;
        color: #1080E8;
        font-weight: 600;
        cursor: pointer;
        filter: drop-shadow(0 10px 8px rgb(0 0 0 / 0.04)) drop-shadow(0 4px 3px rgb(0 0 0 / 0.1));
    }

    .profile_header .navbar #lang:hover {
        background-color: #F1F5F9;
    }

    .profile_header .navbar #lang:focus {
        outline: none;
    }

    .profile_header .navbar .right_nav .profile_button {
        padding: 0;
        border: none;
        margin: 0;
        width: fit-content;
        height: fit-content;
        border-radius: 50%;
        filter: drop-shadow(0 10px 8px rgb(0 0 0 / 0.04)) drop-shadow(0 4px 3px rgb(0 0 0 / 0.1));
    }

    .profile_header .navbar .right_nav .profile_button:hover {
        box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.5);
        cursor: pointer;
    }

    .profile_header .navbar .right_nav .profile_button > img {
        width: 50px;
        height: 50px;
        border-radius: 9999px;
    }

    /* Responsive to mobile device */
    @media only screen and (max-width: 768px) {
        .profile_header {
            padding: 15px 20px 80px 20px !important;
        }

        .profile_header .navbar .label_menu button {
            padding: 0;
        }

        .profile_header .navbar .right_nav {
            gap: 10px;
        }

        .profile_header .navbar #lang {
            padding: 3px 5px;
            font-size: 12px;
            width: 50px
        }

        .profile_header .navbar .right_nav .profile_button > img {
            width: 40px;
            height: 40px;
        }
    }
</style>
