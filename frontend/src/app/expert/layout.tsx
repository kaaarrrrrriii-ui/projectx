import type { ReactNode } from "react";
import SiteHeader from "@/widgets/site-header/site-header";
import ExpertSidebar from "@/widgets/expert-dashboard/expert-sidebar";

export default function ExpertLayout({ children }: { children: ReactNode }) {
  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader />
      <div className="flex flex-1 items-stretch max-[799px]:flex-col">
        <ExpertSidebar />
        {children}
      </div>
    </main>
  );
}
