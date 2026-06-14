import { userNotifications } from '@/services/notifications'

function waitForOneSignal(timeoutMs = 10_000) {
  return new Promise((resolve, reject) => {
    const started = Date.now()
    function check() {
      if (typeof window === 'undefined') return resolve(null)
      if (window.OneSignal && window.OneSignal.User) return resolve(window.OneSignal)
      if (Date.now() - started > timeoutMs) return reject(new Error('OneSignal SDK timeout'))
      setTimeout(check, 200)
    }
    check()
  })
}

export async function enablePushNotifications() {
  try {
    const OneSignal = await waitForOneSignal()
    if (!OneSignal) return null

    const permission = await OneSignal.Notifications.requestPermission()
    if (permission !== true && permission !== 'granted') return null

    let id = OneSignal.User.PushSubscription.id
    if (!id) {
      await new Promise((r) => setTimeout(r, 800))
      id = OneSignal.User.PushSubscription.id
    }
    if (!id) return null

    await userNotifications.registerPushToken(id)
    return id
  } catch (e) {
    console.warn('[onesignal] enablePushNotifications failed:', e)
    return null
  }
}

export async function syncExistingSubscription() {
  try {
    const OneSignal = await waitForOneSignal(4000)
    if (!OneSignal) return
    const optedIn = OneSignal.User.PushSubscription.optedIn
    if (!optedIn) return
    const id = OneSignal.User.PushSubscription.id
    if (id) {
      await userNotifications.registerPushToken(id)
    }
  } catch {
  }
}
