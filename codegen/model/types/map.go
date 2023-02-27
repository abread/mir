package types

import (
	"github.com/dave/jennifer/jen"
)

// already defined in slice.go
//var thisPackage = reflect.TypeOf(Map{}).PkgPath()

// Slice is used to represent repeated fields.
type Map struct {
	Key Type
	Val Type
}

func (s Map) Same() bool {
	return s.Key.Same() && s.Val.Same()
}

func (s Map) PbType() *jen.Statement {
	return jen.Map(s.Key.PbType()).Add(s.Val.PbType())
}

func (s Map) MirType() *jen.Statement {
	return jen.Map(s.Key.MirType()).Add(s.Val.MirType())
}

func (s Map) ToMir(code jen.Code) *jen.Statement {
	if s.Same() {
		return jen.Add(code)
	}

	return jen.Qual(thisPackage, "ConvertMap").Call(code,
		pb2MirFn(s.Key),
		pb2MirFn(s.Val),
	)
}

func (s Map) ToPb(code jen.Code) *jen.Statement {
	if s.Same() {
		return jen.Add(code)
	}

	return jen.Qual(thisPackage, "ConvertMap").Call(code,
		mir2PbFn(s.Key),
		mir2PbFn(s.Val),
	)
}

func mir2PbFn(t Type) jen.Code {
	return jen.Func().Params(jen.Id("x").Add(t.MirType())).Add(t.PbType()).Block(
		jen.Return(t.ToPb(jen.Id("x"))),
	)
}

func pb2MirFn(t Type) jen.Code {
	return jen.Func().Params(jen.Id("x").Add(t.PbType())).Add(t.MirType()).Block(
		jen.Return(t.ToMir(jen.Id("x"))),
	)
}

// ConvertMap is used by the generated code.
func ConvertMap[K1, K2 comparable, V1, V2 any](tm map[K1]V1, keyMapper func(k K1) K2, valMapper func(v V1) V2) map[K2]V2 {
	rm := make(map[K2]V2, len(tm))
	for k, v := range tm {
		rm[keyMapper(k)] = valMapper(v)
	}
	return rm
}
