"use client";

import { useEffect, useState, useCallback } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ChartContainer, ChartTooltip } from '@/components/ui/chart';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { XAxis, YAxis, CartesianGrid, ResponsiveContainer, LineChart, Line } from 'recharts';
import { useStatistics } from '@/components/contexts/StatisticsContext';
import { HabitStats, ProgressPoint } from '@/types';
import { TrendingUp, Calendar, Target, Flame, RotateCcw, Activity } from 'lucide-react';

interface HabitStatsCardProps {
  habitId: string;
  days?: number;
}

export function HabitStatsCard({ habitId, days = 30 }: HabitStatsCardProps) {
  const { fetchHabitStats, fetchHabitProgress } = useStatistics();
  const [stats, setStats] = useState<HabitStats | null>(null);
  const [progress, setProgress] = useState<ProgressPoint[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [customDays, setCustomDays] = useState<number>(days);

  const loadData = useCallback(async (dayRange: number = customDays) => {
    try {
      setLoading(true);
      setError(null);
      const [habitStats, habitProgress] = await Promise.all([
        fetchHabitStats(habitId),
        fetchHabitProgress(habitId, dayRange)
      ]);
      setStats(habitStats);
      setProgress(habitProgress);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load habit statistics');
    } finally {
      setLoading(false);
    }
  }, [habitId, customDays, fetchHabitStats, fetchHabitProgress]);

  useEffect(() => {
    loadData(days);
  }, [habitId, days, loadData]);

  const handleDaysChange = (newDays: number) => {
    setCustomDays(newDays);
    loadData(newDays);
  };

  const handleRefresh = () => {
    loadData(customDays);
  };

  if (loading) {
    return (
      <Card className="animate-pulse">
        <CardHeader>
          <div className="flex justify-between items-start">
            <div className="space-y-2">
              <div className="h-5 bg-gray-200 rounded w-3/4"></div>
              <div className="h-3 bg-gray-200 rounded w-1/2"></div>
            </div>
            <div className="h-8 w-20 bg-gray-200 rounded"></div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {[...Array(4)].map((_, i) => (
                <div key={i} className="space-y-2">
                  <div className="h-6 bg-gray-200 rounded w-8"></div>
                  <div className="h-3 bg-gray-200 rounded w-12"></div>
                </div>
              ))}
            </div>
            <div className="h-40 bg-gray-200 rounded"></div>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error || !stats) {
    return (
      <Card>
        <CardHeader>
          <div className="flex justify-between items-start">
            <CardTitle>Error Loading Habit</CardTitle>
            <Button 
              variant="outline" 
              size="sm" 
              onClick={handleRefresh}
              className="flex items-center gap-1"
            >
              <RotateCcw className="h-3 w-3" />
              Retry
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <p className="text-destructive text-sm">{error || 'No data available'}</p>
        </CardContent>
      </Card>
    );
  }

  const chartConfig = {
    completions: {
      label: "Completions",
      color: "hsl(var(--chart-1))",
    },
  };

  const progressData = progress?.map(point => ({
    date: new Date(point.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
    completions: point.count,
  })) || [];

  const completionRate = Math.round(stats.completionRate * 100);
  const hasProgressData = progressData.length > 0;

  // Calculate some additional metrics
  const totalDaysTracked = progress?.length || 0;
  const averageCompletionsPerDay = totalDaysTracked > 0 ? 
    (progress?.reduce((sum, p) => sum + p.count, 0) || 0) / totalDaysTracked : 0;

  return (
    <Card className="relative overflow-hidden">
      <CardHeader>
        <div className="flex justify-between items-start">
          <div className="space-y-1">
            <CardTitle className="flex items-center gap-2">
              <Activity className="h-4 w-4" />
              {stats.habitName}
            </CardTitle>
            <CardDescription>
              Statistics for the last {customDays} days
            </CardDescription>
          </div>
          <div className="flex gap-2">
            <Select value={customDays.toString()} onValueChange={(value) => handleDaysChange(parseInt(value))}>
              <SelectTrigger className="w-[100px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="7">7d</SelectItem>
                <SelectItem value="30">30d</SelectItem>
                <SelectItem value="90">90d</SelectItem>
                <SelectItem value="365">1y</SelectItem>
              </SelectContent>
            </Select>
            <Button 
              variant="outline" 
              size="sm" 
              onClick={handleRefresh}
              className="flex items-center gap-1"
            >
              <RotateCcw className="h-3 w-3" />
            </Button>
          </div>
        </div>
      </CardHeader>
      
      <CardContent>
        <div className="space-y-6">
          {/* Key Metrics Grid */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="flex items-center space-x-2">
              <Target className="h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-lg font-semibold">{stats.totalEntries}</p>
                <p className="text-xs text-muted-foreground">Total</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-2">
              <Flame className="h-4 w-4 text-orange-500" />
              <div>
                <p className="text-lg font-semibold">{stats.currentStreak}</p>
                <p className="text-xs text-muted-foreground">Current Streak</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-2">
              <TrendingUp className="h-4 w-4 text-green-500" />
              <div>
                <p className="text-lg font-semibold">{stats.longestStreak}</p>
                <p className="text-xs text-muted-foreground">Best Streak</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-2">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <div>
                <Badge variant={completionRate >= 80 ? "default" : completionRate >= 60 ? "secondary" : "destructive"}>
                  {completionRate}%
                </Badge>
                <p className="text-xs text-muted-foreground">Success Rate</p>
              </div>
            </div>
          </div>

          {/* Secondary Metrics */}
          <div className="grid grid-cols-2 gap-4 p-3 bg-muted/20 rounded-lg">
            <div>
              <p className="text-sm font-medium">Days Tracked</p>
              <p className="text-xs text-muted-foreground">{totalDaysTracked} of {customDays} days</p>
            </div>
            <div>
              <p className="text-sm font-medium">Daily Average</p>
              <p className="text-xs text-muted-foreground">{averageCompletionsPerDay.toFixed(1)} per day</p>
            </div>
          </div>

          {/* Progress Chart */}
          {hasProgressData ? (
            <div>
              <div className="flex items-center justify-between mb-3">
                <h4 className="text-sm font-medium">Progress Trend</h4>
                <Badge variant="outline" className="text-xs">
                  {progressData.length} data points
                </Badge>
              </div>
              <ChartContainer config={chartConfig} className="h-[200px]">
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={progressData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis 
                      dataKey="date" 
                      fontSize={11}
                      tickLine={false}
                      axisLine={false}
                    />
                    <YAxis 
                      fontSize={11}
                      tickLine={false}
                      axisLine={false}
                    />
                    <ChartTooltip 
                      content={({ active, payload, label }) => {
                        if (active && payload && payload.length) {
                          return (
                            <div className="rounded-lg border bg-background p-2 shadow-sm">
                              <div className="grid gap-1">
                                <div className="font-medium text-xs">{label}</div>
                                <div className="text-sm">
                                  <span className="font-medium">{payload[0].value}</span>
                                  <span className="text-muted-foreground ml-1">
                                    completion{payload[0].value !== 1 ? 's' : ''}
                                  </span>
                                </div>
                              </div>
                            </div>
                          );
                        }
                        return null;
                      }}
                    />
                    <Line
                      type="monotone"
                      dataKey="completions"
                      stroke="hsl(var(--chart-1))"
                      strokeWidth={2}
                      dot={{ fill: "hsl(var(--chart-1))", strokeWidth: 2, r: 3 }}
                      activeDot={{ r: 5, stroke: "hsl(var(--chart-1))", strokeWidth: 2 }}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </ChartContainer>
            </div>
          ) : (
            <div className="text-center py-8 text-muted-foreground">
              <Activity className="h-8 w-8 mx-auto mb-2 opacity-50" />
              <p className="text-sm">No progress data available for this period</p>
            </div>
          )}

          {/* Additional Info Footer */}
          <div className="text-xs text-muted-foreground space-y-1 pt-2 border-t">
            <div className="flex items-center justify-between">
              <span>Frequency:</span>
              <Badge variant="outline" className="text-xs">{stats.frequency}</Badge>
            </div>
            <div className="flex items-center justify-between">
              <span>Started:</span>
              <span>{new Date(stats.startDate).toLocaleDateString()}</span>
            </div>
            {stats.lastCompleted && (
              <div className="flex items-center justify-between">
                <span>Last completed:</span>
                <span>{new Date(stats.lastCompleted).toLocaleDateString()}</span>
              </div>
            )}
          </div>
        </div>
      </CardContent>
      
      {/* Success indicator gradient */}
      <div 
        className="absolute bottom-0 left-0 right-0 h-1 opacity-60"
        style={{
          background: completionRate >= 80 ? 
            'linear-gradient(90deg, hsl(var(--chart-1)), hsl(var(--chart-2)))' : 
            completionRate >= 60 ? 
            'linear-gradient(90deg, hsl(var(--chart-2)), hsl(var(--chart-3)))' : 
            'linear-gradient(90deg, hsl(var(--destructive)), hsl(var(--muted)))'
        }}
      />
    </Card>
  );
} 