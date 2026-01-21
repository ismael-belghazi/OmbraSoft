import api from './api'

export const seriesService = {
  getAll: () => {
    return api.get('/series')
  },

  getOne: (id) => {
    return api.get(`/series/${id}`)
  },

  create: (series) => {
    return api.post('/series', series)
  }
}
