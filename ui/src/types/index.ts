export enum Frequency {
  HOURLY = "hourly",
  DAILY = "daily", 
  WEEKLY = "weekly",
  BIWEEKLY = "biweekly",
  MONTHLY = "monthly",
  QUARTERLY = "quarterly",
  YEARLY = "yearly"
}

export interface Habit {
  id: string;
  name: string;
  description: string;
  frequency: 'hourly' | 'daily' | 'weekly' | 'biweekly' | 'monthly' | 'quarterly' | 'yearly';
  startDate: string;
}

export interface TrackingEntry {
  id: string;
  habitId: string;
  timestamp: string;
  note: string;
}

export interface CreateHabitRequest {
  name: string;
  description: string;
  frequency: Habit['frequency'];
  startDate: string;
}

export interface CreateTrackingRequest {
  note?: string;
  timestamp?: string;
}

// WebSocket message types
export interface ReminderMessage {
  type: "reminder";
  data: {
    habitId: string;
    habitName: string;
    description: string;
    frequency: string;
    timestamp: string;
  };
}

export interface AuthMessage {
  type: "auth";
  data: {
    userId: string;
  };
}

// Statistics Types
export interface HabitStats {
  habitId: string;
  habitName: string;
  frequency: Habit['frequency'];
  startDate: string;
  totalEntries: number;
  currentStreak: number;
  longestStreak: number;
  completionRate: number;
  lastCompleted: string;
}

export interface ProgressPoint {
  date: string;
  count: number;
}

export interface OverallStats {
  totalHabits: number;
  totalEntries: number;
  entriesToday: number;
  entriesThisWeek: number;
  avgEntriesPerDay: number;
}

export interface HabitCompletionRate {
  habitId: string;
  habitName: string;
  frequency: Habit['frequency'];
  startDate: string;
  actualCompletions: number;
  expectedCompletions: number;
  completionRate: number;
}

export interface DailyCompletion {
  date: string;
  completions: number;
} 