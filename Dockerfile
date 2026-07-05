FROM golang:1.25.9-alpine AS build

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN go build -o bookmark-management cmd/api/main.go

FROM alpine:3.21

RUN addgroup -S appgroup && adduser -S appuser -G appgroup


WORKDIR /app

COPY --from=build /opt/app/bookmark-management /app/bookmark-management
COPY --from=build /opt/app/docs /app/docs


CMD ["/app/bookmark-management"]

