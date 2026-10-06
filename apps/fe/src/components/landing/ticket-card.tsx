"use client";

import Image from "next/image";
import { Calendar, MapPin, Users, Flame, Ban, ArrowRight } from "lucide-react";
import { PublicTicket } from "@/types/ticket";
import { formatVND } from "@/lib/utils/format";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { TicketNotchDivider } from "@/components/ui/ticket-notch";

interface TicketCardProps {
  ticket: PublicTicket;
  onSelect: (ticket: PublicTicket) => void;
  index?: number;
}

export function TicketCard({ ticket, onSelect, index = 1 }: TicketCardProps) {
  const isSoldOut = ticket.is_sold_out || ticket.available_stock <= 0;
  const isLowStock = !isSoldOut && ticket.available_stock < 20;

  return (
    <div className="group relative flex flex-col justify-between font-sans transition-all duration-300 hover:-translate-y-1.5 drop-shadow-xs hover:drop-shadow-md">
      {/* 1. Top Section: Event Poster & Core Information */}
      <div className="rounded-t-2xl border-t border-x border-[#ded4c1] bg-[#fdfcf9] overflow-hidden">
        {/* Banner Image */}
        <div className="relative h-48 w-full overflow-hidden bg-[#ebe4d5]">
          {ticket.image_url ? (
            <Image
              src={ticket.image_url}
              alt={ticket.name}
              fill
              unoptimized
              sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
              className="object-cover transition-transform duration-500 group-hover:scale-105"
            />
          ) : (
            <div className="flex h-full w-full items-center justify-center bg-[#ebe4d5] text-[#857361] text-xs">
              TicketBox
            </div>
          )}

          {/* Vignette Gradient */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/65 via-black/15 to-transparent" />

          {/* Badges on Banner */}
          <div className="absolute top-3 left-3 right-3 flex items-center justify-between">
            <span className="rounded-lg bg-[#faf7f0]/95 px-2.5 py-1 text-xs font-bold text-[#1b4332] shadow-2xs border border-[#ded4c1]/60 backdrop-blur-xs">
              {ticket.category || "Sự kiện"}
            </span>

            {isSoldOut ? (
              <Badge variant="destructive" className="text-xs font-semibold">
                <Ban className="h-3 w-3 mr-1" />
                Hết vé
              </Badge>
            ) : isLowStock ? (
              <Badge className="bg-[#c2841f] text-white text-xs font-semibold">
                <Flame className="h-3 w-3 mr-1 fill-white" />
                Chỉ còn {ticket.available_stock} vé
              </Badge>
            ) : ticket.featured ? (
              <Badge className="bg-[#1b4332] text-[#faf7f0] text-xs font-semibold">
                <Flame className="h-3 w-3 mr-1 fill-[#e9c46a]" />
                Nổi bật
              </Badge>
            ) : null}
          </div>

          {/* Date & Time Badge */}
          <div className="absolute bottom-3 left-3 flex items-center gap-1.5 rounded-lg bg-[#faf7f0]/95 px-2.5 py-1 text-xs font-semibold text-[#2b241e] border border-[#ded4c1] shadow-2xs">
            <Calendar className="h-3.5 w-3.5 text-[#1b4332]" />
            <span>{ticket.date || "15/11/2026"}</span>
            <span className="text-[#b8ab96]">•</span>
            <span>{ticket.time || "19:30"}</span>
          </div>
        </div>

        {/* Event Content */}
        <div className="p-5 space-y-3">
          <h3 className="line-clamp-2 text-base font-bold tracking-tight text-[#2b241e] group-hover:text-[#1b4332] transition-colors leading-snug">
            {ticket.name}
          </h3>

          <div className="flex items-start gap-1.5 text-xs text-[#5c4e40] font-sans">
            <MapPin className="h-3.5 w-3.5 shrink-0 text-[#c2841f] mt-0.5" />
            <span className="line-clamp-1">{ticket.venue || "Địa điểm cập nhật"}</span>
          </div>

          <p className="line-clamp-2 text-xs text-[#6e5d4d] leading-relaxed font-sans">
            {ticket.description}
          </p>

          <div className="pt-2 border-t border-[#ded4c1]/60 flex items-center justify-between text-xs text-[#5c4e40] font-sans">
            <div className="flex items-center gap-1.5">
              <Users className="h-3.5 w-3.5 text-[#1b4332]" />
              <span>Tối đa {ticket.max_booking_per_user ?? 4} vé/người</span>
            </div>
            <span className="text-[#857361] font-mono text-[11px] font-bold">#TK-{String(index).padStart(2, "0")}</span>
          </div>
        </div>
      </div>

      {/* 2. Authentic Ticket Notches & Perforated Tear Line (True Concave Cutout) */}
      <TicketNotchDivider
        fillColor="#fdfcf9"
        borderColor="#ded4c1"
        dashedLineColor="#c8beaf"
        notchRadius={10}
      />

      {/* 3. Bottom Stub: Ticket Price & Booking CTA */}
      <div className="rounded-b-2xl border-b border-x border-[#ded4c1] bg-[#fdfcf9] p-5 pt-3.5 flex items-center justify-between gap-3">
        <div>
          <div className="text-[10px] uppercase tracking-wider text-[#857361] font-bold">
            Giá vé niêm yết
          </div>
          <div className="text-lg font-black text-[#1b4332] tabular-nums">
            {formatVND(ticket.price)}
          </div>
        </div>

        <Button
          onClick={() => onSelect(ticket)}
          disabled={isSoldOut}
          className={`text-xs font-semibold px-4 h-10 rounded-xl transition-all cursor-pointer flex items-center gap-1.5 ${
            isSoldOut
              ? "bg-[#ded4c1] text-[#857361] cursor-not-allowed"
              : "bg-[#1b4332] hover:bg-[#143225] text-[#faf7f0] border border-[#143225] shadow-xs hover:shadow-sm"
          }`}
        >
          {isSoldOut ? (
            <span>Hết vé</span>
          ) : (
            <>
              <span>Đặt vé ngay</span>
              <ArrowRight className="h-3.5 w-3.5" />
            </>
          )}
        </Button>
      </div>
    </div>
  );
}
