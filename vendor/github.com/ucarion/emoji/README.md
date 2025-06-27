# emoji

`emoji` is a Golang package that lets you get information about an emoji,
including whether a string is an emoji at all.

```go
e, ok := emoji.Lookup("a") // ok is false
e, ok := emoji.Lookup("😎") // ok is true, e.Name is "smiling face with sunglasses"
```

`emoji` relies on a dataset that's `go generate`'d from raw data
in [Unicode Technical Standard #51, "Unicode Emoji"](http://unicode.org/reports/tr51/)
. This package is updated to the Emoji standard version 13.1.

## Installation

Add `emoji` as a dependency by running:

```shell
go get github.com/ucarion/emoji
```

## Usage Guide

Text in general is a tricky topic, and emojis are an especially tricky part of
modern text. Here are some recommendations on what you should or should not do
when trying to handle emojis in Go:

* **Iterating over a string will sometimes break apart emojis in Go.** When you
  do a `for` loop over a Go string, you will get each Unicode codepoint (Go
  calls a codepoint a `rune`) in the string. Some (but not all) emojis consist
  of multiple codepoints. For example, this emoji:

  > 🧑🏻‍🤝‍🧑🏿

  Is formally called "people holding hands: light skin tone, dark skin tone".
  It's implemented as a sequence of codepoints, which is why this code:

  ```go
  s := "🧑🏻‍🤝‍🧑🏿"
  for _, r := range s {
    fmt.Printf("%U %s\n", r, string(r))
  }
  ```

  Outputs:

  ```text
  U+1F9D1 🧑
  U+1F3FB 🏻
  U+200D ‍
  U+1F91D 🤝
  U+200D ‍
  U+1F9D1 🧑
  U+1F3FF 🏿
  ```

  This may seem like a strange decision by the folks at Unicode who assign
  emojis, but implementing emojis like this has a nice benefit: if a platform
  doesn't know about the skin tone variants of the "people holding hands" emoji,
  then it can fall back to showing the "sub-emojis", which are separated
  by `U+200D`, (called the "zero-width joiner"):

  > 🧑🏻🤝🧑🏿

  Related to this is another important note:

* **Not all emojis have the same `len`.** As of Emoji 13.1, some emojis are just
  one codepoint long, such as:

  > 😀 (1F600)

  While others can be as long as 10 codepoints long:

  > 👩🏻‍❤️‍💋‍👩🏻 (1F469 1F3FB 200D 2764 FE0F 200D 1F48B 200D 1F469 1F3FB)

  There's nothing stopping emojis from getting longer. For instance, Unicode
  could start adding support for both skin tone and hair color modifiers to all
  existing emojis, which could double, triple, or quadruple the number of
  possible codepoints for emojis that have two, three, or four people in them.

  In Go, `len(s)` tells you how many UTF-8 bytes there are in `s`,
  and `len([]rune(s))` tells you how many codepoints are in `s`. Since emojis
  contain a variable number of codepoints, neither of these have a predictable
  value for emojis.

* **Some emojis are prefixes/suffixes of other emojis.** Returning to the first
  example:

  > 🧑🏻‍🤝‍🧑🏿 (1F9D1 1F3FB 200D 1F91D 200D 1F9D1 1F3FF)

  Because of how emojis like this one are encoded in terms of sub-emojis, it is
  a suffix of some of its sub-emojis:

  > 🧑🏻 ("person: light skin tone", 1F9D1 1F3FB)
  >
  > 🧑 ("person", 1F9D1)

  Between this point and the two previous points, an important fact emerges: *it
  is not straightforward to "extract" the emojis from a string*:

    * You can't just iterate over the `byte` or `rune` contents of a `string`,
      because that would split multi-`rune` emojis.
    * You can't try to look at the `len` of the string in any meaningful way,
      because emojis are variable-width, and the range of possible lengths is
      changing as new versions of Emoji come out.
    * You can't try to look at whether a given *prefix* of a string is an emoji,
      because that will greedily miss longer emojis that are suffixes of shorter
      ones. If you try to look for the *longest* prefix of a string that's an
      emoji, that will not work as new emojis are introduced to the standard.

  If you really need to solve this problem, the technical term for what you need
  is "text segmentation", and in particular, segmenting text into "extended
  grapheme clusters". The standard for this is in a document
  called [UAX #29](https://unicode.org/reports/tr29/), so look for something
  that talks about implementing that.

  Once you have your string segmented into extended grapheme clusters, you can
  then pass each extended grapheme cluster to the `Lookup` function from this
  package, and then do whatever process you like from there. Every emoji,
  including emojis that might be added in the future, forms a single extended
  grapheme cluster, even if it consists of many codepoints.

* **Not all platforms support the same emojis.** For instance, the 10-codepoint
  long emoji in a previous example (👩🏻‍❤️‍💋‍👩🏻) is called "kiss: woman,
  woman, light skin tone", and was added in Emoji 13.1. It is supported in iOS
  14.5, but not macOS 11, its contemporary.

* **Platforms don't always fully support past versions of Emoji.** For example,
  macOS 11 supports 🦤 ("dodo"), which was added in Emoji 13.0, but it does not
  directly support 👩🏻‍🤝‍👨🏼 ("woman and man holding hands: light skin tone,
  medium-light skin tone"), which was added in Emoji 12.0.

  As a result of the previous two points, you should avoid trying to assume that
  an emoji will be drawn in a particular way. This is mostly a fool's errand,
  because emoji support is all over the place. There is no easy way to tell if
  an emoji is supported on a given platform.

  However, if you want to check what version of the Emoji standard an emoji was
  added, you can check its `Introduced` property:

  ```go
  _, e := emoji.Lookup("🦤")
  e.Introduced // 13.0
  
  _, e := emoji.Lookup("👩🏻‍🤝‍👨🏼")
  e.Introduced // 12.0
  ```

  This can be useful if, for example, you want to avoid emitting emojis that are
  definitely not supported by platforms that haven't been updated since the
  release of a particular version of the Emoji standard.

* **Some emojis aren't always displayed as a pictogram.** Some characters we
  consider to be emojis today were added to Unicode before emojis were
  introduced to Unicode. Whereas many new emojis, like 🦔 ("hedgehog", added in
  Emoji 5.0, aka Unicode 10.0), were intended to be presented as emojis from the
  day of their introduction, some older ones, like 🐿 ("chipmunk", added in
  Unicode 7.0, before emojis were added to Unicode), were *retroactively*
  classified as emojis.

  These "retroactive" emojis are said to lack what the Emoji specification
  calls "default emoji presentation". In these cases, it's up to the
  implementation to decide whether to present the character as a pictographic
  emoji, or whether to use some other behavior.

  Unicode has a special character, called U+FE0F "Variation Selector-16" ("
  VS16"), that lets you explicitly mark a character lacking default emoji
  presentation as being intended to be treated as an emoji. Marking a character
  that lacks default emoji presentation with VS16 makes it go from being an "
  unqualified" emoji to being "fully-qualified".

  For example, probably the most commonly-encountered example of a character
  that lacks default emoji presentation is U+263A, "White Smiling Face" (where
  the word "white" means "not filled in"):

  > ☺ (263A)

  It's up to the implementation to decide whether to display that as an emoji.
  Different tools will display that character differently. But if you add VS16:

  > ☺️ (263A FE0F)

  Then the emoji is unambiguously intended for emoji presentation. That said,
  even when an emoji uses VS16, many implementations will still display the
  emoji with a "text" presentation instead of an "emoji" presentation. As noted
  previously, emoji support is all over the place.

  With this `emoji` package, you can get whether an emoji is fully-qualified by
  checking its `Status`:

  ```go
  _, e := emoji.Lookup("☺") // the VS16-less version of the emoji
  e.Status // emoji.Unqualified
  
  _, e := emoji.Lookup("☺️") // the VS16'd version
  e.Status // emoji.FullyQualified
  ```

  (Because some emojis are encoded as a sequence of sub-emojis, there's also
  a `MinimallyQualified` status for emoji sequences where one of the sub-emojis
  is `Unqualified`.)

  All emojis except for those that have a `Status` of `Component` (which is for
  the special skin tone and hair color emoji modifiers) have a
  non-empty `FullyQualifiesAs` property, which gives you the fully-qualified
  rendition of any emoji. For example:

  ```go
  _, e := emoji.Lookup("☺") // the VS16-less version of the emoji
  e.FullyQualifiesAs // ☺️, the VS16'd version
  ```

## Contributing

To update this package to the latest version of the Emoji specification, do the
following:

* Remove the data file in the `data/`, and replace it with the `emoji-test.txt`
  file of the latest Emoji specification.
* Update `col1` and `col2` in `internal/cmd/genemoji/main.go` if needed.
* Run `go generate` and `go fmt emoji_data.go`.
* Update `Version` in `emoji.go` to the new appropriate value.

And that's all! Barring any significant changes to the Emoji data model in
future versions of Unicode, nothing else should require updating.
