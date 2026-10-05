import { getPublicTickets } from "@/lib/api/tickets";
import { Navbar } from "@/components/layout/navbar";
import { Footer } from "@/components/layout/footer";
import { LandingPageClient } from "@/components/landing/landing-page-client";

export const revalidate = 60; // Revalidate every minute for ISR

export default async function HomePage() {
  // Server-side fetch targeting Go backend /api/v1/tickets
  const ticketsData = await getPublicTickets(1, 20);

  return (
    <div className="flex min-h-screen flex-col bg-background text-foreground selection:bg-[#1b4332]/20 selection:text-[#1b4332]">
      <Navbar />
      <main className="flex-1">
        <LandingPageClient
          initialTickets={ticketsData.items}
          pagination={ticketsData.pagination}
        />
      </main>
      <Footer />
    </div>
  );
}
