package alea

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type void struct{}
type processSet map[t.NodeID]void

type cobaltAbbaInst struct {
	finished      bool
	finishSent    bool
	finishSenders map[bit]processSet

	r        uint64 // round number (!= alea agreement round number)
	values   BitSet
	estimate bit

	initSenders map[bit]processSet

	auxSent    bool
	auxSenders map[t.NodeID]BitSet

	confSent    bool
	confSenders map[t.NodeID]BitSet

	coinShareSent bool
	coinShares    map[t.NodeID]threshSigShare
}

func (alea *Alea) applyAgreementMessage(agreementMsg *aleapb.Agreement, source t.NodeID) *events.EventList {
	id := MsgIdFromDomain(agreementMsg.Instance)

	inst, ok := alea.abbaInstances[id]

	if ok {
		return inst.applyMessage(alea, id, agreementMsg.Msg, source)
	}

	// TODO: silently drop?
}

func (inst *cobaltAbbaInst) applyMessage(alea *Alea, id MsgId, abbaMsg *aleapb.CobaltABBA, source t.NodeID) *events.EventList {
	if inst.finished {
		return &events.EventList{}
	}

	switch msg := abbaMsg.Type.(type) {
	case *aleapb.CobaltABBA_Finish:
		return inst.applyFinishMessage(alea, id, msg.Finish, source)
	case *aleapb.CobaltABBA_Init:
		return inst.applyInitMessage(alea, id, msg.Init, source)
	case *aleapb.CobaltABBA_Aux:
		return inst.applyAuxMessage(alea, id, msg.Aux, source)
	case *aleapb.CobaltABBA_Conf:
		return inst.applyConfMessage(alea, id, msg.Conf, source)
	case *aleapb.CobaltABBA_CoinShare:
		return inst.applyCoinShareMessage(alea, id, msg.CoinShare, source)
	default:
		panic(fmt.Errorf("unknown Cobalt ABBA message type: %T", msg))
	}
}

func (inst *cobaltAbbaInst) applyFinishMessage(alea *Alea, id MsgId, msg *aleapb.CobaltFinish, source t.NodeID) *events.EventList {
	v := bit(msg.Value)
	inst.finishSenders[v][source] = void{}

	strongSupportSize := len(alea.config.Membership) - alea.config.F + 1
	weakSupportSize := alea.config.F + 1

	if len(inst.finishSenders[v]) == weakSupportSize && !inst.finishSent {
		// weak support for <FINISH, v>
		inst.finishSent = true
		return alea.broadcastMessage(CobaltAbbaFinish(id, v))
	} else if len(inst.finishSenders[v]) == strongSupportSize {
		// strong support for <FINISH, v>
		inst.finished = true
		return (&events.EventList{}).PushBack(AbaDeliver(id, v))
	} else {
		return &events.EventList{}
	}
}

func (inst *cobaltAbbaInst) applyInitMessage(alea *Alea, id MsgId, msg *aleapb.CobaltInit, source t.NodeID) *events.EventList {
	if inst.r != msg.Round {
		// not the current round, ignore message (will be retransmitted if it just came before anything important)
		return &events.EventList{}
	}

	v := bit(msg.Value)
	inst.initSenders[v][source] = void{}

	strongSupportSize := len(alea.config.Membership) - alea.config.F + 1
	weakSupportSize := alea.config.F + 1

	if len(inst.initSenders[v]) == weakSupportSize {
		// weak support for <INIT, r, v>
		return alea.broadcastMessage(CobaltAbbaInit(id, inst.r, v)) // TODO: may broadcast init more than once when dupe packets are received
	} else if len(inst.initSenders[v]) == strongSupportSize && !inst.auxSent {
		// strong support for <INIT, r, v>
		inst.values = inst.values.WithBit(v)
		inst.auxSent = true
		return alea.broadcastMessage(CobaltAbbaAux(id, inst.r, v))
	} else {
		return &events.EventList{}
	}
}

