package verbs

import(
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abworrall/dicebot/pkg/character"
	"github.com/abworrall/dicebot/pkg/config"
	"github.com/abworrall/dicebot/pkg/dnd5e/encounter"
	"github.com/abworrall/dicebot/pkg/state"
)

// VerbContext is passed by value to all the various verbs. If any
// fields are to be mutated by verbs, they need to be pointers, and be
// initialized during `Setup`.

type VerbContext struct {
	Ctx           context.Context
	state.StateManager

	// State we prepopulate, to be helpful
	Character      *character.Character
	MasqueradeAs    string // User we want to masquerade as
	Encounter      *encounter.Encounter
	
	// Request specific fields
	ExternalUserId  string // External systems provide their own IDs
	User            string // This should be the bot's nickname for the user

	Events       *[]Event  // Basically an audit log, that we shorhorn into the History verb
	Debug           string
}

// LogEvent should be used by verbs when they want to put things in the Log Of Record.
func (vc *VerbContext)LogEvent(description string) {
	*vc.Events = append(*vc.Events, Event{"", time.Now(), vc.User, description})
}


// Populate pulls in  all the interesting information for this user, and
// puts it in the context ready for the verb to use.
func (vc *VerbContext)Setup() {
	name, bs := "bot-state", BotSetup{} // FIXME: breaks layering; we're loading a verb's state
	if err := vc.StateManager.ReadState(vc.Ctx, name, &bs); err != nil {
		log.Printf("ReadState(%s): %v", name, err)
		return
	}

	vc.Events = &[]Event{}

	vc.User = bs.NameClaims[vc.ExternalUserId] // will be nil if they have no claim

	if vc.MasqueradeAs != "" && vc.User == config.Get("adminuser") {
		vc.User = vc.MasqueradeAs
		log.Printf("[masquerading as %s]\n", vc.MasqueradeAs)
	}
	
	if vc.User == "" { return }

	vc.loadContextCharacter()
	vc.loadContextEncounter()
}

func (vc *VerbContext)Teardown() {
	vc.maybeSaveContextCharacter()
	vc.maybeSaveContextEncounter()
	
	if len(*vc.Events) > 0 {
		stateName := "history-state" // FIXME: need a better way to identify singleton state
		h := NewHistory()

		if err := vc.StateManager.ReadState(vc.Ctx, stateName, &h); err != nil {
			// FIXME: breaks layering; we're loading a verb's state
			log.Printf("ReadState(%s): %v", stateName, err)
			return
		}

		h.Events = append (h.Events, *vc.Events...)

		if err := vc.StateManager.WriteState(vc.Ctx, stateName, &h); err != nil {
			log.Printf("%T.WriteState(%s, %T): %v\n", vc.StateManager, stateName, vc.Character, err)
		}
	}
}

// FIXME: replace these with something less awful
func (vc *VerbContext)characterStateName(name string) string { return fmt.Sprintf("char-state-%s", name) }
func (vc *VerbContext)encounterStateName() string { return "encounter-state" }

func (vc *VerbContext)loadContextCharacter() {
	if vc.User == "" {
		return
	}
	vc.Character = vc.loadCharacter(vc.User)
}

func (vc *VerbContext)maybeSaveContextCharacter() {
	// FIXME: ideally this would check for changes before writing the character
	if vc.User == "" || vc.Character == nil {
		return
	}
	vc.maybeSaveCharacter(vc.Character)
}

func (vc *VerbContext)loadContextEncounter() {
	e := encounter.NewEncounter()
	vc.Encounter = &e
	if vc.User == "" {
		return
	}
	if err := vc.StateManager.ReadState(vc.Ctx, vc.encounterStateName(), &e); err != nil {
		log.Printf("ReadState(%q): %v", vc.encounterStateName(), err)
	}
}
func (vc *VerbContext)maybeSaveContextEncounter() {
	// FIXME: ideally this would check for changes before writing the character
	if vc.User == "" || vc.Encounter == nil {
		return
	}
	if err := vc.StateManager.WriteState(vc.Ctx, vc.encounterStateName(), vc.Encounter); err != nil {
		log.Printf("%T.WriteState(%q, %T): %v\n", vc.StateManager, vc.encounterStateName(), vc.Encounter, err)
	}
}


// These routines might be used by other verbs
func (vc *VerbContext)loadCharacter(name string) *character.Character{
	c := character.NewCharacter()
	if err := vc.StateManager.ReadState(vc.Ctx, vc.characterStateName(name), &c); err != nil {
		// This will happen on the first ever read of a new character
		log.Printf("ReadState(%q): %v", vc.characterStateName(name), err)
	}

	return &c
}

func (vc *VerbContext)maybeSaveCharacter(c *character.Character) {
	// FIXME: ideally this would check for changes before writing the character
	if err := vc.StateManager.WriteState(vc.Ctx, vc.characterStateName(c.Name), c); err != nil {
		log.Printf("%T.WriteState(%s, %T): %v\n", vc.StateManager, vc.characterStateName(c.Name), c, err)
	}
}
