package gacha

import (
	"ca-tech-dojo/model/character"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

type Gacha struct {
	characters []*character.Character
	region     []float64
}

// Gachaを生成
// キャラクターのリストとその出現確率の累積値を管理するregionを格納
func New() (Gacha, error) {
	var gacha Gacha
	var err error

	gacha.characters, err = character.All()
	if err != nil {
		return Gacha{}, errors.Wrap(err, "character.All failed")
	}

	var total float64
	for _, c := range gacha.characters {
		total += c.Likelihood
	}

	gacha.region = make([]float64, 1, len(gacha.characters)+1)
	var sum float64
	for _, c := range gacha.characters {
		sum += c.Likelihood / total
		gacha.region = append(gacha.region, sum)
	}
	return gacha, nil
}

// ガチャを引く回数に対して，ガチャの結果を返す
func (g Gacha) Draw(n uint) []*character.Character {
	var characters []*character.Character
	rand.Seed(time.Now().UnixNano())
	for ; n > 0; n-- {
		p := rand.Float64()
		for i := 0; i < len(g.characters); i++ {
			if g.region[i] <= p && p < g.region[i+1] {
				characters = append(characters, g.characters[i])
				break
			}
		}
	}
	return characters
}
