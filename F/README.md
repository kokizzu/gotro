# F
--
    import "github.com/kokizzu/gotro/F"


## Usage

#### func  If

```go
func If(b bool, yes float64) float64
```go
simplified ternary operator (bool ? val : 0), returns second argument, if the
condition (first arg) is true, returns 0 if not

    F.If(true,3.12) // 3.12
    F.If(false,3)   // 0

#### func  IfElse

```go
func IfElse(b bool, yes, no float64) float64
```go
ternary operator (bool ? val1 : val2), returns second argument if the condition
(first arg) is true, third argument if not

    F.IfElse(true,3.12,3.45))  // 3.12

#### func  ToDateStr

```go
func ToDateStr(num float64) string
```go
convert float64 unix to `YYYY-MM-DD`

#### func  ToIsoDateStr

```go
func ToIsoDateStr(num float64) string
```go
convert to ISO-8601 string

    F.ToIsoDateStr(0) // `1970-01-01T00:00:00`

#### func  ToS

```go
func ToS(num float64) string
```go
convert float64 to string

    F.ToS(3.1284)) // `3.1284`

#### func  ToStr

```go
func ToStr(num float64) string
```go
convert float64 to string with 2 digits behind the decimal point

    F.ToStr(3.1284)) // `3.13`
