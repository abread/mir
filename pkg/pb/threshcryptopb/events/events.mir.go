package threshcryptopbevents

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

func SignShare(destModule types.ModuleID, data [][]uint8, origin *types1.SignShareOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ThreshCrypto{
			ThreshCrypto: &types1.Event{
				Type: &types1.Event_SignShare{
					SignShare: &types1.SignShare{
						Data:   data,
						Origin: origin,
					},
				},
			},
		},
	}
}

func SignShareResult(destModule types.ModuleID, signatureShare tctypes.SigShare, origin *types1.SignShareOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ThreshCrypto{
			ThreshCrypto: &types1.Event{
				Type: &types1.Event_SignShareResult{
					SignShareResult: &types1.SignShareResult{
						SignatureShare: signatureShare,
						Origin:         origin,
					},
				},
			},
		},
	}
}

func VerifyShare(destModule types.ModuleID, data [][]uint8, signatureShare tctypes.SigShare, nodeId types.NodeID, origin *types1.VerifyShareOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ThreshCrypto{
			ThreshCrypto: &types1.Event{
				Type: &types1.Event_VerifyShare{
					VerifyShare: &types1.VerifyShare{
						Data:           data,
						SignatureShare: signatureShare,
						NodeId:         nodeId,
						Origin:         origin,
					},
				},
			},
		},
	}
}

func VerifyShareResult(destModule types.ModuleID, ok bool, error string, origin *types1.VerifyShareOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ThreshCrypto{
			ThreshCrypto: &types1.Event{
				Type: &types1.Event_VerifyShareResult{
					VerifyShareResult: &types1.VerifyShareResult{
						Ok:     ok,
						Error:  error,
						Origin: origin,
					},
				},
			},
		},
	}
}

func VerifyFull(destModule types.ModuleID, data [][]uint8, fullSignature tctypes.FullSig, origin *types1.VerifyFullOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ThreshCrypto{
			ThreshCrypto: &types1.Event{
				Type: &types1.Event_VerifyFull{
					VerifyFull: &types1.VerifyFull{
						Data:          data,
						FullSignature: fullSignature,
						Origin:        origin,
					},
				},
			},
		},
	}
}

func VerifyFullResult(destModule types.ModuleID, ok bool, error string, origin *types1.VerifyFullOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ThreshCrypto{
			ThreshCrypto: &types1.Event{
				Type: &types1.Event_VerifyFullResult{
					VerifyFullResult: &types1.VerifyFullResult{
						Ok:     ok,
						Error:  error,
						Origin: origin,
					},
				},
			},
		},
	}
}

func Recover(destModule types.ModuleID, data [][]uint8, signatureShares []tctypes.SigShare, origin *types1.RecoverOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ThreshCrypto{
			ThreshCrypto: &types1.Event{
				Type: &types1.Event_Recover{
					Recover: &types1.Recover{
						Data:            data,
						SignatureShares: signatureShares,
						Origin:          origin,
					},
				},
			},
		},
	}
}

func RecoverResult(destModule types.ModuleID, fullSignature tctypes.FullSig, ok bool, error string, origin *types1.RecoverOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ThreshCrypto{
			ThreshCrypto: &types1.Event{
				Type: &types1.Event_RecoverResult{
					RecoverResult: &types1.RecoverResult{
						FullSignature: fullSignature,
						Ok:            ok,
						Error:         error,
						Origin:        origin,
					},
				},
			},
		},
	}
}
