package agreementpb

import (
	reflect "reflect"
)

func (*Message) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Message_FinishAbba)(nil)),
	}
}
