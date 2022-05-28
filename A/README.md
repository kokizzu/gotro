# A
--
    import "github.com/kokizzu/gotro/A"


## Usage

#### func  FloatExist

```go
func FloatExist(arr []float64, val float64) bool
```
FloatExist check if float exists on array

#### func  IntAppendIfNotExists

```go
func IntAppendIfNotExists(arr []int64, str int64) []int64
```
IntAppendIfNotExists Append if not exists

#### func  IntContains

```go
func IntContains(arr []int64, str int64) bool
```
IntContains Append int64 to array of string if not exists

#### func  IntJoin

```go
func IntJoin(arr []int64, sep string) string
```
IntJoin combine int64s in the array of int64 with the chosen string separator

    m1:= []int64{123,456}
    A.IntJoin(m1,`|`) // 123|456

#### func  IntsAppendIfNotExists

```go
func IntsAppendIfNotExists(arr []int64, ints []int64) []int64
```
IntsAppendIfNotExists Append slices if not exists

#### func  ParseEmail

```go
func ParseEmail(str_emails, each_name string) []string
```
ParseEmail split, add alias, and concat emails with name

#### func  StrAppendIfNotExists

```go
func StrAppendIfNotExists(arr []string, str string) []string
```
StrAppendIfNotExists Append if not exists

#### func  StrContains

```go
func StrContains(arr []string, str string) bool
```
StrContains Append string to array of string if not exists

#### func  StrJoin

```go
func StrJoin(arr []string, sep string) string
```
StrJoin combine strings in the array of string with the chosen string separator

    m1:= []string{`satu`,`dua`}
    A.StrJoin(m1,`-`) // satu-dua

#### func  StrToInt

```go
func StrToInt(arr []string) []int64
```
StrToInt Convert array of string to array of int64

    func main() {
      m:= []string{`1`,`2`}
      L.Print(A.StrToInt(m))//output [1 2]
    }

convert string list to integer list

#### func  StrsAppendIfNotExists

```go
func StrsAppendIfNotExists(arr []string, strs []string) []string
```
StrsAppendIfNotExists Append slices if not exists

#### func  ToJson

```go
func ToJson(arr []interface{}) string
```
ToJson convert map array of string to JSON string type

    m:= []interface{}{123,`abc`}
    L.Print(A.ToJson(m)) // [123,"abc"]

#### func  UIntJoin

```go
func UIntJoin(arr []uint64, sep string) string
```
UIntJoin combine uint64s in the array of int64 with the chosen string separator

    m1:= []uint64{123,456}
    A.UIntJoin(m1,`-`) // 123-456

#### type MSX

```go
type MSX []map[string]interface{}
```

MSX array (slice) of map with string key and any value

    v := A.MSX{}
    v = append(v, map[string]{
      `foo`: 123,
      `bar`: `yay`,
    })

#### type X

```go
type X []interface{}
```

X array (slice) of anything

    v := A.X{}
    v = append(v, any_value)
