package character

import(
	"fmt"
)

// A Character holds info about a typical RPG character
type Character struct {
	Name string
	Str,Int,Wis,Con,Cha,Dex,Per int
}

func (c Character)String() string {
	s := fmt.Sprintf("%-10.10s [str:%2d int:%2d wis:%2d con:%2d cha:%2d dex:%2d per:%2d]",
		c.Name, c.Str, c.Int, c.Wis, c.Con, c.Cha, c.Dex, c.Per)
	return s
}
