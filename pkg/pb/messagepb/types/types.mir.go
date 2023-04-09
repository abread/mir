package messagepbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types4 "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	types6 "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/types"
	types5 "github.com/filecoin-project/mir/pkg/pb/aleapb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/availabilitypb/mscpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/bcbpb/types"
	checkpointpb "github.com/filecoin-project/mir/pkg/pb/checkpointpb"
	isspb "github.com/filecoin-project/mir/pkg/pb/isspb"
	messagepb "github.com/filecoin-project/mir/pkg/pb/messagepb"
	ordererspb "github.com/filecoin-project/mir/pkg/pb/ordererspb"
	pingpongpb "github.com/filecoin-project/mir/pkg/pb/pingpongpb"
	types7 "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Message struct {
	DestModule   types.ModuleID
	TraceContext map[string]string
	Type         Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() messagepb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb messagepb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *messagepb.Message_Iss:
		return &Message_Iss{Iss: pb.Iss}
	case *messagepb.Message_Bcb:
		return &Message_Bcb{Bcb: types1.MessageFromPb(pb.Bcb)}
	case *messagepb.Message_MultisigCollector:
		return &Message_MultisigCollector{MultisigCollector: types2.MessageFromPb(pb.MultisigCollector)}
	case *messagepb.Message_Pingpong:
		return &Message_Pingpong{Pingpong: pb.Pingpong}
	case *messagepb.Message_Checkpoint:
		return &Message_Checkpoint{Checkpoint: pb.Checkpoint}
	case *messagepb.Message_SbMessage:
		return &Message_SbMessage{SbMessage: pb.SbMessage}
	case *messagepb.Message_Vcb:
		return &Message_Vcb{Vcb: types3.MessageFromPb(pb.Vcb)}
	case *messagepb.Message_Abba:
		return &Message_Abba{Abba: types4.MessageFromPb(pb.Abba)}
	case *messagepb.Message_Alea:
		return &Message_Alea{Alea: types5.MessageFromPb(pb.Alea)}
	case *messagepb.Message_AleaAgreement:
		return &Message_AleaAgreement{AleaAgreement: types6.MessageFromPb(pb.AleaAgreement)}
	case *messagepb.Message_ReliableNet:
		return &Message_ReliableNet{ReliableNet: types7.MessageFromPb(pb.ReliableNet)}
	}
	return nil
}

type Message_Iss struct {
	Iss *isspb.ISSMessage
}

func (*Message_Iss) isMessage_Type() {}

func (w *Message_Iss) Unwrap() *isspb.ISSMessage {
	return w.Iss
}

func (w *Message_Iss) Pb() messagepb.Message_Type {
	return &messagepb.Message_Iss{Iss: w.Iss}
}

func (*Message_Iss) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_Iss]()}
}

type Message_Bcb struct {
	Bcb *types1.Message
}

func (*Message_Bcb) isMessage_Type() {}

func (w *Message_Bcb) Unwrap() *types1.Message {
	return w.Bcb
}

func (w *Message_Bcb) Pb() messagepb.Message_Type {
	return &messagepb.Message_Bcb{Bcb: (w.Bcb).Pb()}
}

func (*Message_Bcb) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_Bcb]()}
}

type Message_MultisigCollector struct {
	MultisigCollector *types2.Message
}

func (*Message_MultisigCollector) isMessage_Type() {}

func (w *Message_MultisigCollector) Unwrap() *types2.Message {
	return w.MultisigCollector
}

func (w *Message_MultisigCollector) Pb() messagepb.Message_Type {
	return &messagepb.Message_MultisigCollector{MultisigCollector: (w.MultisigCollector).Pb()}
}

func (*Message_MultisigCollector) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_MultisigCollector]()}
}

type Message_Pingpong struct {
	Pingpong *pingpongpb.Message
}

func (*Message_Pingpong) isMessage_Type() {}

