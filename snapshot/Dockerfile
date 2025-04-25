FROM golang:alpine

RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone
    

# env
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# move to work dir：/build
WORKDIR /build

# copy code
COPY . .

# go build
RUN go mod tidy
RUN go build -o app .

# move to /dist
WORKDIR /dist

# copy file to /dist
RUN cp /build/app .
RUN cp /build/configuration.yaml .
# RUN cp /build/abi/oracle.json .
# RUN cp /build/proof.ucan .

COPY configuration.yaml .
COPY proof.ucan .
COPY --from=snapshot-builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai
# expose server port
EXPOSE 8001

# run
CMD ["/dist/app"]