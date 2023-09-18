package chkpvalidator

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	tcvpbdsl "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb/threshchkpvalidatorpb/dsl"
	tcvpbtypes "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb/threshchkpvalidatorpb/types"
	threshcheckpointpbtypes "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids.
type ModuleConfig struct {
	Self t.ModuleID
}

// NewModule returns a (passive) ChkpValidator module.
func NewModule(mc ModuleConfig, cv ChkpValidator) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	tcvpbdsl.UponValidateCheckpoint(m, func(checkpoint *threshcheckpointpbtypes.StableCheckpoint, epochNr types.EpochNr, memberships []*trantorpbtypes.Membership, origin *tcvpbtypes.ValidateChkpOrigin) error {
		err := cv.Verify(checkpoint, epochNr, memberships)
		tcvpbdsl.CheckpointValidated(m, origin.Module, err, origin)
		return nil
	})

	return m
}
