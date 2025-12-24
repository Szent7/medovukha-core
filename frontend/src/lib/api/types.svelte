<script lang="ts" context="module">
	import { z } from 'zod';

	export const ContainerIdMessageScheme = z.object({
		Id: z.string()
	});
	export type ContainerIdMessage = z.infer<typeof ContainerIdMessageScheme>;

	export const ContainerBaseInfoScheme = z.array(
		z.object({
			id: z.string(),
			names: z.array(z.string()),
			image: z.string(),
			ports: z
				.array(
					z
						.object({
							ip: z.string(),
							privatePort: z.number(),
							publicPort: z.number(),
							type: z.string()
						})
						.optional()
				)
				.nullable(),
			created: z.number(),
			state: z.string(),
			isMedovukha: z.boolean()
		})
	);
	export type ContainerBaseInfo = z.infer<typeof ContainerBaseInfoScheme>;

	export const ImageBaseInfoScheme = z.array(
		z.object({
			id: z.string(),
			tags: z.array(z.string()),
			size: z.number(),
			created: z.number()
		})
	);
	export type ImageBaseInfo = z.infer<typeof ImageBaseInfoScheme>;

	export const NetworkBaseInfoScheme = z.array(
		z.object({
			name: z.string(),
			id: z.string(),
			driver: z.string(),
			enableIPv6: z.boolean(),
			ipamDriver: z.string(),
			subnet: z.array(z.string()),
			gateway: z.array(z.string()),
			attachable: z.boolean(),
			dockerNetwork: z.boolean()
		})
	);
	export type NetworkBaseInfo = z.infer<typeof NetworkBaseInfoScheme>;

	export const VolumeBaseInfoScheme = z.array(
		z.object({
			name: z.string(),
			driver: z.string(),
			mountpoint: z.string(),
			created: z.string()
		})
	);
	export type VolumeBaseInfo = z.infer<typeof VolumeBaseInfoScheme>;
</script>
