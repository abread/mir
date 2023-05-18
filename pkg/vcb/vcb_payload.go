package vcb

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	hasherpbdsl "github.com/filecoin-project/mir/pkg/pb/hasherpb/dsl"
	hasherpbtypes "github.com/filecoin-project/mir/pkg/pb/hasherpb/types"
	mpdsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

type vcbPayloadManager struct {
	txs       []*trantorpbtypes.Transaction
	txIDs     []tt.TxID
	txIDsHash []byte
	sigData   [][]byte
}

func (mgr *vcbPayloadManager) init(m dsl.Module, mc ModuleConfig, params *ModuleParams) {
	mpdsl.UponTransactionIDsResponse(m, func(txIDs []tt.TxID, context *vcbPayloadMgrInputTxs) error {
		mgr.txIDs = txIDs

		idsBytes := make([][]byte, len(txIDs))
		for i, x := range txIDs {
			idsBytes[i] = []byte(x)
		}

		hasherpbdsl.RequestOne(m, mc.Hasher,
			&hasherpbtypes.HashData{Data: sliceutil.Transform(txIDs, func(_ int, txID tt.TxID) []byte {
				return []byte(txID)
			})},
			context,
		)
		return nil
	})
	hasherpbdsl.UponResultOne(m, func(txIDsHash []byte, context *vcbPayloadMgrInputTxs) error {
		mgr.txIDsHash = txIDsHash
		mgr.sigData = SigData(params.InstanceUID, txIDsHash)
		return nil
	})
}

type vcbPayloadMgrInputTxs struct{}

func (mgr *vcbPayloadManager) Input(m dsl.Module, mc *ModuleConfig, txs []*trantorpbtypes.Transaction) {
	if mgr.txs != nil {
		return
	}

	mgr.txs = txs
	mpdsl.RequestTransactionIDs(m, mc.Mempool, txs, &vcbPayloadMgrInputTxs{})
}

func (mgr *vcbPayloadManager) Txs() []*trantorpbtypes.Transaction {
	return mgr.txs
}

func (mgr *vcbPayloadManager) TxIDs() []tt.TxID {
	return mgr.txIDs
}

func (mgr *vcbPayloadManager) SigData() [][]byte {
	return mgr.sigData
}

func SigData(instanceUID []byte, txIDsHash []byte) [][]byte {
	return [][]byte{instanceUID, txIDsHash}
}
