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
RUN go build -o app .


#requires Node 18 or higher
FROM node:22

# 
RUN npm install -g @web3-storage/w3cli

# move to /dist
WORKDIR /dist

# copy file to /dist
COPY --from=backend-builder /build/app .
COPY --from=backend-builder /build/configuration.yaml .
COPY --from=backend-builder /build/abi/power-voting.json .
COPY --from=backend-builder /build/abi/oracle.json .

# expose server port
EXPOSE 9999

# run
CMD ["/dist/app"]