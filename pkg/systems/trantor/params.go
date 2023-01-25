package trantor

import (
	"time"

	"github.com/filecoin-project/mir/pkg/alea"
	issconfig "github.com/filecoin-project/mir/pkg/iss/config"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool"
	"github.com/filecoin-project/mir/pkg/net/libp2p"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	t "github.com/filecoin-project/mir/pkg/types"
)

type Params struct {
	Mempool     *simplemempool.ModuleParams
	Iss         *issconfig.ModuleParams
	Net         libp2p.Params
	Alea        *alea.Params // TODO: extract protocol parameters away or figure out a better way to handle this
	ReliableNet *reliablenet.ModuleParams
}

func DefaultParams(initialMembership map[t.NodeID]t.NodeAddress) Params {
	return Params{
		Mempool:     simplemempool.DefaultModuleParams(),
		Iss:         issconfig.DefaultParams(initialMembership),
		Net:         libp2p.DefaultParams(),
		Alea:        alea.DefaultParams(initialMembership),
		ReliableNet: reliablenet.DefaultModuleParams(nil),
	}
}

func (p *Params) AdjustSpeed(maxProposeDelay time.Duration) *Params {
	p.Iss.AdjustSpeed(maxProposeDelay)
	return p
}
