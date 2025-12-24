<script lang="ts">
    import logo from "@assets/logo.svg";
</script>

<div class="header">
    <img class="logoHeader" src={logo} alt="logo" />
</div>

<style>
    .logoHeader {
        display: block;
        max-height: 100%;
    }

    .header {
        height: 100px;
        -webkit-box-shadow: 0px 10px 10px 0px rgba(0, 0, 0, 0.2);
        -moz-box-shadow: 0px 10px 10px 0px rgba(0, 0, 0, 0.2);
        box-shadow: 0px 10px 10px 0px rgba(0, 0, 0, 0.2);
        margin-top: 15px;
        margin-left: 15px;
        margin-right: 15px;
        border-radius: 15px;
    }
</style>
