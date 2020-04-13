package module

func Legalletter(mtype MType) bool {
	if _, ok := legalletterMap[mtype]; ok {
		return true
	}
	return false
}
