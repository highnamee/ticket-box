"use client";

import { useState, useMemo } from "react";
import { Search, Ticket, XCircle } from "lucide-react";
import { PublicTicket, PaginationMeta } from "@/types/ticket";
import { TicketCard } from "./ticket-card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

interface TicketCatalogProps {
  initialTickets: PublicTicket[];
  pagination: PaginationMeta;
  onSelectTicket: (ticket: PublicTicket) => void;
}

const CATEGORIES = [
  "Tất cả",
  "Âm nhạc",
  "Công nghệ",
  "Workshop",
  "Thể thao",
  "Triển lãm",
  "Hài kịch",
  "Điện ảnh",
];

export function TicketCatalog({
  initialTickets,
  pagination,
  onSelectTicket,
}: TicketCatalogProps) {
  const [selectedCategory, setSelectedCategory] = useState("Tất cả");
  const [searchQuery, setSearchQuery] = useState("");
  const [sortBy, setSortBy] = useState<"default" | "price-asc" | "price-desc" | "stock">(
    "default"
  );
  const [onlyAvailable, setOnlyAvailable] = useState(false);

  // Filter & sort logic
  const filteredTickets = useMemo(() => {
    return initialTickets
      .filter((ticket) => {
        // Category match
        if (
          selectedCategory !== "Tất cả" &&
          ticket.category?.toLowerCase() !== selectedCategory.toLowerCase()
        ) {
          return false;
        }

        // Available only match
        if (onlyAvailable && (ticket.is_sold_out || ticket.available_stock <= 0)) {
          return false;
        }

        // Search match
        if (searchQuery.trim() !== "") {
          const query = searchQuery.toLowerCase();
          const matchName = ticket.name.toLowerCase().includes(query);
          const matchDesc = ticket.description.toLowerCase().includes(query);
          const matchVenue = ticket.venue?.toLowerCase().includes(query);
          return matchName || matchDesc || matchVenue;
        }

        return true;
      })
      .sort((a, b) => {
        if (sortBy === "price-asc") return a.price - b.price;
        if (sortBy === "price-desc") return b.price - a.price;
        if (sortBy === "stock") return a.available_stock - b.available_stock;
        return 0;
      });
  }, [initialTickets, selectedCategory, searchQuery, sortBy, onlyAvailable]);

  return (
    <section id="all-tickets" className="scroll-mt-20 py-16 md:py-24 border-b border-[#ded4c1] bg-[#f5f2eb] font-sans">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 space-y-8">
        {/* Title Bar */}
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
          <div className="space-y-2">
            <span className="text-xs font-bold text-[#1b4332] uppercase tracking-wider">
              Danh Mục Sự Kiện Mở Bán
            </span>
            <h2 className="text-2xl sm:text-3xl font-extrabold tracking-tight text-[#2b241e]">
              Khám Phá Toàn Bộ Vé Sự Kiện
            </h2>
            <p className="text-sm text-[#5c4e40] max-w-xl">
              Dễ dàng tìm kiếm và lựa chọn vé theo sở thích của bạn. Chọn số lượng, giữ chỗ ngay và thanh toán nhanh chóng.
            </p>
          </div>

          <div className="text-xs font-semibold text-[#524538] bg-[#faf7f0] border border-[#ded4c1] px-3.5 py-2 rounded-xl shadow-2xs">
            Hiển thị <strong>{filteredTickets.length}</strong> / {pagination.total_items} sự kiện
          </div>
        </div>

        {/* Filter Controls Box */}
        <div className="space-y-4 rounded-3xl border border-[#ded4c1] bg-[#faf7f0] p-5 sm:p-6 shadow-xs">
          {/* Category Tabs */}
          <div className="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-none">
            {CATEGORIES.map((cat) => {
              const active = selectedCategory === cat;
              return (
                <button
                  key={cat}
                  onClick={() => setSelectedCategory(cat)}
                  className={`px-4 py-2 rounded-xl text-xs font-bold whitespace-nowrap transition-colors cursor-pointer ${
                    active
                      ? "bg-[#1b4332] text-[#faf7f0] shadow-xs"
                      : "bg-[#ebe4d5] text-[#524538] hover:bg-[#ded4c1] hover:text-[#2b241e]"
                  }`}
                >
                  {cat}
                </button>
              );
            })}
          </div>

          {/* Search & Sort Row */}
          <div className="flex flex-col sm:flex-row items-center gap-3 pt-1">
            <div className="relative w-full sm:flex-1">
              <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-[#857361]" />
              <Input
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Tìm theo tên sự kiện, nghệ sĩ biểu diễn, địa điểm..."
                className="pl-9 h-11 bg-[#fdfcf9] border border-[#ded4c1] rounded-xl text-sm text-[#2b241e] placeholder:text-[#857361] focus-visible:border-[#1b4332] shadow-2xs"
              />
              {searchQuery && (
                <button
                  onClick={() => setSearchQuery("")}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-[#857361] hover:text-[#2b241e]"
                >
                  <XCircle className="h-4 w-4" />
                </button>
              )}
            </div>

            <div className="flex items-center gap-2 w-full sm:w-auto text-xs font-semibold">
              {/* Only Available Toggle */}
              <Button
                variant={onlyAvailable ? "default" : "outline"}
                size="sm"
                onClick={() => setOnlyAvailable(!onlyAvailable)}
                className={`h-11 rounded-xl px-4 text-xs font-semibold cursor-pointer ${
                  onlyAvailable
                    ? "bg-[#1b4332] hover:bg-[#143225] text-[#faf7f0] border-transparent"
                    : "border-[#ded4c1] bg-[#fdfcf9] text-[#524538] hover:bg-[#ebe4d5]"
                }`}
              >
                {onlyAvailable ? "✓ Đang còn vé" : "Tất cả trạng thái"}
              </Button>

              {/* Sort By Select */}
              <select
                value={sortBy}
                onChange={(e) =>
                  setSortBy(
                    e.target.value as "default" | "price-asc" | "price-desc" | "stock"
                  )
                }
                className="h-11 rounded-xl bg-[#fdfcf9] border border-[#ded4c1] px-3.5 text-xs font-semibold text-[#524538] focus:outline-none focus:border-[#1b4332] cursor-pointer shadow-2xs"
              >
                <option value="default">Sắp xếp: Mặc định</option>
                <option value="price-asc">Giá: Thấp đến cao</option>
                <option value="price-desc">Giá: Cao đến thấp</option>
                <option value="stock">Ưu tiên sắp hết vé</option>
              </select>
            </div>
          </div>
        </div>

        {/* Tickets Grid or Empty State */}
        {filteredTickets.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
            {filteredTickets.map((ticket, idx) => (
              <TicketCard
                key={ticket.id}
                ticket={ticket}
                index={idx + 1}
                onSelect={onSelectTicket}
              />
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-16 text-center space-y-4 rounded-3xl border border-dashed border-[#c8beaf] bg-[#faf7f0] p-8">
            <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-[#ebe4d5] text-[#1b4332]">
              <Ticket className="h-7 w-7" />
            </div>
            <div className="space-y-1">
              <h3 className="text-base font-bold text-[#2b241e]">
                Không tìm thấy sự kiện phù hợp
              </h3>
              <p className="text-xs text-[#5c4e40] max-w-sm">
                Hãy thử kiểm tra lại từ khóa tìm kiếm hoặc bấm &ldquo;Tất cả&rdquo; để xem các sự kiện khác.
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                setSelectedCategory("Tất cả");
                setSearchQuery("");
                setOnlyAvailable(false);
              }}
              className="text-xs border-[#ded4c1] bg-[#fdfcf9] hover:bg-[#ebe4d5] text-[#524538]"
            >
              Đặt lại bộ lọc
            </Button>
          </div>
        )}
      </div>
    </section>
  );
}
