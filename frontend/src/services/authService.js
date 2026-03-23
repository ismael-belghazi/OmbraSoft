import api from './api'

export const authService = {
  register: (email, password, secretPhrase) =>
    api.post('/auth/register', { email, password, secretPhrase }),

  login: (email, password) =>
    api.post('/auth/login', { email, password }),

  logout: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }
}
