package gacha

import (
	"ca-tech-dojo/model/character"
)

type Gacha struct {
	characters []*character.Character
	region     []float64
}