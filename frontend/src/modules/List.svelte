<div class="secretsList">
  {#if $secrets.length}
    {#each $secrets as secret}
      <div
        on:click="{() => selectSecret(secret)}"
        class="secret row a-center j-between {selectedSecret.name === secret.name ? 'active' : ''}"
        title="{secret.name}"
      >
        <h2>{secret.name}</h2>
        <p>{secret.objects}</p>
      </div>
    {/each}
    <div class="footer row a-center j-end">
      <button on:click="{createSecret}">Add</button>
    </div>
  {:else}
    <Loader/>
  {/if}
</div>

<script>
  import Loader from "components/Loader"
  import { namespaces, secrets } from "store.js"
  import { get } from "svelte/store"
  import axios from "axios"

  let selectedSecret = {}

  secrets.selected.subscribe(value => {
    if (value) {
      selectedSecret = value
    }
  })

  function selectSecret (value) {
    secrets.select(value)
  }

  function createSecret () {
    let name = prompt("Provide name for new secret:")
    if (!name.length) {
      return
    }
    const namespace = get(namespaces.selected)

    axios.post(`${process.env.API_HOST}/vault/${namespace}/${name}`, {}).then((res) => {
      if (res.data) {
        secrets.update()
      }
    }).catch((err) => {
      console.log(err)
    })
  }
</script>

<style>
  .secretsList {
    width: 240px;
    height: 100%;
    border-right: solid #d7d7d7 1px;
    overflow-x: hidden;
    position: relative;
  }

  .secretsList > .secret {
    width: 100%;
    height: 50px;
    padding: 0 10px;
    cursor: pointer;
    transition: background-color 100ms ease;
    border-bottom: solid #ededed 1px;
  }
  .secretsList > .secret:last-child {
    border: none;
  }
  .secretsList > .secret.active {
    background: #f1f1f1;
  }
  .secretsList > .secret:hover {
    background: #f1f1f1;
    transition: background-color 200ms ease;
  }
  .secretsList > .secret > h2 {
    font-weight: normal;
    font-size: 14px;
  }
  .secretsList > .secret > p {
    font-weight: lighter;
    font-size: 12px;
  }

  .secretsList > .footer {
    width: 100%;
    height: 50px;
    position: absolute;
    bottom: 0;
    border-top: solid #d7d7d7 1px;
    padding: 8px;
  }
</style>