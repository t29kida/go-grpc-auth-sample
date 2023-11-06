package server

import (
	"context"

	"github.com/t29kida/go-grpc-auth-sample/internal/pb"
	"github.com/t29kida/go-grpc-auth-sample/internal/server/interceptor"
	"github.com/t29kida/go-grpc-auth-sample/internal/service/auth"
	"github.com/t29kida/go-grpc-auth-sample/internal/service/hash"
)

var _ pb.BackendServiceServer = (*Server)(nil)

type Server struct {
	pb.UnimplementedBackendServiceServer

	auth auth.Auther
	hash hash.Hasher
}

func New(auther auth.Auther, hasher hash.Hasher) *Server {
	return &Server{
		auth: auther,
		hash: hasher,
	}
}

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

func (s *Server) SignUp(_ context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	// TODO: 登録処理内容は適宜変更する

	// TODO: ハッシュ化したパスワードやユーザー情報などをDBに保存する
	_, err := s.hash.CreateHash(req.GetPassword())
	if err != nil {
		return nil, err
	}

	accessToken, err := s.auth.CreateAccessToken()
	if err != nil {
		return nil, err
	}

	return &pb.SignUpResponse{AccessToken: accessToken}, nil
}

func (s *Server) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// TODO: ログイン処理内容は適宜変更する

	// TODO: DBからハッシュ化したパスワードを取得し、比較する
	//password := "hashed_password_from_database"
	//match, err := s.hash.CompareHash(req.GetPassword(), password)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if !match {
	//	return nil, errors.New("user identifier or password is invalid")
	//}

	accessToken, err := s.auth.CreateAccessToken()
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{AccessToken: accessToken}, nil
}

func (s *Server) Greet(_ context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	return &pb.GreetResponse{Message: "Public: Hello " + req.GetName()}, nil
}

func (s *Server) PrivateGreet(_ context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	return &pb.GreetResponse{Message: "Private: Hello " + req.GetName()}, nil
}
