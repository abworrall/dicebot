package dnd5e

import(
	"fmt"
	"log"
)

// Dnd5E wraps up the DnD 5th Edition API objects (see https://www.dnd5eapi.co/docs/#resource-lists)
type Dnd5E struct {
	SpellList
}

var(
	Dnd Dnd5E
)

// InitDnd5E populates the global var with the data it loads from `datadir`
func InitDnd5E(datadir string) {
	Dnd = LoadDnd5E(datadir)
	log.Printf("InitDnd5E: loaded %s\n", Dnd)
}
	
// LoadDnd5E loads up the JSON files it needs to find in `datadir`
func LoadDnd5E(datadir string) Dnd5E {
	return Dnd5E{
		SpellList: LoadSpells(datadir + "/" + "5e-spells.json"),
	}
}

func (dnd Dnd5E)String() string {
	return fmt.Sprintf("API objects: %d spells\n", len(dnd.SpellList))
}
