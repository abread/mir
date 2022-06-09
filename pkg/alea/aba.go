package alea

import (
	"fmt"

	"github.com/hyperledger-labs/mirbft/pkg/events"
	"github.com/hyperledger-labs/mirbft/pkg/pb/aleapb"
	t "github.com/hyperledger-labs/mirbft/pkg/types"
)

func (alea *Alea) applyAbaMessage(abaMsg *aleapb.CobaltABBA, source t.NodeID) *events.EventList {
	// TODO
	panic(fmt.Errorf("TODO"))
}
