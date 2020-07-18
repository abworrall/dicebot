package verbs

import(
	"fmt"
	"strconv"

	"github.com/abworrall/dicebot/pkg/rules"
)

// Character is stateless, in that the verb doesn't have its own state; it simply operates
// on the character's state in the context
type Character struct{}

func (c Character)Help() string {
	return "[set FIELD VALUE], [list], [remove FIELD VALUE]"
}

func (c Character)Process(vc VerbContext, args []string) string {
	if vc.User == "" {
		return "who are you, eh ?"
	}

	if len(args) == 0 {
		return fmt.Sprintf("%s", vc.Character)
	}

	switch args[0] {		
	case "set":
		if len(args) != 3 { return "`set field value`, plz" }
		return c.Set(vc, args[1], args[2])

	case "remove":
		if len(args) != 3 { return "`remove field entry`, plz" }
		return c.Remove(vc, args[1], args[2])

	case "list":
		return "[Useful fields: weapon, armor, shield]\n"+
			"[Less useful fields: race, class, alignment, level, maxhp, buff, hp, str, int, wis, con, cha, dex per]"

	case "setstats":
		if len(args) != 8 { return "`setstats 1 2 3 4 5 6 7`, plz" }
		c.Set(vc, "str", args[1])
		c.Set(vc, "int", args[2])
		c.Set(vc, "wis", args[3])
		c.Set(vc, "con", args[4])
		c.Set(vc, "cha", args[5])
		c.Set(vc, "dex", args[6])
		c.Set(vc, "per", args[7])

	default: return "wat?"
	}

	return ""
}

func (c Character)Set(vc VerbContext, k,v string) string {
	myatoi := func(s string) int {
		i,_ := strconv.Atoi(s)
		return i
	}

	switch k {
	case "name": vc.Character.Name = v
	case "race": vc.Character.Race = v
	case "class": vc.Character.Class = v
	case "subclass": vc.Character.Subclass = v
	case "alignment": vc.Character.Alignment = v

	case "buff":
		if err := vc.Character.AddBuff(v); err != nil {
			return fmt.Sprintf("bad buff: %v", err)
		}

	case "weapon":
		if ! rules.TheRules.IsWeapon(v) {
			return fmt.Sprintf("'%s' is not a known weapon", v)
		}
		vc.Character.Weapons[v] = 1
		vc.Character.CurrWeapon = v

	case "shield": vc.Character.Shield = (myatoi(v) == 1)
		
	case "armor":
		if ! rules.TheRules.IsArmor(v) {
			return fmt.Sprintf("'%s' is not a known kind of armor", v)
		}
		vc.Character.Armor = v
		
	case "str": vc.Character.Str = myatoi(v)
	case "int": vc.Character.Int = myatoi(v)
	case "wis": vc.Character.Wis = myatoi(v)
	case "con": vc.Character.Con = myatoi(v)
	case "cha": vc.Character.Cha = myatoi(v)
	case "dex": vc.Character.Dex = myatoi(v)
	case "per": vc.Character.Per = myatoi(v)

	case "level": vc.Character.Level = myatoi(v)
	case "maxhp": vc.Character.MaxHitpoints = myatoi(v)
	case "hp": vc.Character.CurrHitpoints = myatoi(v)

	default: return fmt.Sprintf("I don't set '%s'", k)
	}

	return fmt.Sprintf("%s set to %s", k, v)
}

func (c Character)Remove(vc VerbContext, k,v string) string {

	switch k {
	case "buff":
		if err := vc.Character.RemoveBuff(v); err != nil {
			return fmt.Sprintf("bad buff: %v", err)
		}
	default: return fmt.Sprintf("I don't remove '%s'", k)
	}

	return fmt.Sprintf("%s has %s removed", k, v)
}
