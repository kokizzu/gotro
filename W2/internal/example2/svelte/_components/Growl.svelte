<script>
    import { Icon } from 'svelte-icons-pack';
    import { FaSolidCircleInfo, FaSolidCircleCheck, FaSolidTriangleExclamation, FaSolidCircleXmark } from 'svelte-icons-pack/fa';

    let message = $state('This is growl, XD');
    let growlType = $state('info'); // info, warning, error, success
    let isShow = $state(false);
    let icon = $state(FaSolidCircleInfo);

    function show( msg, typ, ico ) {
        icon = ico;
        growlType = typ;
        isShow = true;
        message = msg;
        setTimeout( () => {
            isShow = false;
        }, 3000 );
    }

    export const showInfo = function( msg ) {
        console.log( 'grow.showInfo', msg );
        show( msg, 'info', FaSolidCircleInfo );
    };

    export const showWarning = function( msg ) {
        console.log( 'grow.showWarning', msg );
        show( msg, 'warning', FaSolidTriangleExclamation );
    };

    export const showError = function( msg ) {
        console.log( 'grow.showError', msg );
        show( msg, 'error', FaSolidCircleXmark );
    };

    export const showSuccess = function( msg ) {
        console.log( 'grow.showSuccess', msg );
        show( msg, 'success', FaSolidCircleCheck );
    };
</script>

<div class={`growl ${growlType} ${isShow ? '':'hidden'}`}>
    <Icon className='icon_growl' size={20} color='#FFF' src={icon} />
    <span>{message}</span>
</div>

<style>
    .growl.hidden {
        display : none
    }

    .growl {
        position       : fixed;
        display        : flex;
        flex-direction : row;
        align-items    : center;
        gap            : 8px;
        font-size      : 14px;
        bottom         : 20px;
        right          : 20px;
        padding        : 10px 20px;
        border-radius  : 3px;
        box-shadow     : 0px 4px 24px 0px rgba(0, 0, 0, 0.25);
        z-index        : 9999;
        min-width      : 120px;
        max-width      : 350px;
        height         : fit-content;
    }

    :global(.icon_growl) {
        flex-shrink : 0;
    }

    .growl span {
        flex-shrink : 1;
    }

    .growl.info {
        background-color : #1080E8;
        color            : #FFF;
    }

    .growl.success {
        background-color : #059669;
        color            : #FFF;
    }

    .growl.error {
        background-color : #EF4444;
        color            : #FFF;
    }

    .growl.warning {
        background-color : #D97706;
        color            : #FFF;
    }
</style>
