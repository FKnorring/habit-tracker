import { toast } from 'sonner';
import { requestNotificationPermission, showBrowserNotification, getNotificationPermission } from '../lib/notifications';
import { ReminderMessage } from '../types';

export const useMessageHandler = () => {
  const handleMessage = async (event: MessageEvent) => {
    console.log('message', event.data);
    
    try {
      const message = JSON.parse(event.data);
      
      if (message.type === "reminder") {
        await handleReminderMessage(message as ReminderMessage);
      }
    } catch (error) {
      console.error('Error parsing websocket message:', error);
    }
  };

  const handleReminderMessage = async (reminderMessage: ReminderMessage) => {
    const { habitName, frequency } = reminderMessage.data;
    const notificationPermission = getNotificationPermission();
    
    toast(habitName, {
      description: `It's time to track this habit - should be tracked ${frequency}`,
      duration: 10000,
      cancel: notificationPermission === 'default' ? 
        undefined : {
          label: 'Close',
          onClick: () => {
            toast.dismiss();
          },
        },
      action: notificationPermission === 'default' ? {
        label: 'Enable notifications',
        onClick: async () => {
          const granted = await requestNotificationPermission();
          if (granted) {
            toast.success('Browser notifications enabled!');
            showBrowserNotification(habitName, frequency);
          }
        },
      } : undefined,
    });

    // Show browser notification if permission is granted
    if (notificationPermission === 'granted') {
      showBrowserNotification(habitName, frequency);
    }
  };

  return { handleMessage };
}; 