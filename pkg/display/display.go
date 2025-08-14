package display

import (
	"mp/sozzler/pkg/sozzler"
)

type Display interface {
	Error(string)
	List([]*sozzler.Recipe)
	Show(*sozzler.Recipe)
	String(string)
}
