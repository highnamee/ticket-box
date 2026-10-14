import { ApiResponse, PaginatedTicketsData } from "@/types/ticket";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

// Fetches public tickets from Go Backend API (/api/v1/tickets)
export async function getPublicTickets(
  page: number = 1,
  limit: number = 12
): Promise<PaginatedTicketsData> {
  const fallback: PaginatedTicketsData = {
    items: [],
    pagination: {
      page,
      limit,
      total_items: 0,
      total_pages: 0,
    },
  };

  try {
    const res = await fetch(`${API_BASE_URL}/tickets?page=${page}&limit=${limit}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      next: { revalidate: 60 },
    });

    if (!res.ok) {
      return fallback;
    }

    const json: ApiResponse<PaginatedTicketsData> = await res.json();
    if (json.success && json.data) {
      return json.data;
    }

    return fallback;
  } catch (error) {
    console.error("Failed to fetch public tickets:", error);
    return fallback;
  }
}
