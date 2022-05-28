# S
--
    import "github.com/kokizzu/gotro/S"


## Usage

```go
const MaxStrLenCB63 = 11
```

```go
const WebBR = "\n<br/>"
```

```go
var ModCB63 []uint64
```

#### func  AsF

```go
func AsF(str string) (float64, bool)
```
AsF convert to float64 with check

    S.AsF(`1234.5`) // 1234.5, true
    S.AsF(`1abc`) // 0.0, false

#### func  AsI

```go
func AsI(str string) (int64, bool)
```
AsI convert to int64 with check

    S.AsI(`1234`) // 1234, true
    S.AsI(`1abc`) // 0, false

#### func  AsU

```go
func AsU(str string) (uint, bool)
```
AsU convert to uint with check

    S.AsU(`1234`) // 1234, true
    S.AsU(`1abc`) // 0, false

#### func  BT

```go
func BT(str string) string
```
BT add backtick quote in the beginning and the end of string, without escaping.

    S.Q(`coba`) // "`coba`"
    S.Q(`123`)  // "`123`"

#### func  CamelCase

```go
func CamelCase(s string) string
```
CamelCase convert to camelCase

#### func  CharAt

```go
func CharAt(str string, index int) string
```
CharAt get character at specific index, utf-8 safe

    S.CharAt(`Halo 世界`, 5) // `世` // utf-8 example, if characters not shown, it's probably your font/editor/plugin
    S.CharAt(`Halo`, 3) // `o`

#### func  CheckPassword

```go
func CheckPassword(hash string, rawPassword string) error
```
CheckPassword check encrypted password

#### func  Coalesce

```go
func Coalesce(strs ...string) string
```
Coalesce coalesce, return first non-empty string

    S.Coalesce(`1`,`2`) // `1`
    S.Coalesce(``,`2`) // `2`
    S.Coalesce(``,``,`3`) // `3`

#### func  ConcatIfNotEmpty

```go
func ConcatIfNotEmpty(str, sep string) string
```
ConcatIfNotEmpty concat if not empty with additional separator

#### func  Contains

```go
func Contains(str, substr string) bool
```
Contains check whether the input string (first arg) contains a certain sub
string (second arg) or not.

    S.Contains(`komputer`,`om`)) // bool(true)
    S.Contains(`komputer`,`opu`)) // bool(false)

#### func  Count

```go
func Count(str, substr string) int
```
Count count how many specific character (first arg) that the string (second arg)
contains

    S.Count(`komputeer`,`e`))// output int(2)

#### func  DecodeCB63

```go
func DecodeCB63(str string) (int64, bool)
```
DecodeCB63 convert custom base-63 encoding to int64

    S.DecodeCB63(`--0`) // 1, true
    S.DecodeCB64(`(*&#$`) // 0, false

#### func  EncodeCB63

```go
func EncodeCB63(id int64, min_len int) string
```
EncodeCB63 convert integer to custom base-63 encoding that lexicographically
correct, positive integer only

    0       -
    1..10   0..9
    11..36  A..Z
    37      _
    38..63  a..z
    S.EncodeCB63(11,1) // `A`
    S.EncodeCB63(1,3) // `--0`

#### func  EncryptPassword

```go
func EncryptPassword(s string) string
```
EncryptPassword hash password (with salt)

#### func  EndsWith

```go
func EndsWith(str, suffix string) bool
```
EndsWith check whether the input string (first arg) ends with a certain
character (second arg) or not.

    S.EndsWith(`adakah`,`ah`)) // bool(true)
    S.EndsWith(`adakah`,`aka`)) // bool(false)

#### func  Equals

```go
func Equals(strFirst, strSecond string) bool
```
Equals compare two input string (first arg) equal with another input string
(second arg).

    S.Equals(`komputer`,`komputer`)) // bool(true)
    S.Equals(`komputer`,`Komputer`)) // bool(false)

#### func  EqualsIgnoreCase

```go
func EqualsIgnoreCase(strFirst, strSecond string) bool
```
EqualsIgnoreCase compare two input string (first arg) equal with ignoring case
another input string (second arg).

    S.EqualsIgnoreCase(`komputer`,`komputer`)) // bool(true)
    S.EqualsIgnoreCase(`komputer`,`Komputer`)) // bool(true)

#### func  FirstIsLower

```go
func FirstIsLower(s string) bool
```
FirstIsLower check first character is lowercase

#### func  HashPassword

```go
func HashPassword(pass string) string
```
HashPassword hash password with sha256 (without salt)

