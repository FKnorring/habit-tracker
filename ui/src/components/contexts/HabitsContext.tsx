"use client";

import React, { createContext, useContext, useState, ReactNode, useCallback } from 'react';
import { Habit, TrackingEntry } from '@/types';
import { api } from '@/lib/api';
import { toast } from 'sonner';

export interface EnrichedHabit extends Habit {
  trackingEntries?: TrackingEntry[] | null;
}

interface HabitsContextType {
  habits: EnrichedHabit[] | null;
  loading: boolean;
  error: string | null;
  initialized: boolean;
  reminders: Set<string>;
  fetchHabits: () => Promise<void>;
  createHabit: (habit: Parameters<typeof api.createHabit>[0]) => Promise<Habit>;
  updateHabit: (id: string, habit: Parameters<typeof api.updateHabit>[1]) => Promise<void>;
  deleteHabit: (id: string) => Promise<void>;
  enrichHabitWithTracking: (habitId: string) => Promise<void>;
  addTrackingEntry: (habitId: string, entry: Parameters<typeof api.createTrackingEntry>[1]) => Promise<void>;
  addReminder: (habitId: string) => void;
  removeReminder: (habitId: string) => void;
  clearAllReminders: () => void;
}

const HabitsContext = createContext<HabitsContextType | undefined>(undefined);

interface HabitsProviderProps {
  children: ReactNode;
}

export function HabitsProvider({ children }: HabitsProviderProps) {
  const [habits, setHabits] = useState<EnrichedHabit[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [initialized, setInitialized] = useState(false);
  const [reminders, setReminders] = useState<Set<string>>(new Set());

  const fetchHabits = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const fetchedHabits = await api.getHabits();
      setHabits(fetchedHabits);
      setInitialized(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch habits");
      toast.error("Failed to load habits");
    } finally {
      setLoading(false);
    }
  }, []);

  const createHabit = useCallback(async (habitData: Parameters<typeof api.createHabit>[0]) => {
    try {
      const newHabit = await api.createHabit(habitData);
      setHabits((prev) => [...(prev || []), newHabit]);
      toast.success("Habit created successfully");
      return newHabit;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to create habit";
      toast.error(errorMessage);
      throw err;
    }
  }, []);

  const updateHabit = useCallback(async (id: string, habitData: Parameters<typeof api.updateHabit>[1]) => {
    try {
      const updatedHabit = await api.updateHabit(id, habitData);
      setHabits((prev) => 
        prev?.map((habit) => 
          habit.id === id ? { ...habit, ...updatedHabit } : habit
        ) || []
      );
      toast.success("Habit updated successfully");
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to update habit";
      toast.error(errorMessage);
      throw err;
    }
  }, []);

  const deleteHabit = useCallback(async (id: string) => {
    try {
      await api.deleteHabit(id);
      setHabits((prev) => prev?.filter((habit) => habit.id !== id) || []);
      // Remove any reminder for deleted habit
      setReminders((prev) => {
        const newSet = new Set(prev);
        newSet.delete(id);
        return newSet;
      });
      toast.success("Habit deleted successfully");
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to delete habit";
      toast.error(errorMessage);
      throw err;
    }
  }, []);

  const enrichHabitWithTracking = useCallback(async (habitId: string) => {
    try {
      const trackingEntries = await api.getTrackingEntries(habitId);
      setHabits((prev) => 
        prev?.map((habit) => 
          habit.id === habitId 
            ? { ...habit, trackingEntries } 
            : habit
        ) || []
      );
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to load tracking data";
      toast.error(errorMessage);
      throw err;
    }
  }, []);

  const addTrackingEntry = useCallback(async (habitId: string, entryData: Parameters<typeof api.createTrackingEntry>[1]) => {
    try {
      const newEntry = await api.createTrackingEntry(habitId, entryData);
      setHabits((prev) => 
        prev?.map((habit) => {
          if (habit.id === habitId) {
            const updatedEntries = habit.trackingEntries 
              ? [...habit.trackingEntries, newEntry]
              : [newEntry];
            return { ...habit, trackingEntries: updatedEntries };
          }
          return habit;
        }) || []
      );
      // Remove reminder when user tracks the habit
      setReminders((prev) => {
        const newSet = new Set(prev);
        newSet.delete(habitId);
        return newSet;
      });
      toast.success("Progress tracked successfully");
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to track progress";
      toast.error(errorMessage);
      throw err;
    }
  }, []);

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

  const value: HabitsContextType = {
    habits,
    loading,
    error,
    initialized,
    reminders,
    fetchHabits,
    createHabit,
    updateHabit,
    deleteHabit,
    enrichHabitWithTracking,
    addTrackingEntry,
    addReminder,
    removeReminder,
    clearAllReminders,
  };

  return (
    <HabitsContext.Provider value={value}>
      {children}
    </HabitsContext.Provider>
  );
}

export function useHabits() {
  const context = useContext(HabitsContext);
  if (context === undefined) {
    throw new Error('useHabits must be used within a HabitsProvider');
  }
  return context;
} 