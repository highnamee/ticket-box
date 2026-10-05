"use client";

import Link from "next/link";
import { Ticket, ShieldCheck, Clock, QrCode, PhoneCall } from "lucide-react";

export function Footer() {
  return (
    <footer className="border-t border-[#ded4c1] bg-[#faf7f0] font-sans text-xs text-[#5c4e40]">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-14">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Brand Info */}
          <div className="space-y-3.5 md:col-span-1">
            <Link href="/" className="flex items-center gap-2 text-[#2b241e] font-bold text-base">
              <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-[#1b4332] text-[#faf7f0] shadow-xs">
                <Ticket className="h-5 w-5 text-[#e9c46a]" />
              </div>
              <span className="text-xl font-black">Ticket<span className="text-[#1b4332]">Box</span></span>
            </Link>
            <p className="text-xs text-[#6e5d4d] leading-relaxed font-sans">
              Nền tảng phân phối vé sự kiện trực tuyến chính hãng hàng đầu Việt Nam. Mang lại trải nghiệm giữ chỗ nhanh, minh bạch và an tâm tuyệt đối.
            </p>
            <div className="flex items-center gap-2 text-[#1b4332] text-xs font-semibold pt-1">
              <span className="h-2 w-2 rounded-full bg-[#1b4332] animate-pulse" />
              <span>Hệ thống mở bán trực tuyến 24/7</span>
            </div>
          </div>

          {/* Quick Menu */}
          <div className="space-y-2.5">
            <h4 className="text-[#2b241e] font-bold text-xs uppercase tracking-wider">
              Khám Phá Sự Kiện
            </h4>
            <ul className="space-y-2 text-xs">
              <li>
                <a href="#hot-events" className="hover:text-[#1b4332] transition-colors">
                  Sự kiện nổi bật
                </a>
              </li>
              <li>
                <a href="#all-tickets" className="hover:text-[#1b4332] transition-colors">
                  Đại nhạc hội &amp; Concert
                </a>
              </li>
              <li>
                <a href="#all-tickets" className="hover:text-[#1b4332] transition-colors">
                  Hội nghị công nghệ
                </a>
              </li>
              <li>
                <a href="#all-tickets" className="hover:text-[#1b4332] transition-colors">
                  Workshop &amp; Triển lãm
                </a>
              </li>
            </ul>
          </div>

          {/* Platform Commitments */}
          <div className="space-y-2.5">
            <h4 className="text-[#2b241e] font-bold text-xs uppercase tracking-wider">
              Cam Kết Người Mua
            </h4>
            <ul className="space-y-2 text-xs">
              <li className="flex items-center gap-2 text-[#4d4034]">
                <ShieldCheck className="h-4 w-4 text-[#1b4332] shrink-0" />
                <span>100% Vé gốc từ Ban Tổ Chức</span>
              </li>
              <li className="flex items-center gap-2 text-[#4d4034]">
                <Clock className="h-4 w-4 text-[#1b4332] shrink-0" />
                <span>Giữ chỗ tự động 10 phút</span>
              </li>
              <li className="flex items-center gap-2 text-[#4d4034]">
                <QrCode className="h-4 w-4 text-[#1b4332] shrink-0" />
                <span>Mã QR check-in tiện lợi</span>
              </li>
            </ul>
          </div>

          {/* Customer Support */}
          <div className="space-y-2.5">
            <h4 className="text-[#2b241e] font-bold text-xs uppercase tracking-wider">
              Trung Tâm Hỗ Trợ
            </h4>
            <p className="text-xs text-[#6e5d4d] font-sans leading-relaxed">
              Bạn cần giải đáp thắc mắc hoặc hỗ trợ hoàn đổi vé?
            </p>
            <div className="rounded-2xl bg-[#ebe4d5] border border-[#ded4c1] p-4 space-y-1">
              <div className="flex items-center gap-2 text-[#1b4332] font-black text-sm">
                <PhoneCall className="h-4 w-4 text-[#1b4332]" />
                <span>1900-TICKET</span>
              </div>
              <div className="text-[11px] text-[#6e5d4d] font-medium">Thời gian hỗ trợ: 8h00 - 22h00 hàng ngày</div>
            </div>
          </div>
        </div>

        {/* Bottom Bar */}
        <div className="mt-12 pt-6 border-t border-[#ded4c1]/60 flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-[#857361] font-sans">
          <div>© 2026 TicketBox Vietnam. Mọi quyền được bảo lưu.</div>
          <div className="flex items-center gap-6">
            <span className="hover:text-[#2b241e] cursor-pointer">Điều khoản dịch vụ</span>
            <span className="hover:text-[#2b241e] cursor-pointer">Chính sách bảo mật</span>
            <span className="hover:text-[#2b241e] cursor-pointer">Quy chế hoạt động</span>
          </div>
        </div>
      </div>
    </footer>
  );
}
