// Package emoji provides a lookup function to get information about a given
// potential emoji.
//
// Information in this package is derived from Unicode Technical Standard #51
// ("Unicode Emoji"):
//
// http://unicode.org/reports/tr51/
package emoji

//go:generate go run ./internal/cmd/genemoji/... -data data/13.1 -out emoji_data.go

// Version is the edition of Unicode Technical Standard #51 ("Unicode Emoji")
// from which Lookup is derived.
const Version = "13.1"

// Emoji holds information about an emoji.
type Emoji struct {
	// The CLDR short name of the emoji.
	Name string

	// The emoji's qualification status.
	Status Status

	// The edition of Unicode Emoji when the emoji was introduced.
	Introduced string

	// The fully qualified representation of the emoji. Empty if the emoji does
	// not have any fully qualified representation.
	FullyQualifiesAs string
}

// Status is the qualification status of the emoji.
//
// The qualification status of an emoji sequence informs whether an
// implementation must process and display the sequence as an emoji.
// FullyQualified emojis must be processed as emojis; MinimallyQualified and
// Unqualified emojis may or may not be.
type Status int

const (
	// Component indicates the emoji is an emoji component.
	//
	// Emoji components are not intended for independent, direct output. They do
	// not have a fully-qualified representation.
	Component Status = iota

	// FullyQualified indicates the emoji is fully-qualified.
	//
	// Fully-qualified emojis are unambiguously intended for emoji presentation.
	// The Unicode Emoji standard recommends that user input devices only emit
	// fully-qualified emojis.
	FullyQualified

	// MinimallyQualified indicates that the emoji an emoji sequence where the
	// first character is qualified, but the full sequence is not.
	//
	// It is up to the implementation to choose whether to process and display
	// minimally-qualified emojis in the same way as their fully-qualified
	// forms.
	MinimallyQualified

	// Unqualified indicates that the emoji is neither fully-qualified nor
	// minimally-qualified.
	//
	// Most unqualified emojis are codepoints that were introduced prior to the
	// Emoji standard being created. They were retroactively categorized as
	// emojis.
	//
	// It is up to the implementation to choose whether to process and display
	// unqualified emojis in the same way as their fully-qualified forms.
	Unqualified
)

// Lookup finds information about a single emoji. If the emoji is found, its
// information is returned and the boolean is true. Otherwise the returned value
// will be empty and the boolean is false.
//
// Lookup looks up the inputted string in its entirety. As a result, it will not
// find any emoji if the input string consists of multiple emojis.
//
// Lookup only finds emojis that are recommended for general interchange
// ("RGI"), are a minimally-qualified or unqualified version of an RGI emoji, or
// which are emoji components requiring emoji presentation.
func Lookup(s string) (Emoji, bool) {
	e, ok := emojis[s]
	return e, ok
}
