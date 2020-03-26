package filter

func Filter(text string, replace string) (rText string, isReplace bool) {
	return trieFilter.Check(text, replace)
}
