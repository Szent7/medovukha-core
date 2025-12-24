<script lang="ts" context="module">
    export function UnixTimeFormat(unixTimestamp: number): string {
        const date = new Date(unixTimestamp * 1000);
        return date.toLocaleString();
    }
</script>
