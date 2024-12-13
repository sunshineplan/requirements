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
  label: string
  note: string
}

interface Field {
  key?: keyof Requirement
  name: string
  size?: number
  title?: boolean
  searchable?: boolean
}

type FieldMap = {
  [key in keyof Requirement]: Field
}

interface Info {
  brand: string
  username: string
  done: string
  labels: string[]
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
