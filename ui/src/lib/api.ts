import { Habit, TrackingEntry, CreateHabitRequest, CreateTrackingRequest, HabitStats, ProgressPoint, OverallStats, HabitCompletionRate, DailyCompletion } from '@/types';

const API_BASE_URL = 'http://localhost:8080';

class ApiError extends Error {
  constructor(message: string, public status: number) {
    super(message);
    this.name = 'ApiError';
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    throw new ApiError(`HTTP error! status: ${response.status}`, response.status);
  }
  return response.json();
}

export const api = {
  
  async getHabits(): Promise<Habit[]> {
    const response = await fetch(`${API_BASE_URL}/habits`);
    return handleResponse<Habit[]>(response);
  },

  async getHabit(id: string): Promise<Habit> {
    const response = await fetch(`${API_BASE_URL}/habits/${id}`);
    return handleResponse<Habit>(response);
  },

  async createHabit(habit: CreateHabitRequest): Promise<Habit> {
    const response = await fetch(`${API_BASE_URL}/habits`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(habit),
    });
    return handleResponse<Habit>(response);
  },

  async updateHabit(id: string, habit: Partial<Omit<Habit, 'id'>>): Promise<Habit> {
    const response = await fetch(`${API_BASE_URL}/habits/${id}`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(habit),
    });
    return handleResponse<Habit>(response);
  },

  async deleteHabit(id: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/habits/${id}`, {
      method: 'DELETE',
    });
    console.log(response);
    if (!response.ok) {
      throw new ApiError(`HTTP error! status: ${response.status}`, response.status);
    }
  },

  async getTrackingEntries(habitId: string): Promise<TrackingEntry[]> {
    const response = await fetch(`${API_BASE_URL}/habits/${habitId}/tracking`);
    return handleResponse<TrackingEntry[]>(response);
  },

  async createTrackingEntry(habitId: string, entry: CreateTrackingRequest): Promise<TrackingEntry> {
    const response = await fetch(`${API_BASE_URL}/habits/${habitId}/tracking`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(entry),
    });
    return handleResponse<TrackingEntry>(response);
  },

  
  async updateReminder(habitId: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/reminders/${habitId}`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: habitId + "-reminder",
        habitId: habitId,
        lastReminder: new Date().toISOString(),
      }),
    });
    if (!response.ok) {
      throw new ApiError(`HTTP error! status: ${response.status}`, response.status);
    }
  },

  // Statistics endpoints
  async getHabitStats(habitId: string): Promise<HabitStats> {
    const response = await fetch(`${API_BASE_URL}/habits/${habitId}/stats`);
    return handleResponse<HabitStats>(response);
  },

  async getHabitProgress(habitId: string, days: number = 30): Promise<ProgressPoint[]> {
    const response = await fetch(`${API_BASE_URL}/habits/${habitId}/progress?days=${days}`);
    return handleResponse<ProgressPoint[]>(response);
  },

  async getOverallStats(): Promise<OverallStats> {
    const response = await fetch(`${API_BASE_URL}/stats/overview`);
    return handleResponse<OverallStats>(response);
  },

  async getHabitCompletionRates(days: number = 30): Promise<HabitCompletionRate[]> {
    const response = await fetch(`${API_BASE_URL}/stats/completion-rates?days=${days}`);
    return handleResponse<HabitCompletionRate[]>(response);
  },

  async getDailyCompletions(days: number = 30): Promise<DailyCompletion[]> {
    const response = await fetch(`${API_BASE_URL}/stats/daily-completions?days=${days}`);
    return handleResponse<DailyCompletion[]>(response);
  },
}; 