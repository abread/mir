// Code generated by Mir codegen. DO NOT EDIT.

package messagepb

import (
	reflect "reflect"
)

func (*Message) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Message_Iss)(nil)),
		reflect.TypeOf((*Message_Bcb)(nil)),
		reflect.TypeOf((*Message_MultisigCollector)(nil)),
		reflect.TypeOf((*Message_Pingpong)(nil)),
		reflect.TypeOf((*Message_Checkpoint)(nil)),
		reflect.TypeOf((*Message_Orderer)(nil)),
		reflect.TypeOf((*Message_Vcb)(nil)),
		reflect.TypeOf((*Message_Abba)(nil)),
		reflect.TypeOf((*Message_AleaBroadcast)(nil)),
		reflect.TypeOf((*Message_AleaAgreement)(nil)),
		reflect.TypeOf((*Message_AleaDirector)(nil)),
		reflect.TypeOf((*Message_ReliableNet)(nil)),
		reflect.TypeOf((*Message_Threshcheckpoint)(nil)),
	}
}
