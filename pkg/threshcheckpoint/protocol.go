// TODO: Finish writing proper comments in this file.

package threshcheckpoint

import (
	"bytes"
	"runtime"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	apppbdsl "github.com/filecoin-project/mir/pkg/pb/apppb/dsl"
	checkpointpbdsl "github.com/filecoin-project/mir/pkg/pb/checkpointpb/dsl"
	hasherpbdsl "github.com/filecoin-project/mir/pkg/pb/hasherpb/dsl"
	reliablenetpbdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	threshcheckpointpbdsl "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb/dsl"
	threshcheckpointpbmsgs "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb/msgs"
	threshcheckpointpbtypes "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb/types"
	threshcryptopbdsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	trantorpbdsl "github.com/filecoin-project/mir/pkg/pb/trantorpb/dsl"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tsagg"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

type State struct {
	// State snapshot associated with this checkpoint.
	StateSnapshot *trantorpbtypes.StateSnapshot

	// Hash of the state snapshot data associated with this checkpoint.
	StateSnapshotHash []byte

	SigData [][]byte

	SigAggregator *tsagg.ThreshSigAggregator

	// Set of Checkpoint messages that were received ahead of time.
	PendingMessages map[t.NodeID]*threshcheckpointpbtypes.Checkpoint

	// Flag ensuring that the stable checkpoint is only Announced once.
	// Set to true when announcing a stable checkpoint for the first time.
	// When true, stable checkpoints are not Announced anymore.
	Announced bool
}

// NewModule allocates and returns a new instance of the ModuleParams associated with sequence number sn.
func NewModule(
	moduleConfig ModuleConfig,
	params *ModuleParams,
	nodeID t.NodeID,
	logger logging.Logger) modules.PassiveModule {

	m := dsl.NewModule(moduleConfig.Self)

	state := &State{
		StateSnapshot: &trantorpbtypes.StateSnapshot{
			AppData: nil,
			EpochData: &trantorpbtypes.EpochData{
				EpochConfig:        params.EpochConfig,
				ClientProgress:     nil, // This will be filled by a separate event.
				LeaderPolicy:       params.LeaderPolicyData,
				PreviousMembership: params.Membership,
			},
		},
		Announced:       false,
		PendingMessages: make(map[t.NodeID]*threshcheckpointpbtypes.Checkpoint),
	}

	state.SigAggregator = tsagg.New(m, &tsagg.Params{
		OwnNodeID:               nodeID,
		TCModuleID:              moduleConfig.ThreshCrypto,
		Threshold:               params.Threshold,
		MaxVerifyShareBatchSize: runtime.NumCPU(),
		SigData: func() [][]byte {
			return state.SigData
		},
		InitialNodeCount: len(params.Membership.Nodes),
	}, logging.Decorate(logger, "threshcrypto sig aggregator: "))

	apppbdsl.UponSnapshot(m, func(appData []uint8) error {
		// Treat nil data as an empty byte slice.
		if appData == nil {
			appData = []byte{}
		}

		// Save the received app snapshot if there is none yet.
		if state.StateSnapshot.AppData == nil {
			state.StateSnapshot.AppData = appData
			if state.SnapshotReady() {
				if err := processStateSnapshot(m, state, moduleConfig); err != nil {
					return err
				}
			}
		}
		return nil
	})

	hasherpbdsl.UponResultOne(m, func(digest []uint8, _ *struct{}) error {
		// Save the received snapshot hash
		state.StateSnapshotHash = digest

		// Request signature
		state.SigData = serializeCheckpointForSig(params.EpochConfig.EpochNr, params.EpochConfig.FirstSn, state.StateSnapshotHash)

		threshcryptopbdsl.SignShare(m, moduleConfig.ThreshCrypto, state.SigData, &struct{}{})
		return nil
	})

	threshcryptopbdsl.UponSignShareResult(m, func(sig tctypes.SigShare, _ *struct{}) error {

		// Save received own checkpoint signature
		state.SigAggregator.Add(sig, nodeID)

		// Send a checkpoint message to all nodes.
		chkpMessage := threshcheckpointpbmsgs.Checkpoint(moduleConfig.Self, params.EpochConfig.EpochNr, params.EpochConfig.FirstSn, state.StateSnapshotHash, sig)
		sortedMembership := maputil.GetSortedKeys(params.Membership.Nodes)

		reliablenetpbdsl.SendMessage(m, moduleConfig.ReliableNet, "c", chkpMessage, sortedMembership)
		logger.Log(logging.LevelDebug, "Sending checkpoint message",
			"epoch", params.EpochConfig.EpochNr,
			"dataLen", len(state.StateSnapshot.AppData),
			"memberships", len(state.StateSnapshot.EpochData.EpochConfig.Memberships),
		)

		// Apply pending Checkpoint messages
		for s, msg := range state.PendingMessages {
			if err := applyCheckpointReceived(m, params, state, moduleConfig, s, msg.Epoch, msg.Sn, msg.SnapshotHash, msg.SignatureShare, logger); err != nil {
				logger.Log(logging.LevelWarn, "Error applying pending Checkpoint message", "error", err, "msg", msg)
				return err
			}

		}
		state.PendingMessages = nil

		return nil
	})

	dsl.UponStateUpdates(m, func() error {
		if state.Stable() {
			announceStable(m, params, state, moduleConfig)
		}
		return nil
	})

	trantorpbdsl.UponClientProgress(m, func(progress map[tt.ClientID]*trantorpbtypes.DeliveredTXs) error {
		// Save the received client progress if there is none yet.
		if state.StateSnapshot.EpochData.ClientProgress == nil {
			state.StateSnapshot.EpochData.ClientProgress = &trantorpbtypes.ClientProgress{
				Progress: progress,
			}
			if state.SnapshotReady() {
				if err := processStateSnapshot(m, state, moduleConfig); err != nil {
					return err
				}
			}
		}
		return nil
	})

	threshcheckpointpbdsl.UponCheckpointReceived(m, func(from t.NodeID, epoch tt.EpochNr, sn tt.SeqNr, snapshotHash []uint8, sigShare tctypes.SigShare) error {
		return applyCheckpointReceived(m, params, state, moduleConfig, from, epoch, sn, snapshotHash, sigShare, logger)
	})

	return m
}

func processStateSnapshot(m dsl.Module, state *State, mc ModuleConfig) error {

	// Serialize the snapshot.
	snapshotData, err := serializeSnapshotForHash(state.StateSnapshot)
	if err != nil {
		return es.Errorf("failed serializing state snapshot: %w", err)
	}

	// Initiate computing the hash of the snapshot.
	hasherpbdsl.RequestOne(m,
		mc.Hasher,
		snapshotData,
		&struct{}{},
	)

	return nil
}

func announceStable(m dsl.Module, p *ModuleParams, state *State, mc ModuleConfig) {

	// Only announce the stable checkpoint once.
	if state.Announced {
		return
	}
	state.Announced = true

	// Announce the stable checkpoint to the ordering protocol.
	threshcheckpointpbdsl.StableCheckpoint(m, mc.Ord, p.EpochConfig.FirstSn, state.StateSnapshot, state.SigAggregator.FullSig())

	// Stop broadcasting checkpoint messages.
	reliablenetpbdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, "c", maputil.GetSortedKeys(p.Membership.Nodes))
}

