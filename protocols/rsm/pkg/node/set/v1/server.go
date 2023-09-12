// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"context"

	"github.com/atomix/atomix/runtime/pkg/logging"
	streams "github.com/atomix/atomix/runtime/pkg/stream"
	"github.com/gogo/protobuf/proto"
	setprotocolv1 "github.com/vpascoalr/atomix/protocols/rsm/api/set/v1"
	"github.com/vpascoalr/atomix/protocols/rsm/pkg/node"
	"google.golang.org/grpc"
)

var log = logging.GetLogger()

func RegisterServer(node *node.Node) {
	node.RegisterService(func(server *grpc.Server) {
		setprotocolv1.RegisterSetServer(server, NewSetServer(node))
	})
}

var serverCodec = node.NewCodec[*setprotocolv1.SetInput, *setprotocolv1.SetOutput](
	func(input *setprotocolv1.SetInput) ([]byte, error) {
		return proto.Marshal(input)
	},
	func(bytes []byte) (*setprotocolv1.SetOutput, error) {
		output := &setprotocolv1.SetOutput{}
		if err := proto.Unmarshal(bytes, output); err != nil {
			return nil, err
		}
		return output, nil
	})

func NewSetServer(protocol node.Protocol) setprotocolv1.SetServer {
	return &setServer{
		handler: node.NewHandler[*setprotocolv1.SetInput, *setprotocolv1.SetOutput](protocol, serverCodec),
	}
}

type setServer struct {
	handler node.Handler[*setprotocolv1.SetInput, *setprotocolv1.SetOutput]
}

