<script lang="ts">
  import axios from "axios";
  export let choose: number = 0;
  const change = () => {
    axios
      .get("http://" + location.hostname + ":10015/ping")
      .then(function (response) {
        console.log(response.data);
      });
    console.log("pressed");
    if (choose == 0) {
      choose = 1;
    } else {
      choose = 0;
    }
  };
</script>

<button on:click={change}>
  choose is {choose}
</button>
