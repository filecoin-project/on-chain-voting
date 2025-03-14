FROM golang:alpine

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
RUN cp /build/abi/oracle.json .

# expose server port
EXPOSE 18001

# run
CMD ["/dist/app"]