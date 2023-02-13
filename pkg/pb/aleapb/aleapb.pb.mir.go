package aleapb

import (
	reflect "reflect"
)

func (*Message) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Message_FillGapMessage)(nil)),
		reflect.TypeOf((*Message_FillerMessage)(nil)),
	}
}
