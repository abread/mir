package modringpb

import (
	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
)

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_Free) Unwrap() *FreeSubmodule {
	return w.Free
}

func (w *Event_Freed) Unwrap() *FreedSubmodule {
	return w.Freed
}

type FreeSubmoduleOrigin_Type = isFreeSubmoduleOrigin_Type

type FreeSubmoduleOrigin_TypeWrapper[T any] interface {
	FreeSubmoduleOrigin_Type
	Unwrap() *T
}

func (w *FreeSubmoduleOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *FreeSubmoduleOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}
