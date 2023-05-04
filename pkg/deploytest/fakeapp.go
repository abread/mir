/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package deploytest

import (
	"github.com/filecoin-project/mir/pkg/checkpoint"
	"github.com/filecoin-project/mir/pkg/logging"

	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"

	"github.com/filecoin-project/mir/pkg/serializing"
)

// FakeApp represents a dummy stub application used for testing only.
type FakeApp struct {
	// The state of the FakeApp only consists of a counter of processed transactions.
	TransactionsProcessed uint64

	Logger logging.Logger
}

func (fa *FakeApp) ApplyTXs(txs []*trantorpbtypes.Transaction) error {
	for _, tx := range txs {
		fa.TransactionsProcessed++
		fa.Logger.Log(logging.LevelDebug, "Processing (already ordered) transaction", "txData", string(tx.Data), "totalProcessedTxsCount", fa.TransactionsProcessed)
	}
	return nil
}

func (fa *FakeApp) Snapshot() ([]byte, error) {
	return serializing.Uint64ToBytes(fa.TransactionsProcessed), nil
}

func (fa *FakeApp) RestoreState(checkpoint *checkpoint.StableCheckpoint) error {
	fa.TransactionsProcessed = serializing.Uint64FromBytes(checkpoint.Snapshot.AppData)
	return nil
}

func (fa *FakeApp) Checkpoint(_ *checkpoint.StableCheckpoint) error {
	return nil
}

func NewFakeApp(logger logging.Logger) *FakeApp {
	return &FakeApp{
		TransactionsProcessed: 0,
		Logger:                logger,
	}
}
