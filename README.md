# GotRo

GotRo (lit. Gotong Royong, meaning in Indonesia: do it together, mutual cooperation) framework is rewrite of gitlab.com/kokizzu/gokil, that uses fasthttprouter

Design Goal:
- As similar as possible to Elixir's Phoenix Framework
- Opinionated (choose the best dependency), for example by default uses int64 and float64
- 1-letter supporting package (so we can use something like: I.ToS(1234) to convert `int64` to `string`), such as:
  - A - Array
  - B - Boolean
  - C - Character (or Rune)
  - F - Floating Point
  - L - Logging
  - M - Map
  - I - Integer
  - S - String
  - T - Time (and Date)
  - X - Anything (aka `interface{}`)
  - Z - Z-Template Engine, that has syntax similar to ruby string interpolation `#{foo}` or any other that javascript friendly `{/* foo */}`, `[/* bar */]`, `/*! bar */`
- Comment and examples on each type and function, so it can be viewed using godoc, something like: `godoc github.com/kokizzu/gotro/A`