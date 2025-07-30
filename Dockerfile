FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o glog main.go

# ──────────────

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/glog .

# config 파일은 컨테이너 실행시 -v로 주입
# 포트도 실행시 환경변수로 지정 가능해야 함

EXPOSE 8080

ENTRYPOINT ["./glog"]