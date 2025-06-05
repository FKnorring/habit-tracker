import { CalendarDays, Target, TrendingUp, Zap } from 'lucide-react';
import { StatCard } from './StatCard';

interface OverallStats {
  totalHabits: number;
  totalEntries: number;
  entriesToday: number;
  avgEntriesPerDay: number;
}

interface StatCardsProps {
  overallStats: OverallStats | null;
  timeRange: number;
}

export function StatCards({ overallStats, timeRange }: StatCardsProps) {
  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <StatCard
        title="Total Habits"
        value={overallStats?.totalHabits || 0}
        description="Active habits"
        icon={Target}
      />
      
      <StatCard
        title="Total Entries"
        value={overallStats?.totalEntries || 0}
        description="All-time completions"
        icon={Zap}
      />
      
      <StatCard
        title="Today's Progress"
        value={overallStats?.entriesToday || 0}
        description="Completed today"
        icon={CalendarDays}
      />
      
      <StatCard
        title="Daily Average"
        value={overallStats?.avgEntriesPerDay ? overallStats.avgEntriesPerDay.toFixed(1) : '0.0'}
        description={`Last ${timeRange} days`}
        icon={TrendingUp}
      />
    </div>
  );
} 