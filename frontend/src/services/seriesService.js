import api from './api';

export const seriesService = {
  getAll: async () => {
    try {
      const response = await api.get('/series');
      return response.data;
    } catch (error) {
      console.error('Erreur lors de la récupération des séries:', error);
      if (error.response) {
        // L'erreur provient de la réponse du serveur
        throw new Error(`Impossible de récupérer les séries. ${error.response.data.message || 'Veuillez réessayer plus tard.'}`);
      } else if (error.request) {
        // Aucun réponse du serveur
        throw new Error('Aucune réponse du serveur. Veuillez vérifier votre connexion.');
      } else {
        // Erreur dans la configuration de la requête
        throw new Error(`Erreur inconnue: ${error.message}`);
      }
    }
  },

  getOne: async (id) => {
    try {
      const response = await api.get(`/series/${id}`);
      return response.data; 
    } catch (error) {
      console.error(`Erreur lors de la récupération de la série avec ID ${id}:`, error);
      if (error.response) {
        throw new Error(`Impossible de récupérer la série avec ID ${id}. ${error.response.data.message || 'Veuillez réessayer plus tard.'}`);
      } else if (error.request) {
        throw new Error('Aucune réponse du serveur. Veuillez vérifier votre connexion.');
      } else {
        throw new Error(`Erreur inconnue: ${error.message}`);
      }
    }
  },

  create: async (series) => {
    try {
      const response = await api.post('/series', series);
      return response.data;
    } catch (error) {
      console.error('Erreur lors de la création de la série:', error);
      if (error.response) {
        throw new Error(`Impossible de créer la série. ${error.response.data.message || 'Veuillez réessayer plus tard.'}`);
      } else if (error.request) {
        throw new Error('Aucune réponse du serveur. Veuillez vérifier votre connexion.');
      } else {
        throw new Error(`Erreur inconnue: ${error.message}`);
      }
    }
  },

  getChapters: async (url) => {
    try {
      const response = await api.get('/series/chapters', { params: { url } });
      return response.data; 
    } catch (error) {
      console.error('Erreur lors de la récupération des chapitres:', error);
      if (error.response) {
        throw new Error(`Impossible de récupérer les chapitres. ${error.response.data.message || 'Veuillez réessayer plus tard.'}`);
      } else if (error.request) {
        throw new Error('Aucune réponse du serveur. Veuillez vérifier votre connexion.');
      } else {
        throw new Error(`Erreur inconnue: ${error.message}`);
      }
    }
  },
};
