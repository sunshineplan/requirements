import Swal, { type SweetAlertIcon } from 'sweetalert2'
import { loading } from './stores'

export const fire = async (title?: string, html?: string, icon?: SweetAlertIcon) => {
  const swal = Swal.mixin({
    confirmButtonText: '确定',
    customClass: { confirmButton: 'swal btn btn-primary' },
    buttonsStyling: false
  })
  await swal.fire(title, html, icon)
  if (title == '严重错误') throw html
}

export const valid = () => {
  let result = true
  Array.from(document.querySelectorAll('input'))
    .forEach(i => { if (!i.checkValidity()) result = false })
  Array.from(document.querySelectorAll('select'))
    .forEach(i => { if (!i.checkValidity()) result = false })
  Array.from(document.querySelectorAll('textarea'))
    .forEach(i => { if (!i.checkValidity()) result = false })
  return result
}

export const poll = async (signal: AbortSignal) => {
  let resp: Response
  try {
    resp = await fetch('/poll', { signal })
  } catch (e) {
    let message = ''
    if (typeof e === 'string')
      message = e
    else if (e instanceof Error)
      message = e.message
    resp = new Response(message, { 'status': 500 })
  }
  return resp
}

export const post = async (url: string, data?: any) => {
  let resp: Response
  const init: RequestInit = {
    method: 'post',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }
  loading.start()
  try {
    resp = await fetch(url, init)
  } catch (e) {
    let message = ''
    if (typeof e === "string") {
      message = e
    } else if (e instanceof Error) {
      message = e.message
    }
    resp = new Response(message, { "status": 500 })
  }
  loading.end()
  if (resp.status == 401) {
    await fire('错误', '登录状态已变更，请重新登录！', 'error')
    window.location.href = '/'
  } else if (resp.status == 409) {
    await fire('错误', '本地数据与服务器数据不一致，按确定后刷新数据！', 'error')
    window.location.href = '/'
  }
  return resp
}

export const confirm = async (text: string, isDanger?: boolean) => {
  let focusCancel = false, confirmButton = 'swal btn btn-primary'
  if (isDanger) focusCancel = true, confirmButton = 'swal btn btn-danger'
  const confirm = await Swal.fire({
    title: '你确定吗？',
    text,
    icon: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    showCancelButton: true,
    focusCancel,
    customClass: {
      confirmButton,
      cancelButton: 'swal btn btn-primary'
    },
    buttonsStyling: false
  })
  return confirm.isConfirmed
}
