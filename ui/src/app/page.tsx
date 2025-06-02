'use client';

import { useEffect, useState } from 'react';
import { Habit } from '@/types';
import { HabitCard } from '@/components/HabitCard';
import { CreateHabitForm } from '@/components/CreateHabitForm';
import { api } from '@/lib/api';
import { toast } from 'sonner';

export default function Home() {
  const [habits, setHabits] = useState<Habit[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const loadHabits = async () => {
    try {
      const fetchedHabits = await api.getHabits();
      setHabits(fetchedHabits);
    } catch (error) {
      toast.error('Failed to load habits');
      console.error('Error loading habits:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleHabitDelete = (deletedId: string) => {
    setHabits(prev => prev.filter(habit => habit.id !== deletedId));
  };

  useEffect(() => {
    loadHabits();
  }, []);

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
      <div className="max-w-6xl mx-auto">
        <header className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            Habit Tracker
          </h1>
          <CreateHabitForm onHabitCreated={loadHabits} />
        </header>

        {isLoading ? (
          <div className="text-center py-12">
            <p className="text-gray-600">Loading habits...</p>
          </div>
        ) : habits.length === 0 ? (
          <div className="text-center py-12">
            <div className="bg-white rounded-lg shadow-sm p-8 max-w-md mx-auto">
              <h2 className="text-xl font-semibold text-gray-900 mb-2">
                No habits yet
              </h2>
            </div>
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {habits.map((habit) => (
              <HabitCard
                key={habit.id}
                habit={habit}
                onHabitUpdate={loadHabits}
                onHabitDelete={handleHabitDelete}
              />
            ))}
          </div>
        )}
      </div>
    </main>
  );
}
