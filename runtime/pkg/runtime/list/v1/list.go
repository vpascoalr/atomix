// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"context"

	listv1 "github.com/atomix/atomix/api/runtime/list/v1"
	"github.com/vpascoalr/atomix/runtime/pkg/logging"
	runtime "github.com/vpascoalr/atomix/runtime/pkg/runtime/v1"
)

var log = logging.GetLogger()

type ListProxy interface {
	runtime.PrimitiveProxy
	// Size gets the number of elements in the list
	Size(context.Context, *listv1.SizeRequest) (*listv1.SizeResponse, error)
	// Append appends a value to the list
	Append(context.Context, *listv1.AppendRequest) (*listv1.AppendResponse, error)
	// Insert inserts a value at a specific index in the list
	Insert(context.Context, *listv1.InsertRequest) (*listv1.InsertResponse, error)
	// Get gets the value at an index in the list
	Get(context.Context, *listv1.GetRequest) (*listv1.GetResponse, error)
	// Set sets the value at an index in the list
	Set(context.Context, *listv1.SetRequest) (*listv1.SetResponse, error)
	// Remove removes an element from the list
	Remove(context.Context, *listv1.RemoveRequest) (*listv1.RemoveResponse, error)
	// Clear removes all elements from the list
	Clear(context.Context, *listv1.ClearRequest) (*listv1.ClearResponse, error)
	// Events listens for change events
	Events(*listv1.EventsRequest, listv1.List_EventsServer) error
	// Items streams all items in the list
	Items(*listv1.ItemsRequest, listv1.List_ItemsServer) error
}

func NewListServer(rt *runtime.Runtime) listv1.ListServer {
	return &listServer{
		ListsServer: NewListsServer(rt),
		primitives:  runtime.NewPrimitiveRegistry[ListProxy](listv1.PrimitiveType, rt),
	}
}

type listServer struct {
	listv1.ListsServer
	primitives runtime.PrimitiveRegistry[ListProxy]
}

func (s *listServer) Size(ctx context.Context, request *listv1.SizeRequest) (*listv1.SizeResponse, error) {
	log.Debugw("Size",
		logging.Trunc64("SizeRequest", request))
	client, err := s.primitives.Get(request.ID)
	if err != nil {
		log.Warnw("Size",
			logging.Trunc64("SizeRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response, err := client.Size(ctx, request)
	if err != nil {
		log.Debugw("Size",
			logging.Trunc64("SizeRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	log.Debugw("Size",
		logging.Trunc64("SizeResponse", response))
	return response, nil
}

func (s *listServer) Append(ctx context.Context, request *listv1.AppendRequest) (*listv1.AppendResponse, error) {
	log.Debugw("Append",
		logging.Trunc64("AppendRequest", request))
	client, err := s.primitives.Get(request.ID)
	if err != nil {
		log.Warnw("Append",
			logging.Trunc64("AppendRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response, err := client.Append(ctx, request)
	if err != nil {
		log.Debugw("Append",
			logging.Trunc64("AppendRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	log.Debugw("Append",
		logging.Trunc64("AppendResponse", response))
	return response, nil
}

func (s *listServer) Insert(ctx context.Context, request *listv1.InsertRequest) (*listv1.InsertResponse, error) {
	log.Debugw("Insert",
		logging.Trunc64("InsertRequest", request))
	client, err := s.primitives.Get(request.ID)
	if err != nil {
		log.Warnw("Insert",
			logging.Trunc64("InsertRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response, err := client.Insert(ctx, request)
	if err != nil {
		log.Debugw("Insert",
			logging.Trunc64("InsertRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	log.Debugw("Insert",
		logging.Trunc64("InsertResponse", response))
	return response, nil
}

func (s *listServer) Get(ctx context.Context, request *listv1.GetRequest) (*listv1.GetResponse, error) {
	log.Debugw("Get",
		logging.Trunc64("GetRequest", request))
	client, err := s.primitives.Get(request.ID)
	if err != nil {
		log.Warnw("Get",
			logging.Trunc64("GetRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response, err := client.Get(ctx, request)
	if err != nil {
		log.Debugw("Get",
			logging.Trunc64("GetRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	log.Debugw("Get",
		logging.Trunc64("GetResponse", response))
	return response, nil
}

func (s *listServer) Set(ctx context.Context, request *listv1.SetRequest) (*listv1.SetResponse, error) {
	log.Debugw("Set",
		logging.Trunc64("SetRequest", request))
	client, err := s.primitives.Get(request.ID)
	if err != nil {
		log.Warnw("Set",
			logging.Trunc64("SetRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response, err := client.Set(ctx, request)
	if err != nil {
		log.Debugw("Set",
			logging.Trunc64("SetRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	log.Debugw("Set",
		logging.Trunc64("SetResponse", response))
	return response, nil
}

func (s *listServer) Remove(ctx context.Context, request *listv1.RemoveRequest) (*listv1.RemoveResponse, error) {
	log.Debugw("Remove",
		logging.Trunc64("RemoveRequest", request))
	client, err := s.primitives.Get(request.ID)
	if err != nil {
		log.Warnw("Remove",
			logging.Trunc64("RemoveRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response, err := client.Remove(ctx, request)
	if err != nil {
		log.Debugw("Remove",
			logging.Trunc64("RemoveRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	log.Debugw("Remove",
		logging.Trunc64("RemoveResponse", response))
	return response, nil
}

func (s *listServer) Clear(ctx context.Context, request *listv1.ClearRequest) (*listv1.ClearResponse, error) {
	log.Debugw("Clear",
		logging.Trunc64("ClearRequest", request))
	client, err := s.primitives.Get(request.ID)
	if err != nil {
		log.Warnw("Clear",
			logging.Trunc64("ClearRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response, err := client.Clear(ctx, request)
	if err != nil {
		log.Debugw("Clear",
			logging.Trunc64("ClearRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	log.Debugw("Clear",
		logging.Trunc64("ClearResponse", response))
	return response, nil
}

func (s *listServer) Events(request *listv1.EventsRequest, server listv1.List_EventsServer) error {
	log.Debugw("Events",
		logging.Trunc64("EventsRequest", request),
		logging.String("State", "started"))
	client, err := s.primitives.Get(request.ID)
	if err != nil {
		log.Warnw("Events",
			logging.Trunc64("EventsRequest", request),
			logging.Error("Error", err))
		return err
	}
	err = client.Events(request, server)
	if err != nil {
		log.Debugw("Events",
			logging.Trunc64("EventsRequest", request),
			logging.Error("Error", err))
		return err
	}
	return nil
}

func (s *listServer) Items(request *listv1.ItemsRequest, server listv1.List_ItemsServer) error {
	log.Debugw("Items",
		logging.Trunc64("ItemsRequest", request),
		logging.String("State", "started"))
	client, err := s.primitives.Get(request.ID)
	if err != nil {
		log.Warnw("Items",
			logging.Trunc64("ItemsRequest", request),
			logging.Error("Error", err))
		return err
	}
	err = client.Items(request, server)
	if err != nil {
		log.Debugw("Items",
			logging.Trunc64("ItemsRequest", request),
			logging.Error("Error", err))
		return err
	}
	return nil
}

var _ listv1.ListServer = (*listServer)(nil)
