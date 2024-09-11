import { get, writable } from 'svelte/store'

export const name = writable('')

export const username = writable('')

export const component = writable('show')

export const mode = writable('')

export const search = writable('')
export const sort = writable('')
export const desc = writable(true)

const scrollTop = writable(0)
export const saveScrollTop = () => {
  const table = document.querySelector('.table-responsive')
  if (table) scrollTop.set(table.scrollTop)
}
export const scroll = (restore?: Boolean) => {
  if (!restore) {
    sort.set('')
    desc.set(true)
    scrollTop.set(0)
  } else if (get(sort) == 'id' && get(desc)) {
    sort.set('')
  }
  const table = document.querySelector('.table-responsive')
  if (table) {
    table.scrollTop = get(scrollTop)
    table.scrollLeft = 0
  }
}

export const goto = (s: string) => {
  let comp = s
  let url = s
  switch (s) {
    case 'add':
    case 'edit':
    case 'view':
      comp = 'requirement'
      mode.set(s)
      saveScrollTop()
      break
    case 'setting':
      saveScrollTop()
      break
    case 'show':
      url = ''
  }
  component.set(comp)
  window.history.pushState({}, '', '/' + url)
}

export const clear = () => {
  search.set('')
  sort.set('')
  desc.set(true)
  scrollTop.set(0)
}

const createLoading = () => {
  const { subscribe, update } = writable(0)
  return {
    subscribe,
    start: () => update(n => n + 1),
    end: () => update(n => n - 1)
  }
}
export const loading = createLoading()
