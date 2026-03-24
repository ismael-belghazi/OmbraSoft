import { useState, useEffect } from 'react'
import { useAuth } from '../hooks/useAuth'
import { bookmarkService } from '../services/bookmarkService'
import { useNavigate } from 'react-router-dom'
import '../styles/Dashboard.css'

const SERIES_PER_PAGE = 10

export default function Dashboard() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [bookmarks, setBookmarks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [seriesPage, setSeriesPage] = useState(1)

  useEffect(() => {
    if (!user) navigate('/login')
  }, [user, navigate])

  useEffect(() => {
    if (!user) return

    let isMounted = true

    const fetchBookmarks = async () => {
      try {
        setLoading(true)
        const data = await bookmarkService.getAll()
        if (isMounted) {
          const userBookmarks = Array.isArray(data) ? data : []
          userBookmarks.sort(
            (a, b) => (b.lastReadChapter || 0) - (a.lastReadChapter || 0)
          )
          setBookmarks(userBookmarks)
          setError('')
        }
      } catch {
        if (isMounted) setError('Erreur lors du chargement des favoris')
      } finally {
        if (isMounted) setLoading(false)
      }
    }

    fetchBookmarks()
    return () => { isMounted = false }
  }, [user])

  const totalPages = Math.ceil(bookmarks.length / SERIES_PER_PAGE)
  const paginatedBookmarks = bookmarks.slice(
    (seriesPage - 1) * SERIES_PER_PAGE,
    seriesPage * SERIES_PER_PAGE
  )

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
            {paginatedBookmarks.length === 0 ? (
              <p>Aucun favori pour le moment</p>
            ) : (
              <ul className="bookmarks-list">
                {paginatedBookmarks.map((bookmark) => {
                  const lastChapterNumber = bookmark.lastReadChapter || 0
                  const chapter =
                    bookmark.series?.chapters?.find(
                      (c) => Number(c.number) === Number(lastChapterNumber)
                    ) || {}

                  const chapterUrl = chapter.url || bookmark.series?.sourceURL || '#'

                  return (
                    <li key={bookmark.id} className="bookmark-item">
                      <span>{bookmark.series?.title || 'Série'}</span>
                      <span>
                        <a href={chapterUrl} target="_blank" rel="noopener noreferrer">
                          Chapitre {lastChapterNumber}
                        </a>
                      </span>
                    </li>
                  )
                })}
              </ul>
            )}

            {totalPages > 1 && (
              <div className="series-pagination">
                <button disabled={seriesPage === 1} onClick={() => setSeriesPage(seriesPage - 1)}>Précédent</button>
                <span>Page {seriesPage} / {totalPages}</span>
                <button disabled={seriesPage === totalPages} onClick={() => setSeriesPage(seriesPage + 1)}>Suivant</button>
              </div>
            )}
          </section>
        </div>
      )}
    </div>
  )
}