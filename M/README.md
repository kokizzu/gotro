# M
--
    import "gotro/M"


## Usage

#### func  SSKeysStartedWith

```go
func SSKeysStartedWith(m SS, prefix string) []string
```
retrieve all keys started with

#### func  ToJson

```go
func ToJson(hash map[string]interface{}) string
```
convert map[string]interface{} to json

    m :=  map[string]interface{}{`buah`:123,`angka`:`dia`}
    M.ToJson(m) // {"angka":"dia","buah":123}

#### type IAX

```go
type IAX map[int64][]interface{}
```

map with int64 key and array of any value

#### type IB

```go
type IB map[int64]bool
```

map with int64 key and bool value

#### func (IB) Keys

```go
func (hash IB) Keys() []int64
```
get array of int64 keys

    m :=  M.IB{1:true,2:false}
    m.Keys() // []int64{1, 2}

get all integer keys

#### func (IB) KeysConcat

```go
func (hash IB) KeysConcat(with string) string
```
get concatenated integer keys

    m := M.IB{1: true, 2: false, 3:true, 5:false}
    m.KeysConcat(`,`) // `1,2,3,5`

#### type II

```go
type II map[int64]int64
```

map with int64 key and int64 value

#### func (II) Keys

```go
func (hash II) Keys() []int64
```
get array of int64 keys

    m :=  M.II{1:1,2:3}
    m.Keys() // []int64{1, 2}

#### func (II) KeysConcat

```go
func (hash II) KeysConcat(with string) string
```
get concatenated integer keys

    m := M.II{1: 2, 2: 567, 3:6, 5:45}
    m.KeysConcat(`,`) // `1,2,3,5`

#### type IS

```go
type IS map[int64]string
```

map with int64 key and string value

#### type IX

```go
type IX map[int64]interface{}
```

map with int64 key and any value

#### func (IX) Keys

```go
func (hash IX) Keys() []int64
```
get array of int64 keys

    m :=  M.IX{1:1,2:`DIA`}
    m.Keys()) // []int64{1, 2}

#### func (IX) ToSX

```go
func (hash IX) ToSX() SX
```
convert keys to string

    m :=  M.IX{1:1,2:`DUA`}
    m.ToSX() // M.SX{"1": int(1),"2": "DUA"}

convert integer keys to string keys

#### type SAX

```go
type SAX map[string][]interface{}
```

map with string key and array of any value

#### type SB

```go
type SB map[string]bool
```

map with string key and bool value

#### func (SB) IntoJson

```go
func (hash SB) IntoJson() (string, bool)
```
convert to json string with check

#### func (SB) IntoJsonPretty

```go
func (hash SB) IntoJsonPretty() (string, bool)
```
convert to pretty json string with check

#### func (SB) KeysConcat

```go
func (hash SB) KeysConcat(with string) string
```
get concatenated string keys

    m := M.SB{`tes`:true,`coba`:true,`lah`:true}
    m.KeysConcat(`,`) // `coba,lah,tes`

#### func (SB) SortedKeys

```go
func (hash SB) SortedKeys() []string
```
get sorted keys

    m := M.SS{`tes`:true,`coba`:false,`lah`:false}
    m.SortedKeys() // []string{`coba`,`lah`,`tes`}

#### func (SB) ToJson

```go
func (hash SB) ToJson() string
```
convert to json string, silently print error if failed

#### func (SB) ToJsonPretty

```go
func (hash SB) ToJsonPretty() string
```
convert to pretty json string, silently print error if failed

#### type SF

```go
type SF map[string]float64
```

map with string key and float64 value

#### type SI

```go
type SI map[string]int64
```

map with string key and int64 value

#### type SS

```go
type SS map[string]string
```

map with string key and string value

#### func (SS) GetFloat

```go
func (hash SS) GetFloat(key string) float64
```
get float64 type from map

#### func (SS) GetInt

```go
func (hash SS) GetInt(key string) int64
```
get int64 from from map

#### func (SS) GetStr

```go
func (hash SS) GetStr(key string) string
```
get string type from map

#### func (SS) GetUint

```go
func (hash SS) GetUint(key string) uint64
```
get uint from map

#### func (SS) Keys

```go
func (hash SS) Keys() []string
```
get array of string keys

    m :=  M.SS{`satu`:`1`,`dua`:`2`}
    m.Keys() // []string{"satu", "dua"}

#### func (SS) KeysConcat

```go
func (hash SS) KeysConcat(with string) string
```
get concatenated string keys

    m := M.SS{`tes`:`tes`,`coba`:`saja`,`lah`:`lah`}
    m.KeysConcat(`,`) // `coba,lah,tes`

#### func (SS) Merge

```go
func (hash SS) Merge(src SS)
```
merge from another map

#### func (SS) Pretty

```go
func (hash SS) Pretty(sep string) string
```
get pretty printed values

#### func (SS) PrettyFunc

```go
func (hash SS) PrettyFunc(sep string, fun func(string, string) string) string
```
get pretty printed values with filter values

#### func (SS) SortedKeys

```go
func (hash SS) SortedKeys() []string
```
get sorted keys

    m := M.SS{`tes`:`tes`,`coba`:`saja`,`lah`:`lah`}
    m.SortedKeys() // []string{`coba`,`lah`,`tes`}

#### func (SS) ToJson

```go
func (hash SS) ToJson() string
```
convert to json string, silently print error if failed

#### func (SS) ToScylla

