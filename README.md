# Go + gRPCでのログイン認証処理の実装

Go + gRPCを使ったログイン認証処理のサンプルを実装した。  
大まかな実装の流れに重点を置いているため詳細な実装は省略している。

## 1. 環境構築

以下の項目をもとに環境構築を行う。

- Go
  - version: v1.21.3
- protocol bufferコンパイラ
  - version: v25.0
- gRPCプラグイン
  - protoc-gen-go
    - version: v1.31.0
  - protoc-gen-go-grpc
    - version: v1.3.0

### 1.1 Goのインストール

[公式](https://go.dev/)からインストーラーをダウンロードしてインストールする。

### 1.2 protocol bufferコンパイラのインストール

```sh
brew install protobuf
protoc --version  # コンパイラのバージョンが3以上であることを確認する
```

### 1.3 GoのgPRCプラグインのインストール

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
```

## 2. proto定義からGoのコードを生成

```sh
protoc --go_out=internal/pb --go_opt=paths=source_relative --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative backend.proto
```

## 3. ログイン認証処理

Authミドルウェアの[grpc_auth](https://github.com/grpc-ecosystem/go-grpc-middleware)を使用する。

### 3.1. ミドルウェアを追加

grpc.ChainUnaryInterceptor()にgrpc_auth.UnaryServerInterceptor()を追加する。

```go
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_auth.UnaryServerInterceptor(interceptor.AuthFunc),
		),
	)
```

### 3.2. 認証処理を実装

interceptor.AuthFuncを実装する。

```go
package interceptor

// AuthFunc is the pluggable function that performs authentication.
func AuthFunc(ctx context.Context) (context.Context, error) {
  // TODO: この関数内部のロジックは適宜変更する。この関数でDBにアクセスしたい場合は、引数にQuerierインターフェースを渡すといいかも
  tokenString, err := grpc_auth.AuthFromMD(ctx, "bearer")
  if err != nil {
    return ctx, err
  }

  userID, err := parseToken(tokenString)
  if err != nil {
    return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
  }

  return context.WithValue(ctx, userIDKey{}, userID), nil
}
```

### 3.3.

認証処理を実施したいメソッドとしたくないメソッドがある場合には`ServiceAuthFuncOverride`インタフェースを活用する。

```go
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	authFreeMethods := []string{
		"/go_grpc_auth_sample.BackendService/Login",
		"/go_grpc_auth_sample.BackendService/Greet",
	}

	for _, method := range authFreeMethods {
		if method == fullMethodName {
			return ctx, nil
		}
	}

	return interceptor.AuthFunc(ctx)
}
```
