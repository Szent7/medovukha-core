<script lang="ts">
	import { onMount } from 'svelte';
	import FrameElement from '@templates/frameElement.svelte';
	import Sidebar from '@templates/sidebar.svelte';
	import Frame from '@templates/frame.svelte';
	import logo from '@assets/logo_small.svg';
	import { GetNetworkList } from '@lib/api/api.svelte';
	import type { NetworkBaseInfo } from '@lib/api/types.svelte';

	let netList: NetworkBaseInfo = [];
	let loading = true;
	let selectedIds: string[] = [];
	let buttonIds = new Map<string, boolean>([['remove-button', false]]);
	let activateButtons: boolean = false;

	const updateNetworkList = () => {
		const getData = async () => {
			try {
				const data = await GetNetworkList();
				if (data) {
					netList = data;
				}
			} catch (err) {
				console.log(err);
			} finally {
				loading = false;
				//console.log(loading + "\nnetList:" + JSON.stringify(netList));
			}
		};
		getData();
	};

	onMount(updateNetworkList);

	function updateSelected(id: string, checked: boolean) {
		if (checked) {
			selectedIds = [...selectedIds, id];
			const elements = document.querySelectorAll('input[name=checkbox-item]');
			if (elements.length == selectedIds.length) {
				const checkboxAll = document.querySelector(
					'input[name="checkbox-item-all"]'
				) as HTMLInputElement;
				if (checkboxAll !== null) {
					checkboxAll.checked = checked;
				}
			}
			if (!activateButtons) {
				buttonIds.forEach((_: boolean, key: string) => {
					buttonIds.set(key, false);
				});
				activateButtons = true;
			}
		} else {
			selectedIds = selectedIds.filter((item) => item !== id);
			if (selectedIds.length == 0) {
				buttonIds.forEach((_: boolean, key: string) => {
					buttonIds.set(key, true);
				});
				activateButtons = false;
			}
			const checkboxAll = document.querySelector(
				'input[name="checkbox-item-all"]'
			) as HTMLInputElement;
			if (checkboxAll !== null) {
				checkboxAll.checked = checked;
			}
		}
		//console.log("SelectedIds:" + selectedIds);
	}

	function selectAll(checked: boolean) {
		const elements = document.querySelectorAll('input[name=checkbox-item]');
		if (elements !== null) {
			Array.prototype.forEach.call(elements, function (item) {
				updateSelected(item.id, checked);
				item.checked = checked;
			});
		}
	}

	function refreshButtonEvent() {
		const refreshButton = document.getElementById('refresh-button') as HTMLButtonElement | null;
		if (refreshButton !== null) {
			const originalText = refreshButton.innerHTML;
			refreshButton.disabled = true;
			refreshButton.innerHTML = 'Refreshing';

			setTimeout(() => {
				refreshButton.disabled = false;
				refreshButton.innerHTML = originalText;
			}, 1000);
		}
	}
</script>

<svelte:head>
	<title>Networks</title>
	<link rel="icon" type="image/svg+xml" href={logo} />
</svelte:head>

<Sidebar />
<Frame {content} />

{#snippet content()}
	<FrameElement
		--display="flex"
		--justify-content="space-between"
		--align-items="center"
		content={header}
	/>
	{#if loading}
		<FrameElement content={loadingData} />
	{:else if !loading}
		<FrameElement content={networkTable} />
	{/if}
{/snippet}

{#snippet loadingData()}
	<p>Loading data...</p>
{/snippet}

{#snippet header()}
	<img class="logoHeader" src={logo} alt="logo" />
	<h2>Network page</h2>
	<div class="releaseLabel">Alpha v0.0.1</div>
{/snippet}

{#snippet networkTable()}
	<div class="button-block">
		<div class="button-block-left">
			<button
				id="refresh-button"
				onclick={(event) => {
					selectAll(false);
					refreshButtonEvent();
					updateNetworkList();
					event;
				}}>Refresh</button
			>
		</div>
		<div class="button-block-right">
			<button
				class="networklist-button"
				id="remove-button"
				onclick={() => {
					//RemoveImage(selectedIds);
					selectAll(false);
					setTimeout(() => {
						updateNetworkList();
					}, 1000);
				}}
				disabled>Remove</button
			>
		</div>
	</div>
	<table class="networklist-table">
		<thead>
			<tr>
				<th>
					<input
						type="checkbox"
						name="checkbox-item-all"
						onchange={(event) => {
							const target = event.target as HTMLInputElement;
							selectAll(target.checked);
						}}
					/>
				</th>
				<th>Name</th>
				<th>Id</th>
				<th>Driver</th>
				<th>EnableIPv6</th>
				<th>IPAMDriver</th>
				<th>Subnet</th>
				<th>Gateway</th>
				<th>Attachable</th>
			</tr>
		</thead>
		<tbody>
			{#each netList as network}
				<tr>
					<td>
						{#if network.dockerNetwork}
							<input type="checkbox" name="checkbox-item" id={network.id} disabled />
						{:else}
							<input
								type="checkbox"
								name="checkbox-item"
								id={network.id}
								onchange={(event) => {
									const target = event.target as HTMLInputElement;
									updateSelected(network.id, target.checked);
								}}
							/>
						{/if}
					</td>
					<td>{network.name}</td>
					<td>{network.id}</td>
					<td>{network.driver}</td>
					<td>{network.enableIPv6}</td>
					<td>{network.ipamDriver}</td>
					<td>
						{#if network.subnet.length == 0}
							-
						{:else}
							{#each network.subnet as sn}
								{#if sn !== undefined}
									<p>{sn}</p>
								{:else}
									-
								{/if}
							{/each}
						{/if}
					</td>
					<td>
						{#if network.gateway.length == 0}
							-
						{:else}
							{#each network.gateway as gw}
								{#if gw !== undefined}
									<p>{gw}</p>
								{:else}
									-
								{/if}
							{/each}
						{/if}
					</td>
					<td>{network.attachable}</td>
				</tr>
			{/each}
		</tbody>
	</table>
{/snippet}

<style>
	h2 {
		font-size: 24pt;
	}

	.releaseLabel {
		background-color: #ffffff;
		color: #1c1c1c;
		padding: 5px;
		border-radius: 5px;
		font-size: 14pt;
		font-weight: bold;
	}

	.logoHeader {
		display: block;
		max-height: 100px;
	}

	.networklist-table {
		width: 100%;
		text-align: center;
	}

	.networklist-button {
		border-radius: 10px;
	}

	.button-block {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.button-block-right {
		text-align: right;
	}
</style>
