# ─── Stage 1: Build ───────────────────────────────────────────────────────────
FROM node:22-alpine AS builder

# Enable corepack so pnpm is available without a separate install step
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app

# Copy dependency manifests first for better layer caching
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./

# Install all dependencies (including devDependencies needed for the build)
RUN pnpm install --frozen-lockfile

# Copy the rest of the source
COPY . .

# CI/CD passes NEXT_PUBLIC_API_URL (legacy name from the old Next.js era).
# Vite reads VITE_API_URL from vite.config.ts, so we accept both and fall back
# gracefully — the vite.config.ts default is https://backend.pharma-hub.ru.
ARG NEXT_PUBLIC_API_URL
ARG VITE_API_URL=${NEXT_PUBLIC_API_URL}
ENV VITE_API_URL=${VITE_API_URL}

RUN pnpm run build

# ─── Stage 2: Serve ───────────────────────────────────────────────────────────
FROM nginx:1.27-alpine AS runner

# Remove the default nginx site
RUN rm -rf /usr/share/nginx/html/*

# Copy the built assets from the builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Nginx config: serve index.html for every route (SPA fallback)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