func (s *setServer) Size(ctx context.Context, request *setprotocolv1.SizeRequest) (*setprotocolv1.SizeResponse, error) {
	log.Debugw("Size",
		logging.Trunc128("SizeRequest", request))
	input := &setprotocolv1.SetInput{
		Input: &setprotocolv1.SetInput_Size_{
			Size_: request.SizeInput,
		},
	}
	output, headers, err := s.handler.Query(ctx, input, request.Headers)
	if err != nil {
		log.Warnw("Size",
			logging.Trunc128("SizeRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response := &setprotocolv1.SizeResponse{
		Headers:    headers,
		SizeOutput: output.GetSize_(),
	}
	log.Debugw("Size",
		logging.Trunc128("SizeRequest", request),
		logging.Trunc128("SizeResponse", response))
	return response, nil
}

func (s *setServer) Add(ctx context.Context, request *setprotocolv1.AddRequest) (*setprotocolv1.AddResponse, error) {
	log.Debugw("Add",
		logging.Trunc128("AddRequest", request))
	input := &setprotocolv1.SetInput{
		Input: &setprotocolv1.SetInput_Add{
			Add: request.AddInput,
		},
	}
	output, headers, err := s.handler.Propose(ctx, input, request.Headers)
	if err != nil {
		log.Warnw("Add",
			logging.Trunc128("AddRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response := &setprotocolv1.AddResponse{
		Headers:   headers,
		AddOutput: output.GetAdd(),
	}
	log.Debugw("Add",
		logging.Trunc128("AddRequest", request),
		logging.Trunc128("AddResponse", response))
	return response, nil
}

func (s *setServer) Contains(ctx context.Context, request *setprotocolv1.ContainsRequest) (*setprotocolv1.ContainsResponse, error) {
	log.Debugw("Contains",
		logging.Trunc128("ContainsRequest", request))
	input := &setprotocolv1.SetInput{
		Input: &setprotocolv1.SetInput_Contains{
			Contains: request.ContainsInput,
		},
	}
	output, headers, err := s.handler.Query(ctx, input, request.Headers)
	if err != nil {
		log.Warnw("Contains",
			logging.Trunc128("ContainsRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response := &setprotocolv1.ContainsResponse{
		Headers:        headers,
		ContainsOutput: output.GetContains(),
	}
	log.Debugw("Contains",
		logging.Trunc128("ContainsRequest", request),
		logging.Trunc128("ContainsResponse", response))
	return response, nil
}

func (s *setServer) Remove(ctx context.Context, request *setprotocolv1.RemoveRequest) (*setprotocolv1.RemoveResponse, error) {
	log.Debugw("Remove",
		logging.Trunc128("RemoveRequest", request))
	input := &setprotocolv1.SetInput{
		Input: &setprotocolv1.SetInput_Remove{
			Remove: request.RemoveInput,
		},
	}
	output, headers, err := s.handler.Propose(ctx, input, request.Headers)
	if err != nil {
		log.Warnw("Remove",
			logging.Trunc128("RemoveRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response := &setprotocolv1.RemoveResponse{
		Headers:      headers,
		RemoveOutput: output.GetRemove(),
	}
	log.Debugw("Remove",
		logging.Trunc128("RemoveRequest", request),
		logging.Trunc128("RemoveResponse", response))
	return response, nil
}

func (s *setServer) Clear(ctx context.Context, request *setprotocolv1.ClearRequest) (*setprotocolv1.ClearResponse, error) {
	log.Debugw("Clear",
		logging.Trunc128("ClearRequest", request))
	input := &setprotocolv1.SetInput{
		Input: &setprotocolv1.SetInput_Clear{
			Clear: request.ClearInput,
		},
	}
	output, headers, err := s.handler.Propose(ctx, input, request.Headers)
	if err != nil {
		log.Warnw("Clear",
			logging.Trunc128("ClearRequest", request),
			logging.Error("Error", err))
		return nil, err
	}
	response := &setprotocolv1.ClearResponse{
		Headers:     headers,
		ClearOutput: output.GetClear(),
	}
	log.Debugw("Clear",
		logging.Trunc128("ClearRequest", request),
		logging.Trunc128("ClearResponse", response))
	return response, nil
}

func (s *setServer) Events(request *setprotocolv1.EventsRequest, server setprotocolv1.Set_EventsServer) error {
	log.Debugw("Events",
		logging.Trunc128("EventsRequest", request))
	input := &setprotocolv1.SetInput{
		Input: &setprotocolv1.SetInput_Events{
			Events: request.EventsInput,
		},
	}

	stream := streams.NewBufferedStream[*node.StreamProposalResponse[*setprotocolv1.SetOutput]]()
	go func() {
		err := s.handler.StreamPropose(server.Context(), input, request.Headers, stream)
		if err != nil {
			log.Warnw("Events",
				logging.Trunc128("EventsRequest", request),
				logging.Error("Error", err))
			stream.Error(err)
			stream.Close()
		}
	}()

	for {
		result, ok := stream.Receive()
		if !ok {
			return nil
		}

		if result.Failed() {
			log.Warnw("Events",
				logging.Trunc128("EventsRequest", request),
				logging.Error("Error", result.Error))
			return result.Error
		}

		response := &setprotocolv1.EventsResponse{
			Headers:      result.Value.Headers,
			EventsOutput: result.Value.Output.GetEvents(),
		}
		log.Debugw("Events",
			logging.Trunc128("EventsRequest", request),
			logging.Trunc128("EventsResponse", response))
		if err := server.Send(response); err != nil {
			log.Warnw("Events",
				logging.Trunc128("EventsRequest", request),
				logging.Error("Error", err))
			return err
		}
	}
}

func (s *setServer) Elements(request *setprotocolv1.ElementsRequest, server setprotocolv1.Set_ElementsServer) error {
	log.Debugw("Elements",
		logging.Trunc128("ElementsRequest", request))
	input := &setprotocolv1.SetInput{
		Input: &setprotocolv1.SetInput_Elements{
			Elements: request.ElementsInput,
		},
	}

	stream := streams.NewBufferedStream[*node.StreamQueryResponse[*setprotocolv1.SetOutput]]()
	go func() {
		err := s.handler.StreamQuery(server.Context(), input, request.Headers, stream)
		if err != nil {
			log.Warnw("Elements",
				logging.Trunc128("ElementsRequest", request),
				logging.Error("Error", err))
			stream.Error(err)
			stream.Close()
		}
	}()

	for {
		result, ok := stream.Receive()
		if !ok {
			return nil
		}

		if result.Failed() {
			log.Warnw("Elements",
				logging.Trunc128("ElementsRequest", request),
				logging.Error("Error", result.Error))
			return result.Error
		}

		response := &setprotocolv1.ElementsResponse{
			Headers:        result.Value.Headers,
			ElementsOutput: result.Value.Output.GetElements(),
		}
		log.Debugw("Elements",
			logging.Trunc128("ElementsRequest", request),
			logging.Trunc128("ElementsResponse", response))
		if err := server.Send(response); err != nil {
			log.Warnw("Elements",
				logging.Trunc128("ElementsRequest", request),
				logging.Error("Error", err))
			return err
		}
	}
}
