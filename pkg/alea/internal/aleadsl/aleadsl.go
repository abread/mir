package aleadsl

import (
	"fmt"

	adsl "github.com/filecoin-project/mir/pkg/availability/dsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing events.

func UponFillerMessageReceived(m dsl.Module, handler func(from t.NodeID, msg *aleapb.FillerMessage) error) {
	UponMessageReceived(m, func(from t.NodeID, msg *aleapb.Message) error {
		fillerMsgWrapper, ok := msg.Type.(*aleapb.Message_FillerMessage)
		if !ok {
			return nil
		}
		fillerMsg := fillerMsgWrapper.FillerMessage

		return handler(from, fillerMsg)
	})
}

func UponFillGapMessageReceived(m dsl.Module, handler func(from t.NodeID, msg *aleapb.FillGapMessage) error) {
	UponMessageReceived(m, func(from t.NodeID, msg *aleapb.Message) error {
		fillGapMsgWrapper, ok := msg.Type.(*aleapb.Message_FillGapMessage)
		if !ok {
			return nil
		}
		fillGapMsg := fillGapMsgWrapper.FillGapMessage

		return handler(from, fillGapMsg)
	})
}

func UponRequestTransactions(m dsl.Module, handler func(cert *aleapb.Cert, origin *apb.RequestTransactionsOrigin) error) {
	adsl.UponRequestTransactions(m, func(cert *apb.Cert, origin *apb.RequestTransactionsOrigin) error {
		aleabcCertWrapper, ok := cert.Type.(*apb.Cert_Alea)
		if !ok {
			return fmt.Errorf("unexpected certificate type. Expected: %T, got: %T", aleabcCertWrapper, cert.Type)
		}

		return handler(aleabcCertWrapper.Alea, origin)
	})
}

func UponMessageReceived(m dsl.Module, handler func(from t.NodeID, msg *aleapb.Message) error) {
	dsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepb.Message) error {
		cbMsgWrapper, ok := msg.Type.(*messagepb.Message_Alea)
		if !ok {
			return nil
		}

		return handler(from, cbMsgWrapper.Alea)
	})
}
