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

interface ExtendedRequirement extends Requirement {
  [key: string]: string
}

interface Field {
  key: string
  name: string
  size: number
  height: string
  title: boolean
  searchable: boolean
  required: boolean
  enum: string[]
}

interface MainField extends Field {
  key: keyof Requirement
}

type FieldMap = {
  [key in keyof Requirement]: MainField
}

type Type = 'checkbox' | 'date' | 'input' | 'select' | 'textarea'

interface CustomField extends Field {
  type: Type
}

interface Status {
  value: string
  closed: boolean
}

interface Search {
  search: string
  field: string
  sort: string
  desc: boolean
  filter: 'type' | 'status' | ''
  value: string
}
