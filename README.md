Package `mutable_string` provides a utility for manipulating strings in a flexible and efficient manner.
The package allows users to apply multiple string modifications (such as insertions, deletions, and replacements)
to an initial string while deferring the actual application of these modifications.

MutableString is the core data structure of this package, representing a string that can undergo multiple
transformations. The transformations are represented as overlays, where each overlay specifies a range of
characters to be replaced and the replacement text. The MutableString struct keeps track of the initial text
and a list of overlays. The text modifications are not immediately applied to the initial text; instead,
they are stored in the overlays until the Commit method is called. During the Commit operation, all accumulated
overlays are applied to the initial text in the order they were added, resulting in the final transformed string.

The package provides a set of methods for creating and manipulating MutableString instances, including:

- ReplaceRange: Adds an overlay that specifies a range of characters in the initial text to be replaced
  with the provided replacement text. The range is defined by a starting position (inclusive) and an ending
  position (exclusive). If the replacement text is empty, this operation effectively performs a deletion.

- Insert: Adds an overlay that inserts the provided text at the end of the initial text. This operation
  extends the length of the initial text.

- Commit: Applies all accumulated overlays to the initial text. After the Commit
  operation, the initial text is updated to the final transformed string, and the list of overlays is cleared.
  The Commit method ensures that no overlapping overlays are applied; if it detects any overlaps, it returns
  an error.

MutableString is designed to handle batch string transformations efficiently. By deferring the actual application
of modifications, the package reduces the number of intermediate string allocations and copies that would be
needed if each modification were applied immediately. This is particularly useful when dealing with large strings
and a sequence of complex transformations.

Usage:

```
ms := NewMutableString("hello world")
ms.ReplaceRange(Range{Pos: 0, End: 5}, "hi") // Replace "hello" with "hi"
ms.Insert(5, " there")                       // Insert " there" in between "hello" and " world"
ms.Append("!")                               // Insert "!" at the end.
res, err := ms.Commit()                      // Apply the overlays
if err != nil {
	fmt.Println(err)
} else {
	fmt.Println(res) // Output: "hi there world!"
}
```

The package is intended for use cases where batch string manipulation is needed, such as text editors, document
processing systems, and code generation tools. It provides a convenient and memory-efficient way to perform
complex string transformations.
