package abbapbdsl

import (
	abbatypes "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
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
	UponMessageReceived[*types.Message_FinishMessage](m, func(from types1.NodeID, msg *types.FinishMessage) error {
		return handler(from, msg.Value)
	})
}

func UponInitMessageReceived(m dsl.Module, handler func(from types1.NodeID, roundNumber uint64, estimate bool) error) {
	UponMessageReceived[*types.Message_InitMessage](m, func(from types1.NodeID, msg *types.InitMessage) error {
		return handler(from, msg.RoundNumber, msg.Estimate)
	})
}

func UponAuxMessageReceived(m dsl.Module, handler func(from types1.NodeID, roundNumber uint64, value bool) error) {
	UponMessageReceived[*types.Message_AuxMessage](m, func(from types1.NodeID, msg *types.AuxMessage) error {
		return handler(from, msg.RoundNumber, msg.Value)
	})
}

func UponConfMessageReceived(m dsl.Module, handler func(from types1.NodeID, roundNumber uint64, values abbatypes.ValueSet) error) {
	UponMessageReceived[*types.Message_ConfMessage](m, func(from types1.NodeID, msg *types.ConfMessage) error {
		return handler(from, msg.RoundNumber, msg.Values)
	})
}

func UponCoinMessageReceived(m dsl.Module, handler func(from types1.NodeID, roundNumber uint64, coinShare []uint8) error) {
	UponMessageReceived[*types.Message_CoinMessage](m, func(from types1.NodeID, msg *types.CoinMessage) error {
		return handler(from, msg.RoundNumber, msg.CoinShare)
	})
}
