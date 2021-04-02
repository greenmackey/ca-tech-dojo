package reluc

import (
	"ca-tech-dojo/model/character"
	"ca-tech-dojo/model/reluc"

	"github.com/pkg/errors"
)

func SaveCharacters(token string, characters []*character.Character) error {
	rels := make([]reluc.Relationship, 0, len(characters))

	for _, c := range characters {
		rel := reluc.Relationship{
			UserToken:   token,
			CharacterId: c.Id,
		}
		rels = append(rels, rel)
	}

	if err := reluc.BulkCreate(rels); err != nil {
		return errors.Wrap(err, "reluc.BulkCreate failed")
	}

	return nil
}
