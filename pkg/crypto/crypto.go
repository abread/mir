// Package crypto provides an implementation of the MirModule module.
// It supports RSA and ECDSA signatures.
package crypto

import (
	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	cryptopbevents "github.com/filecoin-project/mir/pkg/pb/cryptopb/events"
	cryptopbtypes "github.com/filecoin-project/mir/pkg/pb/cryptopb/types"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

func New(crypto Crypto) *modules.SimpleEventApplier {
	return &modules.SimpleEventApplier{&cryptoEvProc{crypto}}
}

type cryptoEvProc struct {
	crypto Crypto
}

func (c *cryptoEvProc) ApplyEvent(event *eventpbtypes.Event) events.EventList {
	switch e := event.Type.(type) {
	case *eventpbtypes.Event_Init:
		// no actions on init
		return events.EmptyList()
	case *eventpbtypes.Event_Crypto:
		switch e := e.Crypto.Type.(type) {
		case *cryptopbtypes.Event_SignRequest:
			// Compute a signature over the provided data and produce a SignResult event.

			signature, err := c.crypto.Sign(e.SignRequest.Data.Data)
			if err != nil {
				panic(err)
			}
			return events.ListOf(
				cryptopbevents.SignResult(
					e.SignRequest.Origin.Module,
					signature,
					e.SignRequest.Origin,
				))

		case *cryptopbtypes.Event_VerifySigs:
			// Verify a batch of node signatures

			// Convenience variables
			verifyEvent := e.VerifySigs
			errors := make([]error, len(verifyEvent.Data))
			allOK := true

			// Verify each signature.
			for i, data := range verifyEvent.Data {
				errors[i] = c.crypto.Verify(data.Data, verifyEvent.Signatures[i], verifyEvent.NodeIds[i])
				if errors[i] != nil {
					allOK = false
				}
			}

			// Return result event
			return events.ListOf(cryptopbevents.SigsVerified(
				verifyEvent.Origin.Module,
				verifyEvent.Origin,
				verifyEvent.NodeIds,
				errors,
				allOK,
			))

		case *cryptopbtypes.Event_VerifySig:
			err := c.crypto.Verify(
				e.VerifySig.Data.Data,
				e.VerifySig.Signature,
				e.VerifySig.NodeId,
			)

			return events.ListOf(cryptopbevents.SigVerified(
				e.VerifySig.Origin.Module,
				e.VerifySig.Origin,
				e.VerifySig.NodeId,
				err,
			))

		default:
			panic(es.Errorf("unexpected type of crypto event: %T", e))
		}
	default:
		// Complain about all other incoming event types.
		panic(es.Errorf("unexpected type of MirModule event: %T", event.Type))
	}
}
