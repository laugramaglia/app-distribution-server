# --- Build Stage ---
# Use a specific version of Node.js for reproducibility.
FROM node:20-alpine AS build

# Set the working directory inside the container.
WORKDIR /usr/src/app

# You can pass the backend API URL as a build argument from docker-compose.
# This example shows how you might expose it for a Vite-based app.
# You may need to adjust the ENV variable name for your specific framework.
ARG API_URL
ENV VITE_API_URL=${API_URL}

# Copy package.json and the lock file to leverage Docker's layer caching.
COPY package*.json ./

# Install dependencies using 'npm ci' for faster, more reliable builds in CI/CD.
RUN npm ci

# Copy the rest of the application source code.
COPY . .

# Build the application for production. 'build' is the conventional script name.
RUN npm run build

# --- Production Stage ---
# Use a lightweight Nginx image to serve the static files.
FROM nginx:1.25-alpine

# Copy the custom Nginx configuration.
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy the built static files from the build stage.
# The build output is often in a 'dist' or 'build' directory.
COPY --from=build /usr/src/app/dist /usr/share/nginx/html

# Expose port 80 for Nginx.
EXPOSE 80

# The default Nginx command starts the server.
CMD ["nginx", "-g", "daemon off;"]