// Code generated by Mir codegen. DO NOT EDIT.

package checkpointpb

import trantorpb "github.com/filecoin-project/mir/pkg/pb/trantorpb"

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_Checkpoint) Unwrap() *Checkpoint {
	return w.Checkpoint
}

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_EpochConfig) Unwrap() *trantorpb.EpochConfig {
	return w.EpochConfig
}

func (w *Event_StableCheckpoint) Unwrap() *StableCheckpoint {
	return w.StableCheckpoint
}

func (w *Event_EpochProgress) Unwrap() *EpochProgress {
	return w.EpochProgress
}
