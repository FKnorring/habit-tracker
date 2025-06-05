"use client";

import { createContext, ReactNode, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

interface SocketContextType {
  sendJsonMessage: ReturnType<typeof useWebSocket>['sendJsonMessage'];
  readyState: ReadyState;
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