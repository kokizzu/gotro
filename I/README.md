# I
--
    import "github.com/kokizzu/gotro/I"


## Usage

#### func  If

```go
func If(b bool, yes int64) int64
```
If simplified ternary operator (bool ? val : 0), returns second argument, if the
condition (first arg) is true, returns 0 if not

    I.If(true,3) // 3
    I.If(false,3) // 0

#### func  IfElse

```go
func IfElse(b bool, yes, no int64) int64
```
IfElse ternary operator (bool ? val1 : val2), returns second argument if the
condition (first arg) is true, third argument if not

    I.IfElse(true,3,4) // 3
    I.IfElse(false,3,4) // 4

#### func  IfZero

```go
func IfZero(val1, val2 int64) int64
```
IfZero simplified ternary operator (bool ? val1==0 : val2), returns second
argument, if val1 (first arg) is zero, returns val2 if not

    I.IfZero(0,3) // 3
    I.IfZero(4,3) // 4

#### func  IsZero

```go
func IsZero(val1, val2 int) int
```
IsZero simplified ternary operator (bool ? val1==0 : val2), returns second
argument, if val1 (first arg) is zero, returns val2 if not

    I.IsZero(0,3) // 3
    I.IsZero(4,3) // 4

#### func  Max

```go
func Max(a, b int64) int64
```
Max int64 max of two values

    I.Max(int64(3),int64(4)) // 4

#### func  MaxOf

```go
func MaxOf(a, b int) int
```
MaxOf int max of two values

    I.MaxOf(3,4) // 4

#### func  Min

```go
func Min(a, b int64) int64
```
Min int64 min of two values

    I.Min(int64(3),int64(4)) // 3

#### func  MinOf

```go
func MinOf(a, b int) int
```
MinOf int min of two values

    I.MinOf(3,4) // 3

#### func  PadZero

```go
func PadZero(num int64, length int) string
```
PadZero converts int64 (first arg) to string with zero padded with maximum
length I.PadZero(123,5) // `00123`

#### func  Roman

```go
func Roman(num int64) string
```
Roman convert int64 to roman number

    I.ToRoman(16)) // output "XVI"

#### func  ToEnglishNum

```go
func ToEnglishNum(num int64) string
```
ToEnglishNum format ordinal number suffix such as st, nd, rd, and th.

    I.ToEnglishNum(241)) // `241st`
    I.ToEnglishNum(242)) // `242nd`
    I.ToEnglishNum(244)) // `244th`

#### func  ToS

```go
func ToS(num int64) string
```
ToS convert int64 to string

    I.ToS(int64(1234)) // `1234`

#### func  ToStr

```go
func ToStr(num int) string
```
ToStr convert int to string

    I.ToStr(1234) // `1234`

#### func  UIf

```go
func UIf(b bool, yes uint64) uint64
```
UIf simplified ternary operator (bool ? val : 0), returns second argument, if
the condition (first arg) is true, returns 0 if not

    UI.UIf(true,3) // 3
    UI.UIf(false,3) // 0

#### func  UIfElse

```go
func UIfElse(b bool, yes, no uint64) uint64
```
UIfElse ternary operator (bool ? val1 : val2), returns second argument if the
condition (first arg) is true, third argument if not

    UI.UIfElse(true,3,4) // 3
    UI.UIfElse(false,3,4) // 4

#### func  UIfZero

```go
func UIfZero(val1, val2 uint64) uint64
```
UIfZero simplified ternary operator (bool ? val1==0 : val2), returns second
argument, if val1 (first arg) is zero, returns val2 if not

    UI.UIfZero(0,3) // 3
    UI.UIfZero(4,3) // 4

#### func  UIsZero

```go
func UIsZero(val1, val2 uint) uint
```
UIsZero simplified ternary operator (bool ? val1==0 : val2), returns second
argument, if val1 (first arg) is zero, returns val2 if not

    I.UIsZero(0,3) // 3
    I.UIsZero(4,3) // 4

#### func  UMax

```go
func UMax(a, b uint64) uint64
```
UMax uint64 max of two values

    I.UMax(uint64(3),uint64(4)) // 4

#### func  UMaxOf

```go
func UMaxOf(a, b uint) uint
```
UMaxOf uint max of two values

    I.MaxOf(3,4) // 4

#### func  UMin

```go
func UMin(a, b uint64) uint64
```
UMin uint64 min of two values

    I.UMin(uint64(3),uint64(4)) // 3

#### func  UMinOf

```go
func UMinOf(a, b uint) uint
```
UMinOf uint min of two values

    I.MinOf(3,4) // 3

#### func  UToS

```go
func UToS(num uint64) string
```
UToS convert uint64 to string

    I.UToS(uint64(1234)) // `1234`

#### func  UToStr

```go
func UToStr(num uint) string
```
UToStr convert int to string

    I.UToStr(uint(1234)) // `1234`
