# X
--
    import "gotro/X"


## Usage

#### func  ArrToIntArr

```go
func ArrToIntArr(any_arr []interface{}) []int64
```
Convert array of any data type to array of int64

    var m4 []interface{}
    m4 = []interface{}{1}     // // tipe array
    L.ParentDescribe(X.ArrToIntArr(m4)) // []int64{1}

#### func  ArrToStrArr

```go
func ArrToStrArr(any_arr []interface{}) []string
```
convert array of any data type to array of string

    var m4 []interface{}
    m4 = []interface{}{1}     // // tipe array
    L.ParentDescribe(X.ArrToStrArr(m4)) // []string{"1"}

#### func  ToAX

```go
func ToAX(any interface{}) A.X
```

#### func  ToArr

```go
func ToArr(any interface{}) []interface{}
```
convert any data type to array of any

    var m3 interface{}
    m3 = []interface{}{1}   // tipe array
    L.ParentDescribe(X.ToArr(m3)) // []interface {}{int(1),}

#### func  ToBool

```go
func ToBool(any interface{}) bool
```
convert any data type to bool

    var m interface{}
    m = `123`
    L.ParentDescribe(X.ToBool(m)) // bool(true)

#### func  ToByte

```go
func ToByte(any interface{}) byte
```
convert any data type to int8

    var m interface{}
    m = `123`
    L.ParentDescribe(X.ToByte(m)) // byte(123)

#### func  ToF

```go
func ToF(any interface{}) float64
```
Convert any data type to float64

    var m interface{}
    m = `123.5`
    L.ParentDescribe(X.ToF(m)) // float64(123.5)

#### func  ToI

```go
func ToI(any interface{}) int64
```
convert any data type to int64

    var m interface{}
    m = `123`
    L.ParentDescribe(X.ToI(m)) // int64(123)

#### func  ToJson

```go
func ToJson(any interface{}) string
```
convert to standard json text

#### func  ToJson5

```go
func ToJson5(any interface{}) string
```

#### func  ToJsonPretty

```go
func ToJsonPretty(any interface{}) string
```
convert to beautiful json text

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
func ToMSS(any interface{}) M.SS
```

#### func  ToMSX

```go
func ToMSX(any interface{}) M.SX
```

#### func  ToS

```go
func ToS(any interface{}) string
```
convert any data type to string

    var m interface{}
    m = `123`
    L.ParentDescribe(X.ToS(m)) // `123`

#### func  ToTime

```go
func ToTime(any interface{}) time.Time
```
convert any to time

#### func  ToU

```go
func ToU(any interface{}) uint64
```
convert any data type to uint

    var m interface{}
    m = `123`
    L.ParentDescribe(X.ToI(m)) // uint(123)

#### func  ToYaml

```go
func ToYaml(any interface{}) string
```
