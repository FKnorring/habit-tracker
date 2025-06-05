import { ChartContainer, ChartTooltip } from '@/components/ui/chart';
import { PieChart, Pie, Cell, ResponsiveContainer } from 'recharts';
import { ChartCard } from './ChartCard';

interface Habit {
  id: string;
  name: string;
  frequency: string;
}

interface HabitFrequencyProps {
  habits: Habit[] | null;
}

export function HabitFrequency({ habits }: HabitFrequencyProps) {
  // Prepare data for habit frequency pie chart
  const frequencyData = habits?.reduce((acc, habit) => {
    const freq = habit.frequency;
    acc[freq] = (acc[freq] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  // Get unique frequencies and create chart config
  const frequencies = Object.keys(frequencyData || {});
  const chartConfig = frequencies.reduce((config, frequency, index) => {
    config[frequency] = {
      label: frequency.charAt(0).toUpperCase() + frequency.slice(1),
      color: `var(--chart-${(index % 5) + 1})`,
    };
    return config;
  }, {} as Record<string, { label: string; color: string }>);

  const pieData = Object.entries(frequencyData || {}).map(([frequency, count]) => ({
    name: frequency.charAt(0).toUpperCase() + frequency.slice(1),
    value: count,
    fill: `var(--chart-${(frequencies.indexOf(frequency) % 5) + 1})`
  }));

  return (
    <ChartCard
      title="Habit Frequency Distribution"
      description="How your habits are distributed by frequency"
    >
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
    </ChartCard>
  );
} 