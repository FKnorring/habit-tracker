import { ChartContainer, ChartTooltip, ChartTooltipContent } from '@/components/ui/chart';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, ResponsiveContainer } from 'recharts';
import { ChartCard } from './ChartCard';

interface DailyCompletion {
  date: string;
  completions: number;
}

interface DailyCompletionsProps {
  dailyCompletions: DailyCompletion[] | null;
  timeRange: number;
}

export function DailyCompletions({ dailyCompletions, timeRange }: DailyCompletionsProps) {
  const chartConfig = {
    completions: {
      label: "Completions",
      color: "var(--chart-1)",
    },
  };

  // Prepare data for daily completions area chart
  const dailyChartData = dailyCompletions?.map(day => ({
    date: new Date(day.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
    completions: day.completions,
  })) || [];

  return (
    <ChartCard
      title="Daily Completions"
      description={`Your completion trend over the last ${timeRange} days`}
    >
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
              stroke="var(--chart-1)"
              fill="var(--chart-1)"
              fillOpacity={0.4}
            />
          </AreaChart>
        </ResponsiveContainer>
      </ChartContainer>
    </ChartCard>
  );
} 