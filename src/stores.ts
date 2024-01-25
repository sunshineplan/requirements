import { writable } from 'svelte/store'

export const component = writable('show')

export const goHome = () => {
  window.history.pushState({}, '', '/')
  component.set('show')
}

export const mode = writable('')

const createLoading = () => {
  const { subscribe, update } = writable(0)
  return {
    subscribe,
    start: () => update(n => n + 1),
    end: () => update(n => n - 1)
  }
}
export const loading = createLoading()
