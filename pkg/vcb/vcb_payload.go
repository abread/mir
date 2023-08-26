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
	txIDs   []tt.TxID
	batchID string
	sigData [][]byte

	batchStored bool
}

func (mgr *vcbPayloadManager) init(m dsl.Module, mc ModuleConfig, params *ModuleParams) {
	mempoolpbdsl.UponTransactionIDsResponse(m, func(txIDs []tt.TxID, context *vcbPayloadMgrInputTxs) error {
		mgr.txIDs = txIDs

		idsBytes := make([][]byte, len(txIDs))
		for i, x := range txIDs {
			idsBytes[i] = []byte(x)
		}

		mempoolpbdsl.RequestBatchID(m, mc.Mempool, txIDs, context)
		return nil
	})

	mempoolpbdsl.UponBatchIDResponse(m, func(batchID string, context *vcbPayloadMgrInputTxs) error {
		mgr.batchID = batchID
		mgr.sigData = SigData(params.InstanceUID, batchID)

		batchdbpbdsl.StoreBatch(m, mc.BatchDB, batchID, mgr.txs, params.RetentitionIndex, context)
		return nil
	})

	batchdbpbdsl.UponBatchStored(m, func(_ *vcbPayloadMgrInputTxs) error {
		mgr.batchStored = true
		return nil
	})
}

type vcbPayloadMgrInputTxs struct{}

func (mgr *vcbPayloadManager) Input(m dsl.Module, mc *ModuleConfig, txs []*trantorpbtypes.Transaction) {
	if mgr.txs != nil {
		return
	}

	mgr.txs = txs
	mempoolpbdsl.RequestTransactionIDs(m, mc.Mempool, txs, &vcbPayloadMgrInputTxs{})
}

func (mgr *vcbPayloadManager) Txs() []*trantorpbtypes.Transaction {
	return mgr.txs
}

func (mgr *vcbPayloadManager) TxIDs() []tt.TxID {
	return mgr.txIDs
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
