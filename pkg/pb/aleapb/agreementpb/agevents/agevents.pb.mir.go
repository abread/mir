package agevents

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_RequestInput)(nil)),
		reflect.TypeOf((*Event_InputValue)(nil)),
		reflect.TypeOf((*Event_Deliver)(nil)),
		reflect.TypeOf((*Event_StaleMsgsRevcd)(nil)),
	}
}
