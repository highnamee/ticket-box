import { ShieldCheck, Clock, QrCode, Users } from "lucide-react";

export function FeaturesSection() {
  const commitments = [
    {
      icon: ShieldCheck,
      title: "100% Vé Gốc Ban Tổ Chức",
      description:
        "Tất cả vé đều được phân phối trực tiếp từ ban tổ chức sự kiện. Cam kết đúng giá niêm yết, bảo vệ người hâm mộ khỏi tình trạng phe vé thổi giá.",
    },
    {
      icon: Clock,
      title: "Giữ Chỗ Tự Động 10 Phút",
      description:
        "Khi bạn chọn số lượng vé, hệ thống sẽ tự động giữ chỗ cho bạn trong vòng 10 phút. Bạn hoàn toàn thong thả thực hiện thanh toán mà không sợ bị người khác tranh mất.",
    },
    {
      icon: QrCode,
      title: "Vé Điện Tử Check-in Nhanh Chóng",
      description:
        "Không lo quên hoặc thất lạc vé giấy. Mã QR điện tử tích hợp sẵn trên điện thoại, nhân viên tại cổng chỉ cần quét mã là bạn có thể vào cổng ngay.",
    },
    {
      icon: Users,
      title: "Cơ Hội Công Bằng Cho Mọi Người",
      description:
        "Áp dụng hạn mức tối đa từ 2 đến 4 vé cho mỗi tài khoản, đảm bảo vé đến tay những người thực sự yêu thích sự kiện chứ không bị gom hàng đầu cơ.",
    },
  ];

  return (
    <section id="why-choose-us" className="scroll-mt-20 py-16 md:py-24 border-b border-[#ded4c1] bg-[#faf7f0] font-sans">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 space-y-12">
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto space-y-2.5">
          <span className="text-xs font-bold text-[#1b4332] uppercase tracking-wider">
            An Tâm Tuyệt Đối
          </span>
          <h2 className="text-2xl sm:text-3xl font-extrabold tracking-tight text-[#2b241e]">
            Vì Sao Bạn Nên Chọn TicketBox?
          </h2>
          <p className="text-sm text-[#5c4e40] max-w-xl mx-auto">
            Chúng tôi xây dựng nền tảng để mang lại trải nghiệm mua vé dễ dàng, văn minh và tiện lợi nhất cho bạn.
          </p>
        </div>

        {/* 4 User Benefit Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {commitments.map((item, idx) => {
            const Icon = item.icon;
            return (
              <div
                key={idx}
                className="rounded-3xl border border-[#ded4c1] bg-[#fdfcf9] p-6 space-y-3.5 hover:border-[#1b4332]/50 hover:shadow-xs transition-all"
              >
                <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-[#ebe4d5] border border-[#ded4c1] text-[#1b4332] shadow-2xs">
                  <Icon className="h-6 w-6" />
                </div>

                <h3 className="text-base font-bold text-[#2b241e]">
                  {item.title}
                </h3>

                <p className="text-xs text-[#6e5d4d] leading-relaxed font-sans">
                  {item.description}
                </p>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
