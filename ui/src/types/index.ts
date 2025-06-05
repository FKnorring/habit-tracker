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
  frequency: Frequency;
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
  frequency: Frequency;
  startDate: string;
}

export interface CreateTrackingRequest {
  note?: string;
} 