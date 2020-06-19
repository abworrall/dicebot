package dnd5e

import(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Spell struct{
	Index string `json:"index"`
	Name string `json:"name"`
	Desc[] string `json:"desc"`
	Higher[] string `json:"higher_level"`
	Range string `json:"range"`
	Duration string `json:"duration"`
	Level int `json:"level"`
	Classes []struct{
		Name string `json:"name"`
	} `json:"classes"`
}

func (s Spell)String() string {
	return fmt.Sprintf(`--{ %s }--
Level: %d (%s)
Range: %s
Duration: %s
%s
%s`, 
	s.Name,
	s.Level, s.Class(),
	s.Range, 
	s.Duration, 
	s.Desc,
	s.Higher)
}

func (s Spell)Class() string {
	names := []string{}
	for _,c := range s.Classes {
		names = append (names, c.Name)
	}
	return strings.Join(names, ",")
}

type SpellList map[string]Spell

// Find searches the spelllist, returns possible matching spell objects
func (sl SpellList)Find(namelike string) []Spell {
	ret := []Spell{}
	for _,v := range sl {
		if strings.Contains(strings.ToLower(v.Name),strings.ToLower(namelike)) {
			ret = append(ret, v)
		}
	}
	return ret
}

// LoadSpells opens a file that should be a ton of JSON objects that parse into spells
func LoadSpells(filename string) SpellList {
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
