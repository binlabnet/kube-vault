<div class="secretsList">
  {#each secretsList as secret}
    <div
      on:click="{() => selectSecret(secret.name)}"
      class="secret row a-center j-between {selectedSecret === secret.name ? 'active' : ''}"
      title="{secret.name}"
    >
      <h2>{secret.name}</h2>
      <p>{secret.objects}</p>
    </div>
  {/each}
</div>


<script>
  import axios from "axios"
  import { namespace, secret } from "store.js"

  let loading = false
  let secretsList = []
  let selectedSecret

  const unsubscribe = namespace.subscribe(value => {
    if (value) {
      axios.get(`${process.env.API_HOST}/vault/${value}`).then((res) => {
        if (res.data) {
          secretsList = res.data
          selectedSecret = secretsList[0].name
          secret.set(selectedSecret)
        }
      }).catch((err) => {
        console.log(err)
      })
    }
  })

  function selectSecret (name) {
    selectedSecret = name
    secret.set(selectedSecret)
  }
</script>


<style>
  .secretsList {
    width: 240px;
    height: 100%;
    border-right: solid #d7d7d7 1px;
    overflow-x: hidden;
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
</style>