package rules

// This is a simple interface, to allow rules consumers a simple way
// to lookup and display entries fomr any flavor of rule (spell,
// monster, etc)

type Entryer interface {
	Type() string
	Summary() string
	Description() string
}

type Lookuper interface {
	Lookup(s string) []Entryer
	LookupFirst(s string) Entryer
}
