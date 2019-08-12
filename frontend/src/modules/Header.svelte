<header class="row a-center j-between">
  <div class="name row a-base">
    <h1 title="kube-vault">kube-vault</h1>
    <a href="https://github.com/exelban/kube-vault" target="_blank" title="Project page">{version}</a>
  </div>

  {#if $namespaces && $namespaces.length}
    <div class="namespaces">
      <select title="Namespace" bind:value={selectedNamespace}>
        {#each $namespaces as n}
          <option value="{n}">{n}</option>
        {/each}
      </select>
    </div>
  {/if}
</header>

<script>
  import { namespaces } from "store.js"

  const version = `v${process.env.VERSION}`
  let selectedNamespace

  namespaces.selected.subscribe(value => {
    selectedNamespace = value
  })

  $: {
    if (selectedNamespace) {
      namespaces.select(selectedNamespace)
    }
  }
</script>

<style>
  header {
    width: 100%;
    height: 56px;
    background: #ffffff;
    box-shadow: rgba(0, 0, 0, 0.2) 0 2px 4px -1px, rgba(0, 0, 0, 0.137255) 0 4px 5px 0, rgba(0, 0, 0, 0.117647) 0 1px 10px 0;
    padding: 0 20px;
    position: relative;
    z-index: 2;
  }

  header > .name > h1 {
    font-weight: lighter;
    padding: 0;
    margin: 0;
    cursor: default;
  }
  header > .name > a {
    font-weight: lighter;
    font-size: 10px;
    margin-left: 6px;
    cursor: pointer;
  }

  header > .namespaces > select {
    outline: none;
  }
</style>