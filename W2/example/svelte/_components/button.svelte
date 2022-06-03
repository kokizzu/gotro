<script>
	export let title = ''

	// since _components starts with underscore, should not generate .html automatically

	function hitApi() {
		fetch('/api/Health')
			.then(res => res.json())
			.then(data => alert(JSON.stringify(data, null, 2)))
	}
</script>

<button on:click={hitApi}>{title}</button>
