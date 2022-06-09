package alea

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func (alea *Alea) applyAbaMessage(abaMsg *aleapb.CobaltABBA, source t.NodeID) *events.EventList {
	// TODO
	panic(fmt.Errorf("TODO"))
}
