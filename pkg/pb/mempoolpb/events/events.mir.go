// Code generated by Mir codegen. DO NOT EDIT.

package mempoolpbevents

import (
	types5 "github.com/filecoin-project/mir/pkg/availability/multisigcollector/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/mempoolpb/types"
	types4 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	types1 "github.com/filecoin-project/mir/pkg/trantor/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func RequestBatch(destModule types.ModuleID, epoch types1.EpochNr, origin *types2.RequestBatchOrigin) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_RequestBatch{
					RequestBatch: &types2.RequestBatch{
						Epoch:  epoch,
						Origin: origin,
					},
				},
			},
		},
	}
}

func NewBatch(destModule types.ModuleID, txIds []types1.TxID, txs []*types4.Transaction, origin *types2.RequestBatchOrigin) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_NewBatch{
					NewBatch: &types2.NewBatch{
						TxIds:  txIds,
						Txs:    txs,
						Origin: origin,
					},
				},
			},
		},
	}
}

func RequestTransactions(destModule types.ModuleID, txIds []types1.TxID, origin *types2.RequestTransactionsOrigin) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_RequestTransactions{
					RequestTransactions: &types2.RequestTransactions{
						TxIds:  txIds,
						Origin: origin,
					},
				},
			},
		},
	}
}

func TransactionsResponse(destModule types.ModuleID, foundIds []types1.TxID, foundTxs []*types4.Transaction, missingIds []types1.TxID, origin *types2.RequestTransactionsOrigin) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_TransactionsResponse{
					TransactionsResponse: &types2.TransactionsResponse{
						FoundIds:   foundIds,
						FoundTxs:   foundTxs,
						MissingIds: missingIds,
						Origin:     origin,
					},
				},
			},
		},
	}
}

func RequestTransactionIDs(destModule types.ModuleID, txs []*types4.Transaction, origin *types2.RequestTransactionIDsOrigin) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_RequestTransactionIds{
					RequestTransactionIds: &types2.RequestTransactionIDs{
						Txs:    txs,
						Origin: origin,
					},
				},
			},
		},
	}
}

func TransactionIDsResponse(destModule types.ModuleID, txIds []types1.TxID, origin *types2.RequestTransactionIDsOrigin) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_TransactionIdsResponse{
					TransactionIdsResponse: &types2.TransactionIDsResponse{
						TxIds:  txIds,
						Origin: origin,
					},
				},
			},
		},
	}
}

func RequestBatchID(destModule types.ModuleID, txIds []types1.TxID, origin *types2.RequestBatchIDOrigin) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_RequestBatchId{
					RequestBatchId: &types2.RequestBatchID{
						TxIds:  txIds,
						Origin: origin,
					},
				},
			},
		},
	}
}

func BatchIDResponse(destModule types.ModuleID, batchId types5.BatchID, origin *types2.RequestBatchIDOrigin) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_BatchIdResponse{
					BatchIdResponse: &types2.BatchIDResponse{
						BatchId: batchId,
						Origin:  origin,
					},
				},
			},
		},
	}
}

func NewTransactions(destModule types.ModuleID, transactions []*types4.Transaction) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_NewTransactions{
					NewTransactions: &types2.NewTransactions{
						Transactions: transactions,
					},
				},
			},
		},
	}
}

func BatchTimeout(destModule types.ModuleID, batchReqID uint64) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_BatchTimeout{
					BatchTimeout: &types2.BatchTimeout{
						BatchReqID: batchReqID,
					},
				},
			},
		},
	}
}

func NewEpoch(destModule types.ModuleID, epochNr types1.EpochNr, clientProgress *types4.ClientProgress) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_NewEpoch{
					NewEpoch: &types2.NewEpoch{
						EpochNr:        epochNr,
						ClientProgress: clientProgress,
					},
				},
			},
		},
	}
}

func MarkStableProposal(destModule types.ModuleID, txs []*types4.Transaction) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_Mempool{
			Mempool: &types2.Event{
				Type: &types2.Event_MarkStableProposal{
					MarkStableProposal: &types2.MarkStableProposal{
						Txs: txs,
					},
				},
			},
		},
	}
}
