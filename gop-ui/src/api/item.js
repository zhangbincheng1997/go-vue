import axios from 'axios'
import request from '@/utils/request'
import { getToken } from '@/utils/auth'

export function getStatusOptions() {
  return request({
    url: '/item/status',
    method: 'get'
  })
}

export function getItemList(query) {
  return request({
    url: '/item/list',
    method: 'get',
    params: query
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

export function updateText(table, language, id, text) {
  return request({
    url: '/item/text',
    method: 'put',
    data: {
      table: table,
      language: language,
      id: id,
      text: text
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
      text: text
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

export function importData(formData) {
  return request({
    url: '/item/import',
    method: 'post',
    data: formData
  })
}

export function exportData(table, language) {
  return axios({
    baseURL: process.env.VUE_APP_BASE_API,
    url: '/item/export',
    method: 'get',
    params: {
      table: table,
      language: language
    },
    headers: {
      'Authorization': 'Bearer ' + getToken()
    },
    responseType: 'blob'
  })
}
