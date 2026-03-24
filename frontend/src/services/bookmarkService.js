import api from './api';

const handleError = (error, message) => {
  console.error(message, error);
  if (error?.response?.data?.message) {
    throw new Error(error.response.data.message);
  }
  throw new Error(message);
};

const DEFAULT_COVER = '/covers/default-cover.jpg';

const normalizeBookmark = (b) => ({
  id: b.id,

  lastReadChapter: Number(
    b.lastReadChapter ?? b.last_read_chapter ?? 0
  ),

  ...b,

  series: {
    ...b.series,

    lastChapterNumber: Number(
      b.series?.lastChapterNumber ??
      b.series?.last_chapter_number ??
      0
    ),

    cover:
      b.series?.cover ||
      b.series?.cover_image_url ||
      DEFAULT_COVER,
  },
});

export const bookmarkService = {
  getAll: async () => {
    try {
      const res = await api.get('/bookmarks');
      const data = res.data;

      const raw = data?.data || data?.bookmarks || data || [];

      return raw.map(normalizeBookmark);
    } catch (error) {
      handleError(error, 'Erreur lors de la récupération des favoris');
    }
  },

  getOne: async (id) => {
    try {
      const res = await api.get(`/bookmarks/${id}`);
      return normalizeBookmark(res.data?.data || res.data);
    } catch (error) {
      handleError(error, `Erreur lors de la récupération du favori ${id}`);
    }
  },

  create: async (sourceURL) => {
    try {
      if (!sourceURL?.trim()) throw new Error("Lien vide");

      const res = await api.post('/bookmarks', {
        sourceURL: sourceURL.trim()
      });

      return normalizeBookmark(res.data?.data || res.data);
    } catch (error) {
      handleError(error, 'Erreur lors de la création du favori');
    }
  },

  update: async (id, bookmark) => {
    try {
      const res = await api.patch(`/bookmarks/${id}`, bookmark);
      return normalizeBookmark(res.data?.data || res.data);
    } catch (error) {
      handleError(error, `Erreur lors de la mise à jour du favori ${id}`);
    }
  },

  delete: async (id) => {
    try {
      const res = await api.delete(`/bookmarks/${id}`);
      return res.data;
    } catch (error) {
      handleError(error, `Erreur lors de la suppression du favori ${id}`);
    }
  },

  markChapterAsRead: async (bookmarkId, chapterNumber) => {
    try {
      const res = await api.patch(
        `/bookmarks/${bookmarkId}/chapters/${chapterNumber}/read`
      );

      return {
        lastReadChapter: Number(
          res.data?.data?.lastReadChapter ??
          res.data?.data?.last_read_chapter ??
          chapterNumber
        )
      };
    } catch (error) {
      handleError(
        error,
        `Erreur lors du marquage du chapitre ${chapterNumber} comme lu`
      );
    }
  },

  markSeriesAsRead: async (bookmarkId) => {
    try {
      const res = await api.patch(
        `/bookmarks/${bookmarkId}/series/read`
      );

      return {
        lastReadChapter: Number(
          res.data?.data?.lastReadChapter ??
          res.data?.data?.last_read_chapter ??
          0
        )
      };
    } catch (error) {
      handleError(
        error,
        'Erreur lors du marquage de la série comme lue'
      );
    }
  }
};