package mutable_string

import (
	"fmt"
	"sort"
	"strings"
)

type MutableString struct {
	initialText string
	overlays    []overlay
}

func NewMutableString(text string) MutableString {
	return MutableString{
		initialText: text,
		overlays:    []overlay{},
	}
}

/*
Commit applies all accumulated overlays to the initial text, resulting in the final
transformed string. The overlays represent string modifications such as insertions, deletions, and replacements.
Each overlay specifies a range of characters in the initial text to be replaced and the replacement text.

After the Commit operation, the initial text is updated to the final transformed string, and the list of overlays
is cleared, allowing for new overlays to be added for future transformations.

The Commit method ensures that no overlapping overlays are applied; if it detects any overlaps, it returns an
error. Overlapping overlays are detected when the starting position of an overlay is earlier than the ending
position of the previous overlay.

Example usage:

	ms := NewMutableString("hello world")
	ms.ReplaceRange(Range{Pos: 0, End: 5}, "hi") // Replace "hello" with "hi"
	ms.Insert(5, " there")                       // Insert " there" in between "hello" and " world"
	ms.Append("!")                               // Insert "!" at the end.
	err := ms.Commit()                           // Apply the overlays
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ms.initialText) // Output: "hi there world!"
	}

The Commit method is useful when performing batch string transformations, as it allows for multiple modifications
to be applied to the initial text in a single operation, reducing the number of intermediate allocations and copies.
*/
func (ms *MutableString) Commit() error {
	sort.Sort(byStartIndex(ms.overlays))
	rawText := []rune(ms.initialText)

	builder := strings.Builder{}
	rawCursor := 0
	// combinedCursor := 0

	var prevOverlay *overlay

	for i, overlay := range ms.overlays {
		prevOverlay = nil
		if i > 0 {
			prevOverlay = &ms.overlays[i-1]
			if overlay.span.Pos < prevOverlay.span.End {
				return fmt.Errorf("ranges overlap")
			}
		}

		rawToCopy := overlay.span.Pos - rawCursor
		str := string(rawText[rawCursor : rawCursor+rawToCopy])
		builder.Write([]byte(str))
		rawCursor += rawToCopy + overlay.span.End - overlay.span.Pos
		// combinedCursor += rawToCopy

		builder.WriteString(overlay.text)
		// combinedCursor += len(overlay.text)
	}

	if len(rawText) > rawCursor {
		str := string(rawText[rawCursor:])
		builder.Write([]byte(str))
	}

	ms.initialText = builder.String()
	ms.overlays = make([]overlay, 0, INITIAL_CAPACITY)

	return nil
}

// ReplaceRange adds an overlay to replace the specified range with the provided text.
func (ms *MutableString) ReplaceRange(r Range, text string) error {
	if !isValidRange(r, len(ms.initialText)) {
		return fmt.Errorf("invalid or out of bounds range")
	}
	ms.overlays = append(ms.overlays, overlay{span: r, text: text})
	return nil
}

// Append adds an overlay to insert the provided text at the end of the initial text.
func (ms *MutableString) Append(text string) error {
	r := Range{Pos: len(ms.initialText), End: len(ms.initialText)}
	ms.overlays = append(ms.overlays, overlay{span: r, text: text})
	return nil
}

// Insert adds an overlay to insert the provided text at the specified position.
func (ms *MutableString) Insert(pos int, text string) error {
	r := Range{Pos: pos, End: pos}
	if !isValidRange(r, len(ms.initialText)) {
		return fmt.Errorf("invalid or out of bounds range")
	}
	ms.overlays = append(ms.overlays, overlay{span: r, text: text})
	return nil
}

/*
DeleteRange removes a specified range of characters from the MutableString by
replacing it with an empty string. The method takes a Range struct as an argument,
which defines the starting position (inclusive) and ending position (exclusive) of
the range to be deleted.

If the specified range is out of bounds or invalid (e.g., Pos is greater than End),
the behavior is undefined, and the method may return an error during the Commit operation.

Example usage:

	ms := NewMutableString("hello world")
	ms.DeleteRange(Range{Pos: 6, End: 11}) // Delete "world"
	err := ms.Commit()
	if err != nil {
	    fmt.Println(err)
	} else {
	    fmt.Println(ms.initialText) // Output: "hello "
	}
*/
func (ms *MutableString) DeleteRange(r Range) error {
	// Replace the specified range with an empty string to perform deletion.
	if !isValidRange(r, len(ms.initialText)) {
		return fmt.Errorf("invalid or out of bounds range")
	}
	ms.ReplaceRange(r, "")
	return nil
}

// isValidRange checks if the provided range is valid based on the length of the initial text.
// It returns true if the range is valid and false if the range is invalid or out of bounds.
func isValidRange(r Range, initialTextLength int) bool {
	if r.Pos > r.End {
		return false
	}
	if r.Pos < 0 || r.End > initialTextLength {
		return false
	}
	return true
}
