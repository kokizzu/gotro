import { writable } from 'svelte/store';
// translation related

export let isSideMenuOpen = writable(false); // Side Menu

export let langOptions = {
    en: 'EN',
    id: 'ID', // add more languages here
};
