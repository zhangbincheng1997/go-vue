import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

export function getInfo(token) {
  return request({
    url: '/user/info',
    method: 'get',
    params: token
  })
}

export function logout() {
  return request({
    url: '/user/logout',
    method: 'post'
  })
}

export function register(data) {
  return request({
    url: '/user/register',
    method: 'post',
    data: data
  })
}

export function getUserList(data) {
  return request({
    url: '/user/list',
    method: 'get',
    params: data
  })
}

export function updatePassword(data) {
  return request({
    url: '/user/password',
    method: 'put',
    data: data
  })
}

export function updateInfo(data) {
  return request({
    url: '/user/info',
    method: 'put',
    data: data
  })
}
export function updateRole(id, role) {
  return request({
    url: '/user/role',
    method: 'put',
    data: {
      id: id,
      role: role
    }
  })
}

export function deleteUser(id) {
  return request({
    url: '/user',
    method: 'delete',
    data: {
      id: id
    }
  })
}
