package A

/*
 Type Alias: A.X
 Desc: array (slice) of anything
 Usage:
   v := A.X{}
   v = append(v, any_value)
*/
type X []interface{}

/*
 Type Alias: A.MSX
 Desc: array (slice) of map with string key and any value
 Usage:
  v := A.MSX{}
  v = append(v, map[string]{
    `foo`: 123,
    `bar`: `yay`,
  })
*/
type MSX []map[string]interface{}
