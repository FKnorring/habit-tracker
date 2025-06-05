import { NavigationProvider } from "./NavigationContext";
import { SidebarProvider } from "@/components/ui/sidebar";
import { HabitsProvider } from "./HabitsContext";
import { SocketProvider } from "./SocketContext";

export function AppProviders({ children }: { children: React.ReactNode }) {
  return (
    <HabitsProvider>
      <NavigationProvider>
        <SidebarProvider style={
            {
              "--sidebar-width": "calc(var(--spacing) * 72)",
              "--header-height": "calc(var(--spacing) * 12)",
            } as React.CSSProperties
          }>
            <SocketProvider>
              {children}
            </SocketProvider>
          </SidebarProvider>
      </NavigationProvider>
    </HabitsProvider>
  );
}