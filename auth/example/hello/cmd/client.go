package main

import (
	"context"

	"github.com/leostudio/kit/auth/example/hello/pb"
	"github.com/leostudio/kit/log"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	testEmail    = "changeit"
	testPassword = "changeit"
)

var (
	logger = log.Logger()
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloClient(conn)

	conf := oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: "http://localhost:50053/token",
		},
		Scopes: []string{"all"},
	}
	token, err := conf.PasswordCredentialsToken(context.TODO(), testEmail, testPassword)
	if err != nil {
		logger.Fatalf("get access token error: %v", err)
	}

	cred := newPerRPCCredentials(token)

	rsp, err := c.SayHello(context.TODO(), &pb.Req{Name: "world"}, grpc.PerRPCCredentials(cred))
	if err != nil {
		logger.Fatalf("SayHello api error: %v", err)
	}
	logger.Infof("SayHello response:%v", rsp.Message)

	rsp, err = c.SayHelloInsecure(context.TODO(), &pb.Req{Name: "world"})
	if err != nil {
		logger.Fatalf("SayHelloInsecure api error: %v", err)
	}
	logger.Infof("SayHelloInsecure response:%v", rsp.Message)
}

func newPerRPCCredentials(token *oauth2.Token) credentials.PerRPCCredentials {
	return &insecure{token: token}
}

type insecure struct {
	token *oauth2.Token
}

func (p *insecure) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": p.token.Type() + " " + p.token.AccessToken,
	}, nil
}

func (p *insecure) RequireTransportSecurity() bool {
	return false
}