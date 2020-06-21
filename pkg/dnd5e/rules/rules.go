package rules

import(
	"fmt"
	"log"
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

// These routines are helpers for things that want to refer to the rules

// IsWeapon verifies that the index-string passed in will lookup a weapon
func (r Rules)IsWeapon(s string) bool {
	if item,exists := r.EquipmentList[s]; exists {
		return item.EquipmentCategory.Name == "Weapon"
	}
	return false
}

// IsSpell verifies that the index-string passed in will lookup a
// spell
func (r Rules)IsSpell(s string) bool {
	_,exists := r.SpellList[s]
	return exists
}

// IsAllowedSpell verifies that the index-string passed in will lookup
// a spell at the specified level (and type ?)
func (r Rules)IsAllowedSpell(s string, lvl int) bool {
	if spell,exists := r.SpellList[s]; exists {
		return spell.Level == lvl
	}
	return false
}
