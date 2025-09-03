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

RUN go build -o partivo_form_service ./cmd/form-server/

###################
# multi-stage build
###################
FROM scratch

WORKDIR /app
COPY ./conf /app/conf

COPY --from=builder /build/partivo_form_service /app/

# Expose both HTTP and gRPC ports
EXPOSE 8081 8082

ENTRYPOINT ["/app/partivo_form_service"]
CMD ["server", "--config", "conf/config_docker.yaml"]