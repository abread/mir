package ordererpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/pb/ordererpb/types"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing net messages.

func UponMessageReceived[W types.Message_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	dsl1.UponMessageReceived[*types2.Message_Orderer](m, func(from types1.NodeID, msg *types.Message) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}