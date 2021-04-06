package usercharacter

import (
	"ca-tech-dojo/model/character"
	"ca-tech-dojo/model/usercharacter"
)

func ToRelationship(token string, characters []*character.Character) []usercharacter.Relationship {
	rels := make([]usercharacter.Relationship, 0, len(characters))

	for _, c := range characters {
		rel := usercharacter.Relationship{
			UserToken:   token,
			CharacterId: c.Id,
		}
		rels = append(rels, rel)
	}

	return rels
}
