<main class="row">
  {#if loaded}
    <List></List>
    <Secret></Secret>
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
  namespaces.selected.subscribe(value => {
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
</style>