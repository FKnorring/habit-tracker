import { ChartContainer, ChartTooltip } from '@/components/ui/chart';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, ResponsiveContainer } from 'recharts';
import { Badge } from '@/components/ui/badge';
import { ChartCard } from './ChartCard';

interface HabitCompletionRate {
  habitId: string;
  habitName: string;
  actualCompletions: number;
  expectedCompletions: number;
  completionRate: number;
}

interface HabitCompletionProps {
  habitCompletionRates: HabitCompletionRate[] | null;
  timeRange: number;
}

export function HabitCompletion({ habitCompletionRates, timeRange }: HabitCompletionProps) {
  const chartConfig = {
    completions: {
      label: "Completions",
      color: "var(--chart-1)",
    },
    expected: {
      label: "Expected",
      color: "var(--chart-2)",
    },
    actual: {
      label: "Actual",
      color: "var(--chart-3)",
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

  return (
    <ChartCard
      title="Habit Completion Rates"
      description={`Actual vs expected completions for each habit (last ${timeRange} days)`}
    >
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
              fill="var(--chart-2)" 
              name="Expected"
              opacity={0.7}
            />
            <Bar 
              dataKey="actual" 
              fill="var(--chart-1)" 
              name="Actual"
            />
          </BarChart>
        </ResponsiveContainer>
      </ChartContainer>
    </ChartCard>
  );
} 