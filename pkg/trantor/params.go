package trantor

import (
	"time"

	"github.com/filecoin-project/mir/pkg/alea"
	"github.com/filecoin-project/mir/pkg/availability/multisigcollector"
	issconfig "github.com/filecoin-project/mir/pkg/iss/config"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool"
	"github.com/filecoin-project/mir/pkg/net/libp2p"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

type Params struct {
	Protocol string

	Mempool      *simplemempool.ModuleParams
	Iss          *issconfig.ModuleParams
	Net          libp2p.Params
	Alea         *alea.Params // TODO: extract protocol parameters away or figure out a better way to handle this
	ReliableNet  *reliablenet.ModuleParams
	Availability multisigcollector.ModuleParams
}

func DefaultParams(initialMembership *trantorpbtypes.Membership) Params {
	allNodes := maputil.GetSortedKeys(initialMembership.Nodes)
	p := Params{
		Mempool:      simplemempool.DefaultModuleParams(),
		Iss:          issconfig.DefaultParams(initialMembership),
		Net:          libp2p.DefaultParams(),
		Alea:         alea.DefaultParams(initialMembership),
		ReliableNet:  reliablenet.DefaultModuleParams(allNodes),
		Availability: multisigcollector.DefaultParamsTemplate(),
	}

	return p
}

func (p *Params) AdjustBatchSize(batchSize int, txSize int) *Params {
	p.Mempool.MaxTransactionsInBatch = batchSize
	p.Mempool.MaxPayloadInBatch = batchSize * txSize * 105 / 100 // 5% overhead

	// ensure network messages can accommodate the chosen batch size
	batchAdjustedMaxMsgSize := p.Mempool.MaxPayloadInBatch * 105 / 100 // 5% overhead
	if p.Net.MaxMessageSize < batchAdjustedMaxMsgSize {
		p.Net.MaxMessageSize = batchAdjustedMaxMsgSize
	}

	return p
}

func (p *Params) AdjustSpeed(maxProposeDelay time.Duration) *Params {
	p.Iss.AdjustSpeed(maxProposeDelay)
	p.Iss.MaxProposeDelay = 0                // simplemempool will handle this
	p.Mempool.BatchTimeout = maxProposeDelay // TODO: account for processing time
	return p
}
