package alea

import (
	"time"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/queueselectionpolicy"
	"github.com/filecoin-project/mir/pkg/clientprogress"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

// Params sets the values for the parameters of an instance of the protocol.
type Params struct {
	// Unique identifier for this instance of Alea.
	// Must be the same across all replicas.
	InstanceUID []byte

	// The identities of all nodes that execute the protocol.
	// Must not be empty.
	Membership *trantorpbtypes.Membership

	// Maximum number of concurrent VCB instances per queue
	// Must be at least 1
	MaxConcurrentVcbPerQueue int

	// Maximum number of unagreed batches that the broadcast component can have in this node's queue
	// Must be at least 1
	MaxOwnUnagreedBatchCount int

	// Maximum number of concurrent ABBA rounds for which we process messages
	// Must be at least 1
	MaxAbbaRoundLookahead int

	// Maximum number of concurrent agreement rounds for which we process messages
	// Must be at least 1
	MaxAgRoundLookahead int

	// Maximum number of agreement rounds for which we send input before
	// allowing them to progress in the normal path.
	// Must be at least 0, must be less that MaxRoundLookahead
	// Setting this parameter too high will lead to costly retransmissions!
	// Should likely be less than MaxRoundLookahead/2 - 1.
	MaxAgRoundAdvanceInput int

	// Subprotocol duration estimates window size
	EstimateWindowSize int

	// How slower can the F slowest nodes be compared to the rest
	// A factor of 2, means we allow F nodes to take double the time we expect the rest to take.
	// Must be >=1
	MaxExtSlowdownFactor float64

	QueueSelectionPolicyType queueselectionpolicy.QueuePolicyType
	EpochLength              uint64
	RetainEpochs             uint64

	MaxAgStall time.Duration
}

// DefaultParams returns the default configuration for a given membership.
// There is no guarantee that this configuration ensures good performance or security,
// but it will pass the CheckParams test.
// DefaultParams is intended for use during testing and hello-world examples.
// A proper deployment is expected to craft a custom configuration,
// for which DefaultParams can serve as a starting point.
func DefaultParams(membership *trantorpbtypes.Membership) *Params {
	EpochLength := 256
	MaxOwnUnagreedBatchCount := 2

	p := &Params{
		InstanceUID:              []byte{42},
		Membership:               membership,
		MaxAbbaRoundLookahead:    4,
		EstimateWindowSize:       64,
		MaxExtSlowdownFactor:     1.5,
		QueueSelectionPolicyType: queueselectionpolicy.RoundRobin,
		EpochLength:              uint64(EpochLength),
		RetainEpochs:             2,
		MaxAgStall:               10 * time.Second,
	}
	p.Adjust(EpochLength, MaxOwnUnagreedBatchCount)

	return p
}

func (p *Params) Adjust(epochSz int, maxOwnUnagBatchCount int) {
	N := len(p.Membership.Nodes)
	p.MaxOwnUnagreedBatchCount = maxOwnUnagBatchCount
	p.MaxConcurrentVcbPerQueue = (epochSz/N + 1) * 2
	if p.MaxConcurrentVcbPerQueue < p.MaxOwnUnagreedBatchCount*2 {
		p.MaxConcurrentVcbPerQueue = p.MaxOwnUnagreedBatchCount * 2
	}

	p.MaxAgRoundLookahead = epochSz * 2
	p.MaxAgRoundAdvanceInput = epochSz - 1
	p.EpochLength = uint64(epochSz)
}

func (p *Params) Check() error {
	if len(p.Membership.Nodes) == 0 {
		return es.Errorf("cannot start Alea with an empty membership")
	}

	if p.MaxConcurrentVcbPerQueue < 1 {
		return es.Errorf("this implementation requires MaxConcurrentVcbPerQueue >= 1")
	}

	if p.MaxOwnUnagreedBatchCount < 1 {
		return es.Errorf("alea requires TargetOwnUnagreedBatchCount >= 1")
	}

	// TODO: support weighted voting
	// easiest solution is distributing <node weight> keys to each node, and sign with all of them in
	// each operation this is probably best to encode in the threshcrypto module, since weighted
	// threshold crypto schemes seem to exist (see "An Efficient and Secure Weighted Threshold
	// Signcryption Scheme", Chien-Hua Tsai, Journal of Internet Technology, 2019)
	w := ""
	for _, identity := range p.Membership.Nodes {
		if w == "" {
			w = string(identity.Weight)
		}

		if string(identity.Weight) != w {
			return es.Errorf("alea does not support weighted voting (yet): mismatched weights %v != %v", w, identity.Weight)
		}
	}

	return nil
}

func (p *Params) AllNodes() []t.NodeID {
	return maputil.GetSortedKeys(p.Membership.Nodes)
}

// InitialStateSnapshot creates and returns a default initial snapshot that can be used to instantiate ISS.
// The created snapshot corresponds to epoch 0, without any committed transactions,
// and with the initial membership (found in params and used for epoch 0)
// repeated for all the params.ConfigOffset following epochs.
func InitialStateSnapshot(
	appState []byte,
	params *Params,
) (*trantorpbtypes.StateSnapshot, error) {

	// TODO: support membership changes.
	memberships := []*trantorpbtypes.Membership{params.Membership}

	// Create the initial queue selection policy.
	qsp, err := queueselectionpolicy.NewQueuePolicy(params.QueueSelectionPolicyType, params.Membership)
	if err != nil {
		return nil, err
	}

	qspData, err := qsp.Bytes()
	if err != nil {
		return nil, err
	}

	return &trantorpbtypes.StateSnapshot{
		AppData: appState,
		EpochData: &trantorpbtypes.EpochData{
			EpochConfig: &trantorpbtypes.EpochConfig{
				EpochNr:     0,
				FirstSn:     0,
				Length:      params.EpochLength,
				Memberships: memberships,
			},
			ClientProgress:     clientprogress.NewClientProgress().DslStruct(),
			LeaderPolicy:       qspData,
			PreviousMembership: nil,
		},
	}, nil
}
