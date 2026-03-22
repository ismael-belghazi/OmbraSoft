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

  // Redirection si non connecté
  useEffect(() => {
    if (!user) {
      navigate('/login')
    }
  }, [user, navigate])

  // Récupération des bookmarks de l'utilisateur
  useEffect(() => {
    if (!user) return

    let isMounted = true

    const fetchBookmarks = async () => {
      try {
        setLoading(true)
        const response = await bookmarkService.getAll()
        if (isMounted) {
          const userBookmarks = response.data.bookmarks || []

          // Trie par dernier chapitre lu (descendant)
          userBookmarks.sort(
            (a, b) => (b.lastReadChapter || 0) - (a.lastReadChapter || 0)
          )

          setBookmarks(userBookmarks)
          setError('')
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
  }, [user])

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
                    <span>
                      {bookmark.series?.chapters?.length > 0 && bookmark.lastReadChapter ? (
                        <a
                          href={bookmark.series.chapters.find(
                            (c) => c.number === bookmark.lastReadChapter
                          )?.url || bookmark.series.sourceURL}
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          Chapitre {bookmark.lastReadChapter}
                        </a>
                      ) : (
                        'N/A'
                      )}
                    </span>
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