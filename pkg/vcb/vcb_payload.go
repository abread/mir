package vcb

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	batchdbpbdsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/dsl"
	mempoolpbdsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
)

type vcbPayloadManager struct {
	txs     []*trantorpbtypes.Transaction
	batchID string
	sigData [][]byte

	batchStored bool
}

func (mgr *vcbPayloadManager) init(m dsl.Module, mc ModuleConfig, params *ModuleParams) {
	mempoolpbdsl.UponTransactionIDsResponse(m, func(txIDs []tt.TxID, context *vcbPayloadMgrInputTxs) error {
		mempoolpbdsl.RequestBatchID(m, mc.Mempool, txIDs, context)
		return nil
	})

	mempoolpbdsl.UponBatchIDResponse(m, func(batchID string, context *vcbPayloadMgrInputTxs) error {
		mgr.batchID = batchID
		mgr.sigData = SigData(params.InstanceUID, batchID)

		batchdbpbdsl.StoreBatch(m, mc.BatchDB, batchID, mgr.txs, params.EpochNr, context)
		return nil
	})

	batchdbpbdsl.UponBatchStored(m, func(_ *vcbPayloadMgrInputTxs) error {
		mgr.batchStored = true
		return nil
	})
}

type vcbPayloadMgrInputTxs struct{}

func (mgr *vcbPayloadManager) Input(m dsl.Module, mc *ModuleConfig, txIDs []tt.TxID, txs []*trantorpbtypes.Transaction) {
	if mgr.txs != nil {
		return
	}

	mgr.txs = txs

	if txIDs == nil {
		// outside source of txs, must compute ids
		mempoolpbdsl.RequestTransactionIDs(m, mc.Mempool, txs, &vcbPayloadMgrInputTxs{})
	} else {
		// batch came from local node, can use existing ids
		mempoolpbdsl.RequestBatchID(m, mc.Mempool, txIDs, &vcbPayloadMgrInputTxs{})
	}
}

func (mgr *vcbPayloadManager) Txs() []*trantorpbtypes.Transaction {
	return mgr.txs
}

func (mgr *vcbPayloadManager) BatchID() string {
	return mgr.batchID
}

func (mgr *vcbPayloadManager) SigData() [][]byte {
	// Only allow signatures after persisting batch
	if mgr.batchStored {
		return mgr.sigData
	}

	return nil
}

func (mgr *vcbPayloadManager) IsStored() bool {
	return mgr.batchStored
}

func SigData(instanceUID []byte, batchID string) [][]byte {
	return [][]byte{instanceUID, []byte(batchID)}
}
