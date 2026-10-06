"use client";

import { useState } from "react";
import Image from "next/image";
import {
  Calendar,
  ShieldCheck,
  CheckCircle2,
  QrCode,
  Download,
  RotateCcw,
  Clock,
  Ticket,
  Sparkles
} from "lucide-react";
import { PublicTicket } from "@/types/ticket";
import { formatVND } from "@/lib/utils/format";
import {
  Dialog,
  DialogContent,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { TicketNotchDivider } from "@/components/ui/ticket-notch";

interface TicketModalProps {
  ticket: PublicTicket | null;
  isOpen: boolean;
  onClose: () => void;
}

export function TicketModal({ ticket, isOpen, onClose }: TicketModalProps) {
  const [quantity, setQuantity] = useState(1);
  const [isProcessing, setIsProcessing] = useState(false);
  const [bookingConfirmed, setBookingConfirmed] = useState(false);
  const [bookingCode, setBookingCode] = useState<string | null>(null);

  if (!ticket) return null;

  const maxAllowed = Math.min(
    ticket.max_booking_per_user ?? 4,
    ticket.available_stock || 1
  );

  const totalPrice = ticket.price * quantity;

  const handleConfirmBooking = () => {
    setIsProcessing(true);

    setTimeout(() => {
      setIsProcessing(false);
      setBookingConfirmed(true);
      setBookingCode(`TB-${Math.floor(100000 + Math.random() * 900000)}`);
    }, 700);
  };

  const handleResetAndClose = () => {
    setBookingConfirmed(false);
    setIsProcessing(false);
    setQuantity(1);
    setBookingCode(null);
    onClose();
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && handleResetAndClose()}>
      <DialogContent className="w-[94vw] sm:w-[90vw] md:w-[840px] max-w-4xl overflow-hidden p-0 border border-[#ded4c1] bg-[#faf7f0] text-[#2b241e] shadow-2xl rounded-2xl sm:rounded-3xl font-sans max-h-[92vh] flex flex-col">
        {/* Header Bar */}
        <div className="flex items-center justify-between border-b border-[#ded4c1] bg-[#ebe4d5] px-6 py-3.5 text-xs text-[#524538] shrink-0">
          <div className="flex items-center gap-2">
            <span className="flex h-2.5 w-2.5 rounded-full bg-[#1b4332]" />
            <span className="font-bold text-[#2b241e]">
              Quầy Đặt Vé Trực Tuyến • TicketBox
            </span>
          </div>
          <div className="flex items-center gap-2 mr-8">
            <span className="rounded-md bg-[#faf7f0] border border-[#ded4c1] px-2.5 py-0.5 font-bold text-[#1b4332] shadow-2xs">
              {ticket.category || "Sự kiện"}
            </span>
          </div>
        </div>

        {bookingConfirmed ? (
          /* State 2: BOOKING CONFIRMED (VÉ ĐÃ XÁC NHẬN - PHONG CÁCH GIẤY CỔ ĐIỂN) */
          <div className="p-6 sm:p-10 space-y-6 text-center overflow-y-auto">
            <div className="space-y-1.5">
              <div className="inline-flex items-center gap-1.5 text-xs font-bold text-[#1b4332] bg-[#ebe4d5] border border-[#ded4c1] px-3.5 py-1 rounded-full">
                <CheckCircle2 className="h-4 w-4 text-[#1b4332]" />
                <span>Đặt Vé Thành Công!</span>
              </div>
              <DialogTitle className="text-2xl font-black text-[#2b241e] pt-1">
                Vé Điện Tử Của Bạn Đã Sẵn Sàng
              </DialogTitle>
              <DialogDescription className="text-xs text-[#5c4e40] max-w-md mx-auto font-sans leading-relaxed">
                Hệ thống đã giữ chỗ {quantity} vé cho bạn. Bạn có thể sử dụng mã QR dưới đây để làm thủ tục check-in tại cổng sự kiện.
              </DialogDescription>
            </div>

            {/* Retro Ticket Pass Card with True Perforated Notches */}
            <div className="relative mx-auto max-w-xl text-[#2b241e] text-left drop-shadow-xs">
              {/* Pass Top Body */}
              <div className="rounded-t-2xl border-t border-x border-[#ded4c1] bg-[#fdfcf9] p-6 space-y-4">
                <div className="flex items-center justify-between border-b border-dashed border-[#c8beaf] pb-3 text-xs">
                  <div className="flex items-center gap-2 font-bold text-[#1b4332] tracking-wide">
                    <Ticket className="h-4 w-4 text-[#c2841f]" />
                    <span>TICKETBOX • VÉ VÀO CỔNG CHÍNH THỨC</span>
                  </div>
                  <span className="text-[#524538] font-mono font-bold text-xs bg-[#ebe4d5] px-2.5 py-0.5 rounded border border-[#ded4c1]">
                    MÃ: {bookingCode}
                  </span>
                </div>

                {/* Event Name */}
                <div>
                  <div className="text-[10px] text-[#857361] uppercase tracking-wider font-bold">Tên sự kiện:</div>
                  <div className="text-lg font-bold text-[#2b241e] leading-snug">
                    {ticket.name}
                  </div>
                </div>

                {/* Time & Venue Row */}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs bg-[#ebe4d5]/40 p-3 rounded-xl border border-[#ded4c1]">
                  <div>
                    <span className="text-[11px] text-[#857361] font-semibold">Thời gian diễn ra:</span>
                    <div className="font-bold text-[#2b241e] mt-0.5">
                      {ticket.date || "15/11/2026"} • {ticket.time || "19:30"}
                    </div>
                  </div>
                  <div>
                    <span className="text-[11px] text-[#857361] font-semibold">Địa điểm tổ chức:</span>
                    <div className="font-bold text-[#2b241e] mt-0.5 truncate">
                      {ticket.venue || "TP.HCM"}
                    </div>
                  </div>
                </div>

                {/* Price & Quantity summary */}
                <div className="flex items-center justify-between text-xs pt-1 border-t border-[#ded4c1]">
                  <span>Số lượng vé: <strong className="text-[#2b241e] text-sm">{quantity} vé</strong></span>
                  <span className="text-lg font-black text-[#1b4332]">
                    {formatVND(totalPrice)}
                  </span>
                </div>
              </div>

              {/* Perforated Notch Divider (True Cutout) */}
              <TicketNotchDivider
                fillColor="#fdfcf9"
                borderColor="#ded4c1"
                dashedLineColor="#c8beaf"
                notchRadius={10}
              />

              {/* Pass Bottom Stub: Stamp and QR Code */}
              <div className="rounded-b-2xl border-b border-x border-[#ded4c1] bg-[#fdfcf9] p-6 pt-3 flex items-center justify-between">
                <div className="space-y-1.5">
                  <div className="inline-block rounded-lg border-2 border-[#1b4332] bg-[#ebe4d5] px-3 py-1 text-xs font-black text-[#1b4332] uppercase tracking-wider">
                    ✓ ĐÃ XÁC THỰC
                  </div>
                  <div className="text-[11px] text-[#5c4e40] font-semibold">Vé điện tử chính hãng TicketBox</div>
                  <div className="text-[10px] text-[#857361] font-mono">Bảo lưu giữ chỗ 10 phút</div>
                </div>

                <div className="flex flex-col items-center justify-center p-2.5 rounded-xl bg-white border border-[#ded4c1] shadow-2xs">
                  <QrCode className="h-14 w-14 text-[#2b241e]" />
                  <span className="text-[9px] text-[#857361] font-bold mt-1 tracking-wider">MÃ QUÉT TẠI CỔNG</span>
                </div>
              </div>
            </div>

            {/* Actions */}
            <div className="flex flex-col sm:flex-row gap-3 pt-2 max-w-xl mx-auto">
              <Button
                variant="outline"
                onClick={handleResetAndClose}
                className="w-full sm:w-1/2 text-xs border-[#ded4c1] bg-[#faf7f0] hover:bg-[#ebe4d5] text-[#524538] h-11 rounded-xl font-semibold"
              >
                <RotateCcw className="h-3.5 w-3.5 mr-1.5" />
                Xem sự kiện khác
              </Button>
              <Button
                onClick={() => {
                  alert("Đã lưu mã vé QR vào máy thành công!");
                  handleResetAndClose();
                }}
                className="w-full sm:w-1/2 bg-[#1b4332] hover:bg-[#143225] text-[#faf7f0] font-semibold text-xs h-11 rounded-xl shadow-xs border border-[#143225]"
              >
                <Download className="h-3.5 w-3.5 mr-1.5" />
                Lưu vé vào điện thoại
              </Button>
            </div>
          </div>
        ) : (
          /* State 1: SELECTION & ORDERING (2-COLUMN WIDE LAYOUT) */
          <div className="grid grid-cols-1 md:grid-cols-2 overflow-y-auto">
            {/* Left Column: Event Visual & Details */}
            <div className="p-6 sm:p-8 bg-[#f5f2eb] border-b md:border-b-0 md:border-r border-[#ded4c1] space-y-4">
              <div className="relative h-48 sm:h-52 w-full rounded-2xl overflow-hidden border border-[#ded4c1] bg-[#ebe4d5] shadow-2xs">
                {ticket.image_url && (
                  <Image
                    src={ticket.image_url}
                    alt={ticket.name}
                    fill
                    unoptimized
                    className="object-cover"
                  />
                )}
                <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent" />
                <div className="absolute bottom-3 left-3 right-3 text-white text-xs font-semibold drop-shadow-xs">
                  {ticket.venue || "Địa điểm cập nhật"}
                </div>
              </div>

              <div className="space-y-2">
                <DialogTitle className="text-xl font-bold text-[#2b241e] leading-snug">
                  {ticket.name}
                </DialogTitle>
                <p className="text-xs text-[#5c4e40] leading-relaxed font-sans">
                  {ticket.description}
                </p>
              </div>

              {/* Time Pill */}
              <div className="flex items-center gap-2.5 text-xs bg-[#faf7f0] p-3.5 rounded-xl border border-[#ded4c1] shadow-2xs">
                <Calendar className="h-4 w-4 text-[#1b4332] shrink-0" />
                <div>
                  <span className="text-[10px] text-[#857361] font-bold uppercase tracking-wider">Thời gian diễn ra</span>
                  <div className="font-bold text-[#2b241e] text-sm">{ticket.date || "15/11/2026"} • {ticket.time || "19:30"}</div>
                </div>
              </div>
            </div>

            {/* Right Column: Order Slip / Quầy Chọn Vé */}
            <div className="p-6 sm:p-8 flex flex-col justify-between space-y-6 bg-[#faf7f0]">
              <div className="space-y-5">
                <div className="space-y-1">
                  <h3 className="text-base font-bold text-[#2b241e]">
                    Chi Tiết Đặt Vé
                  </h3>
                  <p className="text-xs text-[#5c4e40]">
                    Chọn số lượng vé bạn muốn mua cho sự kiện này
                  </p>
                </div>

                {/* Friendly Quota notice */}
                <div className="rounded-xl border border-[#ded4c1] bg-[#ebe4d5] p-3.5 flex items-start gap-2.5 text-xs text-[#524538]">
                  <ShieldCheck className="h-4 w-4 text-[#1b4332] shrink-0 mt-0.5" />
                  <div className="space-y-0.5">
                    <span className="font-bold text-[#2b241e]">Quy định hạn mức vé:</span>
                    <p className="text-[11px] text-[#5c4e40] leading-relaxed font-sans">
                      Mỗi khách hàng được đặt tối đa <strong className="text-[#1b4332]">{ticket.max_booking_per_user ?? 4} vé</strong> để đảm bảo sự công bằng.
                    </p>
                  </div>
                </div>

                {/* Quantity Selector */}
                <div className="flex items-center justify-between border-t border-[#ded4c1] pt-4">
                  <div>
                    <div className="text-xs font-bold text-[#2b241e]">Số lượng vé</div>
                    <div className="text-[11px] text-[#857361] font-sans">
                      Còn {ticket.available_stock} vé khả dụng
                    </div>
                  </div>

                  <div className="flex items-center gap-3">
                    <Button
                      variant="outline"
                      size="icon-sm"
                      onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                      disabled={quantity <= 1 || isProcessing}
                      className="h-9 w-9 rounded-xl border-[#ded4c1] bg-[#fdfcf9] text-[#2b241e] hover:bg-[#ebe4d5] text-lg font-bold cursor-pointer"
                    >
                      -
                    </Button>
                    <div className="w-10 text-center font-bold text-base text-[#2b241e] tabular-nums">
                      {quantity}
                    </div>
                    <Button
                      variant="outline"
                      size="icon-sm"
                      onClick={() => setQuantity((q) => Math.min(maxAllowed, q + 1))}
                      disabled={quantity >= maxAllowed || isProcessing}
                      className="h-9 w-9 rounded-xl border-[#ded4c1] bg-[#fdfcf9] text-[#2b241e] hover:bg-[#ebe4d5] text-lg font-bold cursor-pointer"
                    >
                      +
                    </Button>
                  </div>
                </div>

                {/* Price Breakdown */}
                <div className="rounded-xl bg-[#ebe4d5]/50 p-4 border border-[#ded4c1] space-y-2 text-xs">
                  <div className="flex items-center justify-between text-[#5c4e40]">
                    <span>Đơn giá:</span>
                    <span className="font-semibold text-[#2b241e]">{formatVND(ticket.price)}</span>
                  </div>
                  <div className="flex items-center justify-between text-[#5c4e40]">
                    <span>Số lượng:</span>
                    <span className="font-semibold text-[#2b241e]">{quantity} vé</span>
                  </div>
                  <div className="flex items-center justify-between text-[#5c4e40]">
                    <span>Phí dịch vụ:</span>
                    <span className="font-semibold text-[#1b4332]">0đ (Miễn phí)</span>
                  </div>
                  <div className="border-t border-[#ded4c1] pt-2 flex items-center justify-between text-sm">
                    <span className="font-bold text-[#2b241e]">Tổng thanh toán:</span>
                    <span className="text-xl font-black text-[#1b4332] tabular-nums">
                      {formatVND(totalPrice)}
                    </span>
                  </div>
                </div>
              </div>

              {/* Action Button */}
              <div className="pt-2">
                <Button
                  onClick={handleConfirmBooking}
                  disabled={ticket.is_sold_out || isProcessing}
                  className="w-full bg-[#1b4332] hover:bg-[#143225] text-[#faf7f0] font-semibold text-sm h-12 rounded-xl shadow-sm border border-[#143225] cursor-pointer flex items-center justify-center gap-2"
                >
                  <Sparkles className="h-4 w-4 text-[#e9c46a]" />
                  <span>{isProcessing ? "Đang xử lý giữ vé..." : "Xác nhận đặt vé ngay"}</span>
                </Button>
                <div className="flex items-center justify-center gap-1.5 text-[11px] text-[#857361] mt-2 font-medium">
                  <Clock className="h-3.5 w-3.5 text-[#1b4332]" />
                  <span>Hệ thống sẽ giữ chỗ 10 phút sau khi xác nhận</span>
                </div>
              </div>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
