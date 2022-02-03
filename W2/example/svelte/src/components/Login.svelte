<script>
	import { onMount } from 'svelte';
	import { profile } from '/src/store/profile.js'

	let error;
	let link;

	onMount(async () => {
		if (profile.user) return;
		let res = await fetch('/api/UserExternalLogin?provider=google')
		console.log(res)
		if (res.status === 200) {
			const json = await res.json()
			link = json.link
			error = json.error
			return
		}
		error = res.statusText
	})
</script>

{#if link}
	<a href="{link}">Login Google</a>
{:else}
	Failed fetch {error}
{/if}
