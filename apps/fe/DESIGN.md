# Frontend Design Guidelines: Retro Paper & Box Office

Tài liệu thiết kế, quy chuẩn thị giác và hệ thống thành phần giao diện cho **TicketBox Web Application** (`apps/fe`).

---

## 🎯 1. Triết Lý Thiết Kế & Bảng Màu (Retro Paper Aesthetic)

TicketBox kết hợp giữa cảm xúc hoài niệm của những tấm vé xem ca nhạc bằng giấy truyền thống (vintage admission ticket) và sự tiện dụng, tốc độ của nền tảng web hiện đại.

### 🎨 Bảng Màu Chủ Đạo (Curated Earth & Paper Palette)

| Phân Loại | Mã Màu Hex | Ứng Dụng Trong Giao Diện |
| :--- | :--- | :--- |
| **Giấy Mộc Ngà** | `#f5f2eb` | Màu nền canvas chính toàn trang, tạo cảm giác dịu mắt và chất liệu giấy cổ điển. |
| **Parchment Nhạt** | `#faf7f0` | Nền thanh điều hướng (navbar), các section xen kẽ, container phiếu đặt vé. |
| **Giấy Vé Thượng Hạng** | `#fdfcf9` | Nền thẻ vé (`TicketCard`), phiếu vé điện tử (`TicketPass`), độ sáng cao tương phản với viền giấy. |
| **Viền Giấy Cắt** | `#ded4c1` | Đường viền 1px bao quanh thẻ, các phân cách layout, nút bấm phụ. |
| **Xanh Rừng Vintage** | `#1b4332` | Màu thương hiệu chủ đạo, nút CTA chính (Đặt vé ngay, Xác nhận), con dấu `✓ ĐÃ XÁC THỰC`. |
| **Nâu Gỗ Trầm Ấm** | `#2b241e` | Màu chữ chính, tiêu đề h1-h3, nhãn thương hiệu và biểu tượng, thay thế màu đen thuần. |
| **Nâu Đất Thứ Cấp** | `#5c4e40`, `#857361` | Mô tả sự kiện, chú thích ngày giờ, địa điểm, văn bản thứ cấp. |
| **Vàng Hổ Phách / Mù Tạt** | `#c2841f`, `#e9c46a` | Biểu tượng ngọn lửa vé hot, ngôi sao điểm xuyết, nhãn số lượng vé sắp hết. |
| **Đường Răng Cưa Xé Vé** | `#c8beaf` | Nét đứt (`border-dashed`) mô phỏng rãnh xé vé thật giữa thân vé và cuống vé. |

---

## 🎟️ 2. Cơ Chế Khuyết Nửa Tròn Thật (Authentic Vector Notch Cutouts)

Điểm nhấn độc đáo của thẻ vé và cuống vé TicketBox là **vết khuyết nửa tròn thật (True Concave Cutout)** thay vì dùng các hình tròn màu giả lập dán đè lên:

