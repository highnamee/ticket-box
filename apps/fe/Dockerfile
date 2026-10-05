# ==========================================
# Stage 1: Base Image
# ==========================================
FROM node:22-alpine AS base

# Install libc6-compat for alpine compatibility if needed
RUN apk add --no-cache libc6-compat
RUN corepack enable && corepack prepare yarn@4.18.1 --activate
WORKDIR /app

# ==========================================
# Stage 2: Install dependencies with Yarn Berry
# ==========================================
FROM base AS deps
COPY .yarnrc.yml package.json yarn.lock ./
RUN yarn install --immutable

# ==========================================
# Stage 3: Build Application
# ==========================================
FROM base AS builder
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

RUN yarn build

# ==========================================
# Stage 4: Production Runner
# ==========================================
FROM base AS runner

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Create non-root system user for security
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy static assets and standalone bundle
COPY --from=builder /app/public ./public

# Set proper permissions for prerender cache
RUN mkdir .next && chown nextjs:nodejs .next

COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

CMD ["node", "server.js"]
