package Z_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/Z"
)

const exampleZ = `hi my name #{name}, my age #{age}`

var exampleMap = M.SX{
	`name`: `Tzuyu`,
	`age`:  21,
}

const exampleGo = `hi my name {{.Name}}, my age {{.Age}}`

var exampleStruct = struct {
	Name string
	Age  int
}{
	Name: `Tzuyu`,
	Age:  21,
}

func expectExampleRendered(res string) {
	const expect = `hi my name Tzuyu, my age 21`
	if expect != res {
		panic(`value not equal: ` + res)
	}
}

func Test_Patterns(t *testing.T) {
	patterns := []string{
		`hi my name /*!name*/, my age /*!age*/`,
		`hi my name /*! name*/, my age /*! age*/`,
		`hi my name /*!name */, my age /*!age */`,
		`hi my name /*! name */, my age /*! age */`,
		`hi my name [/*name*/], my age [/*age*/]`,
		`hi my name [/* name*/], my age [/* age*/]`,
		`hi my name [/*name */], my age [/*age */]`,
		`hi my name [/* name */], my age [/* age */]`,
		`hi my name {/*name*/}, my age {/*age*/}`,
		`hi my name {/* name*/}, my age {/* age*/}`,
		`hi my name {/*name */}, my age {/*age */}`,
		`hi my name {/* name */}, my age {/* age */}`,
		`hi my name #{name}, my age #{age}`,
		`hi my name #{ name}, my age #{ age}`,
		`hi my name #{name }, my age #{age }`,
		`hi my name #{ name }, my age #{ age }`,
	}
	for _, pattern := range patterns {
		t.Run(pattern, func(t *testing.T) {
			res := Z.FromString(pattern, true).Str(exampleMap)
			expectExampleRendered(res)
		})
	}
}

func Test_Template(t *testing.T) {
	const autoReload = true
	const printDebug = true
	const fileName = `dummy.html`
	tc, err := Z.ParseFile(autoReload, printDebug, fileName)
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
	Z.FromString(exampleZ, printDebug).Render(&res, exampleMap)
	expectExampleRendered(res.String())
	str := Z.FromString(`i like #{char}`).Str(M.SX{
		`char`: `Rem`,
	})
	assert.Equal(t, `i like Rem`, str)
}

func BenchmarkParseRenderZTemplate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := Z.FromString(exampleZ, false).Str(exampleMap)
		expectExampleRendered(res)
	}
}

func BenchmarkParseRenderGoTemplate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := template.New("test")
		t, _ = t.Parse(exampleGo)
		res := &bytes.Buffer{}
		_ = t.Execute(res, exampleStruct)
		expectExampleRendered(res.String())
	}
}

func BenchmarkRenderZTemplate(b *testing.B) {
	tc := Z.FromString(exampleZ, false)
	res := &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		tc.Render(res, exampleMap)
		expectExampleRendered(res.String())
		res.Reset()
	}
}

func BenchmarkRenderGoTemplate(b *testing.B) {
	res := &bytes.Buffer{}
	t := template.New("test")
	t, _ = t.Parse(exampleGo)
	for i := 0; i < b.N; i++ {
		_ = t.Execute(res, exampleStruct)
		expectExampleRendered(res.String())
		res.Reset()
	}
}

func TestFalsePositive(t *testing.T) {
	dummy2 := L.ReadFile(`dummy2.txt`)
	dummy2expect := L.ReadFile(`dummy2.rendered`)
	rendered := Z.FromString(dummy2).Str(M.SX{})
	assert.Equal(t, dummy2expect, rendered)
}

func Test_TemplateHelpers(t *testing.T) {
	tc := Z.FromString(`hello #{name}`)
	buf := tc.ByteBuffer(M.SX{`name`: `A`})
	assert.Equal(t, `hello A`, buf.String())
	tc.Print() // coverage for debug print helper
}

func Test_PatternsWithSpacePrefix(t *testing.T) {
	patterns := []string{
		`hi { /* name */ } and [ /* age */ ]`,
	}
	for _, pattern := range patterns {
		t.Run(pattern, func(t *testing.T) {
			res := Z.FromString(pattern).Str(M.SX{`name`: `A`, `age`: 21})
			assert.Equal(t, `hi A and 21`, res)
		})
	}
}

func Test_RenderMissingAndUnusedParam(t *testing.T) {
	tc := Z.FromString(`x #{known} y #{missing}`, true)
	var buf bytes.Buffer
	tc.Render(&buf, M.SX{
		`known`: `K`,
		`extra`: `E`,
	})
	assert.Equal(t, `x K y missing`, buf.String())
}

func Test_ReloadBranches(t *testing.T) {
	t.Run(`in memory`, func(t *testing.T) {
		tc := Z.FromString(`x #{a}`)
		tc.InMemory = true
		nt, err := tc.Reload()
		assert.NoError(t, err)
		assert.Equal(t, tc, nt)
	})

	t.Run(`file not found`, func(t *testing.T) {
		tc := &Z.TemplateChain{Filename: filepath.Join(t.TempDir(), `missing.html`)}
		nt, err := tc.Reload()
		assert.Error(t, err)
		assert.True(t, len(nt.Parts) == 1)
		assert.True(t, strings.Contains(string(nt.Parts[0]), `failed to stat the template`))
	})

	t.Run(`unchanged and modified`, func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, `x.html`)
		assert.NoError(t, os.WriteFile(path, []byte(`hello #{name}`), 0644))

		tc, err := Z.ParseFile(false, false, path)
		assert.NoError(t, err)
		assert.Equal(t, `hello A`, tc.Str(M.SX{`name`: `A`}))

		nt1, err := tc.Reload()
		assert.NoError(t, err)
		assert.Equal(t, tc, nt1) // unchanged branch returns original pointer

		// Force modtime change so Reload takes modified path branch.
		assert.NoError(t, os.WriteFile(path, []byte(`bye #{name}`), 0644))
		future := time.Now().Add(2 * time.Second)
		assert.NoError(t, os.Chtimes(path, future, future))

		nt2, err := tc.Reload()
		assert.NoError(t, err)
		assert.NotEqual(t, tc, nt2)
		assert.Equal(t, `bye B`, nt2.Str(M.SX{`name`: `B`}))
	})
}
