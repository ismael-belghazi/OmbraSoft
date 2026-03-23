import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import '../styles/css.css'

export default function NavBar() {
  const { user, logout, isAuthenticated } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const goToSettings = () => {
    navigate('/notifications')
  }

  return (
    <nav className="navbar">
      <div className="navbar-container">
        <Link to="/" className="navbar-logo">
          OmbraSoft
        </Link>
        <ul className="nav-menu">
          {isAuthenticated ? (
            <>
              <li className="nav-item">
                <Link to="/dashboard" className="nav-link">Accueil</Link>
              </li>
              <li className="nav-item">
                <Link to="/bookmarks" className="nav-link">Mes Favoris</Link>
              </li>
              <li className="nav-item">
                <span className="nav-user">{user?.email}</span>
              </li>
              <li className="nav-item">
                <button onClick={handleLogout} className="nav-link logout-btn">
                  Déconnexion
                </button>
              </li>
              <li className="nav-item">
                <button className="settings-button" onClick={goToSettings}>
                  ⚙️ Paramètres de notifications
                </button>
              </li>
            </>
          ) : (
            <>
              <li className="nav-item">
                <Link to="/login" className="nav-link">Connexion</Link>
              </li>
              <li className="nav-item">
                <Link to="/register" className="nav-link">Inscription</Link>
              </li>
            </>
          )}
        </ul>
      </div>
    </nav>
  )
}
