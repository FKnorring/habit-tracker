'use client';

import { useState } from 'react';
import { Habit, Frequency } from '@/types';
import { useHabits } from '@/components/contexts/HabitsContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

interface EditHabitDialogProps {
  habit: Habit;
  children?: React.ReactNode;
}

const frequencyOptions = [
  { value: Frequency.HOURLY, label: 'Hourly' },
  { value: Frequency.DAILY, label: 'Daily' },
  { value: Frequency.WEEKLY, label: 'Weekly' },
  { value: Frequency.BIWEEKLY, label: 'Biweekly' },
  { value: Frequency.MONTHLY, label: 'Monthly' },
  { value: Frequency.QUARTERLY, label: 'Quarterly' },
  { value: Frequency.YEARLY, label: 'Yearly' },
];

export function EditHabitDialog({ habit, children }: EditHabitDialogProps) {
  const { updateHabit } = useHabits();
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [editName, setEditName] = useState(habit.name);
  const [editDescription, setEditDescription] = useState(habit.description);
  const [editFrequency, setEditFrequency] = useState<Frequency>(habit.frequency);

  const handleEdit = async () => {
    if (isSubmitting) return;
    
    setIsSubmitting(true);
    try {
      await updateHabit(habit.id, {
        name: editName,
        description: editDescription,
        frequency: editFrequency,
      });
      setIsEditDialogOpen(false);
    } catch (error) {
      console.error('Error updating habit:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleEditDialogOpen = (open: boolean) => {
    if (open) {
      setEditName(habit.name);
      setEditDescription(habit.description);
      setEditFrequency(habit.frequency);
    }
    setIsEditDialogOpen(open);
  };

  return (
    <Dialog open={isEditDialogOpen} onOpenChange={handleEditDialogOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Habit</DialogTitle>
          <DialogDescription>
            Update the details for &quot;{habit.name}&quot;
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="edit-name">Name</Label>
            <Input
              id="edit-name"
              placeholder="Habit name"
              value={editName}
              onChange={(e) => setEditName(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="edit-description">Description</Label>
            <Input
              id="edit-description"
              placeholder="Habit description"
              value={editDescription}
              onChange={(e) => setEditDescription(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="edit-frequency">Frequency</Label>
            <Select value={editFrequency} onValueChange={(value: Frequency) => setEditFrequency(value)}>
              <SelectTrigger>
                <SelectValue placeholder="Select frequency" />
              </SelectTrigger>
              <SelectContent>
                {frequencyOptions.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="flex gap-2">
            <Button onClick={handleEdit} disabled={isSubmitting || !editName.trim()}>
              {isSubmitting ? 'Updating...' : 'Update Habit'}
            </Button>
            <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>
              Cancel
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
} 