<script lang="ts">
	import { onMount } from 'svelte';
	import FrameElement from '@templates/frameElement.svelte';
	import Sidebar from '@templates/sidebar.svelte';
	import Frame from '@templates/frame.svelte';
	import logo from '@assets/logo_small.svg';
	import { GetImageList } from '@lib/api/api.svelte';
	import type { ImageBaseInfo } from '@lib/api/types.svelte';
	import { UnixTimeFormat } from '@lib/time.svelte';
	import { SizeConvert } from '@lib/size.svelte';

	let imgList: ImageBaseInfo = [];
	let loading = true;
	let selectedIds: string[] = [];
	let buttonIds = new Map<string, boolean>([['remove-button', false]]);
	let activateButtons: boolean = false;

	const updateImageList = () => {
		const getData = async () => {
			try {
				const data = await GetImageList();
				if (data) {
					imgList = data;
				}
			} catch (err) {
				console.log(err);
			} finally {
				loading = false;
				//console.log(loading + "\nimgList:" + JSON.stringify(imgList));
			}
		};
		getData();
	};

	onMount(updateImageList);

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
	<title>Images</title>
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
		<FrameElement content={imageTable} />
	{/if}
{/snippet}

{#snippet loadingData()}
	<p>Loading data...</p>
{/snippet}

{#snippet header()}
	<img class="logoHeader" src={logo} alt="logo" />
	<h2>Image page</h2>
	<div class="releaseLabel">Alpha v0.0.1</div>
{/snippet}

{#snippet imageTable()}
	<div class="button-block">
		<div class="button-block-left">
			<button
				id="refresh-button"
				onclick={(event) => {
					selectAll(false);
					refreshButtonEvent();
					updateImageList();
					event;
				}}>Refresh</button
			>
		</div>
		<div class="button-block-right">
			<button
				class="imagelist-button"
				id="remove-button"
				onclick={() => {
					//RemoveImage(selectedIds);
					selectAll(false);
					setTimeout(() => {
						updateImageList();
					}, 1000);
				}}
				disabled>Remove</button
			>
		</div>
	</div>
	<table class="imagelist-table">
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
				<th>Id</th>
				<th>Tag</th>
				<th>Size</th>
				<th>Created</th>
			</tr>
		</thead>
		<tbody>
			{#each imgList as image}
				<tr>
					<td>
						<input
							type="checkbox"
							name="checkbox-item"
							id={image.id}
							onchange={(event) => {
								const target = event.target as HTMLInputElement;
								updateSelected(image.id, target.checked);
							}}
						/>
					</td>
					<td>{image.id}</td>
					<td>{image.tags}</td>
					<td>{SizeConvert(image.size)}</td>
					<td>{UnixTimeFormat(image.created)}</td>
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

	.imagelist-table {
		width: 100%;
		text-align: center;
	}

	.imagelist-button {
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
