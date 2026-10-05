import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "TicketBox - Nền Tảng Đặt Vé Thông Minh & Chống Phe Vé",
  description:
    "Nền tảng phân phối vé điện tử thế hệ mới cho concert, hội nghị và sự kiện thể thao với hệ thống Smart Queue và Dynamic QR.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html
      lang="vi"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased scroll-smooth`}
    >
      <body className="min-h-full flex flex-col bg-[#f5f2eb] text-[#2b241e] font-sans">
        {children}
      </body>
    </html>
  );
}
