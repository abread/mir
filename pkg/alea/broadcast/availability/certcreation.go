package availability

import (
	es "github.com/go-errors/errors"
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bccommon"
	"github.com/filecoin-project/mir/pkg/dsl"
	bcpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	availabilitypbdsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/dsl"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	mempoolpbdsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type certCreationState struct {
	nextQueueSlot aleatypes.QueueSlot
}

func includeCertCreation(
	m dsl.Module,
	mc bccommon.ModuleConfig,
	params bccommon.ModuleParams,
	nodeID t.NodeID,
	certDB map[bcpbtypes.Slot]*bcpbtypes.Cert,
) {
	ownQueueIdx := aleatypes.QueueIdx(slices.Index(params.AllNodes, nodeID))
	ownQueueModuleID := bccommon.BcQueueModuleID(mc.Self, ownQueueIdx)
	state := &certCreationState{}

	availabilitypbdsl.UponRequestCert(m, func(origin *availabilitypbtypes.RequestCertOrigin) error {
		// TODO: checkpointing and proper epoch nr handling
		return es.Errorf("alea-bc only supports cert creation requests from alea itself")
	})

	// Alea-specific
	bcpbdsl.UponRequestCert(m, func() error {
		mempoolpbdsl.RequestBatch[struct{}](m, mc.Mempool, tt.EpochNr(0), nil)
		return nil
	})
	mempoolpbdsl.UponNewBatch(m, func(txIds []string, txs []*trantorpbtypes.Transaction, _context *struct{}) error {
		// TODO: propagate txIDs to avoid duplicate computation
		bcqueuepbdsl.InputValue(m, ownQueueModuleID, state.nextQueueSlot, txs)
		state.nextQueueSlot++
		return nil
	})

	bcqueuepbdsl.UponDeliver(m, func(cert *bcpbtypes.Cert) error {
		if cert.Slot.QueueIdx == ownQueueIdx {
			// free slot in own queue immediately
			// this is required to avoid a concurrency issue between FreeSlot and RequestCert (FreeSlot
			// sent at the same time of RequestCert, but InputValue may technically reach queue before FreeSlot)
			// modring is lazy in deleting modules so timers will progress until the slot is reused

			bcqueuepbdsl.FreeSlot(m, ownQueueModuleID, cert.Slot.QueueSlot)
		}

		// register that we received this batch
		certDB[*cert.Slot] = cert

		bcpbdsl.DeliverCert(m, mc.Consumer, cert)
		return nil
	})
}
