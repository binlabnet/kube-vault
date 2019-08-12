<main class="row">
  {#if loaded}
    <List></List>
    <Secret></Secret>
  {:else if error}
    <div class="network_error">
      Cannot connect to server ({process.env.API_HOST})
    </div>
  {:else}
    <Loader/>
  {/if}
</main>

<script>
  import List from "modules/List"
  import Secret from "modules/Secret"
  import Loader from "components/Loader"

  import { namespaces } from "store.js"

  let loaded
  let error = false
  namespaces.selected.subscribe(value => {
    if (value === false) {
      error = true
    }
    loaded = value
  })
</script>

<style>
  main {
    width: 100%;
    height: calc(100% - 56px);
    position: relative;
    z-index: 1;
  }

  main > .network_error {
    width: auto;
    height: 20px;
    padding: 0;
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    margin: auto;
    text-align: center;
    font-size: 15px;
  }
</style>