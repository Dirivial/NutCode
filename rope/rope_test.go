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
		{input: "cool_", at: 13, expected: "hello_I_am_a_cool_rope_data_structure"},
		{input: "_cool", at: len(testInput), expected: "hello_I_am_a_rope_data_structure_cool"},
		{input: "_cool_", at: 1, expected: "h_cool_ello_I_am_a_rope_data_structure"},
		{input: "cool_", at: 0, expected: "cool_hello_I_am_a_rope_data_structure"},
	}
	for _, v := range testInserts {

		rope := New(testInput)
		rope = rope.Insert(v.at, v.input)

		// rope.printRope()
		content := rope.GetContent()
		if content != v.expected {
			t.Fatalf("Content mismatch. Expected=%s, got=%s", v.expected, content)
		}
	}
}

func TestRopeDelete(t *testing.T) {
	testInput := "hello_I_am_a_rope_data_structure"
	testRemoves := []struct {
		start    int
		length   int
		expected string
	}{
		{14, 1, "hello_I_am_a_rpe_data_structure"},
		{13, 5, "hello_I_am_a_data_structure"},
		{len(testInput) - 1, 1, "hello_I_am_a_rope_data_structur"},
		{1, 2, "hlo_I_am_a_rope_data_structure"},
		{0, 1, "ello_I_am_a_rope_data_structure"},
	}
	for _, v := range testRemoves {

		rope := New(testInput)
		rope = rope.Delete(v.start, v.length)

		// rope.printRope()
		content := rope.GetContent()
		if content != v.expected {
			rope.printRope()
			t.Fatalf("Content mismatch. Expected=%s, got=%s", v.expected, content)
		}
	}
}

func TestRopeReport(t *testing.T) {
	testInput := "hello_I_am_a_rope_data_structure"
	testReports := []struct {
		start    int
		length   int
		expected string
	}{
		{14, 1, "r"},
		{13, 5, "_rope"},
		{len(testInput), 1, "e"},
		{1, 2, "he"},
		{1, 1, "h"},
	}
	for _, v := range testReports {

		rope := New(testInput)
		content := rope.Report(v.start, v.length)

		if content != v.expected {
			rope.printRope()
			t.Fatalf("Content mismatch. Expected=%s, got=%s", v.expected, content)
		}
	}
}

func TestRopeSearch(t *testing.T) {
	testInput := "Ahello_I_am_Aa_rope_AdaAAta_structurezA"
	testSearch := []struct {
		start     int
		expected  int
		character rune
	}{
		{1, 1, 'A'},
		{2, 10, 'a'},
		{2, 13, 'A'},
		{11, 14, 'a'},
		{8, 13, 'A'},
		{13, 13, 'A'},
		{14, 21, 'A'},
		{22, 24, 'A'},
		{24, 24, 'A'},
		{25, 25, 'A'},
		{26, 39, 'A'},
		{len(testInput) + 1, -1, 'A'},
		{-1, -1, 'A'},
	}
	for _, v := range testSearch {

		//fmt.Printf("\n ================ %c -- %d ================ \n", v.character, v.start)
		rope := New(testInput)
		content := rope.SearchChar(v.character, v.start)

		if content != v.expected {
			rope.printRope()
			t.Fatalf("Result mismatch. Expected=%d, got=%d, start=%d, char=%c", v.expected, content, v.start, v.character)
		}
		//fmt.Println("")
	}
}