1. **Thành phần `<TicketNotchDivider />`** ([`src/components/ui/ticket-notch.tsx`](file:///Users/user/Code/ticket-box/apps/fe/src/components/ui/ticket-notch.tsx)):
   * Sử dụng vector path SVG với bán kính cung tròn chính xác `R = 10px`.
   * Vùng nằm ngoài cung tròn là **không gian vector trong suốt 100% (alpha = 0)**. Khi thẻ đặt trên bất kỳ bề mặt nền nào (canvas mộc ngà, gradient, hay ảnh nền), nền phía dưới đều hiển thị xuyên qua lỗ khuyết một cách chân thực.
   * Viền 1px (`stroke="#ded4c1"`) uốn cong liên tục theo đường khuyết, kết nối chuẩn xác từng subpixel với viền trái/phải của thân vé trên và cuống vé dưới.
   * Ở giữa hai bên rãnh khuyết là đường nét đứt (`border-dashed border-[#c8beaf]`) tái hiện đường răng cưa xé vé của các quầy vé kinh điển.

2. **Ứng Dụng**:
   * **Thẻ vé danh mục & sự kiện hot** ([`TicketCard`](file:///Users/user/Code/ticket-box/apps/fe/src/components/landing/ticket-card.tsx)): Nối liền giữa phần thông tin sự kiện và phần cuống giá vé + nút CTA.
   * **Vé điện tử xác nhận** ([`TicketModal`](file:///Users/user/Code/ticket-box/apps/fe/src/components/landing/ticket-modal.tsx)): Nối liền giữa thông tin sự kiện đã đặt và cuống kiểm soát có mã QR + con dấu mộc.

---

## 🪟 3. Quầy Đặt Vé Rộng Rãi (Ticket Box Office Modal)

1. **Kích Thước Chuẩn Desktop & Tablet**:
   * Độ rộng thoải mái: `w-[94vw] sm:w-[90vw] md:w-[840px] max-w-4xl`.
   * Chiều cao thông minh: `max-h-[92vh] overflow-y-auto` đảm bảo hiển thị hoàn hảo trên cả laptop nhỏ lẫn màn hình lớn.
   * Đã gỡ bỏ giới hạn hẹp `sm:max-w-sm` trong base dialog của shadcn để modal mở rộng tự nhiên.

2. **Bố Cục 2 Cột Đối Xứng**:
   * **Cột Trái (Poster & Thông Tin Sự Kiện)**: Hình ảnh sự kiện lớn (`h-52`), overlay nhẹ, huy hiệu ngày giờ và địa điểm tổ chức rõ ràng.
   * **Cột Phải (Phiếu Đặt Vé Tại Quầy)**:
     * Tiêu đề & hướng dẫn thân thiện với người dùng.
     * Hộp nhắc nhở hạn mức công bằng (tối đa 2-5 vé/người).
     * Bộ chọn số lượng vé `[-] [ 02 ] [+]` dạng nút nổi với số vé hiển thị dạng tabular-nums.
     * Bảng tính tiền chi tiết (đơn giá, số lượng, phí dịch vụ 0đ, tổng tiền nổi bật).
     * Nút bấm xác nhận màu xanh rừng (`#1b4332`) với lời nhắc bảo lưu vé 10 phút.

3. **Trạng Thái Vé Xác Nhận (Confirmed Ticket Pass)**:
   * Chuyển đổi thành phiếu vé điện tử phong cách retro giấy mộc.
   * Con dấu mộc chữ nhật `✓ ĐÃ XÁC THỰC` màu xanh rừng.
   * Mã đặt vé độc bản (`TB-XXXXXX`).
   * Mã QR kích thước lớn để nhân viên quét tại cổng sự kiện.
   * Nút tải vé về máy hoặc tiếp tục khám phá sự kiện khác.

---

## 🔤 4. Quy Chuẩn Typography (100% Sans-serif)

* **Phông chữ chủ đạo**: `Geist Sans` (`next/font/google`) kết hợp fallbacks: `-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif`.
* **Phông chữ mã số**: `Geist Mono` dùng cho mã vé, số thứ tự sự kiện (`#TK-01`), thời gian và số tiền VND.
* **Nguyên tắc**: Tuyệt đối **không sử dụng font có chân (serif)** như Times New Roman, đảm bảo tiếng Việt hiển thị hiện đại, thanh thoát, sắc nét trên mọi độ phân giải màn hình.

---

## 📱 5. Tính Tương Thích & Trải Nghiệm Người Dùng (UX Commitments)

* **Thiết bị**: Hoàn toàn responsive trên Mobile (375px+), Tablet (768px+) và Desktop (1024px - 1440px+).
* **Hiệu ứng vi tương tác**: Thẻ vé nhấc nhẹ khi di chuột (`hover:-translate-y-1.5`), bóng đổ mềm (`drop-shadow-md`), viền sáng nhẹ, không dùng hiệu ứng chớp nháy neon.
* **Ngôn từ thân thiện**: Hướng tới khán giả yêu nghệ thuật và sự kiện, loại bỏ thuật ngữ kỹ thuật thừa (UUID v7, kiosk terminal specs).
