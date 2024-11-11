class Fields {
  #fields: FieldMap
  constructor(m: FieldMap) {
    this.#fields = m
  }
  name(key: keyof Requirement) {
    return this.#fields[key].name
  }
  size(key: keyof Requirement) {
    return this.#fields[key].size
  }
  title(key: keyof Requirement) {
    return /编号|类型|日期|班组/i.test(fields.name(key))
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
export const fields = new Fields({
  id: {
    name: '编号',
    size: 6,
    searchable: true
  },
  type: {
    name: '类型',
    size: 6,
    searchable: false
  },
  desc: {
    name: '描述',
    size: -1,
    searchable: true
  },
  date: {
    name: '提请日期',
    size: 8,
    searchable: true
  },
  deadline: {
    name: '期限日期',
    size: 8,
    searchable: true
  },
  done: {
    name: '完成日期',
    size: 8,
    searchable: true
  },
  submitter: {
    name: '提交人',
    size: 5,
    searchable: true
  },
  recipient: {
    name: '承接人',
    size: 0,
    searchable: true
  },
  acceptor: {
    name: '受理人',
    size: 5,
    searchable: true
  },
  status: {
    name: '状态',
    size: 5,
    searchable: false
  },
  note: {
    name: '备注',
    size: 9,
    searchable: true
  },
  participating: {
    name: '参与班组',
    size: 6,
    searchable: true
  },
})
