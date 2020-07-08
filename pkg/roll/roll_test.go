package roll

import(
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		s        string
		expected Roll
	}{
		{ "d6",                        Roll{NumDice:1, DiceSize:6, Modifier:0} },
		{"3d4+3",                      Roll{NumDice:3, DiceSize:4, Modifier:3} },
		{"2d8+1 >=7",                  Roll{NumDice:2, DiceSize:8, Modifier:1, Target:7} },
		{ "d20+4 >= 16 for fools",     Roll{NumDice:1, DiceSize:20, Modifier:4, Target:16, Reason:"fools"} },
		{"1d20 >=6 withadvantage",     Roll{NumDice:1, DiceSize:20, Target:6, WithAdvantage:true} },
		{"1d20 >=16 withdisadvantage", Roll{NumDice:1, DiceSize:20, Target:16, WithDisadvantage:true} },
		{"1d20 >=4 for a good reason", Roll{NumDice:1, DiceSize:20, Target:4, Reason:"a good reason"} },
	}

	for i,test := range tests {
		actual := Parse(test.s)
		fmt.Printf("Parse: %s\n", actual)
		actual.input = "" // Simplify our test comparisons
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("[T%d]\nwanted %#v\ngot    %#v", i, test.expected, actual)
		}
	}
}

func TestRoll(t *testing.T) {
	tests := []struct {
		seed     int64    // make the rng determinisitc, per test; -1 means don't reset the seed
		r        Roll
		expected Outcome
	}{
		{0,  Roll{NumDice:1, DiceSize: 6},             Outcome{Dice:[]int{1},       Total: 1} },
		{-1, Roll{NumDice:2, DiceSize: 8},             Outcome{Dice:[]int{3,2},     Total: 5} },
		{-1, Roll{NumDice:3, DiceSize: 4, Modifier:3}, Outcome{Dice:[]int{3,4,1},   Total:11} },
		{-1, Roll{NumDice:1, DiceSize:20, Target:9},   Outcome{Dice:[]int{8},       Total: 8,  Success:false} },
		{-1, Roll{NumDice:1, DiceSize:20, Target:18},  Outcome{Dice:[]int{18},      Total:18,  Success:true} },

		// with{dis}advantage means rolling 1d20 twice, and picking most{least} favorable
		{1, Roll{NumDice:1, DiceSize:20, WithAdvantage:true},    Outcome{Dice:[]int{8}, Ignored:2,  Total:8}},  // [2,8]
		{3, Roll{NumDice:1, DiceSize:20, WithDisadvantage:true}, Outcome{Dice:[]int{9}, Ignored:18, Total:9}}, // [9,18]

		// critical hits, misses, etc. The seeds were picked to generate the desired rolls.
		{11,  Roll{NumDice:1, DiceSize:20, Target:1},  Outcome{Dice:[]int{1},  Total:1,  CriticalMiss:true, Success:false} },
		{103, Roll{NumDice:1, DiceSize:20, Target:1},  Outcome{Dice:[]int{20}, Total:20, CriticalHit:true,  Success:true} },
		{30,
			Roll{NumDice:1, DiceSize:20, Target:100, WithImprovedCritical:true},
			Outcome{Dice:[]int{19}, CriticalHit:true, Success:true, Total:19},
		},
	}

	/*
	for i:=0; i<=10000; i++ {
		rand.Seed(int64(i))
		r := Roll{NumDice:1, DiceSize:20}
		o := r.Do()
		if o.Total == 20 {
			fmt.Printf(" %02d] %s\n", i ,o)
		}
    break
	}
  */

	for i,test := range tests {
		if test.seed >= 0 {
			rand.Seed(test.seed)
		}
		actual := test.r.Do()
		fmt.Printf("Roll: %s\n", actual)
		actual.Roll = Roll{} // Simplify the test comparisons

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("[T%d]\nwanted %#v\ngot    %#v", i, test.expected, actual)
		}
	}
}
