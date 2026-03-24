import React, { useState, useEffect, useCallback, useRef } from 'react';
import ReactDOM from 'react-dom';
import { bookmarkService } from '../services/bookmarkService';
import { seriesService } from '../services/seriesService';
import '../styles/Bookmarks.css';

const DEFAULT_COVER = 'http://localhost:8080/covers/default-cover.jpg';
const SERIES_PER_PAGE = 10; 

const getCoverUrl = (series) => {
  if (!series) return DEFAULT_COVER;
  const coverPath = series.cover || series.cover_image_url || DEFAULT_COVER;
  return coverPath.startsWith('http')
    ? coverPath
    : `http://localhost:8080${coverPath.startsWith('/') ? '' : '/'}${coverPath}`;
};

const setFavicon = () => {
  let link = document.querySelector("link[rel*='icon']");
  if (!link) {
    link = document.createElement('link');
    link.rel = 'icon';
    document.head.appendChild(link);
  }
  link.href = DEFAULT_COVER;
};

const useBookmarks = () => {
  const [bookmarks, setBookmarks] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const isFetching = useRef(false);

  const fetchBookmarks = useCallback(async () => {
    if (isFetching.current) return;
    try {
      isFetching.current = true;
      setLoading(true);
      const response = await bookmarkService.getAll();
      const raw = response?.data || response || [];
      const data = raw.map(b => ({
        ...b,
        lastReadChapter: Number(b.lastReadChapter ?? b.last_read_chapter ?? 0),
        series: {
          ...b.series,
          lastChapterNumber: Number(b.series?.lastChapterNumber ?? b.series?.last_chapter_number ?? 0),
          cover: b.series?.cover || b.series?.cover_image_url || DEFAULT_COVER,
        },
      }));
      setBookmarks(data);
      setError('');
      return data;
    } catch {
      setError('Erreur lors du chargement des favoris');
      return [];
    } finally {
      setLoading(false);
      isFetching.current = false;
    }
  }, []);

  const deleteBookmark = async id => {
    try {
      await bookmarkService.delete(id);
      setBookmarks(prev => prev.filter(b => b.id !== id));
    } catch {
      setError('Erreur lors de la suppression');
    }
  };

  const updateBookmark = (bookmarkId, lastReadChapter) => {
    setBookmarks(prev =>
      prev.map(b => b.id === bookmarkId ? { ...b, lastReadChapter: Number(lastReadChapter) } : b)
    );
  };

  return { bookmarks, loading, error, setError, deleteBookmark, fetchBookmarks, updateBookmark };
};

