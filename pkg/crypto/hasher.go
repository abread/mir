package crypto

import (
	"context"
	"hash"
	"runtime"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/hasherpb"
	hasherpbevents "github.com/filecoin-project/mir/pkg/pb/hasherpb/events"
	hasherpbtypes "github.com/filecoin-project/mir/pkg/pb/hasherpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type HashImpl interface {
	New() hash.Hash
}

type HasherModuleParams struct {
	NumWorkers int
}

func DefaultHasherModuleParams() *HasherModuleParams {
	return &HasherModuleParams{
		NumWorkers: runtime.NumCPU(),
	}
}

func NewHasher(ctx context.Context, params *HasherModuleParams, hashImpl HashImpl) modules.ActiveModule {
	return modules.NewGoRoutinePoolModule(ctx, &hasherEventProc{hashImpl}, params.NumWorkers)
}

type hasherEventProc struct {
	hashImpl HashImpl
}

func (hasher *hasherEventProc) ApplyEvent(_ context.Context, event *eventpb.Event) *events.EventList {
	switch e := event.Type.(type) {
	case *eventpb.Event_Init:
		// no actions on init
		return events.EmptyList()
	case *eventpb.Event_Hasher:
		switch e := e.Hasher.Type.(type) {
		case *hasherpb.Event_Request:
			// Return all computed digests in one common event.
			return events.ListOf(hasherpbevents.Result(
				t.ModuleID(e.Request.Origin.Module),
				hasher.computeDigests(e.Request.Data),
				hasherpbtypes.HashOriginFromPb(e.Request.Origin),
			).Pb())
		case *hasherpb.Event_RequestOne:
			// Return a single computed digests.
			return events.ListOf(hasherpbevents.ResultOne(
				t.ModuleID(e.RequestOne.Origin.Module),
				hasher.computeDigests([]*hasherpb.HashData{e.RequestOne.Data})[0],
				hasherpbtypes.HashOriginFromPb(e.RequestOne.Origin),
			).Pb())
		default:
			panic(es.Errorf("unexpected hasher event type: %T", e))
		}
	default:
		// Complain about all other incoming event types.
		panic(es.Errorf("unexpected type of Hash event: %T", event.Type))
	}
}

func (hasher *hasherEventProc) computeDigests(allData []*hasherpb.HashData) [][]byte {
	// Create a slice for the resulting digests containing one element for each data item to be hashed.
	digests := make([][]byte, len(allData))

	// Hash each data item contained in the event
	for i, data := range allData {

		// One data item consists of potentially multiple byte slices.
		// Add each of them to the hash function.
		h := hasher.hashImpl.New()
		for _, d := range data.Data {
			h.Write(d)
		}

		// Save resulting digest in the result slice
		digests[i] = h.Sum(nil)
	}

	return digests
}
