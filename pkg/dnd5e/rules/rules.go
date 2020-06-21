package rules

import(
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
	"os"
)

var(
	TheRules Rules // A dumb global var (cf. "singleton" :) to hold all the objects
)

// Dnd5E wraps up the DnD 5th Edition API objects (see https://www.dnd5eapi.co/docs/#resource-lists)
// https://github.com/bagelbits/5e-database
type Rules struct {
	EquipmentList
	SpellList
}

func (r Rules)String() string {
	return fmt.Sprintf("Rules{%d spells, %d items}\n", len(r.SpellList), len(r.EquipmentList))
}

// InitRules populates the global var with the data it loads from `datadir`
func Init(datadir string) {
	TheRules = load(datadir)
	log.Printf("InitRules: datadir=%q, got: %s\n", datadir, TheRules)
}


// loadRules loads up the JSON files it needs to find in `datadir`
// TODO: error handling
func load(datadir string) Rules {
	return Rules{
		EquipmentList: loadEquipment(datadir + "/" + "5e-equipment.json"),
		SpellList: loadSpells(datadir + "/" + "5e-spells.json"),
	}
}

// LoadSpells opens a file that should be a ton of JSON objects that parse into spells
func loadSpells(filename string) SpellList {
	sl := map[string]Spell{}

	if jsonF,err := os.Open(filename); err == nil {
		defer jsonF.Close()

		file, _ := ioutil.ReadAll(jsonF)
		spells := []Spell{}
		json.Unmarshal(file, &spells)

		for _,spell := range spells {
			sl[spell.Index] = spell
		}
		log.Printf("%s, loaded %d objects\n", filename, len(sl))
	} else {
		log.Printf("open %s: %v\n", filename, err)
	}
	
	return sl
}

func loadEquipment(filename string) EquipmentList {
	el := map[string]Item{}

	if jsonF,err := os.Open(filename); err == nil {
		defer jsonF.Close()

		file, _ := ioutil.ReadAll(jsonF)
		items := []Item{}
		json.Unmarshal(file, &items)

		for _,item := range items {
			el[item.Index] = item
		}
		log.Printf("%s, loaded %d objects\n", filename, len(el))
	} else {
		log.Printf("open %s: %v\n", filename, err)
	}
	
	return el
}
