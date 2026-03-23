import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../services/api'
import '../styles/css.css'

export default function ForgotPassword() {
  const [email, setEmail] = useState('')
  const [secretPhrase, setSecretPhrase] = useState('')
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setMessage('')
    setLoading(true)

    try {
      const res = await api.post('/auth/forgot-password', { email, secret_phrase: secretPhrase })

      if (res.data.message.includes('validée')) {
        setMessage('Phrase secrète validée ! Vous pouvez maintenant réinitialiser votre mot de passe.')
        setTimeout(() => navigate(`/reset-password?email=${encodeURIComponent(email)}`), 1500)
      } else {
        setMessage(res.data.message)
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Une erreur est survenue')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-container">
      <form onSubmit={handleSubmit} className="auth-form">
        <h2>Réinitialisation du mot de passe</h2>

        {message && <div className="success-message">{message}</div>}
        {error && <div className="error-message">{error}</div>}

        <input
          type="email"
          placeholder="Votre email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />

        <input
          type="text"
          placeholder="Phrase secrète"
          value={secretPhrase}
          onChange={(e) => setSecretPhrase(e.target.value)}
          required
        />

        <button type="submit" disabled={loading}>
          {loading ? 'Vérification...' : 'Valider'}
        </button>
      </form>
    </div>
  )
}
