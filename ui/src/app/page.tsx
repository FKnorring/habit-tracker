import { SiteHeader } from "@/components/SiteHeader";
import { AppSidebar } from "@/components/AppSidebar";
import { SidebarInset } from "@/components/ui/sidebar";
import { MainContent } from "@/components/MainContent";
import { AppProviders } from "@/components/contexts/AppProviders";

export default function Page() {
  return (
    <AppProviders>
      <AppSidebar variant="inset" />
      <SidebarInset>
        <SiteHeader />
        <div className="flex flex-1 flex-col">
          <div className="@container/main flex flex-1 flex-col gap-2">
            <MainContent />
          </div>
        </div>
      </SidebarInset>
    </AppProviders>
  );
}
