<script lang="ts">
	import {APIs, LastUpdatedAt} from "./api";
	import {profile} from "/src/store/profile"
	import Login from "../components/Login.svelte";
	
	let apis = []
	
	for (let z in APIs) {
		let prop = APIs[z];
		apis.push({route: z, prop: prop})
	}

</script>

<main>
	<h1>API List</h1>
	
	{#if !$profile.user}
		<Login/>
	{:else}
		<a href="/dashboard">Dashboard</a>
	{/if}
	
	<h2>List {new Date( LastUpdatedAt * 1000 )}</h2>
	<a id="top" href="api.js">raw JS</a>
	<br/>
	
	<table border="1">
		<thead>
		<tr>
			<th>No</th>
			<th>Name</th>
			<th>In Prop</th>
			<th>Out Prop</th>
			<th>Possible Errors</th>
			<th>Read Tables</th>
			<th>Write Tables</th>
			<th>Stats Tables</th>
			<th>3rd Party Deps</th>
		</tr>
		</thead>
		<tbody>
		{#each apis as { route, prop }, no}
			<tr>
				<td class="r">{1+no}</td>
				<td><a href="#{route}">{route}</a></td>
				<td class="r">{Object.keys( prop.in ).length}</td>
				<td class="r">{Object.keys( prop.out ).length}</td>
				<td class="r">{prop.err.length}</td>
				<td>{prop.read.join(', ')}</td>
				<td>{prop.write.join(', ')}</td>
				<td>{prop.stat.join(', ')}</td>
				<td>{prop.deps.join(', ')}</td>
			</tr>
		{/each}
		</tbody>
	
	</table>
	
	{#each apis as { route, prop }}
		<a id="{route}" href="/api/{route}">{route}</a>
		| <a href="#top">back to top</a>
		<br/>
		<pre>
			{JSON.stringify( prop, null, 3 )}
		</pre>
	{/each}

</main>

<style>
    :root {
        font-family : -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
        Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    }

    pre {
        /*white-space : pre-line;*/
        /*white-space : pre-wrap;*/
        border : 1px solid gray;
    }

    table td, table th {
	     padding-left: 4px;
	     padding-right: 4px;
    }
    
    main {
        text-align : left;
        padding    : 1em;
        margin     : 0 auto;
    }

    h1 {
        color          : #FF3E00;
        text-transform : uppercase;
        font-size      : 4rem;
        font-weight    : 100;
        line-height    : 1.1;
        margin         : 2rem auto;
        max-width      : 14rem;
    }


    @media (min-width : 480px) {
        h1 {
            max-width : none;
        }
    }

    td.r {
        text-align : right
    }
</style>
