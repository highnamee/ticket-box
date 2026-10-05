export type TicketStatus = "ACTIVE" | "INACTIVE" | "ARCHIVED";

export interface PublicTicket {
  id: string;
  name: string;
  description: string;
  price: number;
  available_stock: number;
  max_booking_per_user?: number | null;
  status: TicketStatus;
  is_sold_out: boolean;
  // Enriched presentation fields (parsed/derived or provided)
  category?: string;
  venue?: string;
  date?: string;
  time?: string;
  image_url?: string;
  tags?: string[];
  featured?: boolean;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total_items: number;
  total_pages: number;
}

export interface ApiResponse<T> {
  success: boolean;
  message?: string;
  data: T;
  error?: unknown;
}

export interface PaginatedTicketsData {
  items: PublicTicket[];
  pagination: PaginationMeta;
}
