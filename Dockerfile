# ------------BUILD STAGE--------------------
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

#-------------------GENERATING GO BINARY-------------------------------
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api




#--------------------RUNTIME STAGE---------------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]


