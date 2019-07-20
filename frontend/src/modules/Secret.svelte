<div class="secretView column">
  <div class="settings">
    {$secret}
  </div>

  <div class="secrets">
    {#if items.length}
      {#each items as item}
        <div class="secret row">
          <h3>{item.key}</h3>
          <div>{item.value}</div>
        </div>
      {/each}
    {:else}
      <p class="row a-center j-center">No items</p>
    {/if}
  </div>
</div>


<script>
  import axios from "axios"
  import { namespace, secret } from "store.js"

  let secretItems = {}
  $: items = Object.keys(secretItems).map(k => {
    return {
      key: k,
      value: secretItems[k],
    }
  })

  const unsubscribe = secret.subscribe(value => {
    if (value) {
      axios.get(`${process.env.API_HOST}/vault/${$namespace}/${value}`).then((res) => {
        if (res.data) {
          if (res.data.data) {
            secretItems = res.data.data
          } else {
            secretItems = {}
          }
        }
      }).catch((err) => {
        console.log(err)
      })
    }
  })
</script>


<style>
  .secretView {
    width: calc(100% - 240px);
    height: 100%;
    background: #f6f6f6;
    padding: 14px;
  }

  .secretView > .settings {
    width: 100%;
    height: auto;
    background: #ffffff;
    margin-bottom: 15px;
    padding: 8px;
  }

  .secretView > .secrets {
    width: 100%;
    height: auto;
    background: #ffffff;
    padding: 8px;
    overflow-x: hidden;
  }
  .secretView > .secrets > p {
    width: 100%;
    height: 120px;
  }
  .secretView > .secrets > .secret {
    width: 100%;
    height: auto;
    margin-top: 6px;
    padding: 8px;
    border-bottom: solid #e7e7e7 1px;
    cursor: pointer;
  }
  .secretView > .secrets > .secret:first-child {
    margin: 0;
  }
  .secretView > .secrets > .secret:last-child {
    border-bottom: none;
  }
  .secretView > .secrets > .secret > h3 {
    width: 120px;
    font-weight: normal;
    word-wrap: break-word;
    margin: 0;
    padding: 0;
  }
  .secretView > .secrets > .secret > div {
    width: calc(100% - 120px);
    height: 100%;
    outline: none;
    resize: none;
    word-wrap: break-word;
  }
</style>