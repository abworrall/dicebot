package dnd5e

import(
	"fmt"
	"log"
)

// Dnd5E wraps up the DnD 5th Edition API objects (see https://www.dnd5eapi.co/docs/#resource-lists)
// https://github.com/bagelbits/5e-database
type Dnd5E struct {
	SpellList
}

var(
	// Dnd is a dumb global var (cf. "singleton" :) to hold all the objects
	Dnd Dnd5E
)

// InitDnd5E populates the global var with the data it loads from `datadir`
func InitDnd5E(datadir string) {
	Dnd = loadDnd5E(datadir)
	log.Printf("InitDnd5E: datadir=%q, got: %s\n", datadir, Dnd)
}
	
// LoadDnd5E loads up the JSON files it needs to find in `datadir`
func loadDnd5E(datadir string) Dnd5E {
	return Dnd5E{
		SpellList: LoadSpells(datadir + "/" + "5e-spells.json"),
	}
}

func (dnd Dnd5E)String() string {
	return fmt.Sprintf("Dnd5E{%d spells}\n", len(dnd.SpellList))
}