func applyCheckpointReceived(m dsl.Module,
	p *ModuleParams,
	state *State,
	moduleConfig ModuleConfig,
	from t.NodeID,
	epoch tt.EpochNr,
	sn tt.SeqNr,
	snapshotHash []uint8,
	sigShare tctypes.SigShare,
	logger logging.Logger) error {

	reliablenetpbdsl.MarkRecvd(m, moduleConfig.ReliableNet, moduleConfig.Self, "c", []t.NodeID{from})

	if sn != p.EpochConfig.FirstSn {
		logger.Log(logging.LevelWarn, "invalid sequence number %v, expected %v as first sequence number.\n", sn, p.EpochConfig.FirstSn)
		return nil
	}

	if epoch != p.EpochConfig.EpochNr {
		logger.Log(logging.LevelWarn, "invalid epoch number %v, expected %v.\n", epoch, p.EpochConfig.EpochNr)
		return nil
	}

	// check if from is part of the membership
	if _, ok := p.Membership.Nodes[from]; !ok {
		logger.Log(logging.LevelWarn, "sender %s is not a member.\n", from)
		return nil
	}

	// Notify the protocol about the progress of the from node.
	// If no progress is made for a configured number of epochs,
	// the node is considered to be a straggler and is sent a stable checkpoint to catch uparams.
	checkpointpbdsl.EpochProgress(m, moduleConfig.Ord, from, epoch)

	// If checkpoint is already stable, ignore message.
	if state.Stable() {
		return nil
	}

	// Check snapshot hash
	if state.StateSnapshotHash == nil {
		// The message is received too early, put it aside
		state.PendingMessages[from] = &threshcheckpointpbtypes.Checkpoint{
			Epoch:          epoch,
			Sn:             sn,
			SnapshotHash:   snapshotHash,
			SignatureShare: sigShare,
		}
		return nil
	} else if !bytes.Equal(state.StateSnapshotHash, snapshotHash) {
		// Snapshot hash mismatch
		logger.Log(logging.LevelWarn, "Ignoring Checkpoint message. Mismatching state snapshot hash.", "from", from)
		return nil
	}

	// Ignore duplicate messages.
	state.SigAggregator.Add(sigShare, from)
	return nil
}

func (state *State) SnapshotReady() bool {
	return state.StateSnapshot.AppData != nil &&
		state.StateSnapshot.EpochData.ClientProgress != nil
}

func (state *State) Stable() bool {
	// TODO: support weighting
	return state.SnapshotReady() && state.SigAggregator.FullSig() != nil
}
