// Package threshcrypto provides an implementation of the MirModule module.
// It supports TBLS signatures.
package threshcrypto

import (
	"context"
	"runtime"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	tcEvents "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/events"
	threshcryptopbtypes "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/types"
)

type ModuleParams struct {
	NumWorkers int
}

func DefaultModuleParams() *ModuleParams {
	return &ModuleParams{
		NumWorkers: runtime.NumCPU(),
	}
}

func New(ctx context.Context, params *ModuleParams, threshCrypto ThreshCrypto) modules.ActiveModule {
	return modules.NewGoRoutinePoolModule(ctx, &threshEventProcessor{threshCrypto}, params.NumWorkers)
}

type threshEventProcessor struct {
	threshCrypto ThreshCrypto
}

func (c *threshEventProcessor) ApplyEvent(ctx context.Context, event *eventpbtypes.Event) *events.EventList {
	switch e := event.Type.(type) {
	case *eventpbtypes.Event_Init:
		// no actions on init
		return events.EmptyList()
	case *eventpbtypes.Event_ThreshCrypto:
		return c.applyTCEvent(ctx, e.ThreshCrypto)
	default:
		// Complain about all other incoming event types.
		panic(es.Errorf("unexpected type of event in threshcrypto MirModule: %T", event.Type))
	}
}

// apply a thresholdcryptopbtypes.Event
func (c *threshEventProcessor) applyTCEvent(_ context.Context, event *threshcryptopbtypes.Event) *events.EventList {
	switch e := event.Type.(type) {
	case *threshcryptopbtypes.Event_SignShare:
		// Compute signature share

		sigShare, err := c.threshCrypto.SignShare(e.SignShare.Data)
		if err != nil {
			panic(es.Errorf("could not sign share: %w", err))
		}

		origin := e.SignShare.Origin
		return events.ListOf(
			tcEvents.SignShareResult(origin.Module, sigShare, origin),
		)

	case *threshcryptopbtypes.Event_VerifyShare:
		// Verify signature share

		err := c.threshCrypto.VerifyShare(e.VerifyShare.Data, e.VerifyShare.SignatureShare, e.VerifyShare.NodeId)

		ok := err == nil
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		origin := e.VerifyShare.Origin
		return events.ListOf(
			tcEvents.VerifyShareResult(origin.Module, ok, errStr, origin),
		)

	case *threshcryptopbtypes.Event_VerifyFull:
		// Verify full signature

		err := c.threshCrypto.VerifyFull(e.VerifyFull.Data, e.VerifyFull.FullSignature)

		ok := err == nil
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		origin := e.VerifyFull.Origin
		return events.ListOf(
			tcEvents.VerifyFullResult(origin.Module, ok, errStr, origin),
		)

	case *threshcryptopbtypes.Event_Recover:
		// Recover full signature from shares

		fullSig, err := c.threshCrypto.Recover(e.Recover.Data, e.Recover.SignatureShares)

		ok := err == nil
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		origin := e.Recover.Origin
		return events.ListOf(
			tcEvents.RecoverResult(origin.Module, fullSig, ok, errStr, origin),
		)
	default:
		// Complain about all other incoming event types.
		panic(es.Errorf("unexpected type of threshcrypto event in threshcrypto MirModule: %T", event.Type))
	}
}
