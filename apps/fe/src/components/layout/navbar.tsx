"use client";

import { useState } from "react";
import Link from "next/link";
import { Ticket, Menu, X, PhoneCall, ShieldCheck, Sparkles } from "lucide-react";
import { Button } from "@/components/ui/button";

export function Navbar() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <header className="sticky top-0 z-50 w-full border-b border-[#ded4c1] bg-[#faf7f0]/95 backdrop-blur-md shadow-xs">
      {/* Top Retro Parchment Strip */}
      <div className="bg-[#ebe4d5] border-b border-[#ded4c1] px-4 py-1.5 text-xs text-[#524538]">
        <div className="mx-auto max-w-7xl flex items-center justify-between">
          <div className="flex items-center gap-3">
            <span className="flex items-center gap-1.5 text-[#1b4332] font-semibold">
              <ShieldCheck className="h-4 w-4" />
              Quầy Vé Chính Hãng TicketBox
            </span>
            <span className="hidden sm:inline-block text-[#b8ab96]">•</span>
            <span className="hidden sm:inline-block text-[#615243]">
              100% Vé gốc niêm yết từ Ban Tổ Chức
            </span>
          </div>

          <div className="flex items-center gap-4 text-[#524538]">
            <a
              href="tel:19008425"
              className="flex items-center gap-1 text-[#1b4332] hover:text-[#143225] font-semibold transition-colors"
            >
              <PhoneCall className="h-3.5 w-3.5" />
              <span>Hotline: 1900-TICKET</span>
            </a>
          </div>
        </div>
      </div>

      {/* Main Navbar */}
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        {/* Brand Logo */}
        <Link href="/" className="group flex items-center gap-2.5 transition-transform active:scale-95">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#1b4332] text-[#faf7f0] shadow-sm group-hover:bg-[#143225] transition-colors">
            <Ticket className="h-5 w-5 text-[#e9c46a]" />
          </div>
          <div className="flex flex-col">
            <span className="text-xl font-black tracking-tight text-[#2b241e] font-sans">
              Ticket<span className="text-[#1b4332]">Box</span>
            </span>
            <span className="text-[10px] font-semibold text-[#857361] font-sans tracking-wider uppercase">
              Quầy Vé Sự Kiện Trực Tuyến
            </span>
          </div>
        </Link>

        {/* Desktop Nav Links */}
        <nav className="hidden md:flex items-center gap-8 text-sm font-semibold text-[#524538]">
          <a
            href="#hot-events"
            className="hover:text-[#1b4332] transition-colors"
          >
            Sự Kiện Nổi Bật
          </a>
          <a
            href="#all-tickets"
            className="hover:text-[#1b4332] transition-colors"
          >
            Tất Cả Vé
          </a>
          <a
            href="#why-choose-us"
            className="hover:text-[#1b4332] transition-colors"
          >
            Cam Kết Nền Tảng
          </a>
          <a
            href="#faq"
            className="hover:text-[#1b4332] transition-colors"
          >
            Hỏi Đáp & Trợ Giúp
          </a>
        </nav>

        {/* Action Button */}
        <div className="hidden sm:flex items-center gap-3">
          <a href="#all-tickets">
            <Button
              size="sm"
              className="bg-[#1b4332] hover:bg-[#143225] text-[#faf7f0] font-semibold px-4 h-9 rounded-xl shadow-xs cursor-pointer flex items-center gap-1.5"
            >
              <Sparkles className="h-3.5 w-3.5 text-[#e9c46a]" />
              <span>Chọn Vé Ngay</span>
            </Button>
          </a>
        </div>

        {/* Mobile menu toggle */}
        <div className="flex md:hidden items-center gap-2">
          <button
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="p-2 rounded-lg border border-[#ded4c1] bg-[#faf7f0] text-[#524538] hover:text-[#1b4332]"
            aria-label="Toggle Menu"
          >
            {mobileMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </button>
        </div>
      </div>

      {/* Mobile Drawer Menu */}
      {mobileMenuOpen && (
        <div className="md:hidden border-b border-[#ded4c1] bg-[#faf7f0] p-4 space-y-3 text-sm shadow-md">
          <a
            href="#hot-events"
            onClick={() => setMobileMenuOpen(false)}
            className="block px-3 py-2 rounded-lg text-[#524538] hover:bg-[#ebe4d5] hover:text-[#1b4332] font-semibold"
          >
            Sự Kiện Nổi Bật
          </a>
          <a
            href="#all-tickets"
            onClick={() => setMobileMenuOpen(false)}
            className="block px-3 py-2 rounded-lg text-[#524538] hover:bg-[#ebe4d5] hover:text-[#1b4332] font-semibold"
          >
            Tất Cả Vé
          </a>
          <a
            href="#why-choose-us"
            onClick={() => setMobileMenuOpen(false)}
            className="block px-3 py-2 rounded-lg text-[#524538] hover:bg-[#ebe4d5] hover:text-[#1b4332] font-semibold"
          >
            Cam Kết Nền Tảng
          </a>
          <a
            href="#faq"
            onClick={() => setMobileMenuOpen(false)}
            className="block px-3 py-2 rounded-lg text-[#524538] hover:bg-[#ebe4d5] hover:text-[#1b4332] font-semibold"
          >
            Hỏi Đáp & Trợ Giúp
          </a>
          <div className="pt-2">
            <a href="#all-tickets" onClick={() => setMobileMenuOpen(false)}>
              <Button className="w-full bg-[#1b4332] hover:bg-[#143225] text-[#faf7f0] font-semibold">
                Chọn Vé Ngay
              </Button>
            </a>
          </div>
        </div>
      )}
    </header>
  );
}
