
# Stage 1: Base - Cài đặt môi trường và dependency
FROM golang:1.25.9-alpine AS base
RUN apk add --no-cache git make # Cần cho go mod và chạy lệnh make bên trong
WORKDIR /opt/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Stage 2: Test Execution - Chạy test bên trong Docker
FROM base AS test-exec
# Khai báo các argument để nhận giá trị từ Makefile
ARG COVERAGE_EXCLUDE
RUN mkdir -p /tmp/coverage
# Chạy lệnh test y hệt như trong Makefile nhưng ở trong container
RUN go test ./... -coverprofile=/tmp/coverage/coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
    grep -vE "$COVERAGE_EXCLUDE" /tmp/coverage/coverage.tmp > /tmp/coverage/coverage.out

# Stage 3: Test Output - Dùng để xuất file report ra máy host
FROM scratch AS test
COPY --from=test-exec /tmp/coverage /

# Stage 4: Build - Biên dịch binary
FROM base AS build
RUN go build -o bookmark-management cmd/api/main.go

# Stage 5: Final - Image chạy thực tế (nhẹ và bảo mật)
FROM alpine:3.21 AS final
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=build /opt/app/bookmark-management /app/bookmark-management
COPY --from=build /opt/app/docs /app/docs
USER appuser
CMD ["/app/bookmark-management"]
