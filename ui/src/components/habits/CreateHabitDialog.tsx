'use client';

import { useState } from 'react';
import { CreateHabitRequest } from '@/types';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';

interface CreateHabitDialogProps {
  onHabitCreated: (habitData: CreateHabitRequest) => void | Promise<void>;
  children?: React.ReactNode;
}

export function CreateHabitDialog({ onHabitCreated, children }: CreateHabitDialogProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  // Form state
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [frequency, setFrequency] = useState('');

  const handleSubmit = async () => {
    if (isSubmitting || !name.trim()) return;
    
    setIsSubmitting(true);
    try {
      const newHabit: CreateHabitRequest = {
        name: name.trim(),
        description: description.trim(),
        frequency: frequency.trim(),
        startDate: new Date().toISOString(),
      };
      
      await onHabitCreated(newHabit);
      
      // Reset form and close dialog
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
            <Label htmlFor="frequency">Frequency</Label>
            <Input
              id="frequency"
              placeholder="e.g., Daily, Weekly, 3 times per week"
              value={frequency}
              onChange={(e) => setFrequency(e.target.value)}
            />
          </div>
          <div className="flex gap-2">
            <Button 
              onClick={handleSubmit} 
              disabled={isSubmitting || !name.trim()}
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