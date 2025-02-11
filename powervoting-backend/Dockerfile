FROM golang:alpine AS backend-builder

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

FROM node:22

RUN npm install -g @web3-storage/w3cli

# move to /dist
WORKDIR /dist

# copy file to /dist
RUN cp --from=backend-builder /build/app .
RUN cp --from=backend-builder /build/configuration.yaml .
RUN cp --from=backend-builder /build/abi/power-voting.json .
RUN cp --from=backend-builder /build/abi/oracle.json .

# expose server port
EXPOSE 18007

# run
CMD ["/dist/app"]