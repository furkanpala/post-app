# Build Vue App
FROM node:alpine AS vue-build
WORKDIR /vue

COPY web/post-app/package*.json ./
RUN npm install

COPY web/post-app ./
RUN npm run build

# Build Go Server
FROM golang:1.14 AS go-build
WORKDIR /root

COPY --from=vue-build /vue/dist ./dist

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY go.mod ./
COPY go.sum ./
RUN go build ./cmd/main.go



ENV PORT=4000

CMD ["./main"]