"use client";

import { useEffect, useState } from 'react';
import { useHabits } from '@/components/contexts/HabitsContext';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ChartContainer, ChartTooltip, ChartTooltipContent } from '@/components/ui/chart';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { CalendarDays, Target, TrendingUp, Zap, ChevronDown, ChevronUp, BarChart3 } from 'lucide-react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, ResponsiveContainer, PieChart, Pie, Cell, AreaChart, Area } from 'recharts';
import { HabitStatsCard } from '@/components/HabitStatsCard';

export function Statistics() {
  const {
    overallStats,
    habitCompletionRates,
    dailyCompletions,
    statisticsLoading,
    statisticsError,
    fetchAllStatistics,
    habits
  } = useHabits();

  const [timeRange, setTimeRange] = useState<number>(30);
  const [showDetailedStats, setShowDetailedStats] = useState<boolean>(false);
  const [selectedHabits, setSelectedHabits] = useState<Set<string>>(new Set());

  useEffect(() => {
    fetchAllStatistics(timeRange);
  }, [fetchAllStatistics, timeRange]);

  // Toggle individual habit selection
  const toggleHabitSelection = (habitId: string) => {
    setSelectedHabits(prev => {
      const newSet = new Set(prev);
      if (newSet.has(habitId)) {
        newSet.delete(habitId);
      } else {
        newSet.add(habitId);
      }
      return newSet;
    });
  };

  // Select all habits
  const selectAllHabits = () => {
    if (habits) {
      setSelectedHabits(new Set(habits.map(h => h.id)));
    }
  };

  // Clear all selections
  const clearAllSelections = () => {
    setSelectedHabits(new Set());
  };

  if (statisticsLoading) {
    return (
      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6 px-4 md:px-6">
        <h2 className="text-2xl font-bold">Statistics</h2>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {[...Array(4)].map((_, i) => (
            <Card key={i} className="animate-pulse">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <div className="h-4 bg-gray-200 rounded w-20"></div>
                <div className="h-4 w-4 bg-gray-200 rounded"></div>
              </CardHeader>
              <CardContent>
                <div className="h-8 bg-gray-200 rounded w-16 mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-24"></div>
              </CardContent>
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

  const chartConfig = {
    completions: {
      label: "Completions",
      color: "hsl(var(--chart-1))",
    },
    expected: {
      label: "Expected",
      color: "hsl(var(--chart-2))",
    },
    actual: {
      label: "Actual",
      color: "hsl(var(--chart-3))",
    },
  };

  // Prepare data for habit completion chart
  const completionChartData = habitCompletionRates?.map(habit => ({
    name: habit.habitName.slice(0, 15) + (habit.habitName.length > 15 ? '...' : ''),
    fullName: habit.habitName,
    actual: habit.actualCompletions,
    expected: habit.expectedCompletions,
    rate: Math.round(habit.completionRate * 100),
  })) || [];

  // Prepare data for daily completions area chart
  const dailyChartData = dailyCompletions?.map(day => ({
    date: new Date(day.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
    completions: day.completions,
  })) || [];

  // Prepare data for habit frequency pie chart
  const frequencyData = habits?.reduce((acc, habit) => {
    const freq = habit.frequency;
    acc[freq] = (acc[freq] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  const pieData = Object.entries(frequencyData || {}).map(([frequency, count]) => ({
    name: frequency.charAt(0).toUpperCase() + frequency.slice(1),
    value: count,
    fill: `hsl(var(--chart-${(Object.keys(frequencyData || {}).indexOf(frequency) % 5) + 1}))`
  }));

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
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Habits</CardTitle>
            <Target className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{overallStats?.totalHabits || 0}</div>
            <p className="text-xs text-muted-foreground">Active habits</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Entries</CardTitle>
            <Zap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{overallStats?.totalEntries || 0}</div>
            <p className="text-xs text-muted-foreground">All-time completions</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Today&apos;s Progress</CardTitle>
            <CalendarDays className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{overallStats?.entriesToday || 0}</div>
            <p className="text-xs text-muted-foreground">Completed today</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Daily Average</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {overallStats?.avgEntriesPerDay ? overallStats.avgEntriesPerDay.toFixed(1) : '0.0'}
            </div>
            <p className="text-xs text-muted-foreground">Last {timeRange} days</p>
          </CardContent>
        </Card>
      </div>

      {/* Charts Section */}
      <div className="grid gap-4 md:grid-cols-2">
        {/* Daily Completions Area Chart */}
        <Card>
          <CardHeader>
            <CardTitle>Daily Completions</CardTitle>
            <CardDescription>Your completion trend over the last {timeRange} days</CardDescription>
          </CardHeader>
          <CardContent>
            <ChartContainer config={chartConfig} className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={dailyChartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis 
                    dataKey="date" 
                    fontSize={12}
                    tickLine={false}
                    axisLine={false}
                  />
                  <YAxis 
                    fontSize={12}
                    tickLine={false}
                    axisLine={false}
                  />
                  <ChartTooltip 
                    content={<ChartTooltipContent />}
                  />
                  <Area
                    type="monotone"
                    dataKey="completions"
                    stroke="hsl(var(--chart-1))"
                    fill="hsl(var(--chart-1))"
                    fillOpacity={0.4}
                  />
                </AreaChart>
              </ResponsiveContainer>
            </ChartContainer>
          </CardContent>
        </Card>

        {/* Habit Frequency Distribution */}
        <Card>
          <CardHeader>
            <CardTitle>Habit Frequency Distribution</CardTitle>
            <CardDescription>How your habits are distributed by frequency</CardDescription>
          </CardHeader>
          <CardContent>
            <ChartContainer config={chartConfig} className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={pieData}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={100}
                    paddingAngle={5}
                    dataKey="value"
                  >
                    {pieData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.fill} />
                    ))}
                  </Pie>
                  <ChartTooltip 
                    content={({ active, payload }) => {
                      if (active && payload && payload.length) {
                        const data = payload[0].payload;
                        return (
                          <div className="rounded-lg border bg-background p-2 shadow-sm">
                            <div className="grid grid-cols-2 gap-2">
                              <div className="flex flex-col">
                                <span className="text-[0.70rem] uppercase text-muted-foreground">
                                  Frequency
                                </span>
                                <span className="font-bold text-muted-foreground">
                                  {data.name}
                                </span>
                              </div>
                              <div className="flex flex-col">
                                <span className="text-[0.70rem] uppercase text-muted-foreground">
                                  Count
                                </span>
                                <span className="font-bold">
                                  {data.value}
                                </span>
                              </div>
                            </div>
                          </div>
                        );
                      }
                      return null;
                    }}
                  />
                </PieChart>
              </ResponsiveContainer>
            </ChartContainer>
          </CardContent>
        </Card>
      </div>

      {/* Habit Completion Rates Bar Chart */}
      <Card>
        <CardHeader>
          <CardTitle>Habit Completion Rates</CardTitle>
          <CardDescription>
            Actual vs expected completions for each habit (last {timeRange} days)
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ChartContainer config={chartConfig} className="h-[400px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={completionChartData} margin={{ top: 20, right: 30, left: 20, bottom: 60 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis 
                  dataKey="name" 
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                  angle={-45}
                  textAnchor="end"
                  height={80}
                />
                <YAxis 
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                />
                <ChartTooltip 
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      return (
                        <div className="rounded-lg border bg-background p-3 shadow-sm">
                          <div className="grid gap-2">
                            <div className="font-medium">{data.fullName}</div>
                            <div className="grid grid-cols-2 gap-4">
                              <div className="flex flex-col">
                                <span className="text-[0.70rem] uppercase text-muted-foreground">
                                  Expected
                                </span>
                                <span className="font-bold text-muted-foreground">
                                  {data.expected}
                                </span>
                              </div>
                              <div className="flex flex-col">
                                <span className="text-[0.70rem] uppercase text-muted-foreground">
                                  Actual
                                </span>
                                <span className="font-bold">
                                  {data.actual}
                                </span>
                              </div>
                            </div>
                            <div className="flex items-center gap-2 pt-1">
                              <span className="text-[0.70rem] uppercase text-muted-foreground">
                                Completion Rate
                              </span>
                              <Badge variant={data.rate >= 80 ? "default" : data.rate >= 60 ? "secondary" : "destructive"}>
                                {data.rate}%
                              </Badge>
                            </div>
                          </div>
                        </div>
                      );
                    }
                    return null;
                  }}
                />
                <Bar 
                  dataKey="expected" 
                  fill="hsl(var(--chart-2))" 
                  name="Expected"
                  opacity={0.7}
                />
                <Bar 
                  dataKey="actual" 
                  fill="hsl(var(--chart-1))" 
                  name="Actual"
                />
              </BarChart>
            </ResponsiveContainer>
          </ChartContainer>
        </CardContent>
      </Card>

      {/* Individual Habit Statistics Section */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <BarChart3 className="h-5 w-5" />
                Individual Habit Analytics
              </CardTitle>
              <CardDescription>
                Detailed statistics and progress charts for each habit
              </CardDescription>
            </div>
            <Button
              variant="outline"
              onClick={() => setShowDetailedStats(!showDetailedStats)}
              className="flex items-center gap-2"
            >
              {showDetailedStats ? (
                <>
                  <ChevronUp className="h-4 w-4" />
                  Hide Details
                </>
              ) : (
                <>
                  <ChevronDown className="h-4 w-4" />
                  Show Details
                </>
              )}
            </Button>
          </div>
        </CardHeader>
        {showDetailedStats && (
          <CardContent>
            <div className="space-y-4">
              {/* Habit Selection Controls */}
              {habits && habits.length > 0 && (
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <p className="text-sm text-muted-foreground">
                      Select habits to view detailed analytics:
                    </p>
                    <div className="flex gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={selectAllHabits}
                        disabled={selectedHabits.size === habits.length}
                      >
                        Select All
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={clearAllSelections}
                        disabled={selectedHabits.size === 0}
                      >
                        Clear All
                      </Button>
                    </div>
                  </div>
                  
                  {/* Habit Selection Grid */}
                  <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2">
                    {habits.map((habit) => (
                      <Button
                        key={habit.id}
                        variant={selectedHabits.has(habit.id) ? "default" : "outline"}
                        size="sm"
                        onClick={() => toggleHabitSelection(habit.id)}
                        className="justify-start"
                      >
                        {habit.name}
                      </Button>
                    ))}
                  </div>
                </div>
              )}

              {/* Selected Habit Stats Cards */}
              {selectedHabits.size > 0 ? (
                <div className="grid gap-6 md:grid-cols-1 lg:grid-cols-2">
                  {Array.from(selectedHabits).map((habitId) => (
                    <HabitStatsCard
                      key={habitId}
                      habitId={habitId}
                      days={timeRange}
                    />
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <BarChart3 className="h-12 w-12 text-muted-foreground mx-auto mb-3" />
                  <p className="text-sm text-muted-foreground">
                    Select one or more habits above to view detailed analytics
                  </p>
                </div>
              )}
            </div>
          </CardContent>
        )}
      </Card>
    </div>
  );
} 