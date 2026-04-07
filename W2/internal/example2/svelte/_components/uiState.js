import { writable } from 'svelte/store';
// translation related

export const isSideMenuOpen = writable(false); // Side Menu

export const langOptions = {
    en: 'EN',
    id: 'ID', // add more languages here
};
