
# ZTemplate

a simple javascript syntax-friendly template engine (it means in the javascript editor/IDE, it wont show a lexical/syntax error or wrong syntax highlighting, the autocomplete will still works).

## Syntax

There's 4 tags:
- `/*! mapKey */` -- for HTML or CSS
- `#{mapKey}` -- for HTML, or if quoted can be used for Javascript string
- `[/* mapKey */]` -- for Javascript array
- `{/* mapKey */}` -- for Javascript object

You can switch around, it doesn't matter, but it would be better to use each differently so the frontend guy (if not yourself) knows what's the type.

## Template Example

```
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

```
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

```<html>
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
