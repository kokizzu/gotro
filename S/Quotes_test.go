package S

import (
	"strings"
	"testing"
)

func TestBasicQuotes(t *testing.T) {
	if got := Q(`abc`); got != `'abc'` {
		t.Fatalf("Q mismatch: %q", got)
	}
	if got := QQ(`abc`); got != `"abc"` {
		t.Fatalf("QQ mismatch: %q", got)
	}
	if got := BT(`abc`); got != "`abc`" {
		t.Fatalf("BT mismatch: %q", got)
	}
	if got := ZZ(`  a"b  `); got != `"a&quot;b"` {
		t.Fatalf("ZZ mismatch: %q", got)
	}
	if got := ZB(true); got != `'true'` {
		t.Fatalf("ZB mismatch: %q", got)
	}
	if got := ZI(42); got != `'42'` {
		t.Fatalf("ZI mismatch: %q", got)
	}
	if got := ZU(42); got != `'42'` {
		t.Fatalf("ZU mismatch: %q", got)
	}
}

func TestJsonQuotes(t *testing.T) {
	src := "a\"b'c\\d\r\n"
	if got := ZJJ(src); got != `"a\"b'c\\d\r\n"` {
		t.Fatalf("ZJJ mismatch: %q", got)
	}
	if got := ZJ(src); got != `'a"b\'c\\d\r\n'` {
		t.Fatalf("ZJ mismatch: %q", got)
	}
}

func TestSanitizeHelpers(t *testing.T) {
	src := " <a>'\"\\% "
	if got := Z(src); got != `'&lt;a&gt;&apos;&quot;\\%'` {
		t.Fatalf("Z mismatch: %q", got)
	}
	if got := ZS(src); got != `' &lt;a&gt;&apos;&quot;\\% '` {
		t.Fatalf("ZS mismatch: %q", got)
	}
	if got := ZLIKE(src); got != `'%&lt;a&gt;&apos;&quot;\\\%%'` {
		t.Fatalf("ZLIKE mismatch: %q", got)
	}
	if got := ZJLIKE(src); got != `'%&lt;a&gt;&apos;"\\\%%'` {
		t.Fatalf("ZJLIKE mismatch: %q", got)
	}
	if got := XSS(src); got != `&lt;a&gt;&apos;&quot;\%` {
		t.Fatalf("XSS mismatch: %q", got)
	}
}

func TestUnescapeHelpersAndTrace(t *testing.T) {
	escaped := `&lt;&gt;&amp;&apos;&quot;`
	if got := UZ(escaped); got != "<>&\u2018\u02ba" {
		t.Fatalf("UZ mismatch: %q", got)
	}
	if got := UZRAW(escaped); got != `<>&'"` {
		t.Fatalf("UZRAW mismatch: %q", got)
	}

	zt := ZT("hello", "world")
	if !strings.HasPrefix(zt, "-- ") || !strings.Contains(zt, "hello|world") {
		t.Fatalf("ZT mismatch: %q", zt)
	}
	if zt2 := ZT2(); !strings.HasPrefix(zt2, "-- ") || !strings.HasSuffix(zt2, "\n") {
		t.Fatalf("ZT2 mismatch: %q", zt2)
	}
}
