package rope

import (
	"testing"
)

func TestRopeInit(t *testing.T) {
	cases := []string{"", "Test", "Longer input", "Even longer input. Is this going to be a problem?"}
	for _, i := range cases {
		rope := New(i)
		if rope == nil {
			t.Fail()
		}
	}
}

func TestRopeIndex(t *testing.T) {
	testInput := "Test1! & some other text"
	rope := New(testInput)

	for i, c := range testInput {
		res := rope.Index(i + 1)
		if res != string(c) {
			t.Fatalf("Wrong output '%s', got=%c", res, c)
		}
	}
}

func TestRopeConcat(t *testing.T) {
	testA := "Hello mr"
	testB := ", world!"
	testC := testA + testB

	ropeA := New(testA)
	ropeB := New(testB)
	rope := ropeA.Concat(ropeB)

	for i, c := range testC {
		res := rope.Index(i + 1)
		if res != string(c) {
			t.Fatalf("Wrong character '%s', expected='%c'", res, c)
		}
	}
	if testC != rope.GetContent() {
		t.Fatalf("Content mismatch. Expected=%s, got=%s", testC, rope.GetContent())
	}
}

func TestRopeSplit(t *testing.T) {
	testInput := "hello_I_am_a_rope_data_structure"

	rope := New(testInput)
	ropeHeadWeight := rope.Head.Weight
	secondRope := rope.Split(9)

	weightSum := rope.Head.Weight + secondRope.Head.Weight

	if weightSum != ropeHeadWeight {
		t.Fatalf("Weights of split tree does not add up to original number. Expected=%d, got=%d+%d(=%d)", ropeHeadWeight, rope.Head.Weight, secondRope.Head.Weight, weightSum)
	}

	appendedContent := rope.GetContent() + secondRope.GetContent()
	if appendedContent != testInput {
		t.Fatalf("Original input does not equal split content. Expected=%s, got=%s", testInput, appendedContent)
	}
}

func TestRopeInsert(t *testing.T) {
	testInput := "hello_I_am_a_rope_data_structure"
	testInserts := []struct {
		input    string
		at       int
		expected string
	}{
		{input: "cool_", at: 14, expected: "hello_I_am_a_cool_rope_data_structure"},
		{input: "_cool", at: len(testInput), expected: "hello_I_am_a_rope_data_structure_cool"},
		{input: "_cool_", at: 2, expected: "h_cool_ello_I_am_a_rope_data_structure"},
		{input: "cool_", at: 1, expected: "cool_hello_I_am_a_rope_data_structure"},
	}
	for _, v := range testInserts {

		rope := New(testInput)
		rope = rope.Insert(v.at, v.input)

		rope.printRope()
		content := rope.GetContent()
		if content != v.expected {
			t.Fatalf("Content mismatch. Expected=%s, got=%s", v.expected, content)
		}
	}
}
