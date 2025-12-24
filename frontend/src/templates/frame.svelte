<script lang="ts">
    let { content } = $props();
</script>

<div class="container">
    {@render content?.()}
</div>

<style>
    .container {
        width: auto;
        height: 100%;
        background-color: #1c1c1c;
        display: flex;
        flex-direction: column;
        overflow-y: auto;
    }
</style>
