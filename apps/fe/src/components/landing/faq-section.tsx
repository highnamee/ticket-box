"use client";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

export function FaqSection() {
  const faqs = [
    {
      q: "1. Tôi có thể mua tối đa bao nhiêu vé cho một sự kiện?",
      a: "Tùy thuộc vào quy định của từng ban tổ chức, mỗi tài khoản thông thường được đặt tối đa từ 2 đến 4 vé cho một sự kiện nhằm đảm bảo sự công bằng cho tất cả khán giả.",
    },
    {
      q: "2. Thời gian giữ chỗ vé 10 phút tính như thế nào?",
      a: "Ngay khi bạn bấm nút 'Xác nhận đặt vé', hệ thống sẽ tạm khóa và giữ vé cho bạn trong vòng 10 phút. Bạn có thể thong thả thanh toán mà không sợ người khác mua mất. Sau 10 phút nếu chưa hoàn tất, vé sẽ tự động quay lại kho mở bán.",
    },
    {
      q: "3. Làm sao để sử dụng vé khi đến cổng sự kiện?",
      a: "Sau khi đặt vé thành công, bạn sẽ nhận được mã QR điện tử. Bạn chỉ cần mở mã QR này trên điện thoại để nhân viên tại cổng sự kiện quét mã check-in trong vòng 1 giây, rất tiện lợi và không lo quên vé giấy.",
    },
    {
      q: "4. Nền tảng hỗ trợ những hình thức thanh toán nào?",
      a: "Hệ thống hỗ trợ quét mã VietQR tự động qua tất cả các ứng dụng ngân hàng tại Việt Nam (xác nhận giao dịch tức thì 24/7), thẻ tín dụng/ghi nợ quốc tế (Visa, Mastercard) và các ví điện tử phổ biến.",
    },
    {
      q: "5. Nếu tôi cần trợ giúp khẩn cấp thì liên hệ qua kênh nào?",
      a: "Đội ngũ chăm sóc khách hàng của TicketBox luôn sẵn sàng hỗ trợ bạn qua Hotline 1900-TICKET (từ 8h00 đến 22h00 hàng ngày) hoặc qua email support@ticketbox.vn.",
    },
  ];

  return (
    <section id="faq" className="scroll-mt-20 py-16 md:py-24 border-b border-[#ded4c1] bg-[#f5f2eb] font-sans">
      <div className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8 space-y-8">
        <div className="text-center space-y-2">
          <span className="text-xs font-bold text-[#1b4332] uppercase tracking-wider">
            Giải Đáp Thắc Mắc
          </span>
          <h2 className="text-2xl sm:text-3xl font-extrabold tracking-tight text-[#2b241e]">
            Câu Hỏi Thường Gặp
          </h2>
          <p className="text-sm text-[#5c4e40]">
            Những thông tin cần biết trước khi đặt vé tại TicketBox
          </p>
        </div>

        <div className="rounded-3xl border border-[#ded4c1] bg-[#faf7f0] p-5 sm:p-8 shadow-xs">
          <Accordion defaultValue={["item-0"]} className="w-full space-y-2">
            {faqs.map((faq, idx) => (
              <AccordionItem
                key={idx}
                value={`item-${idx}`}
                className="border-b border-[#ded4c1]/60 last:border-none py-1"
              >
                <AccordionTrigger className="text-left text-sm sm:text-base font-bold text-[#2b241e] hover:text-[#1b4332] hover:no-underline py-3.5">
                  {faq.q}
                </AccordionTrigger>
                <AccordionContent className="text-xs sm:text-sm text-[#5c4e40] leading-relaxed pb-3.5 font-sans">
                  {faq.a}
                </AccordionContent>
              </AccordionItem>
            ))}
          </Accordion>
        </div>
      </div>
    </section>
  );
}
