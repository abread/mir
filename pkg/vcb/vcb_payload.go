package vcb

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	mpdsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	requestpbtypes "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type vcbPayloadManager struct {
	m      dsl.Module
	mc     *ModuleConfig
	params *ModuleParams

	txs       []*requestpbtypes.Request
	txIDs     []t.TxID
	txIDsHash []byte
	sigData   [][]byte
}

func newVcbPayloadManager(m dsl.Module, mc *ModuleConfig, params *ModuleParams) *vcbPayloadManager {
	mgr := &vcbPayloadManager{
		m:      m,
		mc:     mc,
		params: params,
	}

	mpdsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *vcbPayloadMgrInputTxs) error {
		mgr.txIDs = txIDs
		dsl.HashOneMessage(m, mc.Hasher, t.TxIDSlicePb(txIDs), context)
		return nil
	})
	dsl.UponOneHashResult(m, func(txIDsHash []byte, context *vcbPayloadMgrInputTxs) error {
		mgr.txIDsHash = txIDsHash
		mgr.sigData = SigData(mgr.params.InstanceUID, txIDsHash)
		return nil
	})

	return mgr
}

type vcbPayloadMgrInputTxs struct{}

func (mgr *vcbPayloadManager) Input(txs []*requestpbtypes.Request) {
	if mgr.txs != nil {
		return
	}

	mgr.txs = txs
	mpdsl.RequestTransactionIDs(mgr.m, mgr.mc.Mempool, txs, &vcbPayloadMgrInputTxs{})
}

func (mgr *vcbPayloadManager) Txs() []*requestpbtypes.Request {
	return mgr.txs
}

func (mgr *vcbPayloadManager) TxIDs() []t.TxID {
	return mgr.txIDs
}

func (mgr *vcbPayloadManager) SigData() [][]byte {
	return mgr.sigData
}

func SigData(instanceUID []byte, txIDsHash []byte) [][]byte {
	return [][]byte{instanceUID, txIDsHash}
}
