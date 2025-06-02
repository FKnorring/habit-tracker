'use client';

import { useState } from 'react';
import { Habit, TrackingEntry } from '@/types';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { api } from '@/lib/api';
import { toast } from 'sonner';

interface HabitCardProps {
  habit: Habit;
  onHabitUpdate: () => void;
  onHabitDelete: (id: string) => void;
}

export function HabitCard({ habit, onHabitUpdate, onHabitDelete }: HabitCardProps) {
  const [isTrackDialogOpen, setIsTrackDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [trackingNote, setTrackingNote] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [trackingEntries, setTrackingEntries] = useState<TrackingEntry[]>([]);
  const [showEntries, setShowEntries] = useState(false);

  const handleTrack = async () => {
    if (isSubmitting) return;
    
    setIsSubmitting(true);
    try {
      await api.createTrackingEntry(habit.id, { 
        note: trackingNote || undefined 
      });
      toast.success('Progress tracked successfully!');
      setTrackingNote('');
      setIsTrackDialogOpen(false);
      onHabitUpdate();
    } catch (error) {
      toast.error('Failed to track progress');
      console.error('Error tracking habit:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async () => {
    if (isSubmitting) return;
    
    setIsSubmitting(true);
    try {
      await api.deleteHabit(habit.id);
      toast.success('Habit deleted successfully!');
      onHabitDelete(habit.id);
      setIsDeleteDialogOpen(false);
    } catch (error) {
      toast.error('Failed to delete habit');
      console.error('Error deleting habit:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const loadTrackingEntries = async () => {
    try {
      const entries = await api.getTrackingEntries(habit.id);
      setTrackingEntries(entries);
      setShowEntries(true);
    } catch (error) {
      toast.error('Failed to load tracking entries');
      console.error('Error loading tracking entries:', error);
    }
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>{habit.name}</CardTitle>
        <CardDescription>{habit.description}</CardDescription>
        <p className="text-sm text-muted-foreground">
          Frequency: {habit.frequency} | Started: {new Date(habit.startDate).toLocaleDateString()}
        </p>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex gap-2 flex-wrap">
          <Dialog open={isTrackDialogOpen} onOpenChange={setIsTrackDialogOpen}>
            <DialogTrigger asChild>
              <Button variant="default">Track Progress</Button>
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

          <Button variant="outline" onClick={loadTrackingEntries}>
            View History
          </Button>

          <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
            <DialogTrigger asChild>
              <Button variant="destructive">Delete</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Delete Habit</DialogTitle>
                <DialogDescription>
                  Are you sure you want to delete &quot;{habit.name}&quot;? This action cannot be undone.
                </DialogDescription>
              </DialogHeader>
              <div className="flex gap-2">
                <Button variant="destructive" onClick={handleDelete} disabled={isSubmitting}>
                  {isSubmitting ? 'Deleting...' : 'Delete'}
                </Button>
                <Button variant="outline" onClick={() => setIsDeleteDialogOpen(false)}>
                  Cancel
                </Button>
              </div>
            </DialogContent>
          </Dialog>
        </div>

        {showEntries && (
          <div className="space-y-2">
            <h4 className="font-semibold">Recent Entries:</h4>
            {trackingEntries.length === 0 ? (
              <p className="text-muted-foreground">No tracking entries yet.</p>
            ) : (
              <div className="max-h-32 overflow-y-auto space-y-1">
                {trackingEntries
                  .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
                  .map((entry) => (
                    <div key={entry.id} className="text-sm p-2 bg-muted rounded">
                      <p className="font-medium">{new Date(entry.timestamp).toLocaleString()}</p>
                      {entry.note && <p className="text-muted-foreground">{entry.note}</p>}
                    </div>
                  ))}
              </div>
            )}
            <Button variant="ghost" size="sm" onClick={() => setShowEntries(false)}>
              Hide History
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
} 