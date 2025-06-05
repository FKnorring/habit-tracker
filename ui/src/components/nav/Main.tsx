"use client"

import { IconChartBar, IconCirclePlusFilled, IconCoffee } from "@tabler/icons-react"

import { Button } from "@/components/ui/button"
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { useNavigation } from "@/components/contexts/NavigationContext"

export function NavMain() {
  const { activeItem, setActiveItem } = useNavigation()

  return (
    <SidebarGroup>
      <SidebarGroupContent className="flex flex-col gap-2">
        <SidebarMenu>
          <SidebarMenuItem className="flex items-center gap-2">
            <SidebarMenuButton
              tooltip="Habits"
              onClick={() => setActiveItem('habits')}
              className={
                activeItem === 'habits'
                  ? "bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground min-w-8 duration-200 ease-linear"
                  : "min-w-8 duration-200 ease-linear"
              }
            >
              <IconCoffee />
              <span>Habits</span>
            </SidebarMenuButton>
            <Button
              size="icon"
              className="size-8 group-data-[collapsible=icon]:opacity-0"
              variant="outline"
            >
              <IconCirclePlusFilled />
              <span className="sr-only">Quick Create</span>
            </Button>
          </SidebarMenuItem>
          <SidebarMenuItem className="flex items-center gap-2">
            <SidebarMenuButton
              tooltip="Statistics"
              onClick={() => setActiveItem('statistics')}
              className={
                activeItem === 'statistics'
                  ? "bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground min-w-8 duration-200 ease-linear"
                  : "min-w-8 duration-200 ease-linear"
              }
            >
              <IconChartBar />
              <span>Statistics</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
         
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}
