package data

type Item map[string]interface{}

func (i *Item) Valid() bool {
	return i != nil
}