```go
func (hash SS) ToScylla() string
```
convert to scylla based map<text,text>

#### type SX

```go
type SX map[string]interface{}
```

map with string key and any value

#### func (SX) GetAX

```go
func (json SX) GetAX(key string) []interface{}
```
get array of anything value from map

    m :=  M.SX{`tes`:[]interface{}{123,`buah`}}
    m.GetAX(`tes`) // []interface {}{int(123),"buah"}

#### func (SX) GetBool

```go
func (json SX) GetBool(key string) bool
```
get bool type from map (empty string, 0, `f`, `false` are false, other non empty
are considered truthy value) m :=
M.SX{`test`:234.345,`coba`:`buah`,`angka`:float64(123),`salah`:123}

    m.GetBool(`test`)  // bool(true)
    m.GetBool(`coba`)  // bool(true)
    m.GetBool(`angka`) // bool(true)
    m.GetBool(`salah`) // bool(false)

#### func (SX) GetFloat

```go
func (json SX) GetFloat(key string) float64
```
get float64 type from map

    m := M.SX{`test`:234.345,`coba`:`buah`,`dia`:true,`angka`:23435}
    m.GetFloat(`test`)  // float64(234.345)
    m.GetFloat(`dia`)   // int64(1)
    m.GetFloat(`coba`)  // 0
    m.GetFloat(`angka`) // 0

#### func (SX) GetInt

```go
func (json SX) GetInt(key string) int64
```
get int64 type from map

    m := M.SX{`test`:234.345,`coba`:`buah`,`dia`:true,`angka`:int64(23435)}
    m.GetInt(`test`))  // int64(234)
    m.GetInt(`dia`))   // int64(1)
    m.GetInt(`coba`))  // int64(0)
    m.GetInt(`angka`)) // int64(23435)

#### func (SX) GetIntArr

```go
func (json SX) GetIntArr(key string) []int64
```
get array int64 value from map

    m :=  M.SX{`tes`:[]int64{123,234}}
    m.GetIntArr(`tes`)) // []int64{123, 234}

#### func (SX) GetMIB

```go
func (json SX) GetMIB(key string) IB
```
get map string int64 value from map

    m := M.SX{`tes`:M.SB{`satu`:true,`2`:false}}
    m.GetMSB(`tes`) // M.SB{"satu":true, "2":false}

#### func (SX) GetMSB

```go
func (json SX) GetMSB(key string) SB
```
get map string bool value from map

    m := M.SX{`tes`:M.SB{`1`:true,`2`:false}}
    m.GetMSB(`tes`) // M.SB{"1":true, "2":false}

#### func (SX) GetMSF

```go
func (json SX) GetMSF(key string) SF
```
get map string float64 value from map

    m := M.SX{`tes`:M.SF{`satu`:32.45,`2`:12}}
    m.GetMSF(`tes`) // M.SF{"satu":32.45, "2":12}

#### func (SX) GetMSI

```go
func (json SX) GetMSI(key string) SI
```
get map string int64 value from map

    m := M.SX{`tes`:M.SF{`satu`:32,`2`:12}}
    m.GetMSI(`tes`) // M.SF{"satu":32, "2":12}

#### func (SX) GetMSX

```go
func (json SX) GetMSX(key string) SX
```
get map string anything value from map

    m :=  M.SX{`tes`:M.SX{`satu`:234.345,`dua`:`huruf`,`tiga`:123}}
    m.GetMSX(`tes`) // M.SX{"tiga": int(123),"satu": float64(234.345),"dua":  "huruf"}

#### func (SX) GetStr

```go
func (json SX) GetStr(key string) string
```
get string type from map

    m := M.SX{`test`:234.345,`coba`:`buah`,`angka`:int64(123)}
    m.GetStr(`test`)  // `234.345`
    m.GetStr(`coba`)  // `buah`
    m.GetStr(`angka`) // `123`

#### func (SX) GetUint

```go
func (json SX) GetUint(key string) uint64
```
get uint type from map

    m := M.SX{`test`:234.345,`coba`:`buah`,`dia`:true,`angka`:int64(23435)}
    m.GetInt(`test`))  // int64(234)
    m.GetInt(`dia`))   // int64(1)
    m.GetInt(`coba`))  // int64(0)
    m.GetInt(`angka`)) // int64(23435)

#### func (SX) IntoJson

```go
func (hash SX) IntoJson() (string, bool)
```
convert to json string with check

#### func (SX) IntoJsonPretty

```go
func (hash SX) IntoJsonPretty() (string, bool)
```
convert to pretty json string with check

#### func (SX) Keys

```go
func (hash SX) Keys() []string
```
get array of string keys

    m :=  M.SS{`satu`:`1`,`dua`:`2`}
    m.Keys() // []string{"satu", "dua"}

#### func (SX) Pretty

```go
func (hash SX) Pretty(sep string) string
```
get pretty printed values

#### func (SX) Set

```go
func (hash SX) Set(key string, val interface{})
```
set key with any value

#### func (SX) SortedKeys

```go
func (hash SX) SortedKeys() []string
```
get sorted keys

    m := M.SX{`tes`:1,`coba`:12.4,`lah`:false}
    m.SortedKeys() // []string{`coba`,`lah`,`tes`}

#### func (SX) ToJson

```go
func (hash SX) ToJson() string
```
convert to json string, silently print error if failed

#### func (SX) ToJsonPretty

```go
func (hash SX) ToJsonPretty() string
```
convert to pretty json string, silently print error if failed
