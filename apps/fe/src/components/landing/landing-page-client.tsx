"use client";

import { useState } from "react";
import { PublicTicket, PaginationMeta } from "@/types/ticket";
import { HeroSection } from "./hero-section";
import { HotTickets } from "./hot-tickets";
import { TicketCatalog } from "./ticket-catalog";
import { FeaturesSection } from "./features-section";
import { FaqSection } from "./faq-section";
import { TicketModal } from "./ticket-modal";

interface LandingPageClientProps {
  initialTickets: PublicTicket[];
  pagination: PaginationMeta;
}

export function LandingPageClient({
  initialTickets,
  pagination,
}: LandingPageClientProps) {
  const [selectedTicket, setSelectedTicket] = useState<PublicTicket | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleOpenTicketModal = (ticket: PublicTicket) => {
    setSelectedTicket(ticket);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedTicket(null);
  };

  return (
    <div className="flex flex-col min-h-screen">
      {/* Hero Section */}
      <HeroSection />

      {/* Hot Events Section */}
      <HotTickets
        tickets={initialTickets}
        onSelectTicket={handleOpenTicketModal}
      />

      {/* Main Filterable Ticket Catalog */}
      <TicketCatalog
        initialTickets={initialTickets}
        pagination={pagination}
        onSelectTicket={handleOpenTicketModal}
      />

      {/* Technology & Platform Guarantees */}
      <FeaturesSection />

      {/* FAQs */}
      <FaqSection />

      {/* Interactive Reservation Modal */}
      <TicketModal
        ticket={selectedTicket}
        isOpen={isModalOpen}
        onClose={handleCloseModal}
      />
    </div>
  );
}
