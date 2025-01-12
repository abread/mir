// Code generated by Mir codegen. DO NOT EDIT.

package threshcheckpointpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	threshcheckpointpb "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb"
	types1 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/trantor/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() threshcheckpointpb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb threshcheckpointpb.Message_Type) Message_Type {
	if pb == nil {
		return nil
	}
	switch pb := pb.(type) {
	case *threshcheckpointpb.Message_Checkpoint:
		return &Message_Checkpoint{Checkpoint: CheckpointFromPb(pb.Checkpoint)}
	}
	return nil
}

type Message_Checkpoint struct {
	Checkpoint *Checkpoint
}

func (*Message_Checkpoint) isMessage_Type() {}

func (w *Message_Checkpoint) Unwrap() *Checkpoint {
	return w.Checkpoint
}

func (w *Message_Checkpoint) Pb() threshcheckpointpb.Message_Type {
	if w == nil {
		return nil
	}
	if w.Checkpoint == nil {
		return &threshcheckpointpb.Message_Checkpoint{}
	}
	return &threshcheckpointpb.Message_Checkpoint{Checkpoint: (w.Checkpoint).Pb()}
}

func (*Message_Checkpoint) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*threshcheckpointpb.Message_Checkpoint]()}
}

func MessageFromPb(pb *threshcheckpointpb.Message) *Message {
	if pb == nil {
		return nil
	}
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *threshcheckpointpb.Message {
	if m == nil {
		return nil
	}
	pbMessage := &threshcheckpointpb.Message{}
	{
		if m.Type != nil {
			pbMessage.Type = (m.Type).Pb()
		}
	}

	return pbMessage
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*threshcheckpointpb.Message]()}
}

type Checkpoint struct {
	Epoch          types.EpochNr
	Sn             types.SeqNr
	SnapshotHash   []uint8
	SignatureShare tctypes.SigShare
}

func CheckpointFromPb(pb *threshcheckpointpb.Checkpoint) *Checkpoint {
	if pb == nil {
		return nil
	}
	return &Checkpoint{
		Epoch:          (types.EpochNr)(pb.Epoch),
		Sn:             (types.SeqNr)(pb.Sn),
		SnapshotHash:   pb.SnapshotHash,
		SignatureShare: (tctypes.SigShare)(pb.SignatureShare),
	}
}

func (m *Checkpoint) Pb() *threshcheckpointpb.Checkpoint {
	if m == nil {
		return nil
	}
	pbMessage := &threshcheckpointpb.Checkpoint{}
	{
		pbMessage.Epoch = (uint64)(m.Epoch)
		pbMessage.Sn = (uint64)(m.Sn)
		pbMessage.SnapshotHash = m.SnapshotHash
		pbMessage.SignatureShare = ([]uint8)(m.SignatureShare)
	}

	return pbMessage
}

func (*Checkpoint) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*threshcheckpointpb.Checkpoint]()}
}

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() threshcheckpointpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb threshcheckpointpb.Event_Type) Event_Type {
	if pb == nil {
		return nil
	}
	switch pb := pb.(type) {
	case *threshcheckpointpb.Event_StableCheckpoint:
		return &Event_StableCheckpoint{StableCheckpoint: StableCheckpointFromPb(pb.StableCheckpoint)}
	}
	return nil
}

type Event_StableCheckpoint struct {
	StableCheckpoint *StableCheckpoint
}

func (*Event_StableCheckpoint) isEvent_Type() {}

func (w *Event_StableCheckpoint) Unwrap() *StableCheckpoint {
	return w.StableCheckpoint
}

func (w *Event_StableCheckpoint) Pb() threshcheckpointpb.Event_Type {
	if w == nil {
		return nil
	}
	if w.StableCheckpoint == nil {
		return &threshcheckpointpb.Event_StableCheckpoint{}
	}
	return &threshcheckpointpb.Event_StableCheckpoint{StableCheckpoint: (w.StableCheckpoint).Pb()}
}

func (*Event_StableCheckpoint) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*threshcheckpointpb.Event_StableCheckpoint]()}
}

func EventFromPb(pb *threshcheckpointpb.Event) *Event {
	if pb == nil {
		return nil
	}
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *threshcheckpointpb.Event {
	if m == nil {
		return nil
	}
	pbMessage := &threshcheckpointpb.Event{}
	{
		if m.Type != nil {
			pbMessage.Type = (m.Type).Pb()
		}
	}

	return pbMessage
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*threshcheckpointpb.Event]()}
}

type StableCheckpoint struct {
	Sn        types.SeqNr
	Snapshot  *types1.StateSnapshot
	Signature tctypes.FullSig
}

func StableCheckpointFromPb(pb *threshcheckpointpb.StableCheckpoint) *StableCheckpoint {
	if pb == nil {
		return nil
	}
	return &StableCheckpoint{
		Sn:        (types.SeqNr)(pb.Sn),
		Snapshot:  types1.StateSnapshotFromPb(pb.Snapshot),
		Signature: (tctypes.FullSig)(pb.Signature),
	}
}

func (m *StableCheckpoint) Pb() *threshcheckpointpb.StableCheckpoint {
	if m == nil {
		return nil
	}
	pbMessage := &threshcheckpointpb.StableCheckpoint{}
	{
		pbMessage.Sn = (uint64)(m.Sn)
		if m.Snapshot != nil {
			pbMessage.Snapshot = (m.Snapshot).Pb()
		}
		pbMessage.Signature = ([]uint8)(m.Signature)
	}

	return pbMessage
}

func (*StableCheckpoint) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*threshcheckpointpb.StableCheckpoint]()}
}

type InstanceParams struct {
	Membership       *types1.Membership
	LeaderPolicyData []uint8
	EpochConfig      *types1.EpochConfig
	Threshold        uint64
}

func InstanceParamsFromPb(pb *threshcheckpointpb.InstanceParams) *InstanceParams {
	if pb == nil {
		return nil
	}
	return &InstanceParams{
		Membership:       types1.MembershipFromPb(pb.Membership),
		LeaderPolicyData: pb.LeaderPolicyData,
		EpochConfig:      types1.EpochConfigFromPb(pb.EpochConfig),
		Threshold:        pb.Threshold,
	}
}

func (m *InstanceParams) Pb() *threshcheckpointpb.InstanceParams {
	if m == nil {
		return nil
	}
	pbMessage := &threshcheckpointpb.InstanceParams{}
	{
		if m.Membership != nil {
			pbMessage.Membership = (m.Membership).Pb()
		}
		pbMessage.LeaderPolicyData = m.LeaderPolicyData
		if m.EpochConfig != nil {
			pbMessage.EpochConfig = (m.EpochConfig).Pb()
		}
		pbMessage.Threshold = m.Threshold
	}

	return pbMessage
}

func (*InstanceParams) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*threshcheckpointpb.InstanceParams]()}
}
