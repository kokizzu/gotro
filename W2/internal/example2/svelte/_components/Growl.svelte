<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import FaSolidInfoCircle from 'svelte-icons-pack/fa/FaSolidInfoCircle'; // info
    import FaSolidCheckCircle from 'svelte-icons-pack/fa/FaSolidCheckCircle'; // success
    import FaSolidExclamationTriangle from 'svelte-icons-pack/fa/FaSolidExclamationTriangle'; // warning
    import FaSolidTimesCircle from 'svelte-icons-pack/fa/FaSolidTimesCircle'; // error

    let message = 'This is growl, XD';
    let growlType = 'info'; // info, warning, error, success
    let isShow = false;
    let icon = FaSolidInfoCircle;

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
        show( msg, 'info', FaSolidInfoCircle );
    };

    export const showWarning = function( msg ) {
        console.log( 'grow.showWarning', msg );
        show( msg, 'warning', FaSolidExclamationTriangle );
    };

    export const showError = function( msg ) {
        console.log( 'grow.showError', msg );
        show( msg, 'error', FaSolidTimesCircle );
    };

    export const showSuccess = function( msg ) {
        console.log( 'grow.showSuccess', msg );
        show( msg, 'success', FaSolidCheckCircle );
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
