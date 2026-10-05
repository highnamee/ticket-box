import { PaginatedTicketsData, PublicTicket } from "@/types/ticket";

// MOCK_PUBLIC_TICKETS follows Go Backend PublicTicketResponse structure strictly:
// ID (UUID v7), Name, Description, Price, AvailableStock, MaxBookingPerUser, Status, IsSoldOut
export const MOCK_PUBLIC_TICKETS: PublicTicket[] = [
  {
    id: "019468e2-63b7-7bb2-8069-b570e3037f01",
    name: "CyberSound Arena: Electro Symphony 2026",
    description:
      "Đại tiệc âm nhạc điện tử kết hợp dàn nhạc giao hưởng 60 nhạc công với sân khấu visual Hologram 360 độ hoành tráng bậc nhất Đông Nam Á.",
    price: 850000,
    available_stock: 45,
    max_booking_per_user: 4,
    status: "ACTIVE",
    is_sold_out: false,
    category: "Âm nhạc",
    venue: "Sân vận động Quân khu 7, TP. Hồ Chí Minh",
    date: "15/11/2026",
    time: "19:30",
    image_url:
      "https://images.unsplash.com/photo-1470225620780-dba8ba36b745?auto=format&fit=crop&w=1200&q=80",
    tags: ["Hot", "Music", "Hologram"],
    featured: true,
  },
  {
    id: "019468e2-63b7-7bb2-8069-b570e3037f02",
    name: "Vietnam AI & Cloud Summit 2026",
    description:
      "Diễn đàn công nghệ hàng đầu hội tụ hơn 2.500 kỹ sư, nhà sáng lập và chuyên gia AI từ Google, OpenAI, AWS chia sẻ kiến trúc Agentic AI & High-Load Backend.",
    price: 1200000,
    available_stock: 12,
    max_booking_per_user: 2,
    status: "ACTIVE",
    is_sold_out: false,
    category: "Công nghệ",
    venue: "Trung tâm Hội nghị Quốc gia, Hà Nội",
    date: "28/11/2026",
    time: "08:30",
    image_url:
      "https://images.unsplash.com/photo-1540575467063-178a50c2df87?auto=format&fit=crop&w=1200&q=80",
    tags: ["Tech", "VIP", "Limited"],
    featured: true,
  },
  {
    id: "019468e2-63b7-7bb2-8069-b570e3037f03",
    name: "Indie Acoustic Night: Thanh Âm Mùa Thu",
    description:
      "Đêm nhạc acoustic ấm cúng với những nghệ sĩ indie được yêu mến. Trải nghiệm không gian mộc mạc, gần gũi với đồ uống thủ công tặng kèm.",
    price: 350000,
    available_stock: 80,
    max_booking_per_user: 5,
    status: "ACTIVE",
    is_sold_out: false,
    category: "Âm nhạc",
    venue: "Soul Live Project Complex, TP. Hồ Chí Minh",
    date: "05/12/2026",
    time: "20:00",
    image_url:
      "https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?auto=format&fit=crop&w=1200&q=80",
    tags: ["Acoustic", "Chill", "Indie"],
    featured: false,
  },
  {
    id: "019468e2-63b7-7bb2-8069-b570e3037f04",
    name: "Masterclass: High-Performance Go Microservices",
    description:
      "Khóa huấn luyện chuyên sâu 1 ngày về tối ưu hóa concurrency, zero-allocation buffers, Redis caching và xử lý hàng triệu giao dịch vé mỗi giây.",
    price: 2500000,
    available_stock: 8,
    max_booking_per_user: 1,
    status: "ACTIVE",
    is_sold_out: false,
    category: "Workshop",
    venue: "Dreamplex Tech Hub, Quận 1, TP. Hồ Chí Minh",
    date: "12/12/2026",
    time: "09:00",
    image_url:
      "https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?auto=format&fit=crop&w=1200&q=80",
    tags: ["Golang", "Masterclass", "Sắp hết"],
    featured: true,
  },
  {
    id: "019468e2-63b7-7bb2-8069-b570e3037f05",
    name: "Hanoi International Marathon: Đêm Di Sản 2026",
    description:
      "Cung đường chạy đêm đi qua các di tích lịch sử Hồ Gươm, Cầu Long Biên, Hoàng thành Thăng Long với huy chương mạ vàng độc bản.",
    price: 650000,
    available_stock: 120,
    max_booking_per_user: 4,
    status: "ACTIVE",
    is_sold_out: false,
    category: "Thể thao",
    venue: "Quảng trường Đông Kinh Nghĩa Thục, Hà Nội",
    date: "20/12/2026",
    time: "23:00",
    image_url:
      "https://images.unsplash.com/photo-1530549387789-4c1017266635?auto=format&fit=crop&w=1200&q=80",
    tags: ["Sports", "Night Run", "Medal"],
    featured: false,
  },
  {
    id: "019468e2-63b7-7bb2-8069-b570e3037f06",
    name: "Art & Light Exhibition: Lạc Vào Hư Ảo",
    description:
      "Không gian triển lãm nghệ thuật thị giác kỹ thuật số tương tác với hơn 15 phòng trải nghiệm ánh sáng đa chiều và âm thanh vòm binaural.",
    price: 290000,
    available_stock: 0,
    max_booking_per_user: 6,
    status: "ACTIVE",
    is_sold_out: true,
    category: "Triển lãm",
    venue: "Trung tâm Nghệ thuật Đương đại VCCA, Hà Nội",
    date: "Hằng ngày (Đến 31/12/2026)",
    time: "10:00 - 21:00",
    image_url:
      "https://images.unsplash.com/photo-1508997449629-303059a039c0?auto=format&fit=crop&w=1200&q=80",
    tags: ["Art", "Light", "Sold Out"],
    featured: false,
  },
  {
    id: "019468e2-63b7-7bb2-8069-b570e3037f07",
    name: "Saigon Comedy Club: Cười Xuyên Đêm",
    description:
      "Đêm hài độc thoại đỉnh cao cùng 6 diễn viên hài tài năng hàng đầu Việt Nam. Tặng kèm 01 phần thức uống và bắp rang bơ cao cấp.",
    price: 320000,
    available_stock: 28,
    max_booking_per_user: 4,
    status: "ACTIVE",
    is_sold_out: false,
    category: "Hài kịch",
    venue: "Nhà hát Kịch TP.HCM, Quận 1, TP. Hồ Chí Minh",
    date: "25/12/2026",
    time: "20:00",
    image_url:
      "https://images.unsplash.com/photo-1585699324551-f6c309eedeca?auto=format&fit=crop&w=1200&q=80",
    tags: ["Comedy", "Weekend", "Fun"],
    featured: false,
  },
  {
    id: "019468e2-63b7-7bb2-8069-b570e3037f08",
    name: "Cinema Gala Premiere: Vũ Trụ Điện Ảnh 2026",
    description:
      "Suất chiếu đặc biệt thảm đỏ đầu tiên tại Việt Nam trên màn hình IMAX Laser thế hệ mới, giao lưu cùng đạo diễn và dàn diễn viên chính.",
    price: 450000,
    available_stock: 18,
    max_booking_per_user: 2,
    status: "ACTIVE",
    is_sold_out: false,
    category: "Điện ảnh",
    venue: "Cụm rạp IMAX Landmark 81, TP. Hồ Chí Minh",
    date: "30/12/2026",
    time: "18:45",
    image_url:
      "https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?auto=format&fit=crop&w=1200&q=80",
    tags: ["Cinema", "IMAX", "Red Carpet"],
    featured: true,
  },
];

// Returns MOCK_PUBLIC_TICKETS directly following Go BE structure (BE support will plug in later)
export async function getPublicTickets(
  page: number = 1,
  limit: number = 12
): Promise<PaginatedTicketsData> {
  const startIndex = (page - 1) * limit;
  const paginatedItems = MOCK_PUBLIC_TICKETS.slice(startIndex, startIndex + limit);

  return {
    items: paginatedItems,
    pagination: {
      page,
      limit,
      total_items: MOCK_PUBLIC_TICKETS.length,
      total_pages: Math.ceil(MOCK_PUBLIC_TICKETS.length / limit),
    },
  };
}

export function formatVND(amount: number): string {
  return new Intl.NumberFormat("vi-VN", {
    style: "currency",
    currency: "VND",
  }).format(amount);
}
