package mutable_string

import (
	"testing"
)

func assertEmpty(t *testing.T, overlays []overlay) {
	if len(overlays) != 0 {
		t.Errorf("Expected overlays to be empty, but got %v", overlays)
	}
}

func TestCommit(t *testing.T) {
	t.Run("no_overlays", func(t *testing.T) {
		ms := MutableString{
			initialText: "hello world",
			overlays:    []overlay{},
		}
		err := ms.Commit()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		expectedFinalStr := "hello world"
		if ms.initialText != expectedFinalStr {
			t.Errorf("Expected final string %q, but got %q", expectedFinalStr, ms.initialText)
		}
		assertEmpty(t, ms.overlays)
	})

	t.Run("single_overlay", func(t *testing.T) {
		ms := MutableString{
			initialText: "hello world",
			overlays: []overlay{
				{span: Range{Pos: 0, End: 5}, text: "hi"},
			},
		}
		err := ms.Commit()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		expectedFinalStr := "hi world"
		if ms.initialText != expectedFinalStr {
			t.Errorf("Expected final string %q, but got %q", expectedFinalStr, ms.initialText)
		}
		assertEmpty(t, ms.overlays)
	})

	t.Run("multiple_overlapping_overlays", func(t *testing.T) {
		ms := MutableString{
			initialText: "hello world",
			overlays: []overlay{
				{span: Range{Pos: 0, End: 5}, text: "hi"},
				{span: Range{Pos: 3, End: 8}, text: "planet"},
			},
		}
		err := ms.Commit()
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
	})

	t.Run("multiple_non_overlapping_overlays", func(t *testing.T) {
		ms := MutableString{
			initialText: "hello world",
			overlays: []overlay{
				{span: Range{Pos: 0, End: 5}, text: "hi"},
				{span: Range{Pos: 6, End: 11}, text: "planet"},
			},
		}
		err := ms.Commit()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		expectedFinalStr := "hi planet"
		if ms.initialText != expectedFinalStr {
			t.Errorf("Expected final string %q, but got %q", expectedFinalStr, ms.initialText)
		}
		assertEmpty(t, ms.overlays)
	})

	t.Run("abutting_non_overlapping_overlays", func(t *testing.T) {
		ms := MutableString{
			initialText: "hello world",
			overlays: []overlay{
				{span: Range{Pos: 0, End: 5}, text: "hi"},
				{span: Range{Pos: 5, End: 11}, text: " planet"},
			},
		}
		err := ms.Commit()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		expectedFinalStr := "hi planet"
		if ms.initialText != expectedFinalStr {
			t.Errorf("Expected final string %q, but got %q", expectedFinalStr, ms.initialText)
		}
		assertEmpty(t, ms.overlays)
	})

	t.Run("three_overlays_middle_empty", func(t *testing.T) {
		ms := MutableString{
			initialText: "hello world",
			overlays: []overlay{
				{span: Range{Pos: 0, End: 5}, text: "hi"},
				{span: Range{Pos: 5, End: 6}, text: ""},
				{span: Range{Pos: 6, End: 11}, text: "planet"},
			},
		}
		err := ms.Commit()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		expectedFinalStr := "hiplanet"
		if ms.initialText != expectedFinalStr {
			t.Errorf("Expected final string %q, but got %q", expectedFinalStr, ms.initialText)
		}
		assertEmpty(t, ms.overlays)
	})

	t.Run("insert_string_into_empty", func(t *testing.T) {
		ms := MutableString{
			initialText: "",
			overlays: []overlay{
				{span: Range{Pos: 0, End: 0}, text: "hello"},
			},
		}
		err := ms.Commit()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		expectedFinalStr := "hello"
		if ms.initialText != expectedFinalStr {
			t.Errorf("Expected final string %q, but got %q", expectedFinalStr, ms.initialText)
		}
		assertEmpty(t, ms.overlays)
	})

}

func TestInsert(t *testing.T) {
	ms := MutableString{initialText: " world"}
	err := ms.Insert(0, "hello")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = ms.Commit()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := "hello world"
	if ms.initialText != expected {
		t.Errorf("Expected %q, got %q", expected, ms.initialText)
	}
}

func TestDeleteRange(t *testing.T) {
	ms := MutableString{initialText: "hello world"}
	err := ms.DeleteRange(Range{Pos: 5, End: 11})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = ms.Commit()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := "hello"
	if ms.initialText != expected {
		t.Errorf("Expected %q, got %q", expected, ms.initialText)
	}
}

func TestInvalidDeleteRange(t *testing.T) {
	ms := MutableString{initialText: "hello world"}
	err := ms.DeleteRange(Range{Pos: 5, End: 3})
	if err == nil {
		t.Fatalf("Expected error, but got nil")
	}
	expected := "invalid or out of bounds range"
	if err.Error() != expected {
		t.Errorf("Expected error %q, got %q", expected, err.Error())
	}
}
func TestIsValidRange(t *testing.T) {
	testCases := []struct {
		name             string
		r                Range
		initialText      string
		expectedValidity bool
	}{
		{
			name:             "valid_range",
			r:                Range{Pos: 0, End: 5},
			initialText:      "hello world",
			expectedValidity: true,
		},
		{
			name:             "full_range",
			r:                Range{Pos: 0, End: 11},
			initialText:      "hello world",
			expectedValidity: true,
		},
		{
			name:             "empty_range",
			r:                Range{Pos: 5, End: 5},
			initialText:      "hello world",
			expectedValidity: true,
		},
		{
			name:             "negative_start",
			r:                Range{Pos: -2, End: 5},
			initialText:      "hello world",
			expectedValidity: false,
		},
		{
			name:             "out_of_bounds_end",
			r:                Range{Pos: 5, End: 15},
			initialText:      "hello world",
			expectedValidity: false,
		},
		{
			name:             "start_greater_than_end",
			r:                Range{Pos: 6, End: 4},
			initialText:      "hello world",
			expectedValidity: false,
		},
		{
			name:             "empty_text_valid_range",
			r:                Range{Pos: 0, End: 0},
			initialText:      "",
			expectedValidity: true,
		},
		{
			name:             "empty_text_out_of_bounds",
			r:                Range{Pos: 1, End: 2},
			initialText:      "",
			expectedValidity: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidRange(tc.r, len(tc.initialText))
			if result != tc.expectedValidity {
				t.Errorf("Expected validity %v, got %v", tc.expectedValidity, result)
			}
		})
	}
}
