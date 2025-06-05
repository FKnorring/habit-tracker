"use client";

import { useEffect, useState } from 'react';
import { useHabits } from '@/components/contexts/HabitsContext';
import { useStatistics } from '@/components/contexts/StatisticsContext';
import { Card, CardContent } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { StatCards } from './StatCards';
import { DailyCompletions } from './DailyCompletions';
import { HabitFrequency } from './HabitFrequency';
import { HabitCompletion } from './HabitCompletion';
import { IndividualStatistics } from './IndividualStatistics';

export function Statistics() {
  const { habits } = useHabits();
  const {
    overallStats,
    habitCompletionRates,
    dailyCompletions,
    statisticsLoading,
    statisticsError,
    fetchAllStatistics,
  } = useStatistics();

  const [timeRange, setTimeRange] = useState<number>(30);

  useEffect(() => {
    fetchAllStatistics(timeRange);
  }, [fetchAllStatistics, timeRange]);

  if (statisticsLoading) {
    return (
      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6 px-4 md:px-6">
        <h2 className="text-2xl font-bold">Statistics</h2>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {[...Array(4)].map((_, i) => (
            <Card key={i} className="animate-pulse">
              <div className="h-24 bg-gray-200 rounded"></div>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  if (statisticsError) {
    return (
      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6 px-4 md:px-6">
        <h2 className="text-2xl font-bold">Statistics</h2>
        <Card>
          <CardContent className="pt-6">
            <p className="text-destructive">Error loading statistics: {statisticsError}</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6 px-4 md:px-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <h2 className="text-2xl font-bold">Statistics</h2>
        <div className="flex items-center gap-2">
          <Select value={timeRange.toString()} onValueChange={(value) => setTimeRange(parseInt(value))}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Select time range" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="7">Last 7 days</SelectItem>
              <SelectItem value="30">Last 30 days</SelectItem>
              <SelectItem value="90">Last 90 days</SelectItem>
              <SelectItem value="365">Last year</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Overview Cards */}
      <StatCards overallStats={overallStats} timeRange={timeRange} />

      {/* Charts Section */}
      <div className="grid gap-4 md:grid-cols-2">
        <DailyCompletions dailyCompletions={dailyCompletions} timeRange={timeRange} />
        <HabitFrequency habits={habits} />
      </div>

      {/* Habit Completion Rates Bar Chart */}
      <HabitCompletion habitCompletionRates={habitCompletionRates} timeRange={timeRange} />

      {/* Individual Habit Statistics Section */}
      <IndividualStatistics habits={habits} timeRange={timeRange} />
    </div>
  );
} 