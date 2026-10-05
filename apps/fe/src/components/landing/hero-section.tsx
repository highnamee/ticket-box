"use client";

import { Ticket, Flame, ArrowRight, ShieldCheck, Clock, QrCode } from "lucide-react";
import { Button } from "@/components/ui/button";

export function HeroSection() {
  return (
    <section className="relative overflow-hidden py-16 md:py-24 border-b border-[#ded4c1] bg-[#f5f2eb] font-sans">
      <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8 text-center space-y-8 relative">
        {/* Retro Parchment Pill Badge */}
        <div className="inline-flex items-center gap-2 rounded-full border border-[#ded4c1] bg-[#ebe4d5] px-4 py-1.5 text-xs font-semibold text-[#1b4332] shadow-2xs">
          <Ticket className="h-3.5 w-3.5 text-[#c2841f]" />
          <span>Quầy Vé Trực Tuyến • Giữ Chỗ Tức Thì 10 Phút</span>
        </div>

        {/* Main Headline */}
        <div className="mx-auto max-w-3xl space-y-4">
          <h1 className="text-3xl sm:text-5xl md:text-6xl font-black tracking-tight text-[#2b241e] leading-tight">
            Đặt Vé Nhanh Chóng,{" "}
            <span className="text-[#1b4332] underline decoration-[#c2841f]/50 underline-offset-8">
              Trọn Vẹn Khoảnh Khắc
            </span>
          </h1>

          <p className="mx-auto max-w-2xl text-base sm:text-lg text-[#5c4e40] leading-relaxed font-sans">
            Khám phá hàng loạt concert âm nhạc đỉnh cao, hội nghị công nghệ, giải marathon và workshop nghệ thuật. 
            Giữ vé tự động 10 phút, thanh toán an toàn, nhận vé điện tử QR vào cổng tức thì.
          </p>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-col sm:flex-row items-center justify-center gap-3.5 pt-2">
          <a href="#hot-events" className="w-full sm:w-auto">
            <Button
              size="lg"
              className="w-full sm:w-auto bg-[#1b4332] hover:bg-[#143225] text-[#faf7f0] font-semibold px-8 h-12 rounded-xl shadow-xs cursor-pointer flex items-center justify-center gap-2 border border-[#143225]"
            >
              <Flame className="h-4 w-4 text-[#e9c46a]" />
              <span>Xem Vé Nổi Bật</span>
              <ArrowRight className="h-4 w-4 ml-1" />
            </Button>
          </a>

          <a href="#all-tickets" className="w-full sm:w-auto">
            <Button
              variant="outline"
              size="lg"
              className="w-full sm:w-auto border-[#ded4c1] bg-[#faf7f0] hover:bg-[#ebe4d5] text-[#4d4034] font-semibold px-6 h-12 rounded-xl shadow-2xs cursor-pointer"
            >
              <Ticket className="h-4 w-4 mr-2 text-[#c2841f]" />
              <span>Khám Phá Tất Cả Vé</span>
            </Button>
          </a>
        </div>

        {/* 3 Retro Paper Trust Pillars */}
        <div className="mx-auto max-w-4xl pt-8">
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 text-left">
            <div className="rounded-2xl border border-[#ded4c1] bg-[#faf7f0] p-5 space-y-2 shadow-2xs hover:border-[#1b4332]/50 transition-colors">
              <div className="flex items-center gap-2.5 text-[#1b4332] font-bold text-sm">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-[#ebe4d5] text-[#1b4332]">
                  <ShieldCheck className="h-4 w-4" />
                </div>
                <span>100% Vé Gốc Ban Tổ Chức</span>
              </div>
              <p className="text-xs text-[#5c4e40] leading-relaxed font-sans">
                Vé phân phối chính hãng đúng giá niêm yết, hạn chế số lượng mua mỗi người để chống đầu cơ vé.
              </p>
            </div>

            <div className="rounded-2xl border border-[#ded4c1] bg-[#faf7f0] p-5 space-y-2 shadow-2xs hover:border-[#1b4332]/50 transition-colors">
              <div className="flex items-center gap-2.5 text-[#1b4332] font-bold text-sm">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-[#ebe4d5] text-[#1b4332]">
                  <Clock className="h-4 w-4" />
                </div>
                <span>Giữ Chỗ Tự Động 10 Phút</span>
              </div>
              <p className="text-xs text-[#5c4e40] leading-relaxed font-sans">
                Sau khi chọn số lượng vé, hệ thống tự động giữ vé 10 phút để bạn hoàn tất thanh toán thong thả.
              </p>
            </div>

            <div className="rounded-2xl border border-[#ded4c1] bg-[#faf7f0] p-5 space-y-2 shadow-2xs hover:border-[#1b4332]/50 transition-colors">
              <div className="flex items-center gap-2.5 text-[#1b4332] font-bold text-sm">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-[#ebe4d5] text-[#1b4332]">
                  <QrCode className="h-4 w-4" />
                </div>
                <span>Vé QR Check-in Tiện Lợi</span>
              </div>
              <p className="text-xs text-[#5c4e40] leading-relaxed font-sans">
                Nhận vé điện tử hiển thị trực tiếp trên điện thoại, quét mã qua cổng soát vé nhanh chóng trong 1 giây.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
