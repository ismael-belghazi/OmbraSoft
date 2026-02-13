import { useState, useEffect } from 'react';
import { bookmarkService } from '../services/bookmarkService';
import '../styles/css.css';

const useBookmarks = () => {
  const [bookmarks, setBookmarks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const fetchBookmarks = async () => {
    try {
      const response = await bookmarkService.getAll();
      setBookmarks(response.data.bookmarks || []);
    } catch (err) {
      setError('Erreur lors du chargement des favoris');
    } finally {
      setLoading(false);
    }
  };

  const deleteBookmark = async (id) => {
    try {
      await bookmarkService.delete(id);
      setBookmarks(prev => prev.filter(b => b.id !== id));
    } catch (err) {
      setError('Erreur lors de la suppression');
    }
  };

  useEffect(() => {
    fetchBookmarks();
  }, []);

  return { bookmarks, loading, error, deleteBookmark };
};

// Composant principal
export default function Bookmarks() {
  const [newSeries, setNewSeries] = useState({ title: '', sourceSite: '', sourceURL: '' });
  const [adding, setAdding] = useState(false);
  const { bookmarks, loading, error, deleteBookmark } = useBookmarks();

  const handleAdd = async (e) => {
    e.preventDefault();

    if (!newSeries.title.trim() || !newSeries.sourceURL.trim()) {
      setError('Le titre et l\'URL sont obligatoires');
      return;
    }

    const urlRegex = /^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$/i;
    if (!urlRegex.test(newSeries.sourceURL)) {
      setError('L\'URL fournie n\'est pas valide');
      return;
    }

    setAdding(true);
    setError('');

    try {
      const response = await bookmarkService.create(newSeries);
      setNewSeries({ title: '', sourceSite: '', sourceURL: '' });
      alert('Favori ajouté avec succès'); 
    } catch (err) {
      setError('Erreur lors de l\'ajout du favori');
    } finally {
      setAdding(false);
    }
  };

  return (
    <div className="page-container">
      <h1>Mes Favoris</h1>

      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleAdd} className="add-bookmark-form">
        <input
          type="text"
          placeholder="Titre"
          value={newSeries.title}
          onChange={(e) => setNewSeries({ ...newSeries, title: e.target.value })}
          disabled={adding}
        />
        <input
          type="text"
          placeholder="Source"
          value={newSeries.sourceSite}
          onChange={(e) => setNewSeries({ ...newSeries, sourceSite: e.target.value })}
          disabled={adding}
        />
        <input
          type="url"
          placeholder="URL"
          value={newSeries.sourceURL}
          onChange={(e) => setNewSeries({ ...newSeries, sourceURL: e.target.value })}
          disabled={adding}
        />
        <button type="submit" disabled={adding}>
          {adding ? 'Ajout en cours...' : 'Ajouter'}
        </button>
      </form>

      {loading ? (
        <p>Chargement...</p>
      ) : bookmarks.length === 0 ? (
        <p>Aucun favori pour le moment</p>
      ) : (
        <div className="bookmarks-container">
          <table className="bookmarks-table">
            <thead>
              <tr>
                <th>Titre</th>
                <th>Dernier chapitre lu</th>
                <th>Source</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {bookmarks.map((bookmark) => (
                <tr key={bookmark.id}>
                  <td>{bookmark.series?.title || '-'}</td>
                  <td>{bookmark.last_read_chapter || 0}</td>
                  <td>
                    {bookmark.series?.source_url ? (
                      <a
                        href={bookmark.series.source_url}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        Lire
                      </a>
                    ) : '-'}
                  </td>
                  <td>
                    <button
                      onClick={() => deleteBookmark(bookmark.id)}
                      className="delete-btn"
                    >
                      Supprimer
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
