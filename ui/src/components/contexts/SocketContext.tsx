"use client";

import { createContext, ReactNode, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import { toast } from 'sonner';

interface SocketContextType {
  sendJsonMessage: ReturnType<typeof useWebSocket>['sendJsonMessage'];
  readyState: ReadyState;
}

interface ReminderMessage {
  type: "reminder";
  data: {
    habitId: string;
    habitName: string;
    description: string;
    frequency: string;
    timestamp: string;
  };
}

const SocketContext = createContext<SocketContextType | undefined>(undefined);

const socketUrl = 'ws://localhost:8080/ws';

export function SocketProvider({ children }: { children: ReactNode }) {

  const {
    sendJsonMessage,
    readyState,
  } = useWebSocket(socketUrl, {
    onMessage: (event) => {
      console.log('message', event.data);
      
      try {
        const message = JSON.parse(event.data);
        
        if (message.type === "reminder") {
          const reminderMessage = message as ReminderMessage;
          const { habitName, frequency } = reminderMessage.data;
          
          toast(habitName, {
            description: `It's time to track this habit - should be tracked ${frequency}`,
          });
        }
      } catch (error) {
        console.error('Error parsing websocket message:', error);
      }
    },
  });

  useEffect(() => {
    if (readyState === ReadyState.OPEN) {
      sendJsonMessage({
        type: "auth",
        data: {
          userId: "user-123",
        },
      });
    }
  }, [readyState, sendJsonMessage]);

  return <SocketContext.Provider value={{ sendJsonMessage, readyState }}>{children}</SocketContext.Provider>;
}