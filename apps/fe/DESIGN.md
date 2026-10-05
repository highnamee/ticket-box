# Frontend Design Guidelines: Retro Paper & Box Office

Design specifications, visual system, and UI component standards for the **TicketBox Web Application** (`apps/fe`).

---

## 🎯 1. Design Philosophy & Color System (Retro Paper Aesthetic)

TicketBox blends the tactile nostalgia of traditional paper concert admission tickets with the speed, responsiveness, and clarity of a modern web experience.

### 🎨 Curated Color Palette

| Token / Role | Hex Code | Visual Application |
| :--- | :--- | :--- |
| **Warm Ivory / Parchment** | `#f5f2eb` | Primary page canvas background, offering a soft, paper-like tactile feel. |
| **Soft Beige** | `#faf7f0` | Navigation bar background, alternating sections, and ticket reservation container. |
| **Premium Ticket Paper** | `#fdfcf9` | Ticket card surface (`TicketCard`) and electronic ticket pass (`TicketPass`), high brightness against paper borders. |
| **Paper Trim Border** | `#ded4c1` | 1px border framing cards, structural layout dividers, and secondary buttons. |
| **Vintage Forest Green** | `#1b4332` | Core brand accent, primary CTA buttons ("Đặt vé ngay", "Xác nhận"), and `✓ ĐÃ XÁC THỰC` stamp. |
| **Walnut / Espresso Brown**| `#2b241e` | Headings (h1–h3), brand wordmark, key figures, replacing harsh solid black. |
| **Warm Earth Brown** | `#5c4e40`, `#857361` | Event descriptions, date/time labels, venues, and secondary supporting copy. |
| **Warm Ochre / Amber** | `#c2841f`, `#e9c46a` | Hot ticket flames, star accents, low-stock warnings, and ticket pass highlights. |
| **Perforated Tear Line** | `#c8beaf` | Dashed border (`border-dashed`) mimicking the physical perforation between ticket body and stub. |

---

## 🎟️ 2. Authentic Vector Notch Mechanics (True Concave Cutout)

A hallmark of TicketBox tickets is the **True Concave Cutout**, avoiding pseudo-element discs or fake opaque circles pasted over borders:

1. **The `<TicketNotchDivider />` Primitive** ([`src/components/ui/ticket-notch.tsx`](file:///Users/user/Code/ticket-box/apps/fe/src/components/ui/ticket-notch.tsx)):
   * **True Concave Arc**: Drawn via vector SVG path with an exact radius of $R = 10\text{px}$.
   * **100% Vector Transparency**: The space outside the arc has zero fill (`alpha = 0`). Whatever background lies beneath the ticket card (ivory canvas, gradients, or ambient imagery) is genuinely visible through the cutout hole.
   * **Continuous Inward Stroke**: A 1px stroke (`stroke="#ded4c1"`) seamlessly follows the inward curve, aligning with subpixel precision to the outer left and right borders of the card.
   * **Perforated Tear Line**: Centered between the left and right notches is a dashed perforation line (`border-dashed border-[#c8beaf]`), recreating the authentic feel of a torn admission ticket.

2. **Integration**:
   * **Ticket Catalog & Featured Events** ([`TicketCard`](file:///Users/user/Code/ticket-box/apps/fe/src/components/landing/ticket-card.tsx)): Bridges the main event poster/details with the pricing stub and booking CTA.
   * **Electronic Admission Pass** ([`TicketModal`](file:///Users/user/Code/ticket-box/apps/fe/src/components/landing/ticket-modal.tsx)): Separates event admission details from the gate check-in stub containing the dynamic QR code and authentication stamp.

---

## 🪟 3. Wide Box-Office Modal (Ticket Booking Dialog)

1. **Spacious Desktop & Tablet Viewport**:
   * Generous dimensions: `w-[94vw] sm:w-[90vw] md:w-[840px] max-w-4xl`.
   * Adaptive scroll container: `max-h-[92vh] overflow-y-auto`, ensuring seamless readability on both laptops and larger monitors.
   * Hardcoded `sm:max-w-sm` constraints in the base Dialog component have been removed to allow natural fluid expansion.

2. **Two-Column Symmetrical Layout**:
   * **Left Column (Event Poster & Info)**: Large event banner (`h-52`), venue badge, date & time pill, and event description.
   * **Right Column (Box Office Order Slip)**:
     * User-friendly header and booking instructions.
     * Fair quota policy banner (2–5 tickets max per customer).
     * Interactive quantity stepper `[-] [ 02 ] [+]` with tabular numerals.
     * Itemized price breakdown (base price, quantity, service fee 0đ, and prominent total).
     * Forest green CTA button with a 10-minute automated reservation hold notice.

3. **Confirmed State (Electronic Ticket Pass)**:
   * Replaces the order form with an authentic paper admission ticket slip.
   * Vintage forest green ink stamp `✓ ĐÃ XÁC THỰC`.
   * Unique booking reference code (`TB-XXXXXX`).
   * Large high-contrast QR code for gate check-in.
   * Action buttons to download the digital ticket or browse other events.

---

## 🔤 4. Typography Standards (100% Clean Sans-Serif)

* **Primary Font**: `Geist Sans` (`next/font/google`) paired with system sans-serif fallbacks (`-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif`).
* **Monospace Font**: `Geist Mono` for ticket serials, event indices (`#TK-01`), dates, and currency values.
* **Strict Sans-Serif Rule**: No serif fonts (such as Times New Roman) are permitted, ensuring Vietnamese typography remains crisp, clean, and modern across all display densities.

---

## 📱 5. Responsive Design & User Experience Principles

* **Device Support**: Fully responsive across mobile (375px+), tablet (768px+), and desktop (1024px–1440px+).
* **Subtle Micro-Interactions**: Gentle card hover lift (`hover:-translate-y-1.5`), soft shadow elevation (`drop-shadow-md`), and tactile button states without aggressive neon flashes.
* **Audience-Centric Copy**: Clear, transparent, and approachable messaging focused on event discovery and fair ticket access, avoiding technical jargon.
