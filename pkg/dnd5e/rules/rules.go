package rules

import(
	"fmt"
	"log"
)

// Dnd5E wraps up the DnD 5th Edition API objects (see https://www.dnd5eapi.co/docs/#resource-lists)
// https://github.com/bagelbits/5e-database
type Rules struct {
	SpellList
}

var(
	// Dnd is a dumb global var (cf. "singleton" :) to hold all the objects
	TheRules Rules
)

// InitRules populates the global var with the data it loads from `datadir`
func Init(datadir string) {
	TheRules = load(datadir)
	log.Printf("InitRules: datadir=%q, got: %s\n", datadir, TheRules)
}
	
// loadRules loads up the JSON files it needs to find in `datadir`
// TODO: error handling
func load(datadir string) Rules {
	return Rules{
		SpellList: LoadSpells(datadir + "/" + "5e-spells.json"),
	}
}

func (r Rules)String() string {
	return fmt.Sprintf("Rules{%d spells}\n", len(r.SpellList))
}
