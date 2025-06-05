"use client";

import { useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { useHabits } from "@/components/contexts/HabitsContext";
import { Button } from "@/components/ui/button";
import { Plus, Loader2 } from "lucide-react";
import { CreateHabitDialog } from "./CreateHabitDialog";
import { HabitCard } from "./HabitCard";

export function Habits() {
  const { 
    habits, 
    loading, 
    error, 
    initialized,
    fetchHabits, 
    createHabit 
  } = useHabits();

  useEffect(() => {
    if (!initialized) {
      fetchHabits();
    }
  }, [initialized, fetchHabits]);

  const handleHabitCreated = async (habitData: Parameters<typeof createHabit>[0]) => {
    try {
      await createHabit(habitData);
    } catch (error) {
      console.error('Failed to create habit:', error);
    }
  };

  const containerVariants = {
    hidden: {},
    visible: {
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const itemVariants = {
    hidden: {
      opacity: 0,
      y: 20
    },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        duration: 0.5,
        ease: "easeOut"
      }
    }
  };

  return (
    <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6 px-4 md:px-6">
      <div className="flex items-center justify-between">
        <CreateHabitDialog onHabitCreated={handleHabitCreated}>
          <Button>
            <Plus className="h-4 w-4" />
            New Habit
          </Button>
        </CreateHabitDialog>
      </div>

      {error && (
        <div className="text-center py-8">
          <p className="text-destructive mb-4">Error: {error}</p>
          <Button onClick={fetchHabits} variant="outline">
            Try Again
          </Button>
        </div>
      )}

      {loading && (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-8 w-8 animate-spin" />
          <span className="ml-2">Loading habits...</span>
        </div>
      )}

      {!loading && habits && habits.length === 0 ? (
        <div className="text-center py-12">
          <h3 className="text-lg font-semibold mb-2">No habits yet</h3>
          <p className="text-muted-foreground mb-4">
            Get started by creating your first habit to track!
          </p>
          <CreateHabitDialog onHabitCreated={handleHabitCreated}>
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              Create Your First Habit
            </Button>
          </CreateHabitDialog>
        </div>
      ) : (
        <motion.div 
          className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
          variants={containerVariants}
          initial="hidden"
          animate="visible"
        >
          <AnimatePresence>
            {habits?.map((habit) => (
              <motion.div
                key={habit.id}
                variants={itemVariants}
                layout
              >
                <HabitCard
                  habit={habit}
                />
              </motion.div>
            ))}
          </AnimatePresence>
        </motion.div>
      )}
    </div>
  );
}
