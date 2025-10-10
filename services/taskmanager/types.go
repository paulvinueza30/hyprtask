package taskmanager

type Mode string

const (
	Hypr Mode = "hypr"
	All  Mode = "all"
)

var stringToMode = map[string]Mode{
	"all" : All,
	"hypr" : Hypr,
}

