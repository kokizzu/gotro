<script>
  // should generate index.html
  import Button from './_components/button.svelte'
  import Links from './_components/links.svelte'
  
  let arr = [/* arr */];
	let obj = {/* obj */};
</script>

<Button title="something"></Button>

<div>
	{JSON.stringify(arr)}
</div>
<div>
	{JSON.stringify(obj)}
</div>
<Links></Links>
