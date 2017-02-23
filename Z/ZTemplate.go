package Z

// Z-TemplateEngine library

import (
	"bytes"
	"errors"
	"gitlab.com/kokizzu/gokil/L"
	"gitlab.com/kokizzu/gokil/M"
	"gitlab.com/kokizzu/gokil/S"
	"gitlab.com/kokizzu/gokil/X"
	"io/ioutil"
	"os"
	"time"
)

// The Z-Template engine syntax are javascript friendly or similar to ruby's string interpolation:
//  [/* test */]
//  {/* test */}
//  /*! test */
//  #{test}

// To use the library what you need to do is create a TemplateChain, call ParseFile to load a template
// call Render function to render to buffer
// call Reload if you want to rescan from file

type TemplateChain struct {
	Parts       [][]byte
	Keys        []string
	Filename    string
	ModTime     time.Time
	AutoRefresh bool
	PrintDebug  bool
}

func (t TemplateChain) Describe() {
	for k, v := range t.Parts {
		L.Describe(k, string(v))
	}
}

// write to buffer
func (t *TemplateChain) Render(target *bytes.Buffer, values M.SX) {
	if t.AutoRefresh {
		nt, err := t.Reload()
		if err == nil {
			t = nt
		}
	}
	not_found := M.SB{}
	used := M.SB{}
	for idx, key := range t.Keys {
		target.Write(t.Parts[idx])
		val, ok := values[key]
		if !ok {
			target.WriteString(key)
			not_found[key] = true
		} else {
			target.WriteString(X.ToS(val))
		}
		used[key] = true
	}
	if t.PrintDebug {
		not_used_arr := []string{}
		for key := range values {
			if !used[key] {
				not_used_arr = append(not_used_arr, key)
			}
		}
		if len(not_used_arr) > 0 {
			L.Describe(`Unused template parameter on `+t.Filename, not_used_arr)
		}
		not_found_arr := []string{}
		for key := range not_found {
			not_found_arr = append(not_found_arr, key)
		}
		if len(not_found_arr) > 0 {
			L.Describe(`Template parameter not found on `+t.Filename, not_found_arr)
		}
	}
	target.Write(t.Parts[len(t.Parts)-1])
}

