import { Dexie } from 'dexie'
import { getCookie } from 'typescript-cookie'
import { fire, loading, post } from './misc.svelte'

const db = new Dexie('requirement')
db.version(1).stores({
  requirements: 'id'
})
const table = db.table<ExtendedRequirement>('requirements')

const defaultFieldNames: {
  [key in keyof Requirement]: string
} = {
  id: "编号",
  type: "类型",
  title: "标题",
  date: "日期",
  deadline: "截止日期",
  done: "完成日期",
  submitter: "提交人",
  recipient: "承接人",
  acceptor: "受理人",
  status: "状态",
  label: "标签",
  note: "备注"
}

class Fields {
  raw: MainField[]
  custom: CustomField[]
  #fields: FieldMap
  constructor(fields: MainField[], custom: CustomField[]) {
    this.raw = fields
    this.custom = custom
    this.#fields = fields.reduce((m, field) => {
      if (field.key) {
        if (!field.name)
          field.name = defaultFieldNames[field.key as keyof Requirement]
        m[field.key] = field
      }
      return m
    }, {} as FieldMap)
  }
  enable(key: keyof Requirement) {
    return this.#fields[key] != undefined
  }
  name(key: string) {
    return this.#fields[key as keyof Requirement]?.name ||
      defaultFieldNames[key as keyof Requirement] ||
      this.custom.find(field => field.key === key)?.name ||
      key
  }
  height(key: keyof Requirement) {
    return this.#fields[key]?.height
  }
  required(key: keyof Requirement) {
    return this.#fields[key]?.required
  }
  enum(key: keyof Requirement) {
    return this.#fields[key]?.enum || []
  }
  columns(all?: boolean) {
    const columns = <Field[]>[]
    for (const key in this.#fields)
      if (all || this.#fields[key as keyof Requirement].size)
        columns.push(this.#fields[key as keyof Requirement])
    if (this.custom)
      this.custom.forEach(field => {
        if (all || field.size)
          columns.push(field)
      })
    return columns
  }
  searchable() {
    const fields = <(Field)[]>[]
    for (const key in this.#fields)
      if (this.#fields[key as keyof Requirement].searchable)
        fields.push(this.#fields[key as keyof Requirement])
    if (this.custom)
      this.custom.forEach(field => {
        if (field.searchable)
          fields.push(field)
      })
    return fields
  }
}

class Requirements {
  brand = $state('')
  username = $state('')
  component = $state('show')
  fields = new Fields([], [])
  doneValue = $state('')
  statuses = $state.raw<Status[]>([])
  mode = $state('')
  requirement = $state({} as ExtendedRequirement)
  requirements = $state.raw<ExtendedRequirement[]>([])
  search = $state<Search>({ search: '', field: '', sort: '', desc: true, filter: '', value: '' })
  scrollTop = $state(0)
  controller = $state(new AbortController())
  results = $derived.by(() => {
    let array: ExtendedRequirement[] = []
    if (!this.search.search) array = this.requirements
    else if (this.search.field)
      array = this.requirements.filter((i) =>
        i[this.search.field as keyof ExtendedRequirement].includes(this.search.search),
      )
    else
      array = this.requirements.filter((i) =>
        this.fields
          .searchable()
          .some(field => i[field.key].includes(this.search.search)),
      )
    if (this.search.filter && this.search.value)
      array = array.filter((i) =>
        i[this.search.filter as keyof ExtendedRequirement] === this.search.value,
      )
    if (this.search.sort)
      return array.toSorted((a, b) => {
        const v1 = a[this.search.sort as keyof ExtendedRequirement],
          v2 = b[this.search.sort as keyof ExtendedRequirement]
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
    this.requirement = {} as ExtendedRequirement
    this.requirements = []
    await this.clear()
  }
  async init(load?: Boolean): Promise<string[]> {
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
        this.doneValue = res.done
        this.fields = new Fields(res.fields, res.custom)
        this.statuses = this.#parseStatuses(this.fields.enum('status'))
        if (load) {
          const n = await this.load()
          if (!n) {
            await this.fetch()
            await this.load()
          }
        }
        return res.users
      } else await this.reset()
    } else if (resp.status == 409) {
      await this.clear()
      return await this.init(load)
    } else await this.reset()
    return []
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
  async save(r: ExtendedRequirement) {
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
  async done(r: ExtendedRequirement, date: string) {
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
  #parseStatus(s: string): Status {
    const res = s.split(/[:\s]+/);
    const status: Status = { value: res[0].trim(), closed: false }
    if (res.length > 1) {
      const closed = res[1].toLowerCase().trim()
      if (closed == "1" || closed == "t" || closed == "true")
        status.closed = true
    }
    if (status.value === this.doneValue)
      status.closed = true
    return status
  }
  #parseStatuses(array: string[]) {
    const statuses: Status[] = []
    array.forEach(s => statuses.push(this.#parseStatus(s)))
    return statuses
  }
  #isEqual(a: ExtendedRequirement, b: ExtendedRequirement) {
    for (const k in a) {
      const key = k as keyof ExtendedRequirement
      if (a[key] != b[key]) return false
    }
    return true
  }
  isClosed(r: ExtendedRequirement) {
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




