<script lang="ts">
	import { onMount } from 'svelte';
	import FrameElement from '@templates/frameElement.svelte';
	import Sidebar from '@templates/sidebar.svelte';
	import Frame from '@templates/frame.svelte';
	import logo from '@assets/logo_small.svg';
	import { GetContainerList } from '@lib/api/api.svelte';
	import type { ContainerBaseInfo } from '@lib/api/types.svelte';
	import { UnixTimeFormat } from '@lib/time.svelte';
	import {
		KillContainer,
		PauseContainer,
		RemoveContainer,
		StartContainer,
		UnpauseContainer,
		StopContainer,
		RestartContainer
	} from '@lib/containerActions.svelte';

	let conList: ContainerBaseInfo = [];
	let loading = true;
	let selectedIds: string[] = [];
	let buttonIds = new Map<string, boolean>([
		['start-button', false],
		['stop-button', false],
		['kill-button', false],
		['restart-button', false],
		['pause-button', false],
		['resume-button', false],
		['remove-button', false]
	]);
	let stateList = new Map<string, boolean>([
		['running', false],
		['paused', false],
		['exited', false]
	]);
	let uniqueStates: string[] = [];
	const actionsBlacklist = new Map<string, string[]>([
		['running', ['start-button', 'resume-button']],
		['paused', ['start-button', 'stop-button', 'pause-button']],
		['exited', ['stop-button', 'kill-button', 'pause-button', 'resume-button']]
	]);
	let activateButtons: boolean = false;

	const updateContainerList = () => {
		const getData = async () => {
			try {
				const data = await GetContainerList();
				if (data) {
					conList = data;
				}
			} catch (err) {
				console.log(err);
			} finally {
				loading = false;
			}
		};
		getData();
	};

	onMount(updateContainerList);

	function CheckIsMedovukha(id: string): boolean {
		let isMedovukha: boolean = false;
		conList.forEach((container) => {
			if (container.id == id) {
				isMedovukha = container.isMedovukha;
				return;
			}
		});
		return isMedovukha;
	}

	function updateButtons() {
		buttonIds.forEach((value: boolean, key: string) => {
			const button = document.getElementById(key) as HTMLButtonElement | null;
			if (button !== null) {
				button.disabled = value;
			}
		});
	}

	function updateStateList() {
		if (selectedIds.length == 0) {
			stateList.forEach((_: boolean, key: string) => {
				stateList.set(key, false);
				//console.log("stateList: ", key, ":", false);
			});
			buttonIds.forEach((_: boolean, key: string) => {
				buttonIds.set(key, true);
			});
			updateButtons();
			uniqueStates = [];
			return;
		}
		let states: string[] = new Array<string>(selectedIds.length);
		for (let i = 0; i < selectedIds.length; i++) {
			for (let j = 0; j < conList.length; j++) {
				if (selectedIds[i] == conList[j].id) {
					states[i] = conList[j].state;
				}
			}
		}
		const newUniqueStates = Array.from(new Set(states));
		if (newUniqueStates.length != uniqueStates.length) {
			uniqueStates = newUniqueStates;
			buttonIds.forEach((_: boolean, key: string) => {
				buttonIds.set(key, false);
			});
			uniqueStates.forEach((item) => {
				const blackList = actionsBlacklist.get(item);
				if (blackList !== undefined) {
					blackList.forEach((item) => {
						buttonIds.set(item, true);
					});
				}
			});
			updateButtons();
		}
	}

	function updateSelected(id: string, checked: boolean) {
		if (checked) {
			selectedIds = [...selectedIds, id];
			updateStateList();
			const elements = document.querySelectorAll('input[name=checkbox-item]');
			if (elements.length - 1 == selectedIds.length) {
				const checkboxAll = document.querySelector(
					'input[name="checkbox-item-all"]'
				) as HTMLInputElement;
				if (checkboxAll !== null) {
					checkboxAll.checked = checked;
				}
			}
			if (!activateButtons) {
				/*buttonIds.forEach((_: boolean, key: string) => {
                    buttonIds.set(key, false);
                });
                updateButtons();*/
				updateStateList();
				activateButtons = true;
			}
		} else {
			selectedIds = selectedIds.filter((item) => item !== id);
			updateStateList();
			if (selectedIds.length == 0) {
				/*buttonIds.forEach((_: boolean, key: string) => {
                    buttonIds.set(key, true);
                });
                updateButtons();*/
				updateStateList();
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
				if (!CheckIsMedovukha(item.id) && item.checked != checked) {
					updateSelected(item.id, checked);
					item.checked = checked;
				}
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
	<title>Containers</title>
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
		<FrameElement content={containerTable} />
	{/if}
{/snippet}

{#snippet loadingData()}
	<p>Loading data...</p>
{/snippet}

{#snippet header()}
	<img class="logoHeader" src={logo} alt="logo" />
	<h2>Container page</h2>
	<div class="releaseLabel">Alpha v0.0.1</div>
{/snippet}

{#snippet containerTable()}
	<div class="button-block">
		<div class="button-block-left">
			<button
				id="refresh-button"
				onclick={(event) => {
					selectAll(false);
					refreshButtonEvent();
					updateContainerList();
					event;
				}}>Refresh</button
			>
			<a href="/create-container-git" id="create-container-git-button" class="btn">Create from git</a>
		</div>
		<div class="button-block-right">
			<button
				class="containerlist-button"
				id="start-button"
				onclick={() => {
					StartContainer(selectedIds);
					selectAll(false);
					setTimeout(() => {
						updateContainerList();
					}, 1000);
				}}
				disabled>Start</button
			>
			<button
				class="containerlist-button"
				id="stop-button"
				onclick={() => {
					StopContainer(selectedIds);
					selectAll(false);
					setTimeout(() => {
						updateContainerList();
					}, 1000);
				}}
				disabled>Stop</button
			>
			<button
				class="containerlist-button"
				id="kill-button"
				onclick={() => {
					KillContainer(selectedIds);
					selectAll(false);
					setTimeout(() => {
						updateContainerList();
					}, 1000);
				}}
				disabled>Kill</button
			>
			<button
				class="containerlist-button"
				id="restart-button"
				onclick={() => {
					RestartContainer(selectedIds);
					selectAll(false);
					setTimeout(() => {
						updateContainerList();
					}, 1000);
				}}
				disabled>Restart</button
			>
			<button
				class="containerlist-button"
				id="pause-button"
				onclick={() => {
					PauseContainer(selectedIds);
					selectAll(false);
					setTimeout(() => {
						updateContainerList();
					}, 1000);
				}}
				disabled>Pause</button
			>
			<button
				class="containerlist-button"
				id="resume-button"
				onclick={() => {
					UnpauseContainer(selectedIds);
					selectAll(false);
					setTimeout(() => {
						updateContainerList();
					}, 1000);
				}}
				disabled>Resume</button
			>
			<button
				class="containerlist-button"
				id="remove-button"
				onclick={() => {
					RemoveContainer(selectedIds);
					selectAll(false);
					setTimeout(() => {
						updateContainerList();
					}, 1000);
				}}
				disabled>Remove</button
			>
		</div>
	</div>
	<table class="containerlist-table">
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
				<th>State</th>
				<th>Image</th>
				<th>Port bindings</th>
				<th>Created</th>
			</tr>
		</thead>
		<tbody>
			{#each conList as container}
				<tr>
					<td>
						{#if container.isMedovukha}
							<input type="checkbox" name="checkbox-item" id={container.id} disabled />
						{:else}
							<input
								type="checkbox"
								name="checkbox-item"
								id={container.id}
								onchange={(event) => {
									const target = event.target as HTMLInputElement;
									updateSelected(container.id, target.checked);
								}}
							/>
						{/if}
					</td>
					<td>{container.names[0]}</td>
					<td class="state-{container.state}">{container.state}</td>
					<td>{container.image}</td>
					<td>
						{#if container.ports == null || container.ports.length == 0}
							-
						{:else}
							{#each container.ports as port}
								{#if port?.publicPort !== undefined && port?.privatePort !== undefined}
									<p>{port.publicPort}:{port.privatePort}</p>
								{:else}
									-
								{/if}
							{/each}
						{/if}
					</td>
					<td>{UnixTimeFormat(container.created)}</td>
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

	.containerlist-table {
		width: 100%;
		text-align: center;
	}

	.containerlist-button {
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

	.state-running {
		color: aquamarine;
	}

	.state-exited {
		color: red;
	}

	.state-paused {
		color: orange;
	}

	#create-container-git-button {
		color: orangered;
	}
</style>