func (inst *cobaltAbbaInst) applyAuxMessage(alea *Alea, id MsgId, msg *aleapb.CobaltAux, source t.NodeID) *events.EventList {
	if inst.r != msg.Round {
		// not the current round, ignore message (will be retransmitted if it just came before anything important)
		return &events.EventList{}
	}

	v := bit(msg.Value)
	mapPutIfNotPresent(inst.auxSenders, source, BITS_NONE.WithBit(v))
	//inst.auxSenders[source] = inst.auxSenders[source].WithBit(v) // TODO: do I need to deal with byz nodes sending different AUX messages?

	strongSupportSize := len(alea.config.Membership) - alea.config.F
	if inst.recvAuxWithValidValuesCount() >= strongSupportSize && !inst.confSent {
		inst.confSent = true
		return alea.broadcastMessage(CobaltAbbaConf(id, inst.r, inst.values))
	} else {
		return &events.EventList{}
	}
}

func (inst *cobaltAbbaInst) recvAuxWithValidValuesCount() int {
	count := 0
	for _, auxBitsRecvd := range inst.auxSenders {
		if !auxBitsRecvd.Intersection(inst.values).IsEmpty() {
			// we received an AUX message with some acceptable value
			count += 1
		}
	}

	return count
}

func (inst *cobaltAbbaInst) applyConfMessage(alea *Alea, id MsgId, msg *aleapb.CobaltConf, source t.NodeID) *events.EventList {
	if inst.r != msg.Round {
		// not the current round, ignore message (will be retransmitted if it just came before anything important)
		return &events.EventList{}
	}

	config := BitSetFromCobaltValueSetPb(msg.Config)
	mapPutIfNotPresent(inst.confSenders, source, config)
	//inst.confSenders[source] = config // TODO: do I need to deal with byz nodes sending different configs?

	strongSupportSize := len(alea.config.Membership) - alea.config.F
	if inst.recvConfWithValidValuesCount() >= strongSupportSize && !inst.coinShareSent {
		inst.coinShareSent = true
		// TODO: request coin sign (thresh-sign-share Hash(somestring, msgid, r))
		return alea.broadcastMessage(CobaltAbbaCoinShare(id, inst.r, threshSigShare{}))
	}

	return &events.EventList{}
}

func (inst *cobaltAbbaInst) recvConfWithValidValuesCount() int {
	count := 0
	valuesSubsets := inst.values.PowerSet()

	for _, bs := range inst.confSenders {
		for _, bs2 := range valuesSubsets {
			if bs == bs2 {
				count += 1
			}
		}
	}

	return count
}

func (inst *cobaltAbbaInst) applyCoinShareMessage(alea *Alea, id MsgId, msg *aleapb.CobaltCoinShare, source t.NodeID) *events.EventList {
	if inst.r != msg.Round {
		// not the current round, ignore message (will be retransmitted if it just came before anything important)
		return &events.EventList{}
	}

	// passively collect coin shares even if we can't toss it yet
	if share, err := threshSigShareFromDomain(msg.CoinShare); err == nil {
		mapPutIfNotPresent(inst.coinShares, source, share)
	}

	if !inst.coinShareSent {
		// we don't yet fulfill the condition for tossing a coin
		return &events.EventList{}
	}

	tossedCoin := bit(false) // TODO: thresh-sig-combine + modulo arith?

	if inst.values == BITS_BOTH {
		inst.estimate = tossedCoin
	} else if inst.values == BITS_ONE {
		inst.estimate = bit(true)
	} else if inst.values == BITS_ZERO {
		inst.estimate = bit(false)
	}

	if inst.values == bitToBitset(tossedCoin) {
		return alea.broadcastMessage(CobaltAbbaFinish(id, tossedCoin))
	} else {
		inst.r += 1

		// reset round data
		clearMap(inst.initSenders)
		inst.auxSent = false
		clearMap(inst.auxSenders)
		inst.confSent = false
		clearMap(inst.confSenders)
		inst.coinShareSent = false
		clearMap(inst.coinShares)

		return alea.broadcastMessage(CobaltAbbaInit(id, inst.r, inst.estimate))
	}
}

func mapPutIfNotPresent[K comparable, V any](m map[K]V, k K, v V) {
	if _, exists := m[k]; !exists {
		m[k] = v
	}
}

func clearMap[K comparable, V any](m map[K]V) {
	for k := range m {
		delete(m, k)
	}
}
