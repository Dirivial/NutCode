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
