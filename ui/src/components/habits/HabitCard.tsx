"use client";

import { useEffect } from "react";
import { useHabits, EnrichedHabit } from "@/components/contexts/HabitsContext";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { EditHabitDialog } from "./EditHabitDialog";
import { AddTrackingDialog } from "./AddTrackingDialog";
import { TrackingEntriesDialog } from "./TrackingEntriesDialog";
import { DeleteHabitDialog } from "./DeleteHabitDialog";
import { IconPencil, IconX, IconHistory, IconPlus } from "@tabler/icons-react";
import { Badge } from "../ui/badge";
import { Button } from "../ui/button";

interface HabitCardProps {
  habit: EnrichedHabit;
}

export function HabitCard({ habit }: HabitCardProps) {
  const { enrichHabitWithTracking } = useHabits();

  useEffect(() => {
    if (!habit.trackingEntries) {
      enrichHabitWithTracking(habit.id);
    }
  }, [habit.id, habit.trackingEntries, enrichHabitWithTracking]);

  return (
    <Card className="@container/card">
      <CardHeader>
        <CardDescription>
          Times tracked:{" "}
          {habit.trackingEntries === undefined ? "..." : habit.trackingEntries?.length || 0}
        </CardDescription>
        <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
          {habit.name}
        </CardTitle>
        <CardDescription>{habit.description || "No description"}</CardDescription>
        <CardAction className="flex gap-1">

          {/* Edit Habit Button */}
          <EditHabitDialog habit={habit}>
            <Badge
              className="p-1 hover:bg-accent transition-colors"
              variant="outline"
            >
              <IconPencil />
            </Badge>
          </EditHabitDialog>

          {/* Delete Habit Button */}
          <DeleteHabitDialog habit={habit}>
            <Badge
              className="p-1 hover:bg-accent transition-colors"
              variant="outline"
            >
              <IconX />
            </Badge>
          </DeleteHabitDialog>

        </CardAction>
      </CardHeader>
      <CardContent>
        <div className="line-clamp-1 flex gap-2 font-medium capitalize">
          {habit.frequency || "N/A"}
        </div>
        <div className="text-muted-foreground">
          Started: {new Date(habit.startDate).toLocaleDateString()}
        </div>
      </CardContent>
      <CardFooter className="flex gap-2">

        {/* Quick Add Button */}
        <AddTrackingDialog habit={habit}>
          <Button size="sm" className="">
            <IconPlus className="w-4 h-4" />
          </Button>
        </AddTrackingDialog>
        
        {/* View/Manage All Entries Button */}
        <TrackingEntriesDialog habit={habit}>
          <Button variant="outline" size="sm">
            <IconHistory className="w-4 h-4 mr-1" />
            View History
          </Button>
        </TrackingEntriesDialog>

      </CardFooter>
    </Card>
  );
}
