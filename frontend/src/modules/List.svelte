<div class="secretsList column">
  {#if $secrets.length}
    <div class="container">
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
    </div>
    <div class="footer row a-center j-between">
      <p>{$secrets.length} secrets</p>
      <button on:click="{createSecret}">New secret</button>
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
    let name = prompt("Provide a name for the new secret:")
    if (!name || !name.length) {
      return
    }
    const namespace = get(namespaces.selected)

    axios.post(`${process.env.API_HOST}/vault/${namespace}/${name}`, {}).then((res) => {
      if (res.data) {
        secrets.update()
      }
    }).catch((err) => {
      if (err.response && err.response.data) {
        alert(err.response.data)
      }
    })
  }
</script>

<style>
  .secretsList {
    width: 240px;
    height: 100%;
    border-right: solid #d7d7d7 1px;
    position: relative;
  }

  .secretsList > .container {
    width: 100%;
    height: calc(100% - 50px);
    overflow-x: hidden;
  }
  .secretsList > .container > .secret {
    width: 100%;
    height: 50px;
    padding: 0 10px;
    cursor: pointer;
    transition: background-color 100ms ease;
    border-bottom: solid #ededed 1px;
  }
  .secretsList > .container > .secret:last-child {
    border: none;
  }
  .secretsList > .container > .secret.active {
    background: #f1f1f1;
  }
  .secretsList > .container > .secret:hover {
    background: #f1f1f1;
    transition: background-color 200ms ease;
  }
  .secretsList > .container > .secret > h2 {
    font-weight: normal;
    font-size: 14px;
  }
  .secretsList > .container > .secret > p {
    font-weight: lighter;
    font-size: 12px;
  }

  .secretsList > .footer {
    width: 100%;
    height: 50px;
    border-top: solid #d7d7d7 1px;
    padding: 8px;
    background: #ffffff;
  }
</style>