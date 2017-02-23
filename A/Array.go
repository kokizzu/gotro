package A

// Array Support package

// array (slice) of anything
//  v := A.X{}
//  v = append(v, any_value)
type X []interface{}

// array (slice) of map with string key and any value
//  v := A.MSX{}
//  v = append(v, map[string]{
//    `foo`: 123,
//    `bar`: `yay`,
//  })
type MSX []map[string]interface{}
