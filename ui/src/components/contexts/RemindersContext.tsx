"use client";

import React, { createContext, useContext, useState, ReactNode, useCallback } from 'react';

interface RemindersContextType {
  reminders: Set<string>;
  addReminder: (habitId: string) => void;
  removeReminder: (habitId: string) => void;
  clearAllReminders: () => void;
}

const RemindersContext = createContext<RemindersContextType | undefined>(undefined);

interface RemindersProviderProps {
  children: ReactNode;
}

export function RemindersProvider({ children }: RemindersProviderProps) {
  const [reminders, setReminders] = useState<Set<string>>(new Set());

  const addReminder = useCallback((habitId: string) => {
    setReminders((prev) => new Set([...prev, habitId]));
  }, []);

  const removeReminder = useCallback((habitId: string) => {
    setReminders((prev) => {
      const newSet = new Set(prev);
      newSet.delete(habitId);
      return newSet;
    });
  }, []);

  const clearAllReminders = useCallback(() => {
    setReminders(new Set());
  }, []);

  const value: RemindersContextType = {
    reminders,
    addReminder,
    removeReminder,
    clearAllReminders,
  };

  return (
    <RemindersContext.Provider value={value}>
      {children}
    </RemindersContext.Provider>
  );
}

export function useReminders() {
  const context = useContext(RemindersContext);
  if (context === undefined) {
    throw new Error('useReminders must be used within a RemindersProvider');
  }
  return context;
} 