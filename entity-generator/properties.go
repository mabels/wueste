package entity_generator

import (
	"sort"

	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertiesBuilder struct {
	items   []PropertyItem
	_loader SchemaLoader
}

func NewPropertiesBuilder(loader SchemaLoader) *PropertiesBuilder {
	return &PropertiesBuilder{
		_loader: loader,
	}
}

func isOptional(name string, req []string) bool {
	for _, r := range req {
		if r == name {
			return false
		}
	}
	return true
}

func interfaceToStringPtr(i interface{}) *string {
	if i == nil {
		return nil
	}
	s, found := i.(*string)
	if !found {
		return nil
	}
	return s
}

func (b *PropertiesBuilder) FromJson(js JSONProperties, req []string) *PropertiesBuilder {
	for name, p := range js {
		// falsch !!!
		var pn Property
		switch p.Type {
		case "string":
			pn = NewPropertyString(PropertyStringParam{
				PropertyParam: PropertyParam{
					Id:          p.Id,
					Type:        p.Type,
					Description: rusty.OptionalFromPtr(p.Description),
					Format:      rusty.OptionalFromPtr(p.Format),
					Optional:    isOptional(name, req),
				},
				Default: rusty.OptionalFromPtr(interfaceToStringPtr(p.Default)),
				Format:  rusty.OptionalFromPtr(p.Format),
			})
		case "boolean":
			pn = NewPropertyBoolean(PropertyBooleanParam{
				PropertyParam: PropertyParam{
					Id:          p.Id,
					Type:        p.Type,
					Description: rusty.OptionalFromPtr(p.Description),
					Format:      rusty.OptionalFromPtr(p.Format),
					Optional:    isOptional(name, req),
				},
				Default: rusty.OptionalFromPtr(p.Default.(*bool)),
			})
		case "integer":
			var format string
			if p.Format == nil {
				format = "int"
			} else {
				format = *p.Format
			}
			switch format {
			case "int":
				pn = NewPropertyInteger[int](PropertyIntegerParam[int]{
					PropertyLiteralParam: PropertyLiteralParam[int]{
						PropertyParam: PropertyParam{
							Id:          p.Id,
							Type:        p.Type,
							Description: rusty.OptionalFromPtr(p.Description),
							Format:      rusty.Some("int"),
							Optional:    isOptional(name, req),
						},
					},
					Default: rusty.OptionalFromPtr(p.Default.(*int)),
					Maximum: rusty.OptionalFromPtr(p.Maximum.(*int)),
					Minimum: rusty.OptionalFromPtr(p.Minimum.(*int)),
				})
			case "int8":
				pn = NewPropertyInteger[int8](PropertyIntegerParam[int8]{
					PropertyLiteralParam: PropertyLiteralParam[int8]{
						PropertyParam: PropertyParam{
							Id:          p.Id,
							Type:        p.Type,
							Description: rusty.OptionalFromPtr(p.Description),
							Format:      rusty.Some("int8"),
							Optional:    isOptional(name, req),
						},
					},
					Default: rusty.OptionalFromPtr(p.Default.(*int8)),
					Maximum: rusty.OptionalFromPtr(p.Maximum.(*int8)),
					Minimum: rusty.OptionalFromPtr(p.Minimum.(*int8)),
				})
			case "int16":
				pn = NewPropertyInteger[int16](PropertyIntegerParam[int16]{
					PropertyLiteralParam: PropertyLiteralParam[int16]{
						PropertyParam: PropertyParam{
							Id:          p.Id,
							Type:        p.Type,
							Description: rusty.OptionalFromPtr(p.Description),
							Format:      rusty.Some("float32"),
							Optional:    isOptional(name, req),
						},
					},
					Default: rusty.OptionalFromPtr(p.Default.(*int16)),
					Maximum: rusty.OptionalFromPtr(p.Maximum.(*int16)),
					Minimum: rusty.OptionalFromPtr(p.Minimum.(*int16)),
				})
			case "int32":
				pn = NewPropertyInteger[int32](PropertyIntegerParam[int32]{
					PropertyLiteralParam: PropertyLiteralParam[int32]{
						PropertyParam: PropertyParam{
							Id:          p.Id,
							Type:        p.Type,
							Description: rusty.OptionalFromPtr(p.Description),
							Format:      rusty.Some("float32"),
							Optional:    isOptional(name, req),
						},
					},
					Default: rusty.OptionalFromPtr(p.Default.(*int32)),
					Maximum: rusty.OptionalFromPtr(p.Maximum.(*int32)),
					Minimum: rusty.OptionalFromPtr(p.Minimum.(*int32)),
				})
			case "int64":
				pn = NewPropertyInteger[int64](PropertyIntegerParam[int64]{
					PropertyLiteralParam: PropertyLiteralParam[int64]{
						PropertyParam: PropertyParam{
							Id:          p.Id,
							Type:        p.Type,
							Description: rusty.OptionalFromPtr(p.Description),
							Format:      rusty.Some("float32"),
							Optional:    isOptional(name, req),
						},
					},
					Default: rusty.OptionalFromPtr(p.Default.(*int64)),
					Maximum: rusty.OptionalFromPtr(p.Maximum.(*int64)),
					Minimum: rusty.OptionalFromPtr(p.Minimum.(*int64)),
				})
			default:
				panic("unknown format")
			}
		case "number":
			var format string
			if p.Format == nil {
				format = "float64"
			} else {
				format = *p.Format
			}
			switch format {
			case "float32":
				pn = NewPropertyNumber[float32](PropertyNumberParam[float32]{
					PropertyLiteralParam: PropertyLiteralParam[float32]{
						PropertyParam: PropertyParam{
							Id:          p.Id,
							Type:        p.Type,
							Description: rusty.OptionalFromPtr(p.Description),
							Format:      rusty.Some("float32"),
							Optional:    isOptional(name, req),
						},
					},
					Default: rusty.OptionalFromPtr(p.Default.(*float32)),
					Maximum: rusty.OptionalFromPtr(p.Maximum.(*float32)),
					Minimum: rusty.OptionalFromPtr(p.Minimum.(*float32)),
				})
			case "float64":
				pn = NewPropertyNumber[float64](PropertyNumberParam[float64]{
					PropertyLiteralParam: PropertyLiteralParam[float64]{
						PropertyParam: PropertyParam{
							Id:          p.Id,
							Type:        p.Type,
							Description: rusty.OptionalFromPtr(p.Description),
							Format:      rusty.Some("float32"),
							Optional:    isOptional(name, req),
						},
					},
					Default: rusty.OptionalFromPtr(p.Default.(*float64)),
					Maximum: rusty.OptionalFromPtr(p.Maximum.(*float64)),
					Minimum: rusty.OptionalFromPtr(p.Minimum.(*float64)),
				})
			default:
				panic("unknown format")
			}
		case "object":
			pn = NewPropertyObject(PropertyObjectParam{
				FileName:    p.FileName,
				Id:          p.Id,
				Title:       p.Title,
				Schema:      p.Schema,
				Description: rusty.OptionalFromPtr(p.Description),
				Properties:  NewPropertiesBuilder(b._loader).FromJson(p.Properties, p.Required).Build(),
				Required:    p.Required,
				Ref:         rusty.OptionalFromPtr(p.Ref),
			})
		default:
			panic("unknown type")
		}
		b.Add(NewPropertyItem(name, pn))
	}
	return b
}

func (b *PropertiesBuilder) Add(property PropertyItem) *PropertiesBuilder {
	// property.SetOrder(len(b.items))
	b.items = append(b.items, property)
	return b
}

func (p *PropertiesBuilder) Items() []PropertyItem {
	sort.Slice(p.items, func(i, j int) bool {
		return p.items[i].Name() < p.items[j].Name()
	})
	return p.items
}

func (b *PropertiesBuilder) Build() PropertiesObject {
	return b
}

type propertyItem struct {
	name     string
	property Property
	// order    int
}

// Description implements PropertyItem.
func (pi *propertyItem) Name() string {
	return pi.name
}

// func (pi *propertyItem) Order() int {
// 	return pi.order
// }

// func (pi *propertyItem) SetOrder(order int) {
// 	pi.order = order
// }

func (pi *propertyItem) Property() Property {
	return pi.property
}

func NewPropertyItem(name string, property Property) PropertyItem {
	return &propertyItem{
		name:     name,
		property: property,
	}
}
