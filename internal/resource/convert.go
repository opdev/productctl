package resource

import "encoding/json"

// JSONConvert aids in the marshaling of the various like types, Generated
// graphQL types and on-disk types. This is done by round tripping through JSON,
// and so it's critical that the JSON struct tags match across the two types.
//
// The caller is required to validate that the content is marshaled as expected.
func JSONConvert[R any](in any) (R, error) {
	var out R

	b, err := json.Marshal(in)
	if err != nil {
		return out, err
	}

	err = json.Unmarshal(b, &out)
	if err != nil {
		return out, err
	}

	return out, nil
}
