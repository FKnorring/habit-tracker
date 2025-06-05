'use client';

import { useState } from 'react';
import { EnrichedHabit } from '@/components/contexts/HabitsContext';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';

interface TrackingEntriesDialogProps {
  habit: EnrichedHabit;
  children?: React.ReactNode;
}

export function TrackingEntriesDialog({ habit, children }: TrackingEntriesDialogProps) {
  const [isOpen, setIsOpen] = useState(false);

  const handleDialogOpenChange = (open: boolean) => {
    setIsOpen(open);
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
        </DialogHeader>
        
        <div className="space-y-4">

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