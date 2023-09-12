// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/atomix/atomix/api/errors"
	protocol "github.com/atomix/atomix/protocols/rsm/api/v1"
	"github.com/atomix/atomix/runtime/pkg/utils/grpc/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

func newPartition(id protocol.PartitionID, client *ProtocolClient, sessionTimeout time.Duration) *PartitionClient {
	return &PartitionClient{
		client:         client,
		id:             id,
		sessionTimeout: sessionTimeout,
	}
}

type PartitionClient struct {
	client         *ProtocolClient
	id             protocol.PartitionID
	sessionTimeout time.Duration
	conn           *grpc.ClientConn
	resolver       *partitionResolver
	session        atomic.Value
	mu             sync.RWMutex
}

func (p *PartitionClient) ID() protocol.PartitionID {
	return p.id
}

func (p *PartitionClient) GetSession(ctx context.Context) (*SessionClient, error) {
	stored := p.session.Load()
	if stored != nil {
		return stored.(*SessionClient), nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	stored = p.session.Load()
	if stored != nil {
		return stored.(*SessionClient), nil
	}
	if p.conn == nil {
		return nil, errors.NewUnavailable("not connected")
	}

	request := &protocol.OpenSessionRequest{
		Headers: &protocol.PartitionRequestHeaders{
			PartitionID: p.id,
		},
		OpenSessionInput: &protocol.OpenSessionInput{
			Timeout: p.sessionTimeout,
		},
	}

	client := protocol.NewPartitionClient(p.conn)
	response, err := client.OpenSession(ctx, request)
	if err != nil {
		return nil, err
	}

	session := newSessionClient(response.SessionID, p, p.conn, p.sessionTimeout)
	p.session.Store(session)
	return session, nil
}

func (p *PartitionClient) connect(ctx context.Context, config *protocol.PartitionConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	address := fmt.Sprintf("%s:///%d", resolverName, p.id)
	p.resolver = newResolver(config)
	dialOptions := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024 * 1024 * 20),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, resolverName)),
		grpc.WithResolvers(p.resolver),
		grpc.WithContextDialer(p.client.network.Connect),
		grpc.WithChainUnaryInterceptor(
			interceptors.ErrorHandlingUnaryClientInterceptor(),
			interceptors.RetryingUnaryClientInterceptor(interceptors.WithRetryOn(codes.Unavailable))),
		grpc.WithChainStreamInterceptor(
			interceptors.ErrorHandlingStreamClientInterceptor(),
			interceptors.RetryingStreamClientInterceptor(interceptors.WithRetryOn(codes.Unavailable))),
	}
	dialOptions = append(dialOptions, p.client.GRPCDialOptions...)
	conn, err := grpc.DialContext(ctx, address, dialOptions...)
	if err != nil {
		return err
	}
	p.conn = conn
	return nil
}

func (p *PartitionClient) configure(config *protocol.PartitionConfig) error {
	return p.resolver.update(config)
}

func (p *PartitionClient) close(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	session := p.session.Load()
	if session != nil {
		return session.(*SessionClient).close(ctx)
	}
	return nil
}
