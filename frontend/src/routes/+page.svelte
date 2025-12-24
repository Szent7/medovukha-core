<script lang="ts">
	import FrameElement from '@templates/frameElement.svelte';
	import Sidebar from '@templates/sidebar.svelte';
	import Frame from '@templates/frame.svelte';
	import logo from '@assets/logo_small.svg';
</script>

<Sidebar />

<svelte:head>
	<title>Home</title>
	<link rel="icon" type="image/svg+xml" href={logo} />
</svelte:head>
<Frame {content} />

{#snippet header()}
	<img class="logoHeader" src={logo} alt="logo" />
	<h2>Medovukha panel</h2>
	<div class="releaseLabel">Alpha v0.0.1</div>
{/snippet}

{#snippet content()}
	<FrameElement
		--display="flex"
		--justify-content="space-between"
		--align-items="center"
		content={header}
	/>
{/snippet}

<style>
	.logoHeader {
		display: block;
		max-height: 100px;
	}

	.releaseLabel {
		background-color: #ffffff;
		color: #1c1c1c;
		padding: 5px;
		border-radius: 5px;
		font-size: 14pt;
		font-weight: bold;
	}

	h2 {
		font-size: 24pt;
	}
</style>
