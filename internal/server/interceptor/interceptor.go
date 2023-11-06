package interceptor

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryFunc(p any) (err error) {
	return status.Errorf(codes.Unknown, "unexpected error: %v", p)
}

type userIDKey struct{}

func parseToken(_ string) (string, error) {
	userID := "user_id"
	return userID, nil
}

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
