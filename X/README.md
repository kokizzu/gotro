# X
--
    import "github.com/kokizzu/gotro/X"


## Usage

#### func  ArrToIntArr

```go
func ArrToIntArr(any_arr []any) []int64
```
ArrToIntArr Convert array of any data type to array of int64

    var m4 []any
    m4 = []any{1}     // // tipe array
    L.ParentDescribe(X.ArrToIntArr(m4)) // []int64{1}

#### func  ArrToStrArr

```go
func ArrToStrArr(any_arr []any) []string
```
ArrToStrArr convert array of any data type to array of string

    var m4 []any
    m4 = []any{1}     // // tipe array
    L.ParentDescribe(X.ArrToStrArr(m4)) // []string{"1"}

#### func  ToAX

```go
func ToAX(x any) A.X
```
ToAX convert to []any

#### func  ToArr

```go
func ToArr(x any) []any
```
ToArr convert any data type to array of any

    var m3 any
    m3 = []any{1}   // tipe array
    L.ParentDescribe(X.ToArr(m3)) // []interface {}{int(1),}

#### func  ToBool

```go
func ToBool(any any) bool
```
ToBool convert any data type to bool

    var m any
    m = `123`
    L.ParentDescribe(X.ToBool(m)) // bool(true)

#### func  ToByte

```go
func ToByte(x any) byte
```
ToByte convert any data type to int8

    var m any
    m = `123`
    L.ParentDescribe(X.ToByte(m)) // byte(123)

#### func  ToF

```go
func ToF(x any) float64
```
ToF Convert any data type to float64

    var m any
    m = `123.5`
    L.ParentDescribe(X.ToF(m)) // float64(123.5)

#### func  ToI

```go
func ToI(x any) int64
```
ToI convert any data type to int64

    var m any
    m = `123`
    L.ParentDescribe(X.ToI(m)) // int64(123)

#### func  ToJson

```go
func ToJson(any any) string
```
ToJson convert to standard json text

#### func  ToJson5

```go
func ToJson5(x any) string
```
ToJson5 convert to json5

#### func  ToJsonPretty

```go
func ToJsonPretty(any any) string
```
ToJsonPretty convert to beautiful json text

    m:= []interface {}{true,`1`,23,`wabcd`}
    L.Print(K.ToJsonPretty(m))
    // [
    //   true,
    //   "1",
    //   23,
    //   "wabcd"
    // ]

#### func  ToMSS

```go
func ToMSS(any any) M.SS
```
ToMSS convert to map[string]string

#### func  ToMSX

```go
func ToMSX(x any) M.SX
```
ToMSX convert to map[string]any

#### func  ToS

```go
func ToS(x any) string
```
ToS convert any data type to string

    var m any
    m = `123`
    L.ParentDescribe(X.ToS(m)) // `123`

#### func  ToTime

```go
func ToTime(x any) time.Time
```
ToTime convert any to time

#### func  ToU

```go
func ToU(x any) uint64
```
ToU convert any data type to uint

    var m any
    m = `123`
    L.ParentDescribe(X.ToI(m)) // uint(123)

#### func  ToYaml

```go
func ToYaml(any any) string
```
ToYaml convert to yaml text
