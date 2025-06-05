"use client";

import {
  IconChartBar,
  IconCirclePlusFilled,
  IconCoffee,
} from "@tabler/icons-react";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useNavigation } from "@/components/contexts/NavigationContext";
import { useReminders } from "@/components/contexts/RemindersContext";
import { CreateHabitDialog } from "@/components/habits/CreateHabitDialog";
import { cn } from "@/lib/utils";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "../ui/tooltip";

export function NavMain() {
  const { activeItem, setActiveItem } = useNavigation();
  const { reminders } = useReminders();

  return (
    <SidebarGroup>
      <SidebarGroupContent className="flex flex-col gap-2">
        <SidebarMenu>
          <SidebarMenuItem className="flex items-center gap-2">
            <div className="relative flex items-center w-full">
              <SidebarMenuButton
                tooltip="Habits"
                onClick={() => setActiveItem("habits")}
                className={cn(
                  "min-w-8 duration-200 ease-linear flex-1 border",
                  activeItem === "habits" &&
                    "bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground"
                )}
              >
                <IconCoffee />
                <span>Habits</span>
              </SidebarMenuButton>
              {reminders.size > 0 && (
                <Badge
                  variant="destructive"
                  className="absolute -top-2 -right-2 h-5 w-5 p-0 text-xs flex items-center justify-center min-w-5 z-50 rounded-full"
                >
                  {reminders.size}
                </Badge>
              )}
            </div>
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger>
                  <CreateHabitDialog>
                    <Button
                      size="icon"
                      className="size-8 group-data-[collapsible=icon]:opacity-0"
                      variant="outline"
                    >
                      <IconCirclePlusFilled />
                      <span className="sr-only">Quick Create</span>
                    </Button>
                  </CreateHabitDialog>
                </TooltipTrigger>
                <TooltipContent>Quick Create</TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </SidebarMenuItem>
          <SidebarMenuItem className="flex items-center gap-2">
            <SidebarMenuButton
              tooltip="Statistics"
              onClick={() => setActiveItem("statistics")}
              className={cn(
                "min-w-8 duration-200 ease-linear flex-1 border",
                activeItem === "statistics" &&
                  "bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground"
              )}
            >
              <IconChartBar />
              <span>Statistics</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  );
}
