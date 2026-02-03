import os
import requests
from apprise import Apprise, NotifyFormat

# Variables d'environnement
DISCORD_BOT_TOKEN = os.getenv('DISCORD_BOT_TOKEN')
PUSHBULLET_URL = os.getenv('PUSHBULLET_URL') 
TELEGRAM_URL = os.getenv('TELEGRAM_URL')     
BACKEND_URL = os.getenv('BACKEND_URL', 'http://localhost:8080')

def send_notification(user_id: str, title: str, message: str):
    """Envoie une notification à l'utilisateur selon ses préférences"""
    try:
        res = requests.get(f'{BACKEND_URL}/user/{user_id}/notifications')
        res.raise_for_status()
        prefs = res.json()
    except Exception as e:
        print(f"[Erreur] Impossible de récupérer les préférences pour {user_id} : {e}")
        return

    a = Apprise()
    added_service = False

    if prefs.get('push', False):
        discord_webhook = prefs.get('discord_id')
        if discord_webhook and DISCORD_BOT_TOKEN:
            a.add(f'discord://{DISCORD_BOT_TOKEN}@{discord_webhook}')
            added_service = True
        else:
            print(f"[Info] Pas de webhook Discord configuré pour {user_id}")

        if PUSHBULLET_URL:
            a.add(PUSHBULLET_URL)
            added_service = True

        if TELEGRAM_URL:
            a.add(TELEGRAM_URL)
            added_service = True

    if not added_service:
        print(f"[Info] Aucun service configuré pour l'utilisateur {user_id}")
        return

    try:
        success = a.notify(
            body=message,
            title=title,
            notify_format=NotifyFormat.TEXT
        )
        if success:
            print(f"[Succès] Notification envoyée à {user_id}")
        else:
            print(f"[Erreur] Échec de l'envoi pour {user_id}")
    except Exception as e:
        print(f"[Exception] Erreur lors de l'envoi : {e}")


if __name__ == "__main__":
    chapter = 5
    send_notification(
        user_id='user-uuid',
        title='Nouveau chapitre !',
        message=f'Le chapitre {chapter} est sorti.'
    )
