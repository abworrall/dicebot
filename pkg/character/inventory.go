package character

import(
	"fmt"
	"strconv"
	"time"
	"log"
)

type Inventory struct {
	Items []Item
}

type Item struct {
	Time time.Time
	Description string
}

func (a Item)String() string {
	return a.Description
}

func (i *Inventory)String() string {
	str := ""
	for j,item := range i.Items {
		str += fmt.Sprintf("[%02d] %s\n", j+1, item)
	}
	return str
}

func (i *Inventory)MaybeInit() {
	if i.Items == nil {
		log.Printf("MaybeInit\n")
		i.Items = []Item{}
	}
}

func (i *Inventory)Clear() {
	i.Items = []Item{}
}

func (i *Inventory)Append(desc string) {
	i.MaybeInit()

	i.Items = append(i.Items, Item{time.Now(), desc} )
	log.Printf("Yes, appended %s\n%s", desc, i)
}

// Arg is zero-indexed
func (i *Inventory)Remove(n int) {
	i.MaybeInit()
	if n<0 || n>=len(i.Items) { return }
	i.Items = append(i.Items[:n], i.Items[n+1:]...)
}


func (i *Inventory)ParseIndex(args []string) (int, string) {
	i.MaybeInit()

	if len(args) == 0 {
		return -1, "wat"
	} else if n,err := strconv.Atoi(args[0]); err != nil {
		return -1, fmt.Sprintf("'%s' is such nonsense", args[1])
	} else if n<=0 {
		return -1, fmt.Sprintf("yes yes, very clever")
	} else if n > len(i.Items) {
		return -1, fmt.Sprintf("you don't even have %d items, let alone %d", len(i.Items)+1, n)
	} else {
		return n-1, "" // callers will want a slice index
	}
}
