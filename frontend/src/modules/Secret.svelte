<div class="secretView column">
  {#if secret}
    <div class="settings row a-center j-between">
      <h4>{secret.name}</h4>
      <div class="row">
        {#if !newItemMode}
          <button on:click="{deleteSecret}">Delete secret</button>
          <button on:click="{toggleNewItemMode}">Add item</button>
        {/if}
      </div>
    </div>

    {#if newItemMode}
      <div class="newSecret">
        <input type="text" bind:value={newItemKey} placeholder="Key">
        <input type="text" bind:value={newItemValue} placeholder="Value">

        <div class="row a-center j-end">
          <button on:click="{toggleNewItemMode}">Cancel</button>
          <button on:click="{createItem}">Save</button>
        </div>
      </div>
    {/if}

    <div class="secrets">
      {#if !items}
        <Loader/>
      {:else if items.length}
        {#each items as item}
          <div class="secret row">
            <h3>{item.key}</h3>
            {#if editItemMode && editItemKey === item.key}
              <input type="text" bind:value={item.value} placeholder="Put new value here"/>
            {:else}
              <div>{item.value}</div>
            {/if}
            <div class="menu column">
              {#if editItemMode && editItemKey === item.key}
                <button on:click="{() => updateItem(item)}">Save</button>
              {:else}
                <button on:click="{() => toggleItemEditMode(item)}">Edit</button>
              {/if}
              <button on:click="{() => deleteItem(item)}">Remove</button>
            </div>
          </div>
        {/each}
      {:else}
        <p class="row a-center j-center">No items</p>
      {/if}
    </div>
  {:else}
    <Loader/>
  {/if}
</div>

<script>
  import axios from "axios"
  import Loader from "components/Loader"
  import { secrets } from "store.js"

  let secret
  let items

  let newItemMode = false
  let newItemKey = ""
  let newItemValue = ""

  let editItemMode = false
  let editItemKey = ""

  const unsubscribe = secrets.selected.subscribe(value => {
    if (value) {
      secret = value

      axios.get(`${process.env.API_HOST}/vault/${value.namespace}/${value.name}`).then((res) => {
        if (res.data) {
          let data = []
          if (res.data.data) {
            data = res.data.data
          }

          items = Object.keys(data).map(k => {
            return {
              key: k,
              value: data[k],
            }
          })
        }
      }).catch((err) => {
        console.log(err)
      })
    }
  })

  function toggleNewItemMode () {
    newItemMode = !newItemMode
    newItemKey = ""
    newItemValue = ""
  }
  function toggleItemEditMode (value) {
    editItemMode = !editItemMode
    if (value) {
      editItemKey = value.key
      return
    }
    editItemKey = ""
  }

  function createItem () {
    if (!newItemKey.length || !newItemValue.length) {
      return
    }

    const itemsList = items
    itemsList.push({
      key: newItemKey,
      value: newItemValue,
    })

    updateSecret(itemsList)
    toggleNewItemMode()
  }
  function deleteItem (item) {
    if (confirm(`Are you sure you want to delete item ${item.key}?`)) {
      const itemsList = items.filter(i => i !== item)
      updateSecret(itemsList)
    }
  }
  function updateItem (item) {
    const itemsList = items
    itemsList[item.key] = item.value

    updateSecret(itemsList)
    toggleItemEditMode()
  }

  function updateSecret (itemsList) {
    let data = itemsList.reduce((result, i) => {
      const key = Object.keys(i)[0]
      const value = Object.keys(i)[1]
      result[i[key]] = i[value]
      return result
    }, {})

    axios.post(`${process.env.API_HOST}/vault/${secret.namespace}/${secret.name}`, data).then((res) => {
      if (res.data) {
        let data = []
        if (res.data.data) {
          data = res.data.data
        }

        items = Object.keys(data).map(k => {
          return {
            key: k,
            value: data[k],
          }
        })
      }
    }).catch((err) => {
      console.log(err)
    })
  }
  function deleteSecret () {
    if (confirm(`Are you sure you want to delete secret ${secret.name}?`)) {
      axios.delete(`${process.env.API_HOST}/vault/${secret.namespace}/${secret.name}`).then(() => {
        secrets.update()
      }).catch((err) => {
        console.log(err)
      })
    }
  }
</script>

<style>
  .secretView {
    width: calc(100% - 240px);
    height: 100%;
    background: #f6f6f6;
    padding: 14px;
    position: relative;
  }

  .secretView > .settings {
    width: 100%;
    height: auto;
    min-height: 50px;
    background: #ffffff;
    padding: 8px;
  }
  .secretView > .settings > h4 {
    margin: 0;
    font-weight: bold;
  }

  .secretView > .newSecret {
    width: 100%;
    height: auto;
    min-height: 50px;
    background: #ffffff;
    padding: 8px;
    margin-bottom: 15px;
  }
  .secretView > .newSecret > input {
    width: 100%;
    height: 30px;
    margin-top: 4px;
    outline: none;
  }
  .secretView > .newSecret > div {
    padding: 8px;
  }

  .secretView > .secrets {
    width: 100%;
    height: auto;
    min-height: 100px;
    background: #ffffff;
    padding: 8px;
    overflow-x: hidden;
    position: relative;
    margin-top: 15px;
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
    cursor: default;
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
    width: calc(100% - 210px);
    height: 100%;
    outline: none;
    resize: none;
    word-wrap: break-word;
  }
  .secretView > .secrets > .secret > input {
    width: calc(100% - 210px);
    height: 30px;
    outline: none;
    resize: none;
    word-wrap: break-word;
  }
  .secretView > .secrets > .secret > .menu {
    width: 80px;
    height: auto;
    margin: 0 5px;
  }
  .secretView > .secrets > .secret > .menu > button {
    cursor: pointer;
  }
</style>