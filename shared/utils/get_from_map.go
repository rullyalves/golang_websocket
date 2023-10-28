package utils

func GetFromMap[T any](maps map[string]*any, key string) *T {
	fieldPtr := maps[key]
	if fieldPtr != nil {
		x := (*fieldPtr).(T)
		return &x
	}
	return nil
}
