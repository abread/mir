package crypto

import (
	"context"
	"hash"
	"runtime"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	hasherpbevents "github.com/filecoin-project/mir/pkg/pb/hasherpb/events"
	hasherpbtypes "github.com/filecoin-project/mir/pkg/pb/hasherpb/types"
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

func (hasher *hasherEventProc) ApplyEvent(_ context.Context, event *eventpbtypes.Event) *events.EventList {
	switch e := event.Type.(type) {
	case *eventpbtypes.Event_Init:
		// no actions on init
		return events.EmptyList()
	case *eventpbtypes.Event_Hasher:
		switch e := e.Hasher.Type.(type) {
		case *hasherpbtypes.Event_Request:
			// Return all computed digests in one common event.
			return events.ListOf(hasherpbevents.Result(
				e.Request.Origin.Module,
				hasher.computeDigests(e.Request.Data),
				e.Request.Origin,
			))
		case *hasherpbtypes.Event_RequestOne:
			// Return a single computed digests.
			return events.ListOf(hasherpbevents.ResultOne(
				e.RequestOne.Origin.Module,
				hasher.computeDigests([]*hasherpbtypes.HashData{e.RequestOne.Data})[0],
				e.RequestOne.Origin,
			))
		default:
			panic(es.Errorf("unexpected hasher event type: %T", e))
		}
	default:
		// Complain about all other incoming event types.
		panic(es.Errorf("unexpected type of Hash event: %T", event.Type))
	}
}

func (hasher *hasherEventProc) computeDigests(allData []*hasherpbtypes.HashData) [][]byte {
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
