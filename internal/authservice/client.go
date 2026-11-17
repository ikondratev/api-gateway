package authservice

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ikondratev/api-gateway/internal/apperr"
	"github.com/ikondratev/api-gateway/internal/authpb"
	"github.com/ikondratev/api-gateway/internal/settings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn *grpc.ClientConn
	api  authpb.AuthServiceClient
}

func New(s *settings.Settings) (*Client, error) {
	conn, err := grpc.NewClient(
		s.Upstream.AuthService.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial auth service: %w", err)
	}

	return &Client{
		conn: conn,
		api:  authpb.NewAuthServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Register(ctx context.Context, r *authpb.RegisterRequest) (*authpb.TokenPair, error) {
	resp, err := c.api.Register(ctx, r)
	if err != nil {
		return nil, toAppErr(err)
	}
	return resp, err
}

func (c *Client) Login(ctx context.Context, r *authpb.LoginRequest) (*authpb.TokenPair, error) {
	resp, err := c.api.Login(ctx, r)
	if err != nil {
		return nil, toAppErr(err)
	}
	return resp, err
}

func (c *Client) Refresh(ctx context.Context, r *authpb.RefreshRequest) (*authpb.TokenPair, error) {
	resp, err := c.api.Refresh(ctx, r)
	if err != nil {
		return nil, toAppErr(err)
	}
	return resp, err
}

func (c *Client) Logout(ctx context.Context, r *authpb.LogoutRequest) error {
	_, err := c.api.Logout(ctx, r)
	if err != nil {
		return toAppErr(err)
	}
	return nil
}

func toAppErr(e error) error {
	st, ok := status.FromError(e)
	if !ok {
		return apperr.Internal(e)
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return apperr.Business(http.StatusBadRequest, st.Message(), st.Message())
	case codes.AlreadyExists:
		return apperr.Business(http.StatusConflict, "email_taken", st.Message())
	case codes.Unauthenticated:
		return apperr.Business(http.StatusUnauthorized, "unauthenticated", st.Message())
	default:
		return apperr.Internal(e)
	}
}
