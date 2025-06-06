"use client";

import { createContext, ReactNode, useContext, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import { useMessageHandler } from '../../hooks/useMessageHandler';
import { useAuth } from './AuthContext';

interface SocketContextType {
  sendJsonMessage: ReturnType<typeof useWebSocket>['sendJsonMessage'];
  readyState: ReadyState;
}

const SocketContext = createContext<SocketContextType | undefined>(undefined);

const socketUrl = 'ws://localhost:8080/ws';

export function SocketProvider({ children }: { children: ReactNode }) {
  const { handleMessage } = useMessageHandler();
  const { user } = useAuth();

  const {
    sendJsonMessage,
    readyState,
  } = useWebSocket(socketUrl, {
    onMessage: handleMessage,
  }, !!user);

  useEffect(() => {
    if (readyState === ReadyState.OPEN) {
      sendJsonMessage({
        type: "auth",
        data: {
          userId: user?.id,
        },
      });
    }
  }, [readyState, sendJsonMessage]);

  return <SocketContext.Provider value={{ sendJsonMessage, readyState }}>{children}</SocketContext.Provider>;
}

export const useSocket = () => {
  const context = useContext(SocketContext);
  if (context === undefined) {
    throw new Error('useSocket must be used within a SocketProvider');
  }
  return context;
};