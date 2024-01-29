import { writable, get } from 'svelte/store'

export const component = writable('show')

export const mode = writable('')

export const search = writable('')
export const sort = writable('')
export const desc = writable(true)

const scrollTop = writable(0)
export const save = () => {
  const table = document.querySelector('.table-responsive')
  if (table) scrollTop.set(table.scrollTop)
}
export const scroll = (restore?: Boolean) => {
  if (!restore) {
    sort.set('')
    desc.set(true)
    scrollTop.set(0)
  } else if (get(sort) == '编号' && get(desc)) {
    sort.set('')
  }
  const table = document.querySelector('.table-responsive')
  if (table) {
    table.scrollTop = get(scrollTop)
    table.scrollLeft = 0
  }
}

export const goto = (s: string) => {
  if (get(component) == 'show') save()
  let url = ''
  switch (s) {
    case 'add':
    case 'edit':
    case 'view':
      url = s
      mode.set(s)
      component.set('requirement')
      break
    default:
      component.set(s)
  }
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
