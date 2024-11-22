package utils

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var JSON = jsoniter.Config{
	EscapeHTML:                    true,
	MarshalFloatWith6Digits:       true,
	TagKey:                        "json", // Sử dụng tag "json" để map tên trường
	CaseSensitive:                 false,  // Không phân biệt hoa thường
	ObjectFieldMustBeSimpleString: false,
}.Froze()

func init() {
	// Tự động chuyển PascalCase hoặc camelCase sang snake_case
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)
}
