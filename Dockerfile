FROM golang:1.26.1 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y wget ca-certificates && rm -rf /var/lib/apt/lists/*

ARG TAILWIND_VERSION=4.2.2

RUN arch="$(dpkg --print-architecture)" && \
    case "$arch" in \
      amd64) tailwind_arch="x64"; tailwind_sha256="4ab84f2b496c402d3ec4fd25e0e5559fe1184d886dadae8fb4438344ec044c22" ;; \
      arm64) tailwind_arch="arm64"; tailwind_sha256="ad627e77b496cccada4a6e26eafff698ef0829081e575a4baf3af8524bb00747" ;; \
      *) echo "unsupported architecture: $arch" && exit 1 ;; \
    esac && \
    wget -O /usr/local/bin/tailwindcss "https://github.com/tailwindlabs/tailwindcss/releases/download/v${TAILWIND_VERSION}/tailwindcss-linux-${tailwind_arch}" && \
    printf '%s  %s\n' "$tailwind_sha256" /usr/local/bin/tailwindcss | sha256sum -c - && \
    chmod +x /usr/local/bin/tailwindcss

RUN TEMPLUI_PATH="$(go list -mod=mod -m -f '{{.Dir}}' github.com/templui/templui)" && \
    test -n "$TEMPLUI_PATH" && \
    test -d "$TEMPLUI_PATH/components" && \
    printf '%s\n' \
      '@source "./**/*.templ";' \
      '@source "./**/*.js";' \
      "@source \"$TEMPLUI_PATH/components/**/*.templ\";" \
      "@source \"$TEMPLUI_PATH/components/**/*.js\";" \
      > ./assets/css/sources.generated.css && \
    tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --minify

RUN go tool templ generate

RUN CGO_ENABLED=0 GOOS=linux go build -o meshbox-studio ./cmd/meshbox-studio

FROM alpine:3.23.3
WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=build /app/meshbox-studio /app/meshbox-studio
ENV GO_ENV=production
EXPOSE 8090

CMD ["/app/meshbox-studio"]
