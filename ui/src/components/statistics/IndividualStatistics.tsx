import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { BarChart3, ChevronDown, ChevronUp } from 'lucide-react';
import { HabitStatsCard } from './HabitStatsCard';

interface Habit {
  id: string;
  name: string;
  frequency: string;
}

interface IndividualStatisticsProps {
  habits: Habit[] | null;
  timeRange: number;
}

export function IndividualStatistics({ habits, timeRange }: IndividualStatisticsProps) {
  const [showDetailedStats, setShowDetailedStats] = useState<boolean>(false);
  const [selectedHabits, setSelectedHabits] = useState<Set<string>>(new Set());

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

  return (
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
  );
} 