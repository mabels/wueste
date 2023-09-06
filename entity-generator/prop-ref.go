package entity_generator

import (
	"github.com/mabels/wueste/entity-generator/rusty"
)

type PropertyRef interface {
	Property
	Ref() rusty.Optional[string]
}
