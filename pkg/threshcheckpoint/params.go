package threshcheckpoint

import (
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
)

// ModuleParams represents the state associated with a single instance of the checkpoint protocol
// (establishing a single stable checkpoint).
type ModuleParams struct {

	// The IDs of nodes to execute this instance of the checkpoint protocol.
	// Note that it is the Membership of Epoch e-1 that constructs the Membership for Epoch e.
	// (As the starting checkpoint for e is the "finishing" checkpoint for e-1.)
	Membership *trantorpbtypes.Membership

	// EpochConfig to which this checkpoint belongs
	// It contains:.
	// - the Epoch the checkpoint's associated sequence number (SeqNr) is part of.
	// - Sequence number associated with this checkpoint protocol instance.
	//	 This checkpoint encompasses SeqNr sequence numbers,
	//	 i.e., SeqNr is the first sequence number *not* encompassed by this checkpoint.
	//	 One can imagine that the checkpoint represents the state of the system just before SeqNr,
	//	 i.e., "between" SeqNr-1 and SeqNr.
	//   among others
	EpochConfig *trantorpbtypes.EpochConfig

	// LeaderPolicy serialization data.
	LeaderPolicyData []byte

	// Threshold is the number of required signature shares to reconstruct the checkpoint signature.
	Threshold int
}