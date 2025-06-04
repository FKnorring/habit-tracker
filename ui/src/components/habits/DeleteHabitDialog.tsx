"use client";

import { useState } from "react";
import { Habit } from "@/types";
import { useHabits } from "@/components/contexts/HabitsContext";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { toast } from "sonner";

interface DeleteHabitDialogProps {
  habit: Habit;
  children?: React.ReactNode;
}

export function DeleteHabitDialog({
  habit,
  children,
}: DeleteHabitDialogProps) {
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { deleteHabit } = useHabits();

  const handleDelete = async () => {
    if (isSubmitting) return;

    setIsSubmitting(true);
    try {
      await deleteHabit(habit.id);
      toast.success("Habit deleted successfully!");
      setIsDeleteDialogOpen(false);
    } catch (error) {
      toast.error("Failed to delete habit");
      console.error("Error deleting habit:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Habit</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete &quot;{habit.name}&quot;? This
            action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <div className="flex gap-2">
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={isSubmitting}
          >
            {isSubmitting ? "Deleting..." : "Delete"}
          </Button>
          <Button
            variant="outline"
            onClick={() => setIsDeleteDialogOpen(false)}
          >
            Cancel
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
} 