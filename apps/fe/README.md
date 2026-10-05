# TicketBox Frontend

Web frontend for the **TicketBox** ticketing platform, built with [Next.js](https://nextjs.org/) 16 (App Router) and [Tailwind CSS v4](https://tailwindcss.com/).

---

## 📁 Project Structure

```text
apps/fe/
├── public/                 # Static assets (favicons, images)
├── src/
│   ├── app/                # Next.js App Router (pages, root layout, styling)
│   ├── components/         # Landing sections, layout, and UI primitives
│   ├── lib/                # API client and service helpers
│   └── types/              # TypeScript definitions matching Go backend DTOs
├── Dockerfile              # Production multi-stage Docker build
├── next.config.ts          # Next.js configuration
├── package.json            # Scripts and package dependencies
├── AGENTS.md               # Guidelines for AI assistants
└── DESIGN.md               # Visual design system and UX specifications
```

---

## 🛠️ Tech Stack

- **Framework**: [Next.js](https://nextjs.org/) 16 (App Router, Turbopack, React Server Components)
- **Library**: [React](https://react.dev/) 19 & TypeScript 5
- **Styling**: [Tailwind CSS v4](https://tailwindcss.com/) & [tw-animate-css](https://github.com/jamiebuilds/tw-animate-css)
- **Icons & Primitives**: [Lucide React](https://lucide.dev/) & [shadcn/ui](https://ui.shadcn.com/)
- **Package Manager**: [Yarn Berry](https://yarnpkg.com/) 4.x (`nodeLinker: node-modules`)

---

## 🚀 Getting Started

### Prerequisites
- **Node.js**: Version `>= 20.0.0`
- **Yarn Berry**: Version `4.x` (enabled via Corepack)

### Local Development

1. Navigate to the frontend directory:
   ```bash
   cd apps/fe
   ```

2. Enable Corepack and install dependencies:
   ```bash
   corepack enable
   yarn install
   ```

3. (Optional) Configure environment variables:
   ```bash
   cp .env.example .env.local
   ```

4. Start the development server with Turbopack:
   ```bash
   yarn dev
   ```
   Open [http://localhost:3000](http://localhost:3000) in your browser.

---

## 🛠️ Available Scripts

| Command | Description |
| :--- | :--- |
| `yarn dev` | Runs the Next.js development server with Turbopack |
| `yarn build` | Creates an optimized standalone production build |
| `yarn start` | Starts the production server locally |
| `yarn lint` | Runs ESLint checks across the codebase |

---

## 🐳 Docker Deployment

### 1. Full Stack via Docker Compose
From the repository root:
```bash
# Start backend, database, and frontend services
docker compose up -d

# View frontend logs
docker compose logs -f frontend

# Stop services
docker compose down
```

### 2. Standalone Frontend Container
```bash
docker build -t ticket-box-fe:latest apps/fe
docker run --rm -p 3000:3000 ticket-box-fe:latest
```

---

## 🔌 Go Backend Integration (`apps/be`)

The Go backend runs at `http://localhost:8080/api/v1`. 

The frontend integrates with public ticket endpoints (`/tickets`) matching the `PublicTicketResponse` DTO structure. In development, the application includes a comprehensive mock dataset (`MOCK_PUBLIC_TICKETS`) in `src/lib/api/tickets.ts` for offline testing.

To connect to a live backend instance, set `NEXT_PUBLIC_API_URL` in `.env.local`:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

---

## 📚 Documentation Links

- [DESIGN.md](./DESIGN.md): Visual design specifications, retro paper palette, vector notch mechanics, and modal UX standards.
- [AGENTS.md](./AGENTS.md): Architecture conventions and coding guidelines for AI assistants.
