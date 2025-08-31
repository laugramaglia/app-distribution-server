# --- Build Stage ---
# Use a specific Go version for reproducibility.
FROM golang:1.24-alpine AS build

# Set the working directory inside the container.
WORKDIR /src

# Copy go.mod and go.sum to leverage Docker's layer caching.
# This step will only be re-run if these files change.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Install swag and generate swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN /go/bin/swag init -g cmd/server/main.go

# Build the Go application.
# CGO_ENABLED=0 is important for building a static binary that can run in a minimal container.
# -ldflags="-w -s" strips debugging information, making the binary smaller.
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /server ./cmd/server

# --- Production Stage ---
# Use a minimal, non-root "distroless" base image for security and a small footprint.
FROM gcr.io/distroless/static-debian12

# Copy the compiled binary from the build stage.
COPY --from=build /server /server
# Copy the generated swagger documentation from the build stage.
COPY --from=build /src/docs /docs

# Expose the port the application will run on.
EXPOSE 8080

# Command to run the application when the container starts.
CMD ["/server"]