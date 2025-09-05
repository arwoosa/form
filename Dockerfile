FROM golang:1.24-alpine AS builder

# 設置環境變數
ENV GO111MODULE=on \
GOPROXY=https://proxy.golang.org,direct \
CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64

WORKDIR /build

COPY . .
RUN go mod download

RUN go build -o partivo_form ./cmd/form-server/

###################
# multi-stage build
###################
FROM scratch

WORKDIR /app
COPY ./conf /app/conf

COPY --from=builder /build/partivo_form /app/

# Expose both HTTP and gRPC ports
EXPOSE 8081

ENTRYPOINT ["/app/partivo_form"]
CMD ["server", "--config", "conf/config_docker.yaml"]