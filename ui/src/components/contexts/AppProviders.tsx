import { NavigationProvider } from "./NavigationContext";
import { SidebarProvider } from "@/components/ui/sidebar";
import { HabitsProvider } from "./HabitsContext";

export function AppProviders({ children }: { children: React.ReactNode }) {
  return (
    <HabitsProvider>
      <NavigationProvider>
        <SidebarProvider style={
            {
              "--sidebar-width": "calc(var(--spacing) * 72)",
              "--header-height": "calc(var(--spacing) * 12)",
            } as React.CSSProperties
          }>{children}</SidebarProvider>
      </NavigationProvider>
    </HabitsProvider>
  );
}