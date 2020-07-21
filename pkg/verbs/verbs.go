package verbs

import(
	"fmt"
	"log"
	"reflect"
	"sort"
)

func init() {
	// These verbs are always needed; provides data objects that are used by framework
	HandleVerb("bot",     &BotSetup{})
	HandleVerb("history", &History{})

	// Per-character verbs (use state in context)
	HandleVerb("hp",       HitPoints{})
	HandleVerb("char",     Character{})
	HandleVerb("inv",      Inventory{})
	HandleVerb("spells",   Spells{})
	HandleVerb("roll",     Roll{})

	// Encounter data is stored in the context
	HandleVerb("attack",   Encounter{})
	
	// Verbs with their own state (not part of character objects)
	HandleVerb("rules",    Rules{})
	HandleVerb("list",    &Lists{})
	HandleVerb("vow",     &Vows{})
	HandleVerb("insult",  &Insult{})

	// Oddball; has its own state (list of names), but reads/writes all characters
	HandleVerb("party",   &Party{})
}

// A Verber will respond to a bot command
type Verber interface {
	Process(c VerbContext, args []string) string
	Help() string
}

// verbs is the global registry of things the bot can do
var verbs = map[string]Verber{}

// HandleVerb is a how verb registers itself into the bot. This has
// a subtle trick to handle statefulness.
//
// If you register a pointer object (e.g. `HandleVerb("name",
// &MyVerb{})`), then the bot framework will consider it stateful, and
// load/persist that object each time the verb runs. And you should make the
// interface methods act on pointers (e.g. `func (v *MyVerb)Help()
// ...`).
//
// If you just register a regular object (e.g. `HandleVerb("name2",
// MyVerb2{})`, then it is considered stateless, and you get a fresh
// empty one each time the verb runs. And your interface methods
// should act on objects (e.g. `func (v MyVerb2)Help()`).
func HandleVerb(v string, vr Verber) {
	if _,exists := verbs[v]; exists {
		log.Printf("Verb %s already registered, overwriting", v)
	}
	
	verbs[v] = vr
}

func IsStateful(vr Verber) bool {
	return reflect.ValueOf(vr).Kind() == reflect.Ptr 
}

// Help generates a help summary of all commands
func Help() string {
	keys := []string{}
	for k,_ := range verbs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	str := "I can:\n"
	for _,k := range keys {
		if help := verbs[k].Help(); help != "" {
			str += fmt.Sprintf("- db %s %s\n", k, help)
		}
	}

	str += "\nDetails at https://github.com/abworrall/dicebot/blob/master/README.md#verbs\n"

	return str
}

// Act looks to see if this is a verb we handle, and if so handles it.
// Will load/persist stateful objects as needed.
func Act(vc VerbContext, v string, args []string) string {
	stateName := v+"-state"
	if v == "help" { return Help() }

	// Masquerading - command will start with "as userblah ". Permissions verified in `vc.Setup`.
	if v == "as" && len(args) > 1 {
		vc.MasqueradeAs, v, args = args[0], args[1], args[2:]
	}

	vc.Setup()
	
	vr,exists := verbs[v]
	if !exists {
		return fmt.Sprintf("`%s` ? Computer says no", v)
	}

	if IsStateful(vr) {
		if vc.StateManager != nil {
			if err := vc.StateManager.ReadState(vc.Ctx, stateName, vr); err != nil {
				log.Printf("%T.ReadState(%s, %T): %v\n", vc.StateManager, stateName, vr, err)
			}
		} else {
			log.Printf("Asked to load/save verb state, but have no StateManager\n")
		}
	}

	resp := vr.Process(vc, args)
	
	if IsStateful(vr) {
		if vc.StateManager != nil {
			if err := vc.StateManager.WriteState(vc.Ctx, stateName, vr); err != nil {
				log.Printf("%T.WriteState(%s, %T): %v\n", vc.StateManager, stateName, vr, err)
			}
		}
	}

	vc.Teardown()
	
	return resp
}
