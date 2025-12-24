<script lang="ts">
    let { content } = $props();
</script>

<div class="frame-element">
    {@render content?.()}
</div>

<style>
    .frame-element {
        -webkit-box-shadow: 0px 10px 10px 0px rgba(0, 0, 0, 0.2);
        -moz-box-shadow: 0px 10px 10px 0px rgba(0, 0, 0, 0.2);
        box-shadow: 0px 10px 10px 0px rgba(0, 0, 0, 0.2);
        margin-top: 15px;
        margin-left: 15px;
        margin-right: 15px;
        border-radius: 15px;
        overflow-x: auto;
        padding: 10px;
        display: var(--display);
        justify-content: var(--justify-content);
        align-items: var(--align-items);
    }
</style>
