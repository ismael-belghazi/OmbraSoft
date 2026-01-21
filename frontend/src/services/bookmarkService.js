import api from './api'

export const bookmarkService = {
  getAll: () => {
    return api.get('/bookmarks')
  },

  getOne: (id) => {
    return api.get(`/bookmarks/${id}`)
  },

  create: (bookmark) => {
    return api.post('/bookmarks', bookmark)
  },

  update: (id, bookmark) => {
    return api.patch(`/bookmarks/${id}`, bookmark)
  },

  delete: (id) => {
    return api.delete(`/bookmarks/${id}`)
  }
}
