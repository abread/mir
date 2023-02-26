package abbapbdsl

import (
	abbatypes "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing net messages.

func UponMessageReceived[W types.Message_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	dsl1.UponMessageReceived[*types2.Message_Abba](m, func(from types1.NodeID, msg *types.Message) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}

func UponFinishMessageReceived(m dsl.Module, handler func(from types1.NodeID, value bool) error) {
	UponMessageReceived[*types.Message_Finish](m, func(from types1.NodeID, msg *types.FinishMessage) error {
		return handler(from, msg.Value)
	})
}

func UponRoundMessageReceived[W types.RoundMessage_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	UponMessageReceived[*types.Message_Round](m, func(from types1.NodeID, msg *types.RoundMessage) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}

func UponRoundInitMessageReceived(m dsl.Module, handler func(from types1.NodeID, estimate bool) error) {
	UponRoundMessageReceived[*types.RoundMessage_Init](m, func(from types1.NodeID, msg *types.RoundInitMessage) error {
		return handler(from, msg.Estimate)
	})
}

func UponRoundAuxMessageReceived(m dsl.Module, handler func(from types1.NodeID, value bool) error) {
	UponRoundMessageReceived[*types.RoundMessage_Aux](m, func(from types1.NodeID, msg *types.RoundAuxMessage) error {
		return handler(from, msg.Value)
	})
}

func UponRoundConfMessageReceived(m dsl.Module, handler func(from types1.NodeID, values abbatypes.ValueSet) error) {
	UponRoundMessageReceived[*types.RoundMessage_Conf](m, func(from types1.NodeID, msg *types.RoundConfMessage) error {
		return handler(from, msg.Values)
	})
}

func UponRoundCoinMessageReceived(m dsl.Module, handler func(from types1.NodeID, coinShare tctypes.SigShare) error) {
	UponRoundMessageReceived[*types.RoundMessage_Coin](m, func(from types1.NodeID, msg *types.RoundCoinMessage) error {
		return handler(from, msg.CoinShare)
	})
}