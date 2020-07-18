package roll

import(
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var Intn = rand.Intn // Allow RNG to be overridden, though we don't

func New(s string) Outcome {
	return Parse(s).Do()
}

type Roll struct {
	input    string
	Err      error
	
	NumDice  int
	DiceSize int
	Modifier int

	// These options only apply to 1d20 rolls
	Target               int
	WithAdvantage        bool
	WithDisadvantage     bool
	WithImprovedCritical bool  // Means you get a critical hit on 19 as well as 20
	
	Reason   string
}

func (r Roll)IsNil() bool { return r.NumDice == 0 && r.Modifier == 0 }

func (r Roll)String() string {
	s := fmt.Sprintf("%dd%d", r.NumDice, r.DiceSize)
	if r.Modifier != 0 {
		s += fmt.Sprintf("%+d", r.Modifier)
	}
	if r.Target > 0 {
		s += fmt.Sprintf(" >=%d", r.Target)
	}
	if r.WithAdvantage {
		s += " With Advantage"
	} else if r.WithDisadvantage {
		s += " With Disadvantage"
	}
	return s
}

type Outcome struct {
	Roll                // what we were asked to roll - makes this fully self-contained for String()ing

	Dice        []int   // The output dice that made up the total
	Ignored       int   // For rolls with advantage/disadvantage, rejected value goes here
	Total         int   // Total includes the modifier
	Success       bool  // If there was a target, did we meet/exceed it ?
	
	CriticalHit   bool  // rolled 1d20, and got a natural 20
	CriticalMiss  bool  // rolled 1d20, and got a natural 1
}

func (o Outcome)String() string {
	vals := make([]string, len(o.Dice))
	for i,n := range o.Dice {
		vals[i] = fmt.Sprintf("%d", n)
	}

	// If the roll was with{dis}advantage, show the other roll
	if o.Ignored > 0 {
		vals = append(vals, fmt.Sprintf("%d", o.Ignored))
	}

	diceStr := "[" + strings.Join(vals,",") + "]"
	if o.Roll.NumDice == 1 && (o.Roll.WithAdvantage || o.Roll.WithDisadvantage) {
		diceStr = "{" + strings.Join(vals,",") + "}"
	}

	s := fmt.Sprintf("Roll %s: %s", o.Roll, diceStr)
	if o.Roll.Modifier != 0 {
		s += fmt.Sprintf("%+d", o.Roll.Modifier)
	}
	s += fmt.Sprintf(", total:%d", o.Total)

	if o.Roll.Reason != "" {
		s += fmt.Sprintf(" (%s)", o.Roll.Reason)
	}
	
	if o.Roll.Target > 0 {
		if o.Success {
			s += " - SUCCESS!"
		} else {
			s += " - failed"
		}
	}

	if o.Roll.NumDice == 1 && o.Roll.DiceSize == 20 {
		if o.CriticalHit {
			if o.Dice[0] == 19 {
				s += " (IMPROVED CRITICAL HIT !!!)"
			} else {
				s += " (CRITICAL HIT !!!)"
			}
		} else if o.CriticalMiss {
			s += " (critical miss :/ )"
		}
	}

	return s
}

func (r Roll)Do() Outcome {
	o := Outcome{
		Roll: r,
		Dice: make([]int, r.NumDice),
	}

	for i:=0; i<r.NumDice; i++ {
		o.Dice[i] = Intn(r.DiceSize) + 1 // returns a value in the range [0,ord)
	}

	if r.NumDice == 1 && (r.WithAdvantage || r.WithDisadvantage) {
		// 1dN With{Dis}Advantage: roll 2dN, take most{least} favorable
		roll1 := o.Dice[0]
		roll2 := Intn(r.DiceSize)+1

		if r.WithAdvantage {
			if roll1 > roll2 {
				o.Dice[0], o.Ignored = roll1, roll2
			} else {
				o.Dice[0], o.Ignored = roll2, roll1
			}
		} else {
			if roll1 < roll2 {
				o.Dice[0], o.Ignored = roll1, roll2
			} else {
				o.Dice[0], o.Ignored = roll2, roll1
			}
		}		

	}

	if r.NumDice == 1 && r.DiceSize == 20 {
		switch o.Dice[0] {
		case  1: o.CriticalMiss = true
		case 19: o.CriticalHit  = r.WithImprovedCritical  // 19 is a crit hit when this is set
		case 20: o.CriticalHit  = true
		}
	}

	for i:=0; i<r.NumDice; i++ {
		o.Total += o.Dice[i]
	}

	o.Total += r.Modifier

	// Critical hits always succeed; critical misses always miss
	if r.Target > 0 {
		switch {
		case o.CriticalHit:   o.Success = true
		case o.CriticalMiss:  o.Success = false
		default:              o.Success = (o.Total >= r.Target)
		}
	}
	
	return o
}


// See roll_test.go for syntax examples
func Parse(s string) Roll {
	words := strings.Fields(s)
	if len(words) == 0 {
		return Roll{Err:fmt.Errorf("basic parse failure 1")}
	}		

	bits := regexp.MustCompile(`^(\d*)d(\d+)([-+]\d+)?$`).FindStringSubmatch(words[0])
	// e.g. returns [4d6+3, 4, 6, +3] on success
	if len(bits) < 3 || len(bits) > 4 {
		return Roll{Err:fmt.Errorf("basic parse failure 2")}
	}

	// Can swallow errors now; the regexp succeeded, so we know the bits
	// are parseable
	n,_ := strconv.Atoi(bits[1])
	if n<=0 { n = 1 }
	if n>100 { return Roll{Err:fmt.Errorf("that's not how I roll: way too many dice")} }

	ord,_ := strconv.Atoi(bits[2])
	if ord<=0 { ord = 1 }
	if ord>100 { return Roll{Err:fmt.Errorf("that's not how I roll: that dice is way too big")} }

	modifier := 0
	if len(bits) == 4 {
		modifier,_ = strconv.Atoi(bits[3])
	}

	ret := Roll{input:s, NumDice:n, DiceSize:ord, Modifier:modifier}

	// Look for other stuff relating to the roll, or a reason ("for blah")
	word := ""
	words = words[1:]
	for len(words) > 0 {
		word,words = strings.ToLower(words[0]),words[1:]

		if word == "for" {
			// The rest of the input is a reason string
			ret.Reason = strings.Join(words, " ")
			break
		} else if word == "withadvantage" {
			ret.WithAdvantage = true
		} else if word == "withdisadvantage" {
			ret.WithDisadvantage = true
		} else if word == ">=" && len(words) > 0 {
			// This is a bit hacky ... try to handle ">= 17", instead of the expected ">=17"
			ret.Target,_ = strconv.Atoi(words[0])
			words = words[1:]
			
		} else if bits := regexp.MustCompile(`^>=(\d+)$`).FindStringSubmatch(word); len(bits) == 2 {
			ret.Target,_ = strconv.Atoi(bits[1])
		} else {
			return Roll{Err:fmt.Errorf("confusing floating word %q", word)}
		}
	}
	
	return ret
}

// Combine adds two rolls together, and returns a new roll. Currently
// this only works if the dice size is the same.
func (r1 Roll)Combine(r2 Roll) Roll {
	new := r1

	if new.DiceSize != r2.DiceSize {
		return  Roll{Err:fmt.Errorf("can't combine %s with %s; different dice sizes", r1, r2)}
	}

	new.NumDice += r2.NumDice
	new.Modifier += r2.Modifier

	return new
}
