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

	vc.loadCharacter()
	vc.loadEncounter()
}

func (vc *VerbContext)Teardown() {
	vc.maybeSaveCharacter()
	vc.maybeSaveEncounter()
	
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
func (vc *VerbContext)characterStateName() string { return fmt.Sprintf("char-state-%s", vc.User) }
func (vc *VerbContext)encounterStateName() string { return "encounter-state" }

func (vc *VerbContext)loadCharacter() {
	c := character.NewCharacter()
	vc.Character = &c
	if vc.User == "" {
		return
	}
	if err := vc.StateManager.ReadState(vc.Ctx, vc.characterStateName(), &c); err != nil {
		// This will happen on the first ever read of a new character
		log.Printf("ReadState(%q): %v", vc.characterStateName(), err)
	}
}
func (vc *VerbContext)maybeSaveCharacter() {
	// FIXME: ideally this would check for changes before writing the character
	if vc.User == "" || vc.Character == nil {
		return
	} else if err := vc.StateManager.WriteState(vc.Ctx, vc.characterStateName(), vc.Character); err != nil {
		log.Printf("%T.WriteState(%s, %T): %v\n", vc.StateManager, vc.characterStateName(), vc.Character, err)
		return
	}
}

func (vc *VerbContext)loadEncounter() {
	e := encounter.NewEncounter()
	vc.Encounter = &e
	if vc.User == "" {
		return
	}
	if err := vc.StateManager.ReadState(vc.Ctx, vc.encounterStateName(), &e); err != nil {
		log.Printf("ReadState(%q): %v", vc.encounterStateName(), err)
	}
}
func (vc *VerbContext)maybeSaveEncounter() {
	// FIXME: ideally this would check for changes before writing the character
	if vc.User == "" || vc.Encounter == nil {
		return
	}
	if err := vc.StateManager.WriteState(vc.Ctx, vc.encounterStateName(), vc.Encounter); err != nil {
		log.Printf("%T.WriteState(%q, %T): %v\n", vc.StateManager, vc.encounterStateName(), vc.Encounter, err)
	}
}
