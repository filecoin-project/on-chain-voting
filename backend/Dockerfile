FROM golang:alpine AS backend-builder

RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o app .

FROM alpine:latest

WORKDIR /dist

COPY --from=backend-builder /build/app .
COPY --from=backend-builder /build/configuration-backend.yaml .


# Copy oracle abi to instance
COPY --from=backend-builder /build/abi/oracle-powers.json ./
COPY --from=backend-builder /build/abi/oracle.json ./
COPY --from=backend-builder /build/abi/power-voting-conf.json ./
COPY --from=backend-builder /build/abi/power-voting.json ./
COPY --from=backend-builder /build/abi/power-voting-fip.json ./




COPY --from=backend-builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai
CMD ["/dist/app"]