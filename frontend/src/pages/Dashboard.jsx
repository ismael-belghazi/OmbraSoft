import { useState, useEffect } from 'react'
import { useAuth } from '../hooks/useAuth'
import { bookmarkService } from '../services/bookmarkService'
import { useNavigate } from 'react-router-dom'
import '../styles/css.css'

export default function Dashboard() {
  const { user } = useAuth()
  const navigate = useNavigate()
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

  return (
    <div className="page-container">
      <h1>Bienvenue, {user?.email}</h1>
      {error && <div className="error-message">{error}</div>}
      {loading ? (
        <p>Chargement...</p>
      ) : (
        <div className="dashboard-content">
          <section className="stats">
            <div className="stat-card">
              <h3>Favoris</h3>
              <p className="stat-number">{bookmarks.length}</p>
            </div>
          </section>

          <section className="recent-bookmarks">
            <h2>Vos derniers favoris</h2>
            {bookmarks.length === 0 ? (
              <p>Aucun favori pour le moment</p>
            ) : (
              <ul className="bookmarks-list">
                {bookmarks.map((bookmark) => (
                  <li key={bookmark.id} className="bookmark-item">
                    <span>{bookmark.series?.title || 'Série'}</span>
                    <span>Chapitre {bookmark.last_read_chapter}</span>
                  </li>
                ))}
              </ul>
            )}
          </section>
        </div>
      )}
    </div>
  )
}
