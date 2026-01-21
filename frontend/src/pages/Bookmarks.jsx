import { useState, useEffect } from 'react'
import { bookmarkService } from '../services/bookmarkService'
import '../styles/pages.css'

export default function Bookmarks() {
  const [bookmarks, setBookmarks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let isMounted = true

    const fetchBookmarks = async () => {
      try {
        const response = await bookmarkService.getAll()
        if (isMounted) {
          setBookmarks(response.data.bookmarks || [])
        }
      } catch (err) {
        if (isMounted) {
          setError('Erreur lors du chargement des favoris')
        }
      } finally {
        if (isMounted) {
          setLoading(false)
        }
      }
    }

    fetchBookmarks()

    return () => {
      isMounted = false
    }
  }, [])

  const handleDelete = async (id) => {
    if (confirm('Êtes-vous sûr de vouloir supprimer ce favori ?')) {
      try {
        await bookmarkService.delete(id)
        setBookmarks(bookmarks.filter(b => b.id !== id))
      } catch (err) {
        setError('Erreur lors de la suppression')
      }
    }
  }

  return (
    <div className="page-container">
      <h1>Mes Favoris</h1>

      {error && <div className="error-message">{error}</div>}

      {loading ? (
        <p>Chargement...</p>
      ) : (
        <div className="bookmarks-container">
          {bookmarks.length === 0 ? (
            <p>Aucun favori pour le moment</p>
          ) : (
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
                      <a href={bookmark.series?.source_url} target="_blank" rel="noopener noreferrer">
                        Lire
                      </a>
                    </td>
                    <td>
                      <button 
                        onClick={() => handleDelete(bookmark.id)}
                        className="delete-btn"
                      >
                        Supprimer
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}
    </div>
  )
}