export default function Bookmarks() {
  const [chapters, setChapters] = useState({});
  const [loadingChapters, setLoadingChapters] = useState({});
  const [activeOverlay, setActiveOverlay] = useState(null);
  const [newLink, setNewLink] = useState('');
  const [seriesPage, setSeriesPage] = useState(1);

  const { bookmarks, loading, error, setError, deleteBookmark, fetchBookmarks, updateBookmark } = useBookmarks();

  useEffect(() => { fetchBookmarks(); }, [fetchBookmarks]);
  useEffect(() => { setFavicon(); }, [bookmarks]);

  const loadChapters = async (bookmarkId, seriesId) => {
    try {
      setLoadingChapters(prev => ({ ...prev, [bookmarkId]: true }));
      const res = await seriesService.getChaptersBySeriesId(seriesId);
      const data = res?.data || res?.chapters || [];
      data.sort((a, b) => Number(a.number) - Number(b.number));
      setChapters(prev => ({ ...prev, [bookmarkId]: data }));
    } catch {
      setError('Erreur chargement chapitres');
    } finally {
      setLoadingChapters(prev => ({ ...prev, [bookmarkId]: false }));
    }
  };

  const openOverlay = async (bookmarkId, seriesId) => {
    setActiveOverlay(bookmarkId);
    if (!chapters[bookmarkId]) await loadChapters(bookmarkId, seriesId);
  };

  const closeOverlay = () => setActiveOverlay(null);

  const markChapterAsRead = async (bookmarkId, chapterNumber) => {
    try {
      const res = await bookmarkService.markChapterAsRead(bookmarkId, chapterNumber);
      const last = res?.data?.lastReadChapter ?? chapterNumber;
      updateBookmark(bookmarkId, last);
    } catch {
      setError('Erreur mise à jour chapitre');
    }
  };

  const markSeriesAsRead = async bookmarkId => {
    try {
      const res = await bookmarkService.markSeriesAsRead(bookmarkId);
      const last = res?.data?.lastReadChapter ?? 0;
      updateBookmark(bookmarkId, last);
    } catch {
      setError('Erreur mise à jour série');
    }
  };

  const addBookmark = async () => {
    const trimmedLink = newLink.trim();
    if (!trimmedLink) { setError('Le lien de la série est vide'); return; }
    try { await bookmarkService.create(trimmedLink); await fetchBookmarks(); setNewLink(''); }
    catch (err) { setError(err.message || 'Erreur lors de la création du favori'); }
  };

  const activeBookmark = activeOverlay !== null ? bookmarks.find(b => b.id === activeOverlay) : null;

  const paginatedBookmarks = bookmarks.slice(
    (seriesPage - 1) * SERIES_PER_PAGE,
    seriesPage * SERIES_PER_PAGE
  );

  const totalSeriesPages = Math.ceil(bookmarks.length / SERIES_PER_PAGE);

  return (
    <div className="page-container">
      <h1>Mes Favoris</h1>

      <div className="add-bookmark">
        <input type="text" placeholder="Coller le lien de la série..." value={newLink} onChange={e => setNewLink(e.target.value)} />
        <button onClick={addBookmark}>Ajouter</button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {loading ? <p>Chargement...</p> : (
        <>
          <div className="bookmarks-gallery">
            {paginatedBookmarks.map(bookmark => {
              const total = bookmark.series?.lastChapterNumber || 1;
              const read = bookmark.lastReadChapter || 0;
              const progress = Math.min((read / total) * 100, 100);
              const hasNew = total > 0 && read < total;

              return (
                <div key={bookmark.id} className="bookmark-card">
                  <div className="cover-wrapper">
                    <img src={getCoverUrl(bookmark.series)} alt={bookmark.series?.title} />
                    {hasNew && <span className="badge-new">NEW</span>}
                  </div>

                  <div className="card-content">
                    <h2>{bookmark.series?.title}</h2>

                    <div className="progress-bar">
                      <div className="progress-fill" style={{ width: `${progress}%` }} />
                    </div>

                    <div className="actions">
                      <button onClick={() => openOverlay(bookmark.id, bookmark.series?.id)}>Voir chapitres</button>
                      <button onClick={() => markSeriesAsRead(bookmark.id)}>Tout lire</button>
                      <button onClick={() => deleteBookmark(bookmark.id)}>Supprimer</button>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>

          {totalSeriesPages > 1 && (
            <div className="series-pagination">
              <button disabled={seriesPage === 1} onClick={() => setSeriesPage(seriesPage - 1)}>Précédent</button>
              <span>Page {seriesPage} / {totalSeriesPages}</span>
              <button disabled={seriesPage === totalSeriesPages} onClick={() => setSeriesPage(seriesPage + 1)}>Suivant</button>
            </div>
          )}
        </>
      )}

      {activeBookmark && ReactDOM.createPortal(
        <div className="overlay" onClick={closeOverlay}>
          <div className="overlay-content" onClick={e => e.stopPropagation()}>
            <button className="close-btn" onClick={closeOverlay}>✕</button>

            <div className="overlay-left">
              <img src={getCoverUrl(activeBookmark.series)} alt="" />
              <h2>{activeBookmark.series?.title}</h2>

              <div className="progress-bar">
                <div className="progress-fill" style={{ width: `${(activeBookmark.lastReadChapter / (activeBookmark.series?.lastChapterNumber || 1)) * 100}%` }} />
              </div>

              <button onClick={() => markSeriesAsRead(activeBookmark.id)}>Tout lire</button>
            </div>

            <div className="overlay-right">
              {loadingChapters[activeOverlay] ? <p>Chargement...</p> : (
                <ul className="chapters-scroll">
                  {(chapters[activeOverlay] || []).map(ch => {
                    const isRead = Number(ch.number) <= activeBookmark.lastReadChapter;
                    return (
                      <li key={ch.id} className={isRead ? 'chapter-read' : ''}>
                        <a href={ch.url} target="_blank" rel="noreferrer">Chapitre {ch.number}</a>
                        {!isRead && <button onClick={() => markChapterAsRead(activeOverlay, ch.number)}>✔</button>}
                      </li>
                    );
                  })}
                </ul>
              )}
            </div>
          </div>
        </div>,
        document.body
      )}
    </div>
  );
}