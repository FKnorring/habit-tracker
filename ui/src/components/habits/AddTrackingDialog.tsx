'use client';

import { useState } from 'react';
import { Habit } from '@/types';
import { useHabits } from '@/components/contexts/HabitsContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';

interface AddTrackingDialogProps {
  habit: Habit;
  children?: React.ReactNode;
}

export function AddTrackingDialog({ habit, children }: AddTrackingDialogProps) {
  const { addTrackingEntry } = useHabits();
  const [isTrackDialogOpen, setIsTrackDialogOpen] = useState(false);
  const [trackingNote, setTrackingNote] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleTrack = async () => {
    if (isSubmitting) return;
    
    setIsSubmitting(true);
    try {
      await addTrackingEntry(habit.id, { 
        note: trackingNote || undefined 
      });
      setTrackingNote('');
      setIsTrackDialogOpen(false);
    } catch (error) {
      console.error('Error tracking habit:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={isTrackDialogOpen} onOpenChange={setIsTrackDialogOpen}>
      <DialogTrigger asChild>
        {children || <Button variant="default">Track Progress</Button>}
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Track Progress for {habit.name}</DialogTitle>
          <DialogDescription>
            Add a note about your progress (optional)
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="note">Note</Label>
            <Input
              id="note"
              placeholder="How did it go today?"
              value={trackingNote}
              onChange={(e) => setTrackingNote(e.target.value)}
            />
          </div>
          <div className="flex gap-2">
            <Button onClick={handleTrack} disabled={isSubmitting}>
              {isSubmitting ? 'Saving...' : 'Track Progress'}
            </Button>
            <Button variant="outline" onClick={() => setIsTrackDialogOpen(false)}>
              Cancel
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
} 