package main

import (
	"math/rand"
	"time"
)

type Gacha struct {
	characters []*Character
	region     []float64
}

// Gachaを生成
// キャラクターのリストとその出現確率の累積値を管理するregionを格納
func NewGacha() (Gacha, error) {
	var gacha Gacha

	rows, err := db.Query("select id, name, likelihood from characters order by id asc")
	if err != nil {
		return Gacha{}, err
	}

	var total float64

	for rows.Next() {
		var c Character
		err := rows.Scan(&c.Id, &c.Name, &c.likelihood)
		if err != nil {
			return Gacha{}, err
		}
		gacha.characters = append(gacha.characters, &c)
		total += c.likelihood
	}

	gacha.region = []float64{0}
	var sum float64
	for _, c := range gacha.characters {
		c.likelihood /= total
		sum += c.likelihood
		gacha.region = append(gacha.region, sum)
	}
	return gacha, nil
}

// ガチャを引く回数に対して，ガチャの結果を返す
func (g Gacha) Draw(n int) []*Character {
	var chars []*Character
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
