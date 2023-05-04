package trantor

import (
	"time"

	"github.com/filecoin-project/mir/pkg/alea"
	issconfig "github.com/filecoin-project/mir/pkg/iss/config"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool"
	"github.com/filecoin-project/mir/pkg/net/libp2p"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	"github.com/filecoin-project/mir/pkg/threshcrypto"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

type Params struct {
	Mempool      *simplemempool.ModuleParams
	Iss          *issconfig.ModuleParams
	Net          libp2p.Params
	Alea         *alea.Params // TODO: extract protocol parameters away or figure out a better way to handle this
	ReliableNet  *reliablenet.ModuleParams
	ThreshCrypto *threshcrypto.ModuleParams
}

func DefaultParams(initialMembership *trantorpbtypes.Membership) Params {
	allNodes := maputil.GetSortedKeys(initialMembership.Nodes)
	return Params{
		Mempool:      simplemempool.DefaultModuleParams(),
		Iss:          issconfig.DefaultParams(initialMembership),
		Net:          libp2p.DefaultParams(),
		Alea:         alea.DefaultParams(initialMembership),
		ReliableNet:  reliablenet.DefaultModuleParams(allNodes),
		ThreshCrypto: threshcrypto.DefaultModuleParams(),
	}
}

func (p *Params) AdjustSpeed(maxProposeDelay time.Duration) *Params {
	p.Iss.AdjustSpeed(maxProposeDelay)

	return p
}
