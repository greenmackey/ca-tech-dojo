package reluc

import (
	"ca-tech-dojo/model/character"
	"ca-tech-dojo/model/reluc"
)

func ToRelationship(token string, characters []*character.Character) []reluc.Relationship {
	rels := make([]reluc.Relationship, 0, len(characters))

	for _, c := range characters {
		rel := reluc.Relationship{
			UserToken:   token,
			CharacterId: c.Id,
		}
		rels = append(rels, rel)
	}

	return rels
}
