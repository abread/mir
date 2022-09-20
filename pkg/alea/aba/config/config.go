package config

import t "github.com/filecoin-project/mir/pkg/types"

type Config struct {
	InstanceId uint64

	Members        []t.NodeID
	MembersSelfIdx int
	F              int

	SelfModuleID         t.ModuleID
	NetModuleID          t.ModuleID
	ThreshCryptoModuleID t.ModuleID
}

func (c *Config) SelfNodeID() t.NodeID {
	return c.Members[c.MembersSelfIdx]
}
