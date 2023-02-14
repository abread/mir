package vcb

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	mpdsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	"github.com/filecoin-project/mir/pkg/serializing"
	t "github.com/filecoin-project/mir/pkg/types"
)

type vcbPayloadManager struct {
	m      dsl.Module
	mc     *ModuleConfig
	params *ModuleParams

	txs     []*requestpb.Request
	txIDs   []t.TxID
	sigData [][]byte
}

func newVcbPayloadManager(m dsl.Module, mc *ModuleConfig, params *ModuleParams) *vcbPayloadManager {
	mgr := &vcbPayloadManager{
		m:      m,
		mc:     mc,
		params: params,
	}

	mpdsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *vcbPayloadMgrInputTxs) error {
		return mgr.uponTxIDs(txIDs)
	})

	return mgr
}

type vcbPayloadMgrInputTxs struct{}

func (mgr *vcbPayloadManager) Input(txs []*requestpb.Request) {
	if mgr.txs != nil {
		return
	}

	mgr.txs = txs
	mpdsl.RequestTransactionIDs(mgr.m, mgr.mc.Mempool, txs, &vcbPayloadMgrInputTxs{})
}
func (mgr *vcbPayloadManager) uponTxIDs(txIDs []t.TxID) error {
	mgr.txIDs = txIDs
	mgr.sigData = SigData(mgr.params.InstanceUID, txIDs)
	return nil
}

func (mgr *vcbPayloadManager) Txs() []*requestpb.Request {
	return mgr.txs
}

func (mgr *vcbPayloadManager) TxIDs() []t.TxID {
	return mgr.txIDs
}

func (mgr *vcbPayloadManager) SigData() [][]byte {
	return mgr.sigData
}

func SigData(instanceUID []byte, txIDs []t.TxID) [][]byte {
	res := make([][]byte, 0, len(txIDs)+3)

	res = append(res, []byte("github.com/filecoin-project/mir/pkg/vcb"))
	res = append(res, instanceUID)
	res = append(res, serializing.Uint64ToBytes(uint64(len(txIDs))))

	// TODO: use batch ID from mempool to minimize memory usage
	for _, txID := range txIDs {
		res = append(res, txID.Pb())
	}

	return res
}
