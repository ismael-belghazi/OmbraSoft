import api from './api'

const handleError = (error, defaultMessage) => {
  if (!error.response) {
    throw new Error('Aucune réponse du serveur. Vérifiez votre connexion.')
  }

  const message =
    error.response?.data?.message ||
    error.response?.data?.error ||
    defaultMessage

  throw new Error(message)
}

const validateId = (id) => {
  if (!id || typeof id !== 'string') {
    throw new Error('ID invalide')
  }
}

export const seriesService = {
  getAll: async () => {
    try {
      const { data } = await api.get('/series')
      return data
    } catch (error) {
      handleError(error, 'Impossible de récupérer les séries.')
    }
  },

  getOne: async (id) => {
    validateId(id)

    try {
      const { data } = await api.get(`/series/${id}`)
      return data
    } catch (error) {
      handleError(error, `Impossible de récupérer la série ${id}.`)
    }
  },

  create: async (series) => {
    if (!series?.title || !series?.sourceURL) {
      throw new Error('Données série invalides')
    }

    try {
      const { data } = await api.post('/series', series)
      return data
    } catch (error) {
      handleError(error, 'Impossible de créer la série.')
    }
  },

  getChapters: async (url) => {
    if (!url) {
      throw new Error('URL requise')
    }

    try {
      const { data } = await api.get('/series/chapters', {
        params: { url },
      })
      return data
    } catch (error) {
      handleError(error, 'Impossible de récupérer les chapitres.')
    }
  },

  getChaptersBySeriesId: async (seriesId) => {
    validateId(seriesId)

    try {
      const { data } = await api.get(`/series/${seriesId}/chapters`)
      return data
    } catch (error) {
      handleError(error, 'Impossible de récupérer les chapitres.')
    }
  },
}
