FROM golang:alpine as build

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN go build -o bookmark-management cmd/api/main.go

FROM alpine

WORKDIR /app

COPY --from=build /opt/app/bookmark-management /app/bookmark-management
COPY --from=build /opt/app/docs /app/docs


CMD ["/app/bookmark-management"]

