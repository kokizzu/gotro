package Z

import (
	"bytes"
	"testing"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/stretchr/testify/assert"
)

func Test_Template(t *testing.T) {
	const autoReload = true
	const printDebug = true
	const fileName = `dummy.html`
	tc, err := ParseFile(autoReload, printDebug, fileName)
	if L.IsError(err, `filed Z.ParseFile: `+fileName) {
		t.Fail()
	}
	buff := bytes.Buffer{}
	tc.Render(&buff, M.SX{
		`title`:   `this is a title`,
		`aString`: `this is a string`,
		`anArray`: A.X{1, `b`, `c`, 4},
		`aMap`: M.SX{
			`a`: 1,
			`b`: `test`,
			`c`: `something`,
		},
	})
	assert.Equal(t, `<html>
<head>
	<title>this is a title</title>
</head>
<body>
<script>
  const a = 'this is a string';
  const b = [1,"b","c",4];
  const c = {"a":1,"b":"test","c":"something"};
</script>
</body>
</html>
`, buff.String())
}

func Test_TemplateString(t *testing.T) {
	const printDebug = true
	res := bytes.Buffer{}
	const template = `hi my name #{name}, my age #{age}`
	FromString(template, printDebug).Render(&res, M.SX{
		`name`: `Tzuyu`,
		`age`:  21,
	})
	assert.Equal(t, `hi my name Tzuyu, my age 21`, res.String())
	str := FromString(`i like #{char}`).Str(M.SX{
		`char`: `Rem`,
	})
	assert.Equal(t, `i like Rem`, str)
}
