export interface Habit {
  id: string;
  name: string;
  description: string;
  frequency: string;
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
  frequency: string;
  startDate: string;
}

export interface CreateTrackingRequest {
  note?: string;
} 