// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package raft

import (
	"context"
	"sync/atomic"

	"github.com/atomix/atomix/api/errors"
	streams "github.com/atomix/atomix/runtime/pkg/stream"
	raftv1 "github.com/atomix/atomix/stores/raft/api/v1"
	"github.com/gogo/protobuf/proto"
	"github.com/lni/dragonboat/v3"
	rsmv1 "github.com/vpascoalr/atomix/protocols/rsm/api/v1"
	"github.com/vpascoalr/atomix/protocols/rsm/pkg/node"
	"google.golang.org/grpc/metadata"
)

func newPartition(id rsmv1.PartitionID, memberID raftv1.MemberID, host *dragonboat.NodeHost, streams *protocolContext) *Partition {
	partition := &Partition{
		memberID: memberID,
	}
	partition.Partition = node.NewPartition(id, &Executor{
		Partition: partition,
		host:      host,
		streams:   streams,
	})
	return partition
}

type Partition struct {
	node.Partition
	memberID raftv1.MemberID
	ready    int32
	leader   uint64
	term     uint64
}

func (p *Partition) setReady() {
	atomic.StoreInt32(&p.ready, 1)
}

func (p *Partition) getReady() bool {
	return atomic.LoadInt32(&p.ready) == 1
}

func (p *Partition) setLeader(term raftv1.Term, leader raftv1.MemberID) {
	atomic.StoreUint64(&p.term, uint64(term))
	atomic.StoreUint64(&p.leader, uint64(leader))
}

func (p *Partition) getLeader() (raftv1.Term, raftv1.MemberID) {
	return raftv1.Term(atomic.LoadUint64(&p.term)), raftv1.MemberID(atomic.LoadUint64(&p.leader))
}

type Executor struct {
	*Partition
	host    *dragonboat.NodeHost
	streams *protocolContext
}

// Propose proposes a change to the protocol
func (e *Executor) Propose(ctx context.Context, input *rsmv1.ProposalInput, stream streams.WriteStream[*rsmv1.ProposalOutput]) error {
	term, leader := e.getLeader()
	if leader != e.memberID {
		return errors.NewUnavailable("not the leader")
	}

	inputBytes, err := proto.Marshal(input)
	if err != nil {
		return errors.NewInternal(err.Error())
	}

	sequenceNum := e.streams.addStream(term, stream)
	proposal := &raftv1.RaftProposal{
		Term:        term,
		SequenceNum: sequenceNum,
		Data:        inputBytes,
	}

	proposalBytes, err := proto.Marshal(proposal)
	if err != nil {
		return errors.NewInternal("failed to marshal RaftLogEntry: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, defaultClientTimeout)
	defer cancel()
	if _, err := e.host.SyncPropose(ctx, e.host.GetNoOPSession(uint64(e.ID())), proposalBytes); err != nil {
		return wrapError(err)
	}
	return nil
}

// Query queries the state
func (e *Executor) Query(ctx context.Context, input *rsmv1.QueryInput, stream streams.WriteStream[*rsmv1.QueryOutput]) error {
	query := &protocolQuery{
		input:  input,
		stream: stream,
	}
	md, _ := metadata.FromIncomingContext(ctx)
	sync := md["Sync"] != nil
	if sync {
		ctx, cancel := context.WithTimeout(ctx, defaultClientTimeout)
		defer cancel()
		if _, err := e.host.SyncRead(ctx, uint64(e.ID()), query); err != nil {
			return wrapError(err)
		}
	} else {
		if _, err := e.host.StaleRead(uint64(e.ID()), query); err != nil {
			return wrapError(err)
		}
	}
	return nil
}
