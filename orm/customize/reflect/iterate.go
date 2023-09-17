package reflect

import (
	"errors"
	"reflect"
)

func IterateArrayOrSlice(entity any) ([]any, error) {
	val := reflect.ValueOf(entity)
	res := make([]any, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		res = append(res, val.Index(i).Interface())
	}
	return res, nil

}
func IterateMap(entity any) ([]any, []any, error) {
	val := reflect.ValueOf(entity)
	entityType := reflect.TypeOf(entity)
	if reflect.Map != entityType.Kind() {
		return nil, nil, errors.New("非Map")
	}
	MapKeys := make([]any, 0, val.Len())
	MapValues := make([]any, 0, val.Len())
	Keys := val.MapKeys()
	for _, key := range Keys {
		MapKeys = append(MapKeys, key.Interface())
		MapValues = append(MapValues, val.MapIndex(key).Interface())
	}
	return MapKeys, MapValues, nil
}

func main() {

}
