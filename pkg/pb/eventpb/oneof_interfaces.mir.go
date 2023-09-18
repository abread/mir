// Code generated by Mir codegen. DO NOT EDIT.

package eventpb

import (
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"

	abbapb "github.com/filecoin-project/mir/pkg/pb/abbapb"
	agevents "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents"
	bcpb "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	bcqueuepb "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb"
	directorpb "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb"
	apppb "github.com/filecoin-project/mir/pkg/pb/apppb"
	availabilitypb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	batchdbpb "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb"
	batchfetcherpb "github.com/filecoin-project/mir/pkg/pb/batchfetcherpb"
	bcbpb "github.com/filecoin-project/mir/pkg/pb/bcbpb"
	checkpointpb "github.com/filecoin-project/mir/pkg/pb/checkpointpb"
	chkpvalidatorpb "github.com/filecoin-project/mir/pkg/pb/checkpointpb/chkpvalidatorpb"
	cryptopb "github.com/filecoin-project/mir/pkg/pb/cryptopb"
	factorypb "github.com/filecoin-project/mir/pkg/pb/factorypb"
	hasherpb "github.com/filecoin-project/mir/pkg/pb/hasherpb"
	isspb "github.com/filecoin-project/mir/pkg/pb/isspb"
	mempoolpb "github.com/filecoin-project/mir/pkg/pb/mempoolpb"
	ordererpb "github.com/filecoin-project/mir/pkg/pb/ordererpb"
	pprepvalidatorpb "github.com/filecoin-project/mir/pkg/pb/ordererpb/pprepvalidatorpb"
	pingpongpb "github.com/filecoin-project/mir/pkg/pb/pingpongpb"
	reliablenetpb "github.com/filecoin-project/mir/pkg/pb/reliablenetpb"
	testerpb "github.com/filecoin-project/mir/pkg/pb/testerpb"
	threshcheckpointpb "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb"
	threshchkpvalidatorpb "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb/threshchkpvalidatorpb"
	threshcryptopb "github.com/filecoin-project/mir/pkg/pb/threshcryptopb"
	transportpb "github.com/filecoin-project/mir/pkg/pb/transportpb"
	vcbpb "github.com/filecoin-project/mir/pkg/pb/vcbpb"
)

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_Init) Unwrap() *Init {
	return w.Init
}

func (w *Event_Timer) Unwrap() *TimerEvent {
	return w.Timer
}

func (w *Event_Hasher) Unwrap() *hasherpb.Event {
	return w.Hasher
}

func (w *Event_Bcb) Unwrap() *bcbpb.Event {
	return w.Bcb
}

func (w *Event_Mempool) Unwrap() *mempoolpb.Event {
	return w.Mempool
}

func (w *Event_Availability) Unwrap() *availabilitypb.Event {
	return w.Availability
}

func (w *Event_BatchDb) Unwrap() *batchdbpb.Event {
	return w.BatchDb
}

func (w *Event_BatchFetcher) Unwrap() *batchfetcherpb.Event {
	return w.BatchFetcher
}

func (w *Event_ThreshCrypto) Unwrap() *threshcryptopb.Event {
	return w.ThreshCrypto
}

func (w *Event_Checkpoint) Unwrap() *checkpointpb.Event {
	return w.Checkpoint
}

func (w *Event_Factory) Unwrap() *factorypb.Event {
	return w.Factory
}

func (w *Event_Iss) Unwrap() *isspb.Event {
	return w.Iss
}

func (w *Event_Orderer) Unwrap() *ordererpb.Event {
	return w.Orderer
}

func (w *Event_Crypto) Unwrap() *cryptopb.Event {
	return w.Crypto
}

func (w *Event_App) Unwrap() *apppb.Event {
	return w.App
}

func (w *Event_Transport) Unwrap() *transportpb.Event {
	return w.Transport
}

func (w *Event_ChkpValidator) Unwrap() *chkpvalidatorpb.Event {
	return w.ChkpValidator
}

func (w *Event_PprepValiadtor) Unwrap() *pprepvalidatorpb.Event {
	return w.PprepValiadtor
}

func (w *Event_ReliableNet) Unwrap() *reliablenetpb.Event {
	return w.ReliableNet
}

func (w *Event_Vcb) Unwrap() *vcbpb.Event {
	return w.Vcb
}

func (w *Event_Abba) Unwrap() *abbapb.Event {
	return w.Abba
}

func (w *Event_Threshcheckpoint) Unwrap() *threshcheckpointpb.Event {
	return w.Threshcheckpoint
}

func (w *Event_ThreshchkpValidator) Unwrap() *threshchkpvalidatorpb.Event {
	return w.ThreshchkpValidator
}

func (w *Event_AleaBc) Unwrap() *bcpb.Event {
	return w.AleaBc
}

func (w *Event_AleaBcqueue) Unwrap() *bcqueuepb.Event {
	return w.AleaBcqueue
}

func (w *Event_AleaAgreement) Unwrap() *agevents.Event {
	return w.AleaAgreement
}

func (w *Event_AleaDirector) Unwrap() *directorpb.Event {
	return w.AleaDirector
}

func (w *Event_PingPong) Unwrap() *pingpongpb.Event {
	return w.PingPong
}

func (w *Event_TestingString) Unwrap() *wrapperspb.StringValue {
	return w.TestingString
}

func (w *Event_TestingUint) Unwrap() *wrapperspb.UInt64Value {
	return w.TestingUint
}

func (w *Event_Tester) Unwrap() *testerpb.Tester {
	return w.Tester
}

type TimerEvent_Type = isTimerEvent_Type

type TimerEvent_TypeWrapper[T any] interface {
	TimerEvent_Type
	Unwrap() *T
}

func (w *TimerEvent_Delay) Unwrap() *TimerDelay {
	return w.Delay
}

func (w *TimerEvent_Repeat) Unwrap() *TimerRepeat {
	return w.Repeat
}

func (w *TimerEvent_GarbageCollect) Unwrap() *TimerGarbageCollect {
	return w.GarbageCollect
}
