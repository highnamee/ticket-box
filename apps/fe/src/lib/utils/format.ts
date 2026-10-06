/**
 * Utility functions for formatting currencies and numbers.
 */

/**
 * Format numerical amount into Vietnamese Dong currency string (e.g. 850.000 ₫).
 */
export function formatVND(amount: number): string {
  return new Intl.NumberFormat("vi-VN", {
    style: "currency",
    currency: "VND",
  }).format(amount);
}
