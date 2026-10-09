package eventservice

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ikondratev/api-gateway/internal/apperr"
	"github.com/ikondratev/api-gateway/internal/eventpb"
	"github.com/ikondratev/api-gateway/internal/settings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn *grpc.ClientConn
	api  eventpb.EventServiceClient
}

func New(s *settings.Settings) (*Client, error) {
	conn, err := grpc.NewClient(
		s.Upstream.EventSerivce.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial event service %w", err)
	}

	return &Client{
		conn: conn,
		api:  eventpb.NewEventServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateEvent(
	ctx context.Context,
	authorization string,
	idempotencyKey string,
	req *eventpb.CreateEventRequest,
) (*eventpb.CreateEventResponse, error) {
	ctx = metadata.AppendToOutgoingContext(ctx,
		"authorization", authorization,
		"idempotency-key", idempotencyKey,
	)

	resp, err := c.api.CreateEvent(ctx, req)
	if err != nil {
		return nil, toAppErr(err)
	}

	return resp, nil
}

func toAppErr(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return apperr.Internal(err)
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return apperr.Business(http.StatusBadRequest, "invalid_event", st.Message())
	case codes.Unauthenticated:
		return apperr.Business(http.StatusUnauthorized, "unauthorized", "unauthorized")
	case codes.PermissionDenied:
		return apperr.Business(http.StatusForbidden, "forbidden", "forbidden")
	case codes.AlreadyExists:
		return apperr.Business(http.StatusConflict, "conflict", st.Message())
	default:
		return apperr.Internal(err)
	}
}
