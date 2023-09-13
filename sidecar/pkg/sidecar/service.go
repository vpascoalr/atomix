// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package sidecar

import (
	counterv1 "github.com/atomix/atomix/api/runtime/counter/v1"
	countermapv1 "github.com/atomix/atomix/api/runtime/countermap/v1"
	electionv1 "github.com/atomix/atomix/api/runtime/election/v1"
	indexedmapv1 "github.com/atomix/atomix/api/runtime/indexedmap/v1"
	listv1 "github.com/atomix/atomix/api/runtime/list/v1"
	lockv1 "github.com/atomix/atomix/api/runtime/lock/v1"
	mapv1 "github.com/atomix/atomix/api/runtime/map/v1"
	multimapv1 "github.com/atomix/atomix/api/runtime/multimap/v1"
	setv1 "github.com/atomix/atomix/api/runtime/set/v1"
	topicv1 "github.com/atomix/atomix/api/runtime/topic/v1"
	valuev1 "github.com/atomix/atomix/api/runtime/value/v1"
	"github.com/vpascoalr/atomix/runtime/pkg/network"
	counterproxyv1 "github.com/vpascoalr/atomix/runtime/pkg/runtime/counter/v1"
	countermapproxyv1 "github.com/vpascoalr/atomix/runtime/pkg/runtime/countermap/v1"
	electionproxyv1 "github.com/vpascoalr/atomix/runtime/pkg/runtime/election/v1"
	indexedmapproxyv1 "github.com/vpascoalr/atomix/runtime/pkg/runtime/indexedmap/v1"
	listproxyv1 "github.com/vpascoalr/atomix/runtime/pkg/runtime/list/v1"
	lockproxyv1 "github.com/vpascoalr/atomix/runtime/pkg/runtime/lock/v1"
	mapproxyv1 "github.com/vpascoalr/atomix/runtime/pkg/runtime/map/v1"
	multimapproxyv1 "github.com/vpascoalr/atomix/runtime/pkg/runtime/multimap/v1"
	setproxyv1 "github.com/vpascoalr/atomix/runtime/pkg/runtime/set/v1"
	topicproxyv1 "github.com/vpascoalr/atomix/runtime/pkg/runtime/topic/v1"
	runtime "github.com/vpascoalr/atomix/runtime/pkg/runtime/v1"
	valueproxyv1 "github.com/vpascoalr/atomix/runtime/pkg/runtime/value/v1"
	"github.com/vpascoalr/atomix/runtime/pkg/utils/grpc/interceptors"
	"google.golang.org/grpc"
)

func register(server *grpc.Server, runtime *runtime.Runtime) {
	counterv1.RegisterCounterServer(server, counterproxyv1.NewCounterServer(runtime))
	counterv1.RegisterCountersServer(server, counterproxyv1.NewCountersServer(runtime))
	countermapv1.RegisterCounterMapServer(server, countermapproxyv1.NewCounterMapServer(runtime))
	countermapv1.RegisterCounterMapsServer(server, countermapproxyv1.NewCounterMapsServer(runtime))
	electionv1.RegisterLeaderElectionServer(server, electionproxyv1.NewLeaderElectionServer(runtime))
	electionv1.RegisterLeaderElectionsServer(server, electionproxyv1.NewLeaderElectionsServer(runtime))
	indexedmapv1.RegisterIndexedMapServer(server, indexedmapproxyv1.NewIndexedMapServer(runtime))
	indexedmapv1.RegisterIndexedMapsServer(server, indexedmapproxyv1.NewIndexedMapsServer(runtime))
	listv1.RegisterListServer(server, listproxyv1.NewListServer(runtime))
	listv1.RegisterListsServer(server, listproxyv1.NewListsServer(runtime))
	lockv1.RegisterLockServer(server, lockproxyv1.NewLockServer(runtime))
	lockv1.RegisterLocksServer(server, lockproxyv1.NewLocksServer(runtime))
	mapv1.RegisterMapServer(server, mapproxyv1.NewMapServer(runtime))
	mapv1.RegisterMapsServer(server, mapproxyv1.NewMapsServer(runtime))
	multimapv1.RegisterMultiMapServer(server, multimapproxyv1.NewMultiMapServer(runtime))
	multimapv1.RegisterMultiMapsServer(server, multimapproxyv1.NewMultiMapsServer(runtime))
	setv1.RegisterSetServer(server, setproxyv1.NewSetServer(runtime))
	setv1.RegisterSetsServer(server, setproxyv1.NewSetsServer(runtime))
	topicv1.RegisterTopicServer(server, topicproxyv1.NewTopicServer(runtime))
	topicv1.RegisterTopicsServer(server, topicproxyv1.NewTopicsServer(runtime))
	valuev1.RegisterValueServer(server, valueproxyv1.NewValueServer(runtime))
	valuev1.RegisterValuesServer(server, valueproxyv1.NewValuesServer(runtime))
}

type Service struct {
	network.Service
	Options
}

func NewService(runtime *runtime.Runtime, opts ...Option) network.Service {
	var options Options
	options.apply(opts...)
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024*20),
		grpc.UnaryInterceptor(interceptors.ErrorHandlingUnaryServerInterceptor()),
		grpc.StreamInterceptor(interceptors.ErrorHandlingStreamServerInterceptor()))
	register(server, runtime)
	return &Service{
		Options: options,
		Service: network.NewService(server,
			network.WithDriver(options.Network),
			network.WithHost(options.Host),
			network.WithPort(options.Port)),
	}
}
