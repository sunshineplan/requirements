import { writable, get } from 'svelte/store'
import { Dexie } from 'dexie'
import { fire, post } from './misc'

const db = new Dexie('requirement')
db.version(1).stores({
  requirements: 'id'
})

export const requirement = writable(<Requirement>{})

const createRequirements = () => {
  const { subscribe, set } = writable(<Requirement[]>[])
  return {
    subscribe,
    set,
    clear: async () => await db.table('requirements').clear(),
    load: async () => {
      const array = await db.table<Requirement>('requirements').toArray()
      requirements.set(array.reverse())
      return array.length
    },
    fetch: async () => {
      const resp = await fetch('/get')
      if (resp.ok) await db.table('requirements').bulkAdd(await resp.json())
      else await fire('Fatal', await resp.text(), 'error')
    },
    save: async (r: Requirement) => {
      let resp: Response | undefined = undefined
      if (r.id) {
        if (isEqual(get(requirement), r)) return 0
        resp = await post('/edit', { old: get(requirement), new: r })
      }
      else resp = await post('/add', r)
      if (resp.ok) {
        const res = await resp.json()
        if (res.status == 1) {
          if (res.reload) {
            await requirements.clear()
            await requirements.fetch()
            await requirements.load()
          } else {
            if (r.id) await db.table('requirements').update(r.id, r)
            else {
              r.id = res.id
              await db.table('requirements').add(r)
            }
            await requirements.load()
          }
        } else {
          await fire('Error', res.message, 'error')
          return <number>res.error
        }
      } else await fire('Fatal', await resp.text(), 'error')
      return 0
    },
    done: async (r: Requirement) => {
      if (r.status == '已完成') return 0
      r.status = '已完成'
      return requirements.save(r)
    },
    delete: async (requirement: Requirement) => {
      const resp = await post('/delete/' + requirement.id)
      if (resp.ok) {
        const res = await resp.json()
        if (res.reload) {
          await requirements.clear()
          await requirements.fetch()
        } else await db.table('requirements').where('id').equals(requirement.id).delete()
        await requirements.load()
      } else await fire('Fatal', await resp.text(), 'error')
    },
    submitters: async () => {
      return [...new Set((await db.table<Requirement>('requirements').toArray()).map(i => i.submitter))].sort()
    },
    recipients: async () => {
      return [...new Set((await db.table<Requirement>('requirements').toArray()).map(i => i.recipient))].sort()
    },
    acceptors: async () => {
      return [...new Set((await db.table<Requirement>('requirements').toArray()).map(i => i.acceptor))].sort()
    }
  }
}
export const requirements = createRequirements()

export const info = async (load?: Boolean): Promise<Info> => {
  const resp = await fetch('/info')
  if (resp.ok) {
    const res = await resp.json()
    if (res.username) {
      if (load) {
        const n = await requirements.load()
        if (!n) {
          await requirements.fetch()
          await requirements.load()
        }
      }
      return { username: res.username, participants: res.participants, types: res.types, users: res.users }
    } else await reset()
  } else if (resp.status == 409) {
    await requirements.clear()
    return await info(load)
  } else await reset()
  return <Info>{ username: '' }
}

const reset = async () => {
  requirement.set(<Requirement>{})
  requirements.set([])
  await requirements.clear()
}

const isEqual = (a: Requirement, b: Requirement): boolean => {
  for (const k in a) {
    const key = k as keyof Requirement
    if (a[key] != b[key]) return false
  }
  return true
}
