FROM golang:alpine as builder

WORKDIR /go/src/BaiduPCS-Go
COPY . .

FROM alpine
COPY --from=builder /go/src/BaiduPCS-Go/BaiduPCS-Go /
WORKDIR /
RUN mkdir /upload
ENTRYPOINT ["/BaiduPCS-Go"]
EXPOSE 9181
