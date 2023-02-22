// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: raft/v1beta1/raft.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Sent from a candidate to all peers in the quorum to elect a new Raft leader.
type VoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`                                       // the term of the election
	Candidate    uint32 `protobuf:"varint,2,opt,name=candidate,proto3" json:"candidate,omitempty"`                             // the PID of the candidate requesting the vote
	LastLogIndex uint64 `protobuf:"varint,3,opt,name=last_log_index,json=lastLogIndex,proto3" json:"last_log_index,omitempty"` // the last log index of the candidate's log
	LastLogTerm  uint64 `protobuf:"varint,4,opt,name=last_log_term,json=lastLogTerm,proto3" json:"last_log_term,omitempty"`    // the log of the last entry in the candidate's log
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_v1beta1_raft_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_raft_v1beta1_raft_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_raft_v1beta1_raft_proto_rawDescGZIP(), []int{0}
}

func (x *VoteRequest) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteRequest) GetCandidate() uint32 {
	if x != nil {
		return x.Candidate
	}
	return 0
}

func (x *VoteRequest) GetLastLogIndex() uint64 {
	if x != nil {
		return x.LastLogIndex
	}
	return 0
}

func (x *VoteRequest) GetLastLogTerm() uint64 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

// Sent from peers in the quorum in response to a vote request to bring the candidate's
// state up to date or to elect the candidate as leader for the term.
type VoteReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Remote  uint32 `protobuf:"varint,1,opt,name=remote,proto3" json:"remote,omitempty"`   // the PID of the voter
	Term    uint64 `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`       // the current term of voter
	Granted bool   `protobuf:"varint,3,opt,name=granted,proto3" json:"granted,omitempty"` // if the vote is granted or not
}

func (x *VoteReply) Reset() {
	*x = VoteReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_v1beta1_raft_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteReply) ProtoMessage() {}

func (x *VoteReply) ProtoReflect() protoreflect.Message {
	mi := &file_raft_v1beta1_raft_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteReply.ProtoReflect.Descriptor instead.
func (*VoteReply) Descriptor() ([]byte, []int) {
	return file_raft_v1beta1_raft_proto_rawDescGZIP(), []int{1}
}

func (x *VoteReply) GetRemote() uint32 {
	if x != nil {
		return x.Remote
	}
	return 0
}

func (x *VoteReply) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteReply) GetGranted() bool {
	if x != nil {
		return x.Granted
	}
	return false
}

// Sent from the leader to the peers in the quorum to update their logs, or if no
// entries are sent, as a heartbeat message.
type AppendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         uint64      `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`                                       // the term of the leader
	Leader       uint32      `protobuf:"varint,2,opt,name=leader,proto3" json:"leader,omitempty"`                                   // the PID of the leader
	PrevLogIndex uint64      `protobuf:"varint,3,opt,name=prev_log_index,json=prevLogIndex,proto3" json:"prev_log_index,omitempty"` // the index of the previous log entry in the leader's log
	PrevLogTerm  uint64      `protobuf:"varint,4,opt,name=prev_log_term,json=prevLogTerm,proto3" json:"prev_log_term,omitempty"`    // the term of the previous log entry in the leader's log
	LeaderCommit uint64      `protobuf:"varint,5,opt,name=leader_commit,json=leaderCommit,proto3" json:"leader_commit,omitempty"`   // the index of the last commited entry in the leader's log
	Entries      []*LogEntry `protobuf:"bytes,6,rep,name=entries,proto3" json:"entries,omitempty"`                                  // the entries to be appended to the follower's log
}

func (x *AppendRequest) Reset() {
	*x = AppendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_v1beta1_raft_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendRequest) ProtoMessage() {}

func (x *AppendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_raft_v1beta1_raft_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendRequest.ProtoReflect.Descriptor instead.
func (*AppendRequest) Descriptor() ([]byte, []int) {
	return file_raft_v1beta1_raft_proto_rawDescGZIP(), []int{2}
}

func (x *AppendRequest) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendRequest) GetLeader() uint32 {
	if x != nil {
		return x.Leader
	}
	return 0
}

func (x *AppendRequest) GetPrevLogIndex() uint64 {
	if x != nil {
		return x.PrevLogIndex
	}
	return 0
}

func (x *AppendRequest) GetPrevLogTerm() uint64 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

func (x *AppendRequest) GetLeaderCommit() uint64 {
	if x != nil {
		return x.LeaderCommit
	}
	return 0
}

func (x *AppendRequest) GetEntries() []*LogEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

// Sent from followers back to the leader to acknowledge the append entries or heartbeat
// and to update the leader with their local state.
type AppendReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Remote      uint32 `protobuf:"varint,1,opt,name=remote,proto3" json:"remote,omitempty"`                              // the PID of the follower
	Term        uint64 `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`                                  // the term of the follower
	Success     bool   `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`                            // if the operation was successful
	Index       uint64 `protobuf:"varint,4,opt,name=index,proto3" json:"index,omitempty"`                                // the last index of the follower's log
	CommitIndex uint64 `protobuf:"varint,5,opt,name=commit_index,json=commitIndex,proto3" json:"commit_index,omitempty"` // the commit index of the follower's log
}

func (x *AppendReply) Reset() {
	*x = AppendReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_v1beta1_raft_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendReply) ProtoMessage() {}

func (x *AppendReply) ProtoReflect() protoreflect.Message {
	mi := &file_raft_v1beta1_raft_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendReply.ProtoReflect.Descriptor instead.
func (*AppendReply) Descriptor() ([]byte, []int) {
	return file_raft_v1beta1_raft_proto_rawDescGZIP(), []int{3}
}

func (x *AppendReply) GetRemote() uint32 {
	if x != nil {
		return x.Remote
	}
	return 0
}

func (x *AppendReply) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *AppendReply) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *AppendReply) GetCommitIndex() uint64 {
	if x != nil {
		return x.CommitIndex
	}
	return 0
}

var File_raft_v1beta1_raft_proto protoreflect.FileDescriptor

var file_raft_v1beta1_raft_proto_rawDesc = []byte{
	0x0a, 0x17, 0x72, 0x61, 0x66, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x72,
	0x61, 0x66, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x72, 0x61, 0x66, 0x74, 0x2e,
	0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x16, 0x72, 0x61, 0x66, 0x74, 0x2f, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x89, 0x01, 0x0a, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74,
	0x65, 0x72, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x12, 0x24, 0x0a, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4c,
	0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x22, 0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x5f,
	0x6c, 0x6f, 0x67, 0x5f, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b,
	0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x51, 0x0a, 0x09, 0x56,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04,
	0x74, 0x65, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x22, 0xdc,
	0x01, 0x0a, 0x0d, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04,
	0x74, 0x65, 0x72, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x0e,
	0x70, 0x72, 0x65, 0x76, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x12, 0x22, 0x0a, 0x0d, 0x70, 0x72, 0x65, 0x76, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x74,
	0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x70, 0x72, 0x65, 0x76, 0x4c,
	0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x23, 0x0a, 0x0d, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x6c,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x30, 0x0a, 0x07, 0x65,
	0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x72,
	0x61, 0x66, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x8c, 0x01,
	0x0a, 0x0b, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a,
	0x06, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x72,
	0x65, 0x6d, 0x6f, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x32, 0x9a, 0x01, 0x0a,
	0x04, 0x52, 0x61, 0x66, 0x74, 0x12, 0x43, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x56, 0x6f, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x56,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x0d, 0x41, 0x70,
	0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x1b, 0x2e, 0x72, 0x61,
	0x66, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e,
	0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x61, 0x6c, 0x69, 0x6f, 0x2f, 0x65, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x72, 0x61, 0x66, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_raft_v1beta1_raft_proto_rawDescOnce sync.Once
	file_raft_v1beta1_raft_proto_rawDescData = file_raft_v1beta1_raft_proto_rawDesc
)

func file_raft_v1beta1_raft_proto_rawDescGZIP() []byte {
	file_raft_v1beta1_raft_proto_rawDescOnce.Do(func() {
		file_raft_v1beta1_raft_proto_rawDescData = protoimpl.X.CompressGZIP(file_raft_v1beta1_raft_proto_rawDescData)
	})
	return file_raft_v1beta1_raft_proto_rawDescData
}

var file_raft_v1beta1_raft_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_raft_v1beta1_raft_proto_goTypes = []interface{}{
	(*VoteRequest)(nil),   // 0: raft.v1beta1.VoteRequest
	(*VoteReply)(nil),     // 1: raft.v1beta1.VoteReply
	(*AppendRequest)(nil), // 2: raft.v1beta1.AppendRequest
	(*AppendReply)(nil),   // 3: raft.v1beta1.AppendReply
	(*LogEntry)(nil),      // 4: raft.v1beta1.LogEntry
}
var file_raft_v1beta1_raft_proto_depIdxs = []int32{
	4, // 0: raft.v1beta1.AppendRequest.entries:type_name -> raft.v1beta1.LogEntry
	0, // 1: raft.v1beta1.Raft.RequestVote:input_type -> raft.v1beta1.VoteRequest
	2, // 2: raft.v1beta1.Raft.AppendEntries:input_type -> raft.v1beta1.AppendRequest
	1, // 3: raft.v1beta1.Raft.RequestVote:output_type -> raft.v1beta1.VoteReply
	3, // 4: raft.v1beta1.Raft.AppendEntries:output_type -> raft.v1beta1.AppendReply
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_raft_v1beta1_raft_proto_init() }
func file_raft_v1beta1_raft_proto_init() {
	if File_raft_v1beta1_raft_proto != nil {
		return
	}
	file_raft_v1beta1_log_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_raft_v1beta1_raft_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raft_v1beta1_raft_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raft_v1beta1_raft_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raft_v1beta1_raft_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_raft_v1beta1_raft_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_raft_v1beta1_raft_proto_goTypes,
		DependencyIndexes: file_raft_v1beta1_raft_proto_depIdxs,
		MessageInfos:      file_raft_v1beta1_raft_proto_msgTypes,
	}.Build()
	File_raft_v1beta1_raft_proto = out.File
	file_raft_v1beta1_raft_proto_rawDesc = nil
	file_raft_v1beta1_raft_proto_goTypes = nil
	file_raft_v1beta1_raft_proto_depIdxs = nil
}
