package trantor

import (
	"github.com/filecoin-project/mir/pkg/alea/agreement"
	"github.com/filecoin-project/mir/pkg/alea/broadcast"
	"github.com/filecoin-project/mir/pkg/alea/director"
	"github.com/filecoin-project/mir/pkg/availability/batchdb/fakebatchdb"
	"github.com/filecoin-project/mir/pkg/availability/multisigcollector"
	"github.com/filecoin-project/mir/pkg/batchfetcher"
	"github.com/filecoin-project/mir/pkg/checkpoint"
	"github.com/filecoin-project/mir/pkg/checkpoint/chkpvalidator"
	"github.com/filecoin-project/mir/pkg/iss"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool"
	ordererscommon "github.com/filecoin-project/mir/pkg/orderers/common"
	"github.com/filecoin-project/mir/pkg/orderers/common/pprepvalidator"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	"github.com/filecoin-project/mir/pkg/threshcheckpoint"
	threshchkpvalidator "github.com/filecoin-project/mir/pkg/threshcheckpoint/chkpvalidator"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleConfig struct {
	App           t.ModuleID
	Availability  t.ModuleID
	BatchDB       t.ModuleID
	BatchFetcher  t.ModuleID
	Checkpointing t.ModuleID
	ChkpValidator t.ModuleID
	Crypto        t.ModuleID
	ThreshCrypto  t.ModuleID
	Hasher        t.ModuleID
	Mempool       t.ModuleID
	Net           t.ModuleID
	ReliableNet   t.ModuleID
	Null          t.ModuleID
	Timer         t.ModuleID

	// TODO: move these outside trantor(?)
	ISS                t.ModuleID
	Ordering           t.ModuleID
	PPrepValidator     t.ModuleID
	PPrepValidatorChkp t.ModuleID

	AleaDirector  t.ModuleID
	AleaBroadcast t.ModuleID
	AleaAgreement t.ModuleID
}

func DefaultModuleConfig() ModuleConfig {
	return ModuleConfig{
		App:           "app",
		Availability:  "availability",
		BatchDB:       "batchdb",
		BatchFetcher:  "batchfetcher",
		Checkpointing: "checkpoint",
		ChkpValidator: "chkpvalidator",
		Crypto:        "crypto",
		ThreshCrypto:  "threshcrypto",
		Hasher:        "hasher",
		Mempool:       "mempool",
		Net:           "net",
		ReliableNet:   "reliablenet",
		Null:          "null",
		Timer:         "timer",

		Ordering:           "ordering",
		ISS:                "iss",
		PPrepValidator:     "pprepvalidator",
		PPrepValidatorChkp: "pprepvalidatorchkp",

		AleaDirector:  "iss", // fills in the role of iss in most modules
		AleaBroadcast: "availability",
		AleaAgreement: "ordering",
	}
}

func (mc ModuleConfig) ConfigureISS() iss.ModuleConfig {
	return iss.ModuleConfig{
		Self:               mc.ISS,
		App:                mc.BatchFetcher,
		Availability:       mc.Availability,
		BatchDB:            mc.BatchDB,
		Checkpoint:         mc.Checkpointing,
		ChkpValidator:      mc.ChkpValidator,
		Net:                mc.Net,
		Ordering:           mc.Ordering,
		PPrepValidator:     mc.PPrepValidator,
		PPrepValidatorChkp: mc.PPrepValidatorChkp,
		Timer:              mc.Timer,
	}
}

func (mc ModuleConfig) ConfigureAleaDirector() director.ModuleConfig {
	return director.ModuleConfig{
		Self:          mc.AleaDirector,
		BatchFetcher:  mc.BatchFetcher,
		Checkpoint:    mc.Checkpointing,
		ChkpValidator: mc.ChkpValidator,
		AleaBroadcast: mc.AleaBroadcast,
		AleaAgreement: mc.AleaAgreement,
		BatchDB:       mc.BatchDB,
		Mempool:       mc.Mempool,
		Net:           mc.Net,
		ReliableNet:   mc.ReliableNet,
		Hasher:        mc.Hasher,
		ThreshCrypto:  mc.ThreshCrypto,
		Timer:         mc.Timer,
		Null:          mc.Null,
	}
}

func (mc ModuleConfig) ConfigureCheckpointing() checkpoint.ModuleConfig {
	return checkpoint.ModuleConfig{
		Self:   mc.Checkpointing,
		App:    mc.BatchFetcher,
		Crypto: mc.Crypto,
		Hasher: mc.Hasher,
		Net:    mc.Net,
		Ord:    mc.ISS,
	}
}

func (mc ModuleConfig) ConfigureChkpValidator() chkpvalidator.ModuleConfig {
	return chkpvalidator.ModuleConfig{
		Self: mc.ChkpValidator,
	}
}

func (mc ModuleConfig) ConfigureThreshCheckpointing() threshcheckpoint.ModuleConfig {
	return threshcheckpoint.ModuleConfig{
		Self:         mc.Checkpointing,
		App:          mc.BatchFetcher,
		Hasher:       mc.Hasher,
		ThreshCrypto: mc.ThreshCrypto,
		ReliableNet:  mc.ReliableNet,
		Ord:          mc.AleaDirector, // <=> iss
	}
}

func (mc ModuleConfig) ConfigureThreshChkpValidator() threshchkpvalidator.ModuleConfig {
	return threshchkpvalidator.ModuleConfig{
		Self: mc.ChkpValidator,
	}
}

func (mc ModuleConfig) ConfigureOrdering() ordererscommon.ModuleConfig {
	return ordererscommon.ModuleConfig{
		Self:           mc.Ordering,
		App:            mc.BatchFetcher,
		Ava:            "", // Ava not initialized yet. It will be set at sub-module instantiation within the factory.
		Crypto:         mc.Crypto,
		Hasher:         mc.Hasher,
		Net:            mc.Net,
		Ord:            mc.ISS,
		PPrepValidator: mc.PPrepValidator,
		Timer:          mc.Timer,
	}
}

func (mc ModuleConfig) ConfigureAleaAgreement() agreement.ModuleConfig {
	return agreement.ModuleConfig{
		Self:         mc.AleaAgreement,
		AleaDirector: mc.AleaDirector,
		Hasher:       mc.Hasher,
		ReliableNet:  mc.ReliableNet,
		Net:          mc.Net,
		ThreshCrypto: mc.ThreshCrypto,
	}
}

func (mc ModuleConfig) ConfigurePreprepareValidator() pprepvalidator.ModuleConfig {
	return pprepvalidator.ModuleConfig{
		Self: mc.PPrepValidator,
	}
}

func (mc ModuleConfig) ConfigurePreprepareValidatorChkp() pprepvalidator.ModuleConfig {
	return pprepvalidator.ModuleConfig{
		Self: mc.PPrepValidatorChkp,
	}
}

func (mc ModuleConfig) ConfigureSimpleMempool() simplemempool.ModuleConfig {
	return simplemempool.ModuleConfig{
		Self:   mc.Mempool,
		Hasher: mc.Hasher,
		Timer:  mc.Timer,
	}
}

func (mc ModuleConfig) ConfigureFakeBatchDB() fakebatchdb.ModuleConfig {
	return fakebatchdb.ModuleConfig{
		Self: mc.BatchDB,
	}
}

func (mc ModuleConfig) ConfigureMultisigCollector() multisigcollector.ModuleConfig {
	return multisigcollector.ModuleConfig{
		Self:    mc.Availability,
		Mempool: mc.Mempool,
		BatchDB: mc.BatchDB,
		Net:     mc.Net,
		Crypto:  mc.Crypto,
	}
}

func (mc ModuleConfig) ConfigureAleaBroadcast() broadcast.ModuleConfig {
	return broadcast.ModuleConfig{
		Self:         mc.AleaBroadcast,
		AleaDirector: mc.AleaDirector,
		BatchDB:      mc.BatchDB,
		Mempool:      mc.Mempool,
		Net:          mc.Net,
		ReliableNet:  mc.ReliableNet,
		ThreshCrypto: mc.ThreshCrypto,
		Timer:        mc.Timer,
	}
}

func (mc ModuleConfig) ConfigureBatchFetcher() batchfetcher.ModuleConfig {
	return batchfetcher.ModuleConfig{
		Self:         mc.BatchFetcher,
		Availability: mc.Availability,
		Checkpoint:   mc.Checkpointing,
		Destination:  mc.App,
		Mempool:      mc.Mempool,
	}
}

func (mc ModuleConfig) ConfigureReliableNet() reliablenet.ModuleConfig {
	return reliablenet.ModuleConfig{
		Self:  mc.ReliableNet,
		Net:   mc.Net,
		Timer: mc.Timer,
	}
}
