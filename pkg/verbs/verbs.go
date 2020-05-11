package verbs

import(
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
)

func init() {
	HandleVerb("save",    SavingThrow{})
	HandleVerb("roll",    Roll{})
	HandleVerb("party",  &Party{})
	HandleVerb("insult", &Insult{})
}

type VerbContext struct {
	Ctx           context.Context
	StateManager

	// Request specific fields
	User          string
	Debug         string
}

// A Verber will respond to a bot command
type Verber interface {
	Process(c VerbContext, args []string) string
	Help() string
}

// A StateManager is a thing that can load/persist a verb's state. The caller
// should place one in the VerbContext.
type StateManager interface {
	ReadState(ctx context.Context, key string, ptr interface{}) error
	WriteState(ctx context.Context, key string, ptr interface{}) error
}

// verbs is the global registry of things the bot can do
var verbs = map[string]Verber{}

// HandleVerb is a how verb registers itself into the bot. This has
// a subtle trick to handle statefulness.
//
// If you register a pointer object (e.g. `HandleVerb("name",
// &MyVerb{})`), then the bot framework will consider it stateful, and
// load/persist it each time the verb runs. And you should make the
// interface methods act on pointers (e.g. `func (v *MyVerb)Help()
// ...`).
//
// If you just register a regular object (e.g. `HandleVerb("name2",
// MyVerb2{})`, then it is considered stateless, and your interface
// methods should act on objects (e.g. `func (v MyVerb2)Help()`).
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
	return str
}

// Act looks to see if this is a verb we handle, and if so handles it.
// Will load/persist stateful objects as needed.
func Act(vc VerbContext, v string, args []string) string {
	stateName := v+"-state"
	if v == "help" { return Help() }

	vr,exists := verbs[v]
	if !exists {
		return fmt.Sprintf("I don't `%s`", v)
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

	return resp
}

// Other verbs need access to the party info, so here is a helper to hide the details.
func LoadParty(vc VerbContext) *Party {
	stateName := "party-state"
	p := Party{}

	if err := vc.StateManager.ReadState(vc.Ctx, stateName, &p); err != nil {
		log.Printf("ReadState(%s): %v", stateName, err)
	}

	return &p
}
