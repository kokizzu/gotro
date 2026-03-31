package I

import "testing"

func TestIfAndMinMaxHelpers(t *testing.T) {
	if If(true, 7) != 7 || If(false, 7) != 0 {
		t.Fatalf("If mismatch")
	}
	if IfElse(true, 7, 8) != 7 || IfElse(false, 7, 8) != 8 {
		t.Fatalf("IfElse mismatch")
	}
	if IfZero(0, 4) != 4 || IfZero(5, 4) != 5 {
		t.Fatalf("IfZero mismatch")
	}
	if IsZero(0, 4) != 4 || IsZero(5, 4) != 5 {
		t.Fatalf("IsZero mismatch")
	}
	if UIf(true, 7) != 7 || UIf(false, 7) != 0 {
		t.Fatalf("UIf mismatch")
	}
	if UIfElse(true, 7, 8) != 7 || UIfElse(false, 7, 8) != 8 {
		t.Fatalf("UIfElse mismatch")
	}
	if UIfZero(0, 4) != 4 || UIfZero(5, 4) != 5 {
		t.Fatalf("UIfZero mismatch")
	}
	if UIsZero(0, 4) != 4 || UIsZero(5, 4) != 5 {
		t.Fatalf("UIsZero mismatch")
	}

	if Min(2, 3) != 2 || Max(2, 3) != 3 {
		t.Fatalf("Min/Max mismatch")
	}
	if UMin(2, 3) != 2 || UMax(2, 3) != 3 {
		t.Fatalf("UMin/UMax mismatch")
	}
	if MinOf(2, 3) != 2 || MaxOf(2, 3) != 3 {
		t.Fatalf("MinOf/MaxOf mismatch")
	}
	if UMinOf(2, 3) != 2 || UMaxOf(2, 3) != 3 {
		t.Fatalf("UMinOf/UMaxOf mismatch")
	}
}

func TestFormattingAndOrdinalHelpers(t *testing.T) {
	if ToS(123) != `123` || ToStr(456) != `456` {
		t.Fatalf("ToS/ToStr mismatch")
	}
	if UToS(123) != `123` || UToStr(456) != `456` {
		t.Fatalf("UToS/UToStr mismatch")
	}
	if PadZero(123, 5) != `00123` {
		t.Fatalf("PadZero mismatch")
	}
	if PadZero(123, 2) != `123` {
		t.Fatalf("PadZero should not trim")
	}

	cases := map[int64]string{
		-1:  ``,
		1:   `1st`,
		2:   `2nd`,
		3:   `3rd`,
		4:   `4th`,
		11:  `11th`,
		12:  `12th`,
		13:  `13th`,
		21:  `21st`,
		112: `112th`,
	}
	for in, want := range cases {
		if got := ToEnglishNum(in); got != want {
			t.Fatalf("ToEnglishNum(%d) mismatch: want=%q got=%q", in, want, got)
		}
	}
}

func TestRoman(t *testing.T) {
	cases := map[int64]string{
		0:   ``,
		4:   `IV`,
		9:   `IX`,
		16:  `XVI`,
		944: `CMXLIV`,
	}
	for in, want := range cases {
		if got := Roman(in); got != want {
			t.Fatalf("Roman(%d) mismatch: want=%q got=%q", in, want, got)
		}
	}
}
