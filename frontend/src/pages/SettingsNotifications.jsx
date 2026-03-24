import { useState, useEffect } from 'react';
import api from '../services/api';
import '../styles/NotificationSettings.css';

export default function NotificationSettings() {
  const [push, setPush] = useState(true);
  const [discordWebhook, setDiscordWebhook] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const BOT_ID = 'VOTRE_BOT_ID_ICI';
  const BOT_INVITE_LINK = `https://discord.com/oauth2/authorize?client_id=${BOT_ID}&scope=bot&permissions=8`;

  useEffect(() => {
    const fetchSettings = async () => {
      try {
        setLoading(true);
        const res = await api.get('/user/notifications');
        setPush(res.data.push ?? true);
        setDiscordWebhook(res.data.discord_id || '');
      } catch (err) {
        console.error(err);
        setError('Impossible de récupérer les préférences.');
      } finally {
        setLoading(false);
      }
    };

    fetchSettings();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setMessage('');
    setError('');

    try {
      await api.post('/user/notifications', { push, discord_id: discordWebhook });
      setMessage('Préférences sauvegardées !');
    } catch (err) {
      setError(err.response?.data?.error || 'Impossible de sauvegarder les préférences.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="notification-form">
      <h2>Préférences notifications</h2>

      {error && <div className="error-message">{error}</div>}
      {message && <div className="success-message">{message}</div>}

      <label>
        <input
          type="checkbox"
          checked={push}
          onChange={(e) => setPush(e.target.checked)}
          disabled={loading}
        />
        Notifications push (Discord via Apprise)
      </label>

      {push && (
        <label>
          Webhook Discord :
          <input
            type="text"
            value={discordWebhook}
            onChange={(e) => setDiscordWebhook(e.target.value)}
            placeholder="Entrez votre webhook Discord"
            disabled={loading}
          />
        </label>
      )}

      <button type="submit" disabled={loading}>
        {loading ? 'Sauvegarde...' : 'Sauvegarder'}
      </button>

      <div className="dev-section">
        <span>En cours de développement :</span> certaines fonctionnalités avancées de notification ne sont pas encore disponibles.
      </div>

      <h2>Ajouter le bot Discord</h2>
      <p>Utilisez ce lien pour ajouter le bot à votre serveur avec les permissions prédéfinies :</p>
      <a href={BOT_INVITE_LINK} target="_blank" rel="noopener noreferrer">
        Ajouter le bot à votre serveur
      </a>
    </form>
  );
}