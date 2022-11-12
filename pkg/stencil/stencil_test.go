package stencil

import (
	"testing"
)

func TestStencilWrite(t *testing.T) {
	st := NewStencil("TEST", "")

	st.DrawText()
}

// func TestSingleCharacter(t *testing.T) {
// 	// Given
// 	input := "*-"

// 	// When
// 	got, err := NewStencil().AssignChars(input, "")
// 	if err != nil {
// 		t.Fatal("Failure")
// 	}

// 	// Then
// 	want := []string{"********", "********", "********"}
// 	AssertCharRow(t, got, want)
// }

// func TestSingleCharacterM(t *testing.T) {
// 	// Given
// 	input := "m"

// 	// When
// 	got, err := NewStencil().AssignChars(input, "")
// 	if err != nil {
// 		t.Fatal("Failure")
// 	}

// 	// Then
// 	want := []string{"    ", "|??|", "|  |"}
// 	AssertCharRow(t, got, want)
// }

// func AssertCharRow(t *testing.T, got []string, want []string) {
// 	t.Helper()

// 	for i, _ := range got {
// 		if len(got[i]) != len(want[i]) {
// 			t.Errorf("Lengths do not match on row: %d", i)
// 		}
// 	}
// }
