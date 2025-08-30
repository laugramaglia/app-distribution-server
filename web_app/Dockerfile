# --- Build Stage ---
FROM node:20-alpine AS build

WORKDIR /usr/src/app

ARG API_URL
ENV NEXT_PUBLIC_API_URL=${API_URL}

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# --- Production Stage ---
FROM node:20-alpine AS production

WORKDIR /usr/src/app

# Copia solo lo necesario para ejecutar
COPY --from=build /usr/src/app/.next ./.next
COPY --from=build /usr/src/app/node_modules ./node_modules
COPY --from=build /usr/src/app/package.json ./package.json
COPY --from=build /usr/src/app/next.config.ts ./next.config.ts
# COPY --from=build /usr/src/app/public ./public   # si existe

EXPOSE 3000
CMD ["npm", "start"]
