package rules

import(
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// loadRules loads up the JSON files it needs to find in `datadir`
// TODO: error handling
func load(datadir string) Rules {
	return Rules{
		EquipmentList:   loadEquipment   (datadir + "/" + "5e-equipment.json"),
		SpellList:       loadSpells      (datadir + "/" + "5e-spells.json"),
		MonsterList:     loadMonsters    (datadir + "/" + "5e-monsters.json"),
		BuffList:        loadBuffs       (datadir + "/" + "5e-features.json"),
		SpellDamageList: loadSpellsDamage(datadir + "/" + "5e-spells-damage.json"),
	}
}

// TODO: write this boilerplate once, with some type cleverness

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

func loadMonsters(filename string) MonsterList {
	ml := map[string]Monster{}

	if jsonF,err := os.Open(filename); err == nil {
		defer jsonF.Close()

		file, _ := ioutil.ReadAll(jsonF)
		monsters := []Monster{}
		json.Unmarshal(file, &monsters)

		for _,monster := range monsters {
			monster.PostLoadFixups()
			ml[monster.Index] = monster
		}
		log.Printf("%s, loaded %d objects\n", filename, len(ml))
	} else {
		log.Printf("open %s: %v\n", filename, err)
	}
	
	return ml
}

func loadBuffs(filename string) BuffList {
	bl := map[string]Buff{}

	if jsonF,err := os.Open(filename); err == nil {
		defer jsonF.Close()

		file, _ := ioutil.ReadAll(jsonF)
		buffs := []Buff{}
		json.Unmarshal(file, &buffs)

		for _,buff := range buffs {
			bl[buff.Index] = buff
		}
		log.Printf("%s, loaded %d objects\n", filename, len(bl))
	} else {
		log.Printf("open %s: %v\n", filename, err)
	}
	
	return bl
}

// LoadSpells opens a file that should be a ton of JSON objects that parse into spells
func loadSpellsDamage(filename string) SpellDamageList {
	sl := map[string]SpellDamage{}

	if jsonF,err := os.Open(filename); err == nil {
		defer jsonF.Close()

		file, _ := ioutil.ReadAll(jsonF)
		spelldamages := []SpellDamage{}
		err := json.Unmarshal(file, &spelldamages)
		if err != nil {
			log.Printf("%s, err: %v\n", filename, err)
		}
		
		for _,spelldamage := range spelldamages {
			sl[spelldamage.Index] = spelldamage
			log.Printf("%s\n", spelldamage)
			
		}
		log.Printf("%s, loaded %d objects\n", filename, len(sl))
	} else {
		log.Printf("open %s: %v\n", filename, err)
	}
	
	return sl
}
