'use client';

import { useState } from 'react';
import { useHabits, EnrichedHabit } from '@/components/contexts/HabitsContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { IconPlus } from '@tabler/icons-react';

interface TrackingEntriesDialogProps {
  habit: EnrichedHabit;
  children?: React.ReactNode;
}

export function TrackingEntriesDialog({ habit, children }: TrackingEntriesDialogProps) {
  const { addTrackingEntry } = useHabits();
  const [isOpen, setIsOpen] = useState(false);
  const [isAddingEntry, setIsAddingEntry] = useState(false);
  const [newEntryNote, setNewEntryNote] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleAddEntry = async () => {
    if (isSubmitting) return;
    
    setIsSubmitting(true);
    try {
      await addTrackingEntry(habit.id, { 
        note: newEntryNote || undefined 
      });
      setNewEntryNote('');
      setIsAddingEntry(false);
    } catch (error) {
      console.error('Error tracking habit:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDialogOpenChange = (open: boolean) => {
    setIsOpen(open);
    if (!open) {
      setIsAddingEntry(false);
      setNewEntryNote('');
    }
  };

  const trackingEntries = habit.trackingEntries || [];
  const sortedEntries = trackingEntries
    .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());

  return (
    <Dialog open={isOpen} onOpenChange={handleDialogOpenChange}>
      <DialogTrigger asChild>
        {children}
      </DialogTrigger>
      <DialogContent className="max-w-2xl max-h-[90vh]">
        <DialogHeader>
          <DialogTitle>Tracking History - {habit.name}</DialogTitle>
          <DialogDescription>
            View all tracking entries and add new ones
          </DialogDescription>
        </DialogHeader>
        
        <div className="space-y-4">
          {/* Add New Entry Section */}
          <div className="border-b pb-4">
            {!isAddingEntry ? (
              <Button 
                onClick={() => setIsAddingEntry(true)}
                className="w-full"
              >
                <IconPlus className="w-4 h-4 mr-2" />
                Add New Entry
              </Button>
            ) : (
              <div className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="new-entry-note">Note (optional)</Label>
                  <Input
                    id="new-entry-note"
                    placeholder="How did it go today?"
                    value={newEntryNote}
                    onChange={(e) => setNewEntryNote(e.target.value)}
                  />
                </div>
                <div className="flex gap-2">
                  <Button onClick={handleAddEntry} disabled={isSubmitting}>
                    {isSubmitting ? 'Adding...' : 'Add Entry'}
                  </Button>
                  <Button 
                    variant="outline" 
                    onClick={() => {
                      setIsAddingEntry(false);
                      setNewEntryNote('');
                    }}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            )}
          </div>

          {/* Entries List */}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <h4 className="font-semibold">
                All Entries ({trackingEntries.length})
              </h4>
            </div>
            
            {trackingEntries.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <p>No tracking entries yet.</p>
                <p className="text-sm">Add your first entry to get started!</p>
              </div>
            ) : (
              <div className="h-[300px] w-full overflow-y-auto rounded-md border p-4">
                <div className="space-y-3">
                  {sortedEntries.map((entry) => (
                    <div key={entry.id} className="p-3 bg-muted rounded-lg">
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <p className="font-medium text-sm">
                            {new Date(entry.timestamp).toLocaleString()}
                          </p>
                          {entry.note && (
                            <p className="text-muted-foreground text-sm mt-1">
                              {entry.note}
                            </p>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
} 