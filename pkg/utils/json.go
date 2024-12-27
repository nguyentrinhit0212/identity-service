package utils

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var JSON = jsoniter.Config{
	EscapeHTML:                    true,
	MarshalFloatWith6Digits:       true,
	TagKey:                        "json",
	CaseSensitive:                 false,
	ObjectFieldMustBeSimpleString: false,
}.Froze()

func init() {
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)
}
