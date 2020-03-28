FROM golang:1.13 as builder

WORKDIR /go/src/BaiduPCS-Go
COPY . .

RUN go build

FROM scratch
COPY --from=builder /go/src/BaiduPCS-Go/BaiduPCS-Go /
WORKDIR /
CMD ./BaiduPCS-Go
EXPOSE 9181