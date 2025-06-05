export const requestNotificationPermission = async (): Promise<boolean> => {
  if ('Notification' in window) {
    if (Notification.permission === 'default') {
      const permission = await Notification.requestPermission();
      return permission === 'granted';
    }
    return Notification.permission === 'granted';
  }
  return false;
};

export const showBrowserNotification = (habitName: string, frequency: string): void => {
  if ('Notification' in window && Notification.permission === 'granted') {
    new Notification(habitName, {
      body: `It's time to track this habit - should be tracked ${frequency}`,
      icon: '/favicon.ico', 
      tag: 'habit-reminder', 
    });
  }
};

export const isNotificationSupported = (): boolean => {
  return 'Notification' in window;
};

export const getNotificationPermission = (): NotificationPermission | null => {
  return 'Notification' in window ? Notification.permission : null;
}; 