interface Requirement {
  id: string
  type: string
  desc: string
  date: string
  deadline: string
  done: string
  submitter: string
  recipient: string
  acceptor: string
  status: string
  note: string
  participating: string
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
  name: string
  username: string
  done: string
  participants: string[]
  types: string[]
  users: string[]
}

interface Status {
  value: string
  closed: boolean
}

