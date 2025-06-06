import { NavigationProvider } from "./NavigationContext";
import { SidebarProvider } from "@/components/ui/sidebar";
import { HabitsProvider } from "./HabitsContext";
import { RemindersProvider } from "./RemindersContext";
import { StatisticsProvider } from "./StatisticsContext";
import { SocketProvider } from "./SocketContext";
import { AuthProvider } from "./AuthContext";

export function AppProviders({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      <RemindersProvider>
        <StatisticsProvider>
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
        </StatisticsProvider>
      </RemindersProvider>
    </AuthProvider>
  );
}