FROM golang:1.21.6-alpine as builder

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係のコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY main.go ./

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 実行環境
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルダーからバイナリをコピー
COPY --from=builder /app/main .

# 環境変数NODE_IPを設定（Kubernetesで上書き可能）
ENV NODE_IP=127.0.0.1

# アプリケーションの実行
CMD ["./main"]
