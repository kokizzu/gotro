
# ZTemplate

a simple javascript syntax-friendly template engine (it means in the javascript editor/IDE, it wont show a lexical/syntax error or wrong syntax highlighting, the autocomplete will still works). You can use [svelte-mpa](//github.com/kokizzu/svelte-mpa) if you need reactivity on each page.

## Syntax

There's 4 tags:
- `/*! mapKey */` -- for HTML or CSS
- `#{mapKey}` -- for HTML, or if quoted can be used for Javascript string
- `[/* mapKey */]` -- for Javascript array
- `{/* mapKey */}` -- for Javascript object

You can switch around, it doesn't matter, but it would be better to use each differently so the frontend guy (if not yourself) knows what's the type.

## Template Example

```html
<html>
<head>
	<title>/*! title */</title>
</head>
<body>
<script>
  const a = +'#{ aNumber }';
  const b = [/* anArray */];
  const c = {/* aMap */};
  const d = '#{aString}';
</script>
</body>
</html>
```

## Usage Example

```go
const autoReload = true
const printDebug = true
const fileName = `dummy.html`
tc, err := ParseFile(autoReload, printDebug, fileName)
if L.IsError(err, `filed Z.ParseFile: `+fileName) {
	return err
}
buff := bytes.Buffer{}
tc.Render(&buff,M.SX{
	`title`: `this is a title`,
	`aNumber`: 123,
	`anArray`: A.X{1,`b`,`c`,4}, // == []interface{}
	`aMap`: M.SX{ // == map[string]interface{}
		`a`:1,
		`b`:`test`,
		`c`:`something`,
	},
	`aString`: `ayaya`,
})
// buff.String() will contain the output below
```

## Example Output

```html
<html>
<head>
	<title>this is a title</title>
</head>
<body>
<script>
  const a = +'123';
  const b = [1,"b","c",4];
  const c = {"a":1,"b":"test","c":"something"};
  const d = 'ayaya';
</script>
</body>
</html>
```

## Example from string

note that it's preferable for `FromString` to be reused (declared as global or cached) so it doesn't have to parse the template again and again.

```go
const printDebug = true
res := bytes.Buffer{}
FromString(`hi my name #{name}, my age #{age}`,printDebug).Render(&res,M.SX{
	`name`: `Tzuyu`,
	`age`: 21,
})
res.String() == `hi my name Tzuyu, my age 21`
str := FromString(`i like #{char}`).Str(M.SX{
	`char`: `Rem`,
})
str == `i like Rem`
```

## How Fast is it?

Because of limited set of feature compared to Go's standard text/template.

```shell
# ~12x faster compared to text/template for parsing and rendering
BenchmarkParseRenderZTemplate-32   1697598   742.1 ns/op
BenchmarkParseRenderGoTemplate-32   152878  9237   ns/op

# ~4.8x faster compared to text/template for rendering only
BenchmarkRenderZTemplate-32        8543964   139.4 ns/op
BenchmarkRenderGoTemplate-32       1736002   678.9 ns/op
```

## Performance Optimization Tips

- build/load the template once (`Z.ParseFile` or `Z.FromString`), call `.Render` multiple times
- preallocate `bytes.Buffer` so it won't need to resize, eg. using [valyala/bytebufferpool](//github.com/valyala/bytebufferpool)
- for static values, use another caching library like `bigcache` so no need to call `.Render` again and again, the render itself is fast `O(n)` where `n` is the bytes, but the serialization to JSON string might be slow.
