// Package threshcrypto provides an implementation of the MirModule module.
// It supports TBLS signatures.
package threshcrypto

import (
	"context"
	"fmt"
	"runtime"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/threshcryptopb"
	tcEvents "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/events"
	threshcryptopbtypes "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/types"
	t "github.com/filecoin-project/mir/pkg/types"
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

func (c *threshEventProcessor) ApplyEvent(ctx context.Context, event *eventpb.Event) *events.EventList {
	switch e := event.Type.(type) {
	case *eventpb.Event_Init:
		// no actions on init
		return events.EmptyList()
	case *eventpb.Event_ThreshCrypto:
		return c.applyTCEvent(ctx, e.ThreshCrypto)
	default:
		// Complain about all other incoming event types.
		panic(fmt.Errorf("unexpected type of event in threshcrypto MirModule: %T", event.Type))
	}
}

// apply a thresholdcryptopb.Event
func (c *threshEventProcessor) applyTCEvent(_ context.Context, event *threshcryptopb.Event) *events.EventList {
	switch e := event.Type.(type) {
	case *threshcryptopb.Event_SignShare:
		// Compute signature share

		sigShare, err := c.threshCrypto.SignShare(e.SignShare.Data)
		if err != nil {
			panic(fmt.Errorf("could not sign share: %w", err))
		}

		origin := threshcryptopbtypes.SignShareOriginFromPb(e.SignShare.Origin)
		return events.ListOf(
			tcEvents.SignShareResult(origin.Module, sigShare, origin).Pb(),
		)

	case *threshcryptopb.Event_VerifyShare:
		// Verify signature share

		err := c.threshCrypto.VerifyShare(e.VerifyShare.Data, e.VerifyShare.SignatureShare, t.NodeID(e.VerifyShare.NodeId))

		ok := err == nil
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		origin := threshcryptopbtypes.VerifyShareOriginFromPb(e.VerifyShare.Origin)
		return events.ListOf(
			tcEvents.VerifyShareResult(origin.Module, ok, errStr, origin).Pb(),
		)

	case *threshcryptopb.Event_VerifyFull:
		// Verify full signature

		err := c.threshCrypto.VerifyFull(e.VerifyFull.Data, e.VerifyFull.FullSignature)

		ok := err == nil
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		origin := threshcryptopbtypes.VerifyFullOriginFromPb(e.VerifyFull.Origin)
		return events.ListOf(
			tcEvents.VerifyFullResult(origin.Module, ok, errStr, origin).Pb(),
		)

	case *threshcryptopb.Event_Recover:
		// Recover full signature from shares

		fullSig, err := c.threshCrypto.Recover(e.Recover.Data, e.Recover.SignatureShares)

		ok := err == nil
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		origin := threshcryptopbtypes.RecoverOriginFromPb(e.Recover.Origin)
		return events.ListOf(
			tcEvents.RecoverResult(origin.Module, fullSig, ok, errStr, origin).Pb(),
		)
	default:
		// Complain about all other incoming event types.
		panic(fmt.Errorf("unexpected type of threshcrypto event in threshcrypto MirModule: %T", event.Type))
	}
}