#### func  If

```go
func If(b bool, yes string) string
```
If simplified ternary operator (bool ? val : 0), returns second argument, if the
condition (first arg) is true, returns empty string if not

    S.If(true,`a`) // `a`
    S.If(false,`a`) // ``

#### func  IfElse

```go
func IfElse(b bool, yes, no string) string
```
IfElse ternary operator (bool ? val1 : val2), returns second argument if the
condition (first arg) is true, third argument if not

    S.IfElse(true,`a`,`b`) // `a`
    S.IfElse(false,`a`,`b`) // `b`

#### func  IfEmpty

```go
func IfEmpty(str1, str2 string) string
```
IfEmpty coalesce, return first non-empty string

    S.IfEmpty(``,`2`) // `2`
    S.IfEmpty(`1`,`2`) // `1`

#### func  IndexOf

```go
func IndexOf(str, sub string) int
```
IndexOf get first index of S.IndexOf(`abcdcd`,`c) // 2, -1 if not exists

#### func  JsonAsArr

```go
func JsonAsArr(str string) (res []interface{}, ok bool)
```
JsonAsArr convert JSON object to []interface{} with check

    json_str := `[1,2,['test'],'a']`
    arr, ok := S.JsonAsArr(json_str)

#### func  JsonAsFloatArr

```go
func JsonAsFloatArr(str string) (res []float64, ok bool)
```
JsonAsFloatArr convert JSON object to []float64 with check

    json_str := `[1,2,3]`
    arr, ok := S.JsonAsFloatArr(json_str)

#### func  JsonAsIntArr

```go
func JsonAsIntArr(str string) (res []int64, ok bool)
```
JsonAsIntArr convert JSON object to []int64 with check

    json_str := `[1,2,3]`
    arr, ok := S.JsonAsIntArr(json_str)

#### func  JsonAsMap

```go
func JsonAsMap(str string) (res map[string]interface{}, ok bool)
```
JsonAsMap convert JSON object to map[string]interface{} with check

    json_str := `{"test":123,"bla":[1,2,3,4]}`
    map1, ok := S.JsonAsMap(json_str)

#### func  JsonAsStrArr

```go
func JsonAsStrArr(str string) (res []string, ok bool)
```
JsonAsStrArr convert JSON object to []string with check

    json_str := `["a","b","c"]`
    arr, ok := S.JsonAsStrArr(json_str)

#### func  JsonToArr

```go
func JsonToArr(str string) (res []interface{})
```
JsonToArr convert JSON object to []interface{}, silently print and return empty
slice of interface if failed

    json_str := `[1,2,['test'],'a']`
    arr := S.JsonToArr(json_str)

#### func  JsonToIntArr

```go
func JsonToIntArr(str string) (res []int64)
```
JsonToIntArr convert JSON object to []int64, silently print and return empty
slice of interface if failed

    json_str := `[1,2,['test'],'a']`
    arr := S.JsonToArr(json_str)

#### func  JsonToMap

```go
func JsonToMap(str string) (res map[string]interface{})
```
JsonToMap convert JSON object to map[string]interface{}, silently print and
return empty map if failed

    json_str := `{"test":123,"bla":[1,2,3,4]}`
    map1 := S.JsonToMap(json_str)

#### func  JsonToObjArr

```go
func JsonToObjArr(str string) (res []map[string]interface{})
```
JsonToObjArr convert JSON object to []map[string]interface{}, silently print and
return empty slice of interface if failed

    json_str := `[{"x":"foo"},{"y":"bar"}]`
    arr := S.JsonToObjArr(json_str)

#### func  JsonToStrArr

```go
func JsonToStrArr(str string) (res []string)
```
JsonToStrArr convert JSON object to []string, silently print and return empty
slice of interface if failed

    json_str := `["123","456",789]`
    arr := S.JsonToStrArr(json_str)

#### func  JsonToStrStrMap

```go
func JsonToStrStrMap(str string) (res map[string]string)
```
JsonToStrStrMap convert JSON object to map[string]string, silently print and
return empty map if failed

    json_str := `{"test":123,"bla":[1,2,3,4]}`
    map1 := S.JsonToMap(json_str)

#### func  LastIndexOf

```go
func LastIndexOf(str, sub string) int
```
LastIndexOf get last index of

    S.LastIndexOf(`abcdcd`,`c`) // 4, -1 if not exists

#### func  Left

```go
func Left(str string, n int) string
```
Left substring at most n characters

#### func  LeftN

```go
func LeftN(str string, n int) string
```
LeftN substring at most n characters

#### func  LeftOf

```go
func LeftOf(str, substr string) string
```
LeftOf substring before first `substr`

#### func  LeftOfLast

```go
func LeftOfLast(str, substr string) string
```
LeftOfLast substring before last `substr`

#### func  LowerFirst

```go
func LowerFirst(s string) string
```
LowerFirst convert to lower only first char

#### func  MergeMailContactEmails

```go
func MergeMailContactEmails(each_name, str_emails string) []string
```
MergeMailContactEmails return formatted array of mail contact <usr@email>

#### func  Mid

```go
func Mid(str string, left int, length int) string
```
Mid substring at set left right n characters

#### func  PadLeft

```go
func PadLeft(s string, padStr string, lenStr int) string
```
PadLeft append padStr to left until length is lenStr

#### func  PadRight

```go
func PadRight(s string, padStr string, lenStr int) string
```
PadRight append padStr to right until length is lenStr

#### func  PascalCase

```go
func PascalCase(s string) string
```
PascalCase convert to PascalCase source: https://github.com/iancoleman/strcase

#### func  Q

```go
func Q(str string) string
```
Q add single quote in the beginning and the end of string, without escaping.

    S.Q(`coba`) // `'coba'`
    S.Q(`123`)  // `'123'`

#### func  QQ

```go
func QQ(str string) string
```
QQ add double quote in the beginning and the end of string, without escaping.

    S.Q(`coba`) // `"coba"`
    S.Q(`123`)  // `"123"`

#### func  RandomCB63

```go
func RandomCB63(len int64) string
```
RandomCB63 random CB63 n-times, the result is n*MaxStrLenCB63 bytes

#### func  RandomPassword

```go
func RandomPassword(strlen int64) string
```
RandomPassword create a random password

#### func  RemoveCharAt

```go
func RemoveCharAt(str string, index int) string
```
RemoveCharAt remove character at specific index, utf-8 safe

    S.RemoveCharAt(`Halo 世界`, 5) // `Halo 界` --> utf-8 example, if characters not shown, it's probably your font/editor/plugin
    S.RemoveCharAt(`Halo`, 3) // `Hal`

#### func  RemoveLastN

```go
func RemoveLastN(str string, n int) string
```
RemoveLastN remove last n character, not UTF-8 friendly

#### func  Repeat

```go
func Repeat(str string, count int) string
```
repeat string

#### func  Replace

```go
func Replace(haystack, needle, gold string) string
```
Replace replace all substring with another substring

    S.Replace(`bisa`,`is`,`us`) // `busa`

#### func  Right

```go
func Right(str string, n int) string
```
Right substring at right most n characters

#### func  RightOf

```go
func RightOf(str, substr string) string
```
RightOf substring after first `substr`

#### func  RightOfLast

```go
func RightOfLast(str, substr string) string
```
RightOfLast substring after last `substr`

#### func  SnakeCase

```go
func SnakeCase(s string) string
```
SnakeCase convert to snake case source: https://github.com/iancoleman/strcase

#### func  Split

```go
func Split(str, sep string) []string
```
Split split a string (first arg) by characters (second arg) into array of
strings (output).

    S.Split(`biiiissssa`,func(ch rune) bool { return ch == `i` }) // output []string{"b", "", "", "", "ssssa"}

#### func  SplitFunc

```go
func SplitFunc(str string, fun func(rune) bool) []string
```
SplitFunc split a string (first arg) based on a function

#### func  SplitN

```go
func SplitN(str string, n int) []string
```
SplitN split to substrings with maximum n characters

#### func  StartsWith

```go
func StartsWith(str, prefix string) bool
```
StartsWith check whether the input string (first arg) starts with a certain
character (second arg) or not.

    S.StartsWith(`adakah`,`ad`) // bool(true)
    S.StartsWith(`adakah`,`bad`) // bool(false)

#### func  ToF

```go
func ToF(str string) float64
```
ToF convert string to float64, returns 0 and silently print error if not valid

    S.ToF(`1234.5`) // 1234.5
    S.ToF(`1a`) // 0.0

#### func  ToI

```go
func ToI(str string) int64
```
ToI convert string to int64, returns 0 and silently print error if not valid

    S.ToI(`1234`) // 1234
    S.ToI(`1a`) // 0

#### func  ToInt

```go
func ToInt(str string) int
```
ToInt convert string to int, returns 0 and silently print error if not valid

    S.ToInt(`1234`) // 1234
    S.ToInt(`1a`) // 0

#### func  ToLower

```go
func ToLower(str string) string
```
ToLower change the characters in string to lowercase

    S.ToLower(`BIsa`) // "bisa"

#### func  ToTitle

```go
func ToTitle(str string) string
```
ToTitle Change first letter for every word to uppercase

    S.ToTitle(`Disa dasi`)) // output "Disa Dasi"

#### func  ToU

```go
func ToU(str string) uint64
```
ToU convert string to uint64, returns 0 and silently print error if not valid

    S.ToU(`1234`) // 1234
    S.ToU(`1a`) // 0

#### func  ToUpper

```go
func ToUpper(str string) string
```
ToUpper change the characters in string to uppercase S.ToUpper(`bisa`) // "BISA"

#### func  Trim

```go
func Trim(str string) string
```
Trim erase spaces from left and right

    S.Trim(` withtrim:  `) // `withtrim:`

#### func  TrimChars

```go
func TrimChars(str, chars string) string
```
TrimChars remove chars from beginning and end

    S.TrimChars(`aoaaffoa`,`ao`) // `ff`

#### func  UZ

```go
func UZ(str string) string
```
UZ replace <, >, and & back, quot and apos to alternative utf8

#### func  UZRAW

```go
func UZRAW(str string) string
```
UZRAW replace <, >, and & back, quot and apos to real html

#### func  UpperFirst

```go
func UpperFirst(s string) string
```
UpperFirst convert to lower only first char

#### func  ValidateEmail

```go
func ValidateEmail(str string) string
```
ValidateEmail return empty string if str is not a valid email

#### func  ValidateFilename

```go
func ValidateFilename(str string) string
```
ValidateFilename validate file name

#### func  ValidateMailContact

```go
func ValidateMailContact(str string) string
```
ValidateMailContact return valid version of mail contact (part before
<usr@email>)

#### func  ValidatePhone

```go
func ValidatePhone(str string) string
```
ValidatePhone remove invalid characters of a phone number

#### func  XSS

```go
func XSS(str string) string
```
XSS replace <, >, ', ", % but without giving single quote

#### func  Z

```go
func Z(str string) string
```
Z trim, replace <, >, ', " and gives single quote

    S.Z(`<>'"`) // `&lt;&gt;&apos;&quot;

#### func  ZB

```go
func ZB(b bool) string
```
ZB give ' to boolean value

    S.ZB(true)  // `'true'`
    S.ZB(false) // `'false'`

#### func  ZI

```go
func ZI(num int64) string
```
ZI give ' to int64 value

    S.ZI(23)) // '23'
    S.ZI(03)) // '3'

