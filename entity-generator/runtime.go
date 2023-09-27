package entity_generator

// import (
// 	"fmt"

// 	"github.com/mabels/wueste/entity-generator/rusty"
// )

type PropertyCtx struct {
	Registry *SchemaRegistry
	// BaseDir  string
}

// type PropertyRuntime struct {
// 	FileName rusty.Optional[string]
// 	Parent   rusty.Optional[Property]
// 	Ref      rusty.Optional[string]
// 	BaseDir  rusty.Optional[string]
// 	Of       Property
// }

// // func NewRuntime(b SchemaLoader) PropertyMeta {
// // 	return PropertyMeta{
// // 		Loader: b,
// // 	}
// // }

// func walkConnectRuntime(curr rusty.Result[Property], parent rusty.Optional[Property]) rusty.Result[Property] {
// 	p := curr.Ok()
// 	switch p.Type() {
// 	case OBJECT:
// 		po := p.(PropertyObject)
// 		po.Runtime().Parent = parent
// 		po.Runtime().Of = po
// 		for _, item := range po.Items() {
// 			myparent := rusty.Some(p)
// 			item.Property().Runtime().Parent = myparent
// 			ret := walkConnectRuntime(rusty.Ok(item.Property()), myparent)
// 			if ret.IsErr() {
// 				return ret
// 			}
// 		}
// 	case ARRAY:
// 		pa := p.(PropertyArray)
// 		pa.Runtime().Parent = parent
// 		pa.Runtime().Of = pa
// 		ret := walkConnectRuntime(rusty.Ok(pa.Items()), rusty.Some(p))
// 		if ret.IsErr() {
// 			return ret
// 		}
// 	case STRING, NUMBER, INTEGER, BOOLEAN:
// 		p.Runtime().Of = p
// 		p.Runtime().Parent = parent
// 	default:
// 		return rusty.Err[Property](fmt.Errorf("unknown type %s", p.Type()))
// 	}
// 	return curr
// }

// func ConnectRuntime(p rusty.Result[Property]) rusty.Result[Property] {
// 	if p.IsErr() {
// 		return rusty.Err[Property](p.Err())
// 	}
// 	return walkConnectRuntime(p, p.Ok().Runtime().Parent)
// }

// func (p *PropertyRuntime) SetFileName(name string) {
// 	p.FileName = rusty.Some(name)
// }

// func (p *PropertyRuntime) SetRef(name string) {
// 	p.Ref = rusty.Some(name)
// }

// func (p *PropertyRuntime) Assign(b PropertyRuntime) *PropertyRuntime {
// 	// p.Registry = b.Registry
// 	if b.Ref.IsSome() {
// 		p.Ref = b.Ref
// 	}
// 	if p.FileName.IsNone() {
// 		p.FileName = b.FileName
// 	}
// 	return p
// }

// func (p *PropertyRuntime) Clone() *PropertyRuntime {
// 	return (&PropertyRuntime{}).Assign(*p)
// }

// func (p *PropertyRuntime) ToPropertyObject() rusty.Result[PropertyObject] {
// 	var pi Property = p.Of
// 	po, ok := pi.(PropertyObject)
// 	if !ok {
// 		return rusty.Err[PropertyObject](fmt.Errorf("not a PropertyObject"))
// 	}
// 	return rusty.Ok[PropertyObject](po)
// }
