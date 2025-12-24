<script lang="ts" context="module">
	import axios from 'axios';
	import {
		ContainerBaseInfoScheme,
		ImageBaseInfoScheme,
		NetworkBaseInfoScheme,
		VolumeBaseInfoScheme,
		type ContainerBaseInfo,
		type ContainerIdMessage,
		type ImageBaseInfo,
		type NetworkBaseInfo,
		type VolumeBaseInfo
	} from './types.svelte';

	export async function PingBack() {
		axios.get('http://localhost:10015/ping').then(function (response) {
			console.log(response.data);
		});
	}
	//Containers
	export async function GetContainerList() {
		try {
			const response = await axios.get('http://localhost:10015/rest/v1/getContainerList');
			const containers: ContainerBaseInfo = ContainerBaseInfoScheme.parse(response.data);
			//console.log("responseData:" + containers);
			return containers;
		} catch (error) {
			console.error('Error GET ContainerList:', error);
		}
	}

	export async function PauseContainerById(id: string) {
		const container: ContainerIdMessage = {
			Id: id
		};
		try {
			const response = await axios.post(
				'http://localhost:10015/rest/v1/pauseContainerById',
				container
			);
			/*const containers: ContainerBaseInfo = ContainerBaseInfoScheme.parse(
                response.data,
            );
            console.log("responseData:" + containers);
            return containers;*/
		} catch (error) {
			console.error('Error POST PauseContainerById:', error);
		}
	}

	export async function UnpauseContainerById(id: string) {
		const container: ContainerIdMessage = {
			Id: id
		};
		try {
			const response = await axios.post(
				'http://localhost:10015/rest/v1/unpauseContainerById',
				container
			);
			/*const containers: ContainerBaseInfo = ContainerBaseInfoScheme.parse(
                response.data,
            );
            console.log("responseData:" + containers);
            return containers;*/
		} catch (error) {
			console.error('Error POST UnpauseContainerById:', error);
		}
	}

	export async function KillContainerById(id: string) {
		const container: ContainerIdMessage = {
			Id: id
		};
		try {
			const response = await axios.post(
				'http://localhost:10015/rest/v1/killContainerById',
				container
			);
			/*const containers: ContainerBaseInfo = ContainerBaseInfoScheme.parse(
                response.data,
            );
            console.log("responseData:" + containers);
            return containers;*/
		} catch (error) {
			console.error('Error POST KillContainerById:', error);
		}
	}

	export async function StartContainerById(id: string) {
		const container: ContainerIdMessage = {
			Id: id
		};
		try {
			const response = await axios.post(
				'http://localhost:10015/rest/v1/startContainerById',
				container
			);
			/*const containers: ContainerBaseInfo = ContainerBaseInfoScheme.parse(
                response.data,
            );
            console.log("responseData:" + containers);
            return containers;*/
		} catch (error) {
			console.error('Error POST StartContainerById:', error);
		}
	}

	export async function StopContainerById(id: string) {
		const container: ContainerIdMessage = {
			Id: id
		};
		try {
			const response = await axios.post(
				'http://localhost:10015/rest/v1/stopContainerById',
				container
			);
			/*const containers: ContainerBaseInfo = ContainerBaseInfoScheme.parse(
                response.data,
            );
            console.log("responseData:" + containers);
            return containers;*/
		} catch (error) {
			console.error('Error POST StopContainerById:', error);
		}
	}

	export async function RestartContainerById(id: string) {
		const container: ContainerIdMessage = {
			Id: id
		};
		try {
			const response = await axios.post(
				'http://localhost:10015/rest/v1/restartContainerById',
				container
			);
			/*const containers: ContainerBaseInfo = ContainerBaseInfoScheme.parse(
                response.data,
            );
            console.log("responseData:" + containers);
            return containers;*/
		} catch (error) {
			console.error('Error POST RestartContainerById:', error);
		}
	}

	export async function RemoveContainerById(id: string) {
		const container: ContainerIdMessage = {
			Id: id
		};
		try {
			const response = await axios.post(
				'http://localhost:10015/rest/v1/removeContainerById',
				container
			);
			/*const containers: ContainerBaseInfo = ContainerBaseInfoScheme.parse(
                response.data,
            );
            console.log("responseData:" + containers);
            return containers;*/
		} catch (error) {
			console.error('Error POST RemoveContainerById:', error);
		}
	}

	//Images
	export async function GetImageList() {
		try {
			const response = await axios.get('http://localhost:10015/rest/v1/getImageList');
			const images: ImageBaseInfo = ImageBaseInfoScheme.parse(response.data);
			//console.log("responseData:" + containers);
			return images;
		} catch (error) {
			console.error('Error GET ImageList:', error);
		}
	}

	//Networks
	export async function GetNetworkList() {
		try {
			const response = await axios.get('http://localhost:10015/rest/v1/getNetworkList');
			const networks: NetworkBaseInfo = NetworkBaseInfoScheme.parse(response.data);
			//console.log("responseData:" + containers);
			return networks;
		} catch (error) {
			console.error('Error GET NetworkList:', error);
		}
	}

	//Volumes
	export async function GetVolumeList() {
		try {
			const response = await axios.get('http://localhost:10015/rest/v1/getVolumeList');
			const volumes: VolumeBaseInfo = VolumeBaseInfoScheme.parse(response.data);
			//console.log("responseData:" + containers);
			return volumes;
		} catch (error) {
			console.error('Error GET VolumeList:', error);
		}
	}
</script>
