package M

import (
	"reflect"
	"sync"

	"github.com/mitchellh/mapstructure"
	msgpack2 "github.com/shamaton/msgpack/v2"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
)

var structTypeMutex = sync.Mutex{}
var structTypeCache = map[string]*StructMapper{}

type FieldTag string

const (
	RawFieldName   FieldTag = ``
	SnakeFieldName FieldTag = `SNAKE`
	CamelFieldName FieldTag = `CAMEL`
	AllFieldName   FieldTag = `ALL`
)

type StructMapper struct {
	StructName    string
	Offset2key    map[uintptr]string
	Key2offset    map[string]uintptr
	key2fieldName map[string]string
}

func (sm *StructMapper) MapToStruct(m SX, s any) {
	value := reflect.ValueOf(s).Elem()
	if !value.IsValid() || value.Kind() != reflect.Struct {
		L.Print(`StructMapper.MapToStruct: invalid type`, s)
		return
	}

	sTyp := value.Type()
	structName := sTyp.String()
	if sm.StructName != structName {
		L.Print(`StructMapper.StructToMap: different struct type`, sm.StructName, structName)
		return
	}

	for k, v := range m {
		fieldName := sm.key2fieldName[k]
		fPtr := value.FieldByName(fieldName)
		if !fPtr.CanSet() || !fPtr.IsValid() {
			continue
		}
		if v == nil {
			continue
		}
		fPtr.Set(reflect.ValueOf(v))
	}
}

func (sm *StructMapper) StructToMap(s any) (m SX) {
	m = SX{}
	value := reflect.ValueOf(s).Elem()
	if !value.IsValid() || value.Kind() != reflect.Struct {
		L.Print(`StructMapper.StructToMap: invalid type`, s)
		return
	}

	sTyp := value.Type()
	structName := sTyp.String()
	if sm.StructName != structName {
		L.Print(`StructMapper.StructToMap: different struct type`, sm.StructName, structName)
		return
	}

	for i := 0; i < value.NumField(); i++ {
		field := sTyp.Field(i)
		if S.FirstIsLower(field.Name) { // skip unexported
			continue
		}
		val := value.Field(i)
		m.Set(sm.Offset2key[field.Offset], val.Interface())
	}
	return
}

// ParseStruct convert struct to structMapper
func ParseStruct(s any, tag FieldTag) (sm *StructMapper) {
	sm = &StructMapper{}
	sm.Offset2key = map[uintptr]string{}
	sm.Key2offset = map[string]uintptr{}
	sm.key2fieldName = map[string]string{}
	value := reflect.ValueOf(s).Elem()
	if !value.IsValid() {
		L.Print(`ParseStruct: invalid: invalid type`, s)
		return
	}
	sm.StructName = value.Type().String()

	sTyp := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := sTyp.Field(i)
		if S.FirstIsLower(field.Name) { // skip unexported
			continue
		}
		switch tag {
		case SnakeFieldName:
			key := S.SnakeCase(field.Name)
			sm.Offset2key[field.Offset] = key
			sm.Key2offset[key] = field.Offset
			sm.key2fieldName[key] = field.Name
		case CamelFieldName:
			key := S.CamelCase(field.Name)
			sm.Offset2key[field.Offset] = key
			sm.Key2offset[key] = field.Offset
			sm.key2fieldName[key] = field.Name
		case RawFieldName:
			sm.Offset2key[field.Offset] = field.Name
			sm.Key2offset[field.Name] = field.Offset
			sm.key2fieldName[field.Name] = field.Name
		case AllFieldName:
			sm.Offset2key[field.Offset] = field.Name
			sm.Key2offset[field.Name] = field.Offset
			sm.key2fieldName[field.Name] = field.Name
			key := S.SnakeCase(field.Name)
			sm.Key2offset[key] = field.Offset
			sm.key2fieldName[key] = field.Name
			key = S.CamelCase(field.Name)
			sm.Key2offset[key] = field.Offset
			sm.key2fieldName[key] = field.Name
		default:
			key := S.LeftOf(field.Tag.Get(string(tag)), `,`)
			if key != `` {
				sm.Offset2key[field.Offset] = key
				sm.Key2offset[key] = field.Offset
				sm.key2fieldName[key] = field.Name
			}
		}
	}
	return
}

// FromStruct convert any struct to map
func FromStruct(srcStructPtr any) SX {
	return StructMap(srcStructPtr).StructToMap(srcStructPtr)
}

// StructMap get or create a struct mapper
func StructMap(structPtr any) *StructMapper {
	structType := reflect.TypeOf(structPtr).String()
	sm, ok := structTypeCache[structType]
	if !ok {
		sm = ParseStruct(structPtr, AllFieldName)
		structTypeMutex.Lock()
		structTypeCache[structType] = sm
		structTypeMutex.Unlock()
	}
	return sm
}

// ToStruct convert to struct
func (m SX) ToStruct(targetStructPtr any) {
	StructMap(targetStructPtr).MapToStruct(m, targetStructPtr)
}

// FastestMapToStruct only for exact match of field name and map key
func FastestMapToStruct(m any, s any) {
	b, _ := msgpack.Marshal(m)
	_ = msgpack.Unmarshal(b, s)
}

// FastestStructToMap using struct's field name as map key
func FastestStructToMap(s any) (m map[string]any) {
	_ = mapstructure.Decode(s, &m)
	return
}

// FastestStructToStruct
func FastestStructToStruct(src any, dst any) {
	b, _ := msgpack2.Marshal(src)
	_ = msgpack.Unmarshal(b, dst)
}
