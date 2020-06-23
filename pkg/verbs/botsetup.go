package verbs

import(
	"fmt"
)

// BotSetup contains data that the rest of the bot uses.
// Most notably, users need to claim their display names.
type BotSetup struct {
	NameClaims map[string]string // Key is whatever foreign opaque string we start with
}

func (bs *BotSetup)Help() string {
	return "claim USER"
}

func (bs *BotSetup)MaybeInit() {
	if bs.NameClaims == nil {
		bs.NameClaims = map[string]string{}
	}
}

func (bs *BotSetup)Process(vc VerbContext, args []string) string {
	bs.MaybeInit()

	if len(args) == 0 {
		return fmt.Sprintf("%#v", bs)
	}
	
	switch args[0] {
	case "debug":        return bs.Debug(vc)
	case "claim":        return bs.Claim(vc, args[1])
	case "-deleteclaim": return bs.ClaimDelete(args[1])
	case "-flushclaims": bs.NameClaims = map[string]string{}
	}

	return ""
}

func (bs *BotSetup)Claim(vc VerbContext, name string) string {
	if vc.ExternalUserId == "" {
		return "I can't tell who you are - you need to add me to your friend list (and agree to the ToU)"
	}
	if _,exists := bs.NameClaims[vc.ExternalUserId]; !exists {
		bs.NameClaims[vc.ExternalUserId] = name
		return fmt.Sprintf("%s has been claimed by %s", name, vc.ExternalUserId)
	} else {
		return fmt.Sprintf("naughty ! you already have a claim")
	}
}

func (bs *BotSetup)ClaimDelete(name string) string {
	for k,v := range bs.NameClaims {
		if v == name {
			delete(bs.NameClaims, k)
			return fmt.Sprintf("%s is unclaimed", name)
		}
	}
	return fmt.Sprintf("no claim found for '%s'", name)
}

func (bs *BotSetup)Debug(vc VerbContext) string {
	return fmt.Sprintf("vc: %#v\n\nc: %#v\n\nencounter: %#v\n\nsetup: %#v", vc, vc.Character, vc.Encounter, bs)
}
