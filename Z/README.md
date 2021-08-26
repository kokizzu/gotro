
# ZTemplate

a simple javascript syntax highlighting-friendly template engine.

## Template Example

```
<html>
<head>
	<title>/*! title */</title>
</head>
<body>
<script>
  const a = '#{ aString }';
  const b = [/* anArray */];
  const c = {/* aMap */ };
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
	`aString`: `this is a string`,
	`anArray`: A.X{1,`b`,`c`,4},
	`aMap`: M.SX{
		`a`:1,
		`b`:`test`,
		`c`:`something`,
	},
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
  const a = 'this is a string';
  const b = [1,"b","c",4];
  const c = {"a":1,"b":"test","c":"something"};
</script>
</body>
</html>
```