#### func  ZJ

```go
func ZJ(str string) string
```
ZJ single quote a json string

    hai := `{'test':123,"bla":[1,2,3,4]}`
    S.ZJ(hai) // "{'test':123,\"bla\":[1,2,3,4]}"

#### func  ZJJ

```go
func ZJJ(str string) string
```
ZJJ double quote a json string

    hai := `{'test':123,"bla":[1,2,3,4]}`
    S.ZJJ(hai) // "{'test':123,\"bla\":[1,2,3,4]}"

#### func  ZJLIKE

```go
func ZJLIKE(str string) string
```
ZJLIKE ZLIKE but for json (not replacing double quote)

#### func  ZLIKE

```go
func ZLIKE(str string) string
```
ZLIKE replace <, >, ', ", % and gives single quote and %

    S.ZLIKE(`coba<`))  // output '%coba&lt;%'
    S.ZLIKE(`"coba"`)) // output '%&quot;coba&quot;%'

#### func  ZS

```go
func ZS(str string) string
```
ZS replace <, >, ', " and gives single quote (without trimming)

    S.Z(`<>'"`) // `&lt;&gt;&apos;&quot;

#### func  ZT

```go
func ZT(strs ...string) string
```
ZT trace function, location of the caller code, replacement for ZC

#### func  ZT2

```go
func ZT2() string
```
ZT2 trace function, location of 2nd level caller, parameterless, with newline

#### func  ZU

```go
func ZU(num uint64) string
```
ZU give ' to uint value

    S.ZI(23)) // '23'
    S.ZI(03)) // '3'

#### func  ZZ

```go
func ZZ(str string) string
```
ZZ replace ` and give double quote (for table names)

    S.ZZ(`coba"`) // `"coba&quot;"`