func (w *Message_Pingpong) Unwrap() *pingpongpb.Message {
	return w.Pingpong
}

func (w *Message_Pingpong) Pb() messagepb.Message_Type {
	return &messagepb.Message_Pingpong{Pingpong: w.Pingpong}
}

func (*Message_Pingpong) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_Pingpong]()}
}

type Message_Checkpoint struct {
	Checkpoint *checkpointpb.Message
}

func (*Message_Checkpoint) isMessage_Type() {}

func (w *Message_Checkpoint) Unwrap() *checkpointpb.Message {
	return w.Checkpoint
}

func (w *Message_Checkpoint) Pb() messagepb.Message_Type {
	return &messagepb.Message_Checkpoint{Checkpoint: w.Checkpoint}
}

func (*Message_Checkpoint) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_Checkpoint]()}
}

type Message_SbMessage struct {
	SbMessage *ordererspb.SBInstanceMessage
}

func (*Message_SbMessage) isMessage_Type() {}

func (w *Message_SbMessage) Unwrap() *ordererspb.SBInstanceMessage {
	return w.SbMessage
}

func (w *Message_SbMessage) Pb() messagepb.Message_Type {
	return &messagepb.Message_SbMessage{SbMessage: w.SbMessage}
}

func (*Message_SbMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_SbMessage]()}
}

type Message_Vcb struct {
	Vcb *types3.Message
}

func (*Message_Vcb) isMessage_Type() {}

func (w *Message_Vcb) Unwrap() *types3.Message {
	return w.Vcb
}

func (w *Message_Vcb) Pb() messagepb.Message_Type {
	return &messagepb.Message_Vcb{Vcb: (w.Vcb).Pb()}
}

func (*Message_Vcb) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_Vcb]()}
}

type Message_Abba struct {
	Abba *types4.Message
}

func (*Message_Abba) isMessage_Type() {}

func (w *Message_Abba) Unwrap() *types4.Message {
	return w.Abba
}

func (w *Message_Abba) Pb() messagepb.Message_Type {
	return &messagepb.Message_Abba{Abba: (w.Abba).Pb()}
}

func (*Message_Abba) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_Abba]()}
}

type Message_Alea struct {
	Alea *types5.Message
}

func (*Message_Alea) isMessage_Type() {}

func (w *Message_Alea) Unwrap() *types5.Message {
	return w.Alea
}

func (w *Message_Alea) Pb() messagepb.Message_Type {
	return &messagepb.Message_Alea{Alea: (w.Alea).Pb()}
}

func (*Message_Alea) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_Alea]()}
}

type Message_AleaAgreement struct {
	AleaAgreement *types6.Message
}

func (*Message_AleaAgreement) isMessage_Type() {}

func (w *Message_AleaAgreement) Unwrap() *types6.Message {
	return w.AleaAgreement
}

func (w *Message_AleaAgreement) Pb() messagepb.Message_Type {
	return &messagepb.Message_AleaAgreement{AleaAgreement: (w.AleaAgreement).Pb()}
}

func (*Message_AleaAgreement) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_AleaAgreement]()}
}

type Message_ReliableNet struct {
	ReliableNet *types7.Message
}

func (*Message_ReliableNet) isMessage_Type() {}

func (w *Message_ReliableNet) Unwrap() *types7.Message {
	return w.ReliableNet
}

func (w *Message_ReliableNet) Pb() messagepb.Message_Type {
	return &messagepb.Message_ReliableNet{ReliableNet: (w.ReliableNet).Pb()}
}

func (*Message_ReliableNet) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_ReliableNet]()}
}

func MessageFromPb(pb *messagepb.Message) *Message {
	return &Message{
		DestModule:   (types.ModuleID)(pb.DestModule),
		TraceContext: pb.TraceContext,
		Type:         Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *messagepb.Message {
	return &messagepb.Message{
		DestModule:   (string)(m.DestModule),
		TraceContext: m.TraceContext,
		Type:         (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message]()}
}
