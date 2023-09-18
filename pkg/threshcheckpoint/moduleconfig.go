package threshcheckpoint

import t "github.com/filecoin-project/mir/pkg/types"

type ModuleConfig struct {
	Self t.ModuleID

	App          t.ModuleID
	Hasher       t.ModuleID
	ThreshCrypto t.ModuleID
	ReliableNet  t.ModuleID
	Ord          t.ModuleID
}