// reload from file
func (t *TemplateChain) Reload() (*TemplateChain, error) {
	dup := TemplateChain{
		Parts:       [][]byte{},
		Keys:        []string{},
		ModTime:     t.ModTime,
		Filename:    t.Filename,
		AutoRefresh: t.AutoRefresh,
		PrintDebug:  t.PrintDebug,
	}
	// check for changes
	info, err := os.Stat(dup.Filename)
	errinfo := `failed to stat the template`
	//	L.Describe(t.Filename, err)
	if L.IsError(err, errinfo, dup.Filename) {
		dup.Parts = [][]byte{[]byte(errinfo)}
		return &dup, errors.New(errinfo)
	}
	// when not modified
	mod_time := info.ModTime()
	if dup.ModTime.Sub(mod_time).Nanoseconds() == 0 && mod_time.UnixNano() != 0 {
		return t, nil
	}
	// when modified
	dup.ModTime = mod_time
	// read the actual file
	bs, err := ioutil.ReadFile(dup.Filename)
	errinfo = `template not found`
	if L.IsError(err, errinfo, dup.Filename) {
		dup.Parts = [][]byte{[]byte(errinfo)}
		return &dup, errors.New(errinfo)
	}
	// clear parts
	dup.Parts = [][]byte{}
	dup.Keys = []string{}
	//	L.Describe(`start parsing`, t.Filename)
	// split into Parts and Keys
	start := 0
	p := 0
	key := ``
	for p < len(bs)-3 {
		end := 0
		ch := bs[p]
		// example: var y = '#{key}'
		// example: var y = '#{ key }'
		if ch == '#' { // start with #{
			if ch1 := bs[p+1]; ch1 == '{' {
				for z := p + 2; z < len(bs); z += 1 { // find the }
					if ec := bs[z]; ec == '}' {
						key = string(bs[p+2 : z])
						end = z + 1
						break
					}
				}
			}
		} else if ch == '/' { // start with `/*!`
			// example: var y {
			//   /*! key */
			//   /*!key*/
			// }
			if ch1 := bs[p+1]; ch1 == '*' {
				if ch2 := bs[p+2]; ch2 == '!' {
					for z := p + 3; z < len(bs)-1; z += 1 { // find the `*/`
						if ec := bs[z]; ec == '*' {
							if ec1 := bs[z+1]; ec1 == '/' {
								key = string(bs[p+3 : z])
								end = z + 2
								break
							}
						}
					}
				}
			}
		} else if ch == '{' { // start with {
			// example: var y = {/* key */}
			// example: var y = {/* key */ }
			// example: var y = {/*key*/}
			// example: var y = {/*key*/ }
			ch1 := bs[p+1]
			if ch1 == '/' { // start with `{/*`
				if ch2 := bs[p+2]; ch2 == '*' {
					for z := p + 3; z < len(bs)-2; z += 1 { // find the `*/}`
						if ec := bs[z]; ec == '*' {
							if ec1 := bs[z+1]; ec1 == '/' {
								ec2 := bs[z+2]
								if ec2 == '}' {
									key = string(bs[p+3 : z])
									end = z + 3
									break
								} else if ec2 == ' ' {
									if ec3 := bs[z+3]; ec3 == '}' {
										key = string(bs[p+3 : z])
										end = z + 4
										break
									}
								}
							}
						}
					}
				}
				// example: var y = { /* key */ }
				// example: var y = { /* key */}
				// example: var y = { /*key*/ }
				// example: var y = { /*key*/}
			} else if ch1 == ' ' { // start with `{ /*`
				if ch2 := bs[p+2]; ch2 == '/' {
					if ch3 := bs[p+3]; ch3 == '*' {
						for z := p + 4; z < len(bs)-3; z += 1 { // find the `*/ }`
							if ec := bs[z]; ec == '*' {
								if ec1 := bs[z+1]; ec1 == '/' {
									ec2 := bs[z+2]
									if ec2 == ' ' {
										if ec3 := bs[z+3]; ec3 == '}' {
											key = string(bs[p+4 : z])
											end = z + 4
											break
										}
									} else if ec2 == '}' {
										key = string(bs[p+4 : z])
										end = z + 3
										break
									}
								}
							}
						}
					}
				}
			}
		} else if ch == '[' { // start with [
			// example: var y = [/* key */]
			// example: var y = [/* key */ ]
			// example: var y = [/*key*/]
			// example: var y = [/*key*/ ]
			ch1 := bs[p+1]
			if ch1 == '/' { // start with `[/*`
				if ch2 := bs[p+2]; ch2 == '*' {
					for z := p + 3; z < len(bs)-2; z += 1 { // find the `*/]`
						if ec := bs[z]; ec == '*' {
							if ec1 := bs[z+1]; ec1 == '/' {
								ec2 := bs[z+2]
								if ec2 == ']' {
									key = string(bs[p+3 : z])
									end = z + 3
									break
								} else if ec2 == ' ' {
									if ec3 := bs[z+3]; ec3 == ']' {
										key = string(bs[p+3 : z])
										end = z + 4
										break
									}
								}
							}
						}
					}
				}
				// example: var y = [ /* key */ ]
				// example: var y = [ /* key */]
				// example: var y = [ /*key*/ ]
				// example: var y = [ /*key*/]
			} else if ch1 == ' ' { // start with `[ /*`
				if ch2 := bs[p+2]; ch2 == '/' {
					if ch3 := bs[p+3]; ch3 == '*' {
						for z := p + 4; z < len(bs)-3; z += 1 { // find the `*/ ]`
							if ec := bs[z]; ec == '*' {
								if ec1 := bs[z+1]; ec1 == '/' {
									ec2 := bs[z+2]
									if ec2 == ' ' {
										if ec3 := bs[z+3]; ec3 == ']' {
											key = string(bs[p+4 : z])
											end = z + 4
											break
										}
									} else if ec2 == ']' {
										key = string(bs[p+4 : z])
										end = z + 3
										break
									}
								}
							}
						}
					}
				}
			}
		}
		if end > 0 {
			//			L.Describe(`part`, string(bs[start:p]))
			//			L.Describe(`key`, string(key))
			dup.Parts = append(dup.Parts, bs[start:p])
			// trim whitespace before and after key
			dup.Keys = append(dup.Keys, S.Trim(key))
			start = end
			p = end
		} else {
			p += 1
		}
	}
	if len(dup.Parts) == 0 {
		//		L.Describe(`part (all) not shown`, len(bs[:]))
		dup.Parts = append(dup.Parts, bs[:])
	} else {
		//		L.Describe(`part`, string(bs[start:]))
		dup.Parts = append(dup.Parts, bs[start:])
	}
	//	L.Describe(`end parsing`, t.Filename, len(t.Parts), len(t.Keys))
	return &dup, nil
}

// parse a file and cache it
func ParseFile(auto_reload, print_debug bool, filename string) (*TemplateChain, error) {
	res := TemplateChain{}
	res.AutoRefresh = auto_reload
	res.PrintDebug = print_debug
	res.Filename = filename
	return res.Reload()
}
