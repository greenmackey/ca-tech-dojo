package gacha

import (
	"ca-tech-dojo/db"
	"ca-tech-dojo/model/character"
	"log"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

// Gachaを生成
// キャラクターのリストとその出現確率の累積値を管理するregionを格納
func NewGacha() (Gacha, error) {
	var gacha Gacha

	q := "SELECT id, name, likelihood FROM characters ORDER BY id ASC"
	rows, err := db.DB.Query(q)
	if err != nil {
		return Gacha{}, errors.Wrapf(err, "query failed in %s", "NewGacha")
	}
	log.Print(q)

	var total float64

	for rows.Next() {
		var c character.Character
		err := rows.Scan(&c.Id, &c.Name, &c.Likelihood)
		if err != nil {
			return Gacha{}, errors.Wrapf(err, "scan failed in %s", "NewGacha")
		}
		gacha.characters = append(gacha.characters, &c)
		total += c.Likelihood
	}

	gacha.region = []float64{0}
	var sum float64
	for _, c := range gacha.characters {
		c.Likelihood /= total
		sum += c.Likelihood
		gacha.region = append(gacha.region, sum)
	}
	return gacha, nil
}

// ガチャを引く回数に対して，ガチャの結果を返す
func (g Gacha) Draw(n uint) []*character.Character {
	var chars []*character.Character
	rand.Seed(time.Now().UnixNano())
	for ; n > 0; n-- {
		p := rand.Float64()
		for i := 0; i < len(g.characters); i++ {
			if g.region[i] <= p && p < g.region[i+1] {
				chars = append(chars, g.characters[i])
				break
			}
		}
	}
	return chars
}
