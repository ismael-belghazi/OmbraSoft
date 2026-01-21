import api from './api'

export const authService = {
  register: (email, password) => {
    return api.post('/auth/register', { email, password })
  },

  login: (email, password) => {
    return api.post('/auth/login', { email, password })
  },

  logout: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }
}
