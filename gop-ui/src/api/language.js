import request from '@/utils/request'

export function getToken() {
  return request({
    url: '/item/token',
    method: 'get',
  })
}

export function getStatusOptions() {
  return request({
    url: '/item/status',
    method: 'get'
  })
}

export function getList(query) {
  return request({
    url: '/item/list',
    method: 'get',
    params: query
  })
}

export function add(table, language, item) {
  return request({
    url: '/item',
    method: 'post',
    data: {
      table: table,
      language: language,
      item: item
    }
  })
}

export function remove(table, language, ids) {
  return request({
    url: '/item',
    method: 'delete',
    data: {
      table: table,
      language: language,
      ids: ids
    }
  })
}

export function update(table, language, item) {
  return request({
    url: '/item',
    method: 'put',
    data: {
      table: table,
      language: language,
      item: item
    }
  })
}

export function updateStatus(table, language, ids, status) {
  return request({
    url: '/item/status',
    method: 'put',
    data: {
      table: table,
      language: language,
      ids: ids,
      status: status
    }
  })
}

export function updateText(table, language, id, text) {
  return request({
    url: '/item/text',
    method: 'put',
    data: {
      table: table,
      language: language,
      id: id,
      text: text,
    }
  })
}

export function updateRecordText(table, language, id, text) {
  return request({
    url: '/item/record/text',
    method: 'put',
    data: {
      table: table,
      language: language,
      id: id,
      text: text,
    }
  })
}

export function getRecordList(table, language, id) {
  return request({
    url: '/item/record',
    method: 'get',
    params: {
      table: table,
      language: language,
      id: id,
    }
  })
}

export function importData(formData) {
  return request({
    url: '/item/import',
    method: 'post',
    data: formData
  })
}

export function exportText(table, language) {
  return request({
    url: '/item/export',
    method: 'get',
    params: {
      table: table,
      language: language,
      method: 'exportText'
    }
  })
}

export function exportImageExcel(table, language) {
  return request({
    url: '/item/export',
    method: 'get',
    params: {
      table: table,
      language: language,
      method: 'exportImageExcel'
    }
  })
}

export function exportImageZip(table, language) {
  return request({
    url: '/item/export',
    method: 'get',
    params: {
      table: table,
      language: language,
      method: 'exportImageZip'
    }
  })
}
