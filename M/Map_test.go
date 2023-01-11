package M

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSX_FromJson(t *testing.T) {
	t.Run(`nil map`, func(t *testing.T) {
		var m SX = nil
		assert.True(t, m.FromJson(`{"a":123,"b":"abc"}`))
		if m.GetInt(`a`) != 123 || m[`b`] != `abc` {
			t.Error(`invalid value`)
		}
	})

	t.Run(`default case`, func(t *testing.T) {
		m := SX{}
		assert.True(t, m.FromJson(`{"a":123,"b":"abc"}`))
		if m.GetInt(`a`) != 123 || m[`b`] != `abc` {
			t.Error(`invalid value`)
		}
	})

	t.Run(`invalid json`, func(t *testing.T) {
		m := SX{}
		assert.False(t, m.FromJson(`{"a":123,"b":"abc"`))
		if len(m) != 2 { // goccy/go-json will still parse the valid json part
			t.Error(`invalid value`)
		}
	})

	t.Run(`overwrites`, func(t *testing.T) {
		m := SX{}
		m["a"] = 234
		assert.True(t, m.FromJson(`{"a":123,"b":"abc"}`))
		if m.GetInt(`a`) != 123 || m[`b`] != `abc` {
			t.Error(`invalid value`)
		}
	})

	t.Run(`not overwrite`, func(t *testing.T) {
		m := SX{}
		m["a"] = 234
		assert.True(t, m.FromJson(`{"b":"abc"}`))
		if m.GetInt(`a`) != 234 || m[`b`] != `abc` {
			t.Error(`invalid value`)
		}
	})

	t.Run(`empty string`, func(t *testing.T) {
		var m SX
		assert.False(t, m.FromJson(``))
		if len(m) != 0 {
			t.Error(`invalid value`)
		}
	})

	t.Run(`inside struct`, func(t *testing.T) {
		x := struct {
			Foo SX
		}{}
		assert.True(t, x.Foo.FromJson(`{"a":123,"b":"abc"}`))
		if x.Foo.GetInt(`a`) != 123 || x.Foo[`b`] != `abc` {
			t.Error(`invalid value`)
		}
	})
}
