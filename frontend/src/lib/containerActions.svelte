<script lang="ts" context="module">
	import {
		KillContainerById,
		PauseContainerById,
		StartContainerById,
		UnpauseContainerById,
		RemoveContainerById,
		StopContainerById,
		RestartContainerById
	} from './api/api.svelte';

	export function StartContainer(ids: string[]) {
		ids.forEach((id: string): void => {
			StartContainerById(id);
		});
	}

	export function PauseContainer(ids: string[]) {
		ids.forEach((id: string): void => {
			PauseContainerById(id);
		});
	}

	export function UnpauseContainer(ids: string[]) {
		ids.forEach((id: string): void => {
			UnpauseContainerById(id);
		});
	}

	export function KillContainer(ids: string[]) {
		ids.forEach((id: string): void => {
			KillContainerById(id);
		});
	}

	export function StopContainer(ids: string[]) {
		ids.forEach((id: string): void => {
			StopContainerById(id);
		});
	}

	export function RestartContainer(ids: string[]) {
		ids.forEach((id: string): void => {
			RestartContainerById(id);
		});
	}

	export function RemoveContainer(ids: string[]) {
		ids.forEach((id: string): void => {
			RemoveContainerById(id);
		});
	}
</script>
