"use client";

import React, { createContext, useContext, useState, ReactNode, useCallback } from 'react';
import { HabitStats, ProgressPoint, OverallStats, HabitCompletionRate, DailyCompletion } from '@/types';
import { api } from '@/lib/api';
import { toast } from 'sonner';

interface StatisticsContextType {
  overallStats: OverallStats | null;
  habitCompletionRates: HabitCompletionRate[] | null;
  dailyCompletions: DailyCompletion[] | null;
  statisticsLoading: boolean;
  statisticsError: string | null;
  
  fetchOverallStats: () => Promise<void>;
  fetchHabitCompletionRates: (days?: number) => Promise<void>;
  fetchDailyCompletions: (days?: number) => Promise<void>;
  fetchHabitStats: (habitId: string) => Promise<HabitStats>;
  fetchHabitProgress: (habitId: string, days?: number) => Promise<ProgressPoint[]>;
  fetchAllStatistics: (days?: number) => Promise<void>;
}

const StatisticsContext = createContext<StatisticsContextType | undefined>(undefined);

interface StatisticsProviderProps {
  children: ReactNode;
}

export function StatisticsProvider({ children }: StatisticsProviderProps) {
  const [overallStats, setOverallStats] = useState<OverallStats | null>(null);
  const [habitCompletionRates, setHabitCompletionRates] = useState<HabitCompletionRate[] | null>(null);
  const [dailyCompletions, setDailyCompletions] = useState<DailyCompletion[] | null>(null);
  const [statisticsLoading, setStatisticsLoading] = useState(false);
  const [statisticsError, setStatisticsError] = useState<string | null>(null);

  const fetchOverallStats = useCallback(async () => {
    try {
      setStatisticsLoading(true);
      setStatisticsError(null);
      const stats = await api.getOverallStats();
      setOverallStats(stats);
    } catch (err) {
      setStatisticsError(err instanceof Error ? err.message : "Failed to fetch overall stats");
      toast.error("Failed to load overall stats");
    } finally {
      setStatisticsLoading(false);
    }
  }, []);

  const fetchHabitCompletionRates = useCallback(async (days?: number) => {
    try {
      setStatisticsLoading(true);
      setStatisticsError(null);
      const rates = await api.getHabitCompletionRates(days);
      setHabitCompletionRates(rates);
    } catch (err) {
      setStatisticsError(err instanceof Error ? err.message : "Failed to fetch habit completion rates");
      toast.error("Failed to load habit completion rates");
    } finally {
      setStatisticsLoading(false);
    }
  }, []);

  const fetchDailyCompletions = useCallback(async (days?: number) => {
    try {
      setStatisticsLoading(true);
      setStatisticsError(null);
      const completions = await api.getDailyCompletions(days);
      setDailyCompletions(completions);
    } catch (err) {
      setStatisticsError(err instanceof Error ? err.message : "Failed to fetch daily completions");
      toast.error("Failed to load daily completions");
    } finally {
      setStatisticsLoading(false);
    }
  }, []);

  const fetchHabitStats = useCallback(async (habitId: string) => {
    try {
      setStatisticsLoading(true);
      setStatisticsError(null);
      const stats = await api.getHabitStats(habitId);
      return stats;
    } catch (err) {
      setStatisticsError(err instanceof Error ? err.message : "Failed to fetch habit stats");
      toast.error("Failed to load habit stats");
      throw err;
    } finally {
      setStatisticsLoading(false);
    }
  }, []);

  const fetchHabitProgress = useCallback(async (habitId: string, days?: number) => {
    try {
      setStatisticsLoading(true);
      setStatisticsError(null);
      const progress = await api.getHabitProgress(habitId, days);
      return progress;
    } catch (err) {
      setStatisticsError(err instanceof Error ? err.message : "Failed to fetch habit progress");
      toast.error("Failed to load habit progress");
      throw err;
    } finally {
      setStatisticsLoading(false);
    }
  }, []);

  const fetchAllStatistics = useCallback(async (days?: number) => {
    try {
      setStatisticsLoading(true);
      setStatisticsError(null);
      await Promise.all([
        fetchOverallStats(),
        fetchHabitCompletionRates(days),
        fetchDailyCompletions(days),
      ]);
    } catch (err) {
      setStatisticsError(err instanceof Error ? err.message : "Failed to fetch statistics");
      toast.error("Failed to load statistics");
    } finally {
      setStatisticsLoading(false);
    }
  }, [fetchOverallStats, fetchHabitCompletionRates, fetchDailyCompletions]);

  const value: StatisticsContextType = {
    overallStats,
    habitCompletionRates,
    dailyCompletions,
    statisticsLoading,
    statisticsError,
    fetchOverallStats,
    fetchHabitCompletionRates,
    fetchDailyCompletions,
    fetchHabitStats,
    fetchHabitProgress,
    fetchAllStatistics,
  };

  return (
    <StatisticsContext.Provider value={value}>
      {children}
    </StatisticsContext.Provider>
  );
}

export function useStatistics() {
  const context = useContext(StatisticsContext);
  if (context === undefined) {
    throw new Error('useStatistics must be used within a StatisticsProvider');
  }
  return context;
} 