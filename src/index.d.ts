interface Requirement {
  id: string
  type: string
  title: string
  date: string
  deadline: string
  done: string
  submitter: string
  recipient: string
  acceptor: string
  status: string
  note: string
  group: string
}

interface Field {
  name: string
  size: number
  searchable: boolean
}

type FieldMap = {
  [key in keyof Requirement]: Field
}

interface Info {
  brand: string
  username: string
  done: string
  groups: string[]
  types: string[]
  users: string[]
}

interface Status {
  value: string
  closed: boolean
}

interface Search {
  search: string
  field: keyof Requirement | ''
  sort: keyof Requirement | ''
  desc: boolean
  filter: '' | 'type' | 'status'
  value: string
}
