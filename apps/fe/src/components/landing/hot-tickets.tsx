"use client";

import { Flame, Sparkles } from "lucide-react";
import { PublicTicket } from "@/types/ticket";
import { TicketCard } from "./ticket-card";

interface HotTicketsProps {
  tickets: PublicTicket[];
  onSelectTicket: (ticket: PublicTicket) => void;
}

export function HotTickets({ tickets, onSelectTicket }: HotTicketsProps) {
  // Filter top 3 featured or low-stock tickets
  const hotTickets = tickets
    .filter((t) => t.featured || (t.available_stock > 0 && t.available_stock < 30))
    .slice(0, 3);

  if (hotTickets.length === 0) return null;

  return (
    <section id="hot-events" className="scroll-mt-20 py-16 md:py-24 border-b border-[#ded4c1] bg-[#faf7f0] font-sans">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4 mb-10">
          <div className="space-y-2">
            <div className="inline-flex items-center gap-1.5 text-xs font-bold text-[#1b4332] bg-[#ebe4d5] border border-[#ded4c1] px-3.5 py-1 rounded-full shadow-2xs">
              <Flame className="h-3.5 w-3.5 fill-[#c2841f] text-[#c2841f]" />
              <span>Sự Kiện Hot Được Đặt Nhiều Nhất</span>
            </div>
            <h2 className="text-2xl sm:text-3xl font-extrabold tracking-tight text-[#2b241e]">
              Vé Nổi Bật Tuần Này
            </h2>
            <p className="text-sm text-[#5c4e40] max-w-xl font-sans">
              Các đại nhạc hội và hội nghị đỉnh cao đang có số lượng vé mở bán giới hạn. Hãy nhanh tay chọn vé trước khi cháy vé!
            </p>
          </div>

          <div className="flex items-center gap-2 text-xs font-semibold text-[#524538] bg-[#ebe4d5] border border-[#ded4c1] px-3.5 py-2 rounded-xl shadow-2xs">
            <Sparkles className="h-4 w-4 text-[#c2841f]" />
            <span>Cập nhật số lượng vé liên tục</span>
          </div>
        </div>

        {/* 3 Hot Tickets Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
          {hotTickets.map((ticket, idx) => (
            <TicketCard
              key={ticket.id}
              ticket={ticket}
              index={idx + 1}
              onSelect={onSelectTicket}
            />
          ))}
        </div>
      </div>
    </section>
  );
}
