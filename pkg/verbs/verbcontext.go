package verbs

import(
	"context"
	"fmt"
	"log"

	"github.com/abworrall/dicebot/pkg/character"
	"github.com/abworrall/dicebot/pkg/config"
)

type VerbContext struct {
	Ctx           context.Context
	StateManager

	// State we prepopulate, to be helpful. Breaks layering.
	Character      *character.Character
	MasqueradeAs    string // User we want to masquerade as
	
	// Request specific fields
	ExternalUserId  string // External systems provide their own IDs
	User            string // This should be the bot's nickname for the user
	Debug           string
}

// Populate pulls in  all the interesting information for this user, and
// puts it in the context ready for the verb to use.
func (vc *VerbContext)Setup() {
	name, bs := "bot-state", BotSetup{} // FIXME: breaks layering; we're loading a verb's state
	if err := vc.StateManager.ReadState(vc.Ctx, name, &bs); err != nil {
		log.Printf("ReadState(%s): %v", name, err)
		return
	}

	vc.User = bs.NameClaims[vc.ExternalUserId] // will be nil if they have no claim

	if vc.MasqueradeAs != "" && vc.User == config.Get("adminuser") {
		vc.User = vc.MasqueradeAs
		log.Printf("[masquerading as %s]\n", vc.MasqueradeAs)
	}
	
	if vc.User == "" { return }

	// FIXME: Something for unknown users ?

	vc.loadCharacter()
}

func (vc *VerbContext)Teardown() {
	vc.maybeSaveCharacter()
}


func (vc *VerbContext)loadCharacter() {
	c := character.NewCharacter()
	vc.Character = &c

	if vc.User == "" {
		return
	} else if err := vc.StateManager.ReadState(vc.Ctx, vc.characterStateName(), &c); err != nil {
		// This will happen on the first ever read of a new character
		log.Printf("ReadState(%s): %v", vc.characterStateName(), err)
		return
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

func (vc *VerbContext)characterStateName() string {
	return fmt.Sprintf("char-state-%s", vc.User)
}
