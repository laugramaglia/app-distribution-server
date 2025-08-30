# --- Build Stage ---
# Use a specific version of Node.js for reproducibility.
FROM node:20-alpine AS build

# Set the working directory inside the container.
WORKDIR /usr/src/app

# You can pass the backend API URL as a build argument from docker-compose.
# This example shows how you might expose it for a Vite-based app.
# You may need to adjust the ENV variable name for your specific framework.
ARG API_URL
ENV NEXT_PUBLIC_API_URL=${API_URL}

# Copy package.json and the lock file to leverage Docker's layer caching.
COPY package*.json ./

# Install dependencies using 'npm ci' for faster, more reliable builds in CI/CD.
RUN npm ci

# Copy the rest of the application source code.
COPY . .

# Build the application for production. 'build' is the conventional script name.
RUN npm run build

# --- Production Stage ---
# Use a lightweight Node.js image to run the app.
FROM node:20-alpine

WORKDIR /usr/src/app

# Copy the built app from the build stage
COPY --from=build /usr/src/app/.next ./.next
COPY --from=build /usr/src/app/public ./public
COPY --from=build /usr/src/app/package.json ./package.json
COPY --from=build /usr/src/app/next.config.ts ./next.config.ts

# Expose port 3000 for the Next.js app. The docker-compose will map it to 8080.
EXPOSE 3000

# The command to start the app.
CMD ["npm", "start"]
