# Frontend Design Guidelines

Initial design notes and UI guidelines for the **Ticket Box** frontend. This document will be expanded as application features develop.

---

## 🎯 1. Design Direction

- **Theme**: Modern, clean, dark-mode first with high-contrast accessibility.
- **Aesthetic**: Minimalist surfaces, subtle borders, and clear hierarchy tailored for event ticketing.
- **Mobile-First**: Fully responsive layouts adaptable across mobile, tablet, and desktop screens.

---

## 🎨 2. Color Foundation

Theming is powered by Tailwind CSS v4 using semantic tokens in the **OKLCH** color space (configured in `src/app/globals.css`):

- **Background & Card**: `bg-background`, `bg-card`, `bg-muted`
- **Text**: `text-foreground`, `text-muted-foreground`, `text-primary`
- **Borders**: `border-border`, `ring-ring`
- **Actions & Status**: `bg-primary`, `bg-secondary`, `bg-destructive`

---

## 🔤 3. Typography

Loaded via `next/font/google`:

- **Geist Sans**: Primary typeface for all UI text, headings, buttons, and form labels.
- **Geist Mono**: Monospace typeface for digital ticket codes, identifiers, and timestamps.

---

## 🧩 4. Component Strategy

- Use **shadcn/ui** primitives added on-demand via `npx shadcn@latest add <component>`.
- Keep component styles consistent with Tailwind CSS v4 variables rather than arbitrary hex values.
