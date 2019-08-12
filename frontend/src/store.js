import { writable, get } from "svelte/store"
import axios from "axios"


export const namespaces = (() => {
  const { subscribe, set } = writable(undefined)

  const selected = (() => {
    const { subscribe, set, update } = writable(undefined)
    return {
      subscribe,
      set,
      update,
    }
  })()

  function fetchNamespaces () {
    axios.get(`${process.env.API_HOST}/namespaces`).then((res) => {
      if (res.data) {
        set(res.data)
        const defaultNamespace = res.data.findIndex(n => n === process.env.NAMESPACE)
        selected.set(res.data[defaultNamespace])
      }
    }).catch(() => {
      selected.set(false)
    })
  }
  fetchNamespaces()

  return {
    selected,
    subscribe,
    get,
    select: (value) => {
      selected.set(value)
    },
    update: () => fetchNamespaces,
  }
})()

export const secrets = (() => {
  const { subscribe, set } = writable([])

  const selected = (() => {
    const { subscribe, set, update } = writable(undefined)
    return {
      subscribe,
      set,
      update,
    }
  })()

  function fetchSecrets () {
    if (get(namespaces.selected)) {
      axios.get(`${process.env.API_HOST}/vault/${get(namespaces.selected)}`).then((res) => {
        if (res.data) {
          res.data.forEach(s => {
            s.namespace = get(namespaces.selected)
          })
          set(res.data)
          selected.set(res.data[0])
        }
      }).catch((err) => {
        console.log(err)
      })
    }
  }
  namespaces.selected.subscribe(() => {
    fetchSecrets()
  })

  return {
    selected,
    subscribe,
    select: (value) => {
      selected.set(value)
    },
    update: () => fetchSecrets(),
  }
})()