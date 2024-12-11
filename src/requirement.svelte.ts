import { Dexie } from 'dexie'
import { getCookie } from 'typescript-cookie'
import { fire, loading, post } from './misc.svelte'

const db = new Dexie('requirement')
db.version(1).stores({
  requirements: 'id'
})
const table = db.table<Requirement>('requirements')

class Fields {
  #fields: FieldMap
  constructor(fields: Field[]) {
    this.#fields = fields.reduce((m, field) => {
      if (field.key) { m[field.key] = field }
      return m
    }, {} as FieldMap)
  }
  name(key: keyof Requirement) {
    return this.#fields[key].name
  }
  size(key: keyof Requirement) {
    return this.#fields[key].size
  }
  title(key: keyof Requirement) {
    return !(<(keyof Requirement)[]>['id', 'type', 'date', 'deadline', 'done']).includes(key)
  }
  columns(all?: boolean) {
    const columns = <(keyof Requirement)[]>[]
    for (const key in this.#fields)
      if (all || this.#fields[key as keyof Requirement].size)
        columns.push(key as keyof Requirement)
    return columns
  }
  searchable() {
    const fields = <(keyof Requirement)[]>[]
    for (const key in this.#fields)
      if (this.#fields[key as keyof Requirement].searchable)
        fields.push(key as keyof Requirement)
    return fields
  }
}

class Requirements {
  brand = $state('')
  username = $state('')
  component = $state('show')
  fields = new Fields([])
  mode = $state('')
  requirement = $state({} as Requirement)
  requirements = $state.raw<Requirement[]>([])
  types = $state<string[]>([])
  statuses = $state.raw<Status[]>([])
  search = $state<Search>({ search: '', field: '', sort: '', desc: true, filter: '', value: '' })
  scrollTop = $state(0)
  controller = $state(new AbortController())
  results = $derived.by(() => {
    let array: Requirement[] = []
    if (!this.search.search) array = this.requirements
    else if (this.search.field)
      array = this.requirements.filter((i) =>
        i[this.search.field as keyof Requirement].includes(this.search.search),
      )
    else
      array = this.requirements.filter((i) =>
        this.fields
          .searchable()
          .some((field) => i[field].includes(this.search.search)),
      )
    if (this.search.filter && this.search.value)
      array = array.filter((i) =>
        i[this.search.filter as keyof Requirement] === this.search.value,
      )
    if (this.search.sort)
      return array.toSorted((a, b) => {
        const v1 = a[this.search.sort as keyof Requirement],
          v2 = b[this.search.sort as keyof Requirement]
        let res = 0
        if (v1 < v2) res = 1
        else if (v1 > v2) res = -1
        if (this.search.desc) return res
        else return -res
      })
    else return array.sort()
  })
  async clear() { await table.clear() }
  async reset() {
    this.username = ''
    this.requirement = {} as Requirement
    this.requirements = []
    await this.clear()
  }
  async init(load?: Boolean): Promise<Info> {
    loading.start()
    let resp: Response
    try {
      resp = await fetch('/info')
    } catch (e) {
      console.error(e)
      resp = new Response(null, { "status": 500 })
    }
    loading.end()
    if (resp.ok) {
      const res = await resp.json()
      this.brand = res.brand
      if (res.username) {
        this.username = res.username
        this.fields = new Fields(res.fields)
        this.types = res.types
        this.statuses = res.statuses
        if (load) {
          const n = await this.load()
          if (!n) {
            await this.fetch()
            await this.load()
          }
        }
        return {
          done: res.done,
          groups: res.groups,
          users: res.users
        } as Info
      } else await this.reset()
    } else if (resp.status == 409) {
      await this.clear()
      return await this.init(load)
    } else await this.reset()
    return {} as Info
  }
  async load() {
    const array = await table.toArray()
    this.requirements = array.reverse()
    return array.length
  }
  async fetch() {
    const resp = await fetch('/get')
    if (resp.ok) await table.bulkAdd(await resp.json())
    else await fire('Fatal', await resp.text(), 'error')
  }
  async save(r: Requirement) {
    let resp: Response | undefined = undefined
    if (r.id) {
      if (this.#isEqual(this.requirement, r)) return 0
      resp = await post('/edit', { old: this.requirement, new: r })
    }
    else resp = await post('/add', r)
    if (resp.ok) {
      const res = await resp.json()
      if (res.status == 1) {
        if (res.reload) {
          await this.clear()
          await this.fetch()
          await this.load()
        } else {
          if (r.id) await table.update(r.id, r)
          else {
            r.id = res.id
            await table.add(r)
          }
          await this.load()
        }
      } else {
        await fire('Error', res.message, 'error')
        return <number>res.error
      }
    } else await fire('Fatal', await resp.text(), 'error')
    return 0
  }
  async done(r: Requirement, date: string) {
    const resp = await post(`/done?date=${date}`, r)
    if (resp.ok) {
      const res = await resp.json()
      if (res.status == 1) {
        if (res.reload) {
          await this.clear()
          await this.fetch()
          await this.load()
        } else {
          if (res.value) {
            r.status = res.value
            r.done = res.done
            await table.update(r.id, r)
            await this.load()
          }
        }
      } else {
        await fire('Error', res.message, 'error')
        return <number>res.error
      }
    } else await fire('Fatal', await resp.text(), 'error')
    return 0
  }
  async delete(r: Requirement) {
    const resp = await post('/delete/' + r.id)
    if (resp.ok) {
      const res = await resp.json()
      if (res.reload) {
        await this.clear()
        await this.fetch()
      } else await table.where('id').equals(r.id).delete()
      await this.load()
    } else await fire('Fatal', await resp.text(), 'error')
  }
  #isEqual(a: Requirement, b: Requirement) {
    for (const k in a) {
      const key = k as keyof Requirement
      if (a[key] != b[key]) return false
    }
    return true
  }
  isClosed(r: Requirement) {
    const status = this.statuses.find(e => e.value === r.status)
    return status ? status.closed : false
  }
  async submitters() {
    return [...new Set((await table.toArray()).map(i => i.submitter))].sort()
  }
  async recipients() {
    return [...new Set((await table.toArray()).map(i => i.recipient))].sort()
  }
  async acceptors() {
    return [...new Set((await table.toArray()).map(i => i.acceptor))].sort()
  }
  goto(s: string) {
    let comp = s
    let url = s
    switch (s) {
      case 'add':
      case 'edit':
      case 'view':
        comp = 'requirement'
        this.mode = s
      case 'setting':
        this.saveScrollTop()
        break
      case 'show':
        url = ''
    }
    this.component = comp
    window.history.pushState({}, '', '/' + url)
  }
  scroll(restore?: Boolean) {
    if (!restore) {
      this.search.sort = ''
      this.search.desc = true
      this.scrollTop = 0
    } else if (this.search.sort == 'id' && this.search.desc) {
      this.search.sort = ''
    }
    const table = document.querySelector('.table-responsive')
    if (table) {
      table.scrollTop = this.scrollTop
      table.scrollLeft = 0
    }
  }
  saveScrollTop() {
    const table = document.querySelector('.table-responsive')
    if (table) this.scrollTop = table.scrollTop
  }
  clearSearch() {
    this.search = { search: '', field: '', sort: '', desc: true, filter: '', value: '' }
    this.scrollTop = 0
  }
  async subscribe(reset?: boolean) {
    if (reset)
      this.controller = new AbortController()
    let resp: Response
    try {
      resp = await fetch('/poll', { signal: this.controller.signal })
    } catch (e) {
      if (e instanceof DOMException && e.name === 'AbortError') return
      console.error(e)
      resp = new Response(null, { status: 500 })
    }
    if (resp.ok) {
      const last = await resp.text()
      if (last && getCookie('last') != last) await this.init()
      await this.subscribe()
    } else if (resp.status == 401) {
      await this.init()
    } else {
      await new Promise((sleep) => setTimeout(sleep, 30000))
      await this.subscribe()
    }
  }
}
export const requirements = new Requirements




