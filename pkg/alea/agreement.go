package alea

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func (alea *Alea) applyAgreementMessage(agreementMsg *aleapb.Agreement, source t.NodeID) *events.EventList {
	switch msg := agreementMsg.Type.(type) {
	case *aleapb.Agreement_Aba:
		return alea.applyAbaMessage(msg.Aba, source)
	case *aleapb.Agreement_FillGap:
		// TODO
		panic(fmt.Errorf("TODO"))
	case *aleapb.Agreement_Filler:
		// TODO
		panic(fmt.Errorf("TODO"))
	default:
		panic(fmt.Errorf("unknown Alea Agreement message type: %T", msg))
	}
}
