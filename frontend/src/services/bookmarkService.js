import api from './api'

export const bookmarkService = {
  getAll: async () => {
    try {
      const response = await api.get('/bookmarks');
      return response.data;
    } catch (error) {
      console.error('Erreur lors de la récupération des favoris:', error);
      throw new Error('Impossible de récupérer les favoris. Veuillez vérifier votre connexion ou réessayer plus tard.');
    }
  },

  getOne: async (id) => {
    try {
      const response = await api.get(`/bookmarks/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Erreur lors de la récupération du favori avec ID ${id}:`, error);
      throw new Error(`Impossible de récupérer le favori avec ID ${id}. Veuillez vérifier l'ID ou réessayer plus tard.`);
    }
  },

  create: async (bookmark) => {
    try {
      const response = await api.post('/bookmarks', bookmark);
      return response.data;
    } catch (error) {
      console.error('Erreur lors de la création du favori:', error);
      throw new Error('Impossible de créer le favori. Veuillez vérifier les données envoyées ou réessayer plus tard.');
    }
  },

  update: async (id, bookmark) => {
    try {
      const response = await api.patch(`/bookmarks/${id}`, bookmark);
      return response.data;
    } catch (error) {
      console.error(`Erreur lors de la mise à jour du favori avec ID ${id}:`, error);
      throw new Error(`Impossible de mettre à jour le favori avec ID ${id}. Veuillez réessayer plus tard.`);
    }
  },

  delete: async (id) => {
    try {
      const response = await api.delete(`/bookmarks/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Erreur lors de la suppression du favori avec ID ${id}:`, error);
      throw new Error(`Impossible de supprimer le favori avec ID ${id}. Veuillez réessayer plus tard.`);
    }
  }
}
