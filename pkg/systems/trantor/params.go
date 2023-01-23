package trantor

import (
	"time"

	"github.com/filecoin-project/mir/pkg/alea"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	"github.com/filecoin-project/mir/pkg/util/issutil"

	"github.com/filecoin-project/mir/pkg/mempool/simplemempool"
	"github.com/filecoin-project/mir/pkg/net/libp2p"
	t "github.com/filecoin-project/mir/pkg/types"
)

type Params struct {
	Mempool     *simplemempool.ModuleParams
	Iss         *issutil.ModuleParams
	Alea        *alea.Params // TODO: extract protocol parameters away or figure out a better way to handle this
	Net         libp2p.Params
	ReliableNet *reliablenet.ModuleParams
}

func DefaultParams(initialMembership map[t.NodeID]t.NodeAddress) Params {
	return Params{
		Mempool:     simplemempool.DefaultModuleParams(),
		Iss:         issutil.DefaultParams(initialMembership),
		Alea:        alea.DefaultParams(initialMembership),
		Net:         libp2p.DefaultParams(),
		ReliableNet: reliablenet.DefaultModuleParams(nil),
	}
}

func (p *Params) AdjustSpeed(maxProposeDelay time.Duration) *Params {
	p.Iss.AdjustSpeed(maxProposeDelay)
	return p
}
