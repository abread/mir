// Code generated by Mir codegen. DO NOT EDIT.

package threshcheckpointpb

import (
	reflect "reflect"
)

func (*Message) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Message_Checkpoint)(nil)),
	}
}

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_StableCheckpoint)(nil)),
	}
}
