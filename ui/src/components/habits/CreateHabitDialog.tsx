'use client';

import { useState } from 'react';
import { CreateHabitRequest, Frequency } from '@/types';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { useHabits } from '@/components/contexts/HabitsContext';

interface CreateHabitDialogProps {
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

export function CreateHabitDialog({ children }: CreateHabitDialogProps) {
  const { createHabit } = useHabits();
  const [isOpen, setIsOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [frequency, setFrequency] = useState<Frequency | ''>('');

  const handleSubmit = async () => {
    if (isSubmitting || !name.trim() || !frequency) return;
    
    setIsSubmitting(true);
    try {
      const newHabit: CreateHabitRequest = {
        name: name.trim(),
        description: description.trim(),
        frequency: frequency as Frequency,
        startDate: new Date().toISOString(),
      };
      
      await createHabit(newHabit);
      
      
      resetForm();
      setIsOpen(false);
    } catch (error) {
      console.error('Error creating habit:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const resetForm = () => {
    setName('');
    setDescription('');
    setFrequency('');
  };

  const handleOpenChange = (open: boolean) => {
    if (open) {
      resetForm();
    }
    setIsOpen(open);
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        {children}
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create New Habit</DialogTitle>
          <DialogDescription>
            Add a new habit to track your progress
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Name *</Label>
            <Input
              id="name"
              placeholder="e.g., Drink 8 glasses of water"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Input
              id="description"
              placeholder="Brief description of your habit"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="frequency">Frequency *</Label>
            <Select value={frequency} onValueChange={(value: Frequency) => setFrequency(value)}>
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
            <Button 
              onClick={handleSubmit} 
              disabled={isSubmitting || !name.trim() || !frequency}
            >
              {isSubmitting ? 'Creating...' : 'Create Habit'}
            </Button>
            <Button variant="outline" onClick={() => setIsOpen(false)}>
              Cancel
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
} 