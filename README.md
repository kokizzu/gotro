# Garuda
Garuda framework is rewrite of gitlab.com/kokizzu/gokil, that uses fasthttprouter

Design Goal:
- As similar as possible to Elixir's Phoenix Framework
- Opinionated (choose the best)
- 1-letter supporting package (so we can use something like: I.ToS(1234) to convert `int64` to `string`), such as:
  - A - Array
  - B - Boolean
  - C - Character (or Rune)
  - F - Floating Point
  - M - Map
  - I - Integer
  - S - String
  - T - Time (and Date)
  - X - Anything
