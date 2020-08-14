package util

func AppendSuffix(list []string, suffix string) []string {
	var res []string
	for _, n := range list {
		res = append(res, n+suffix)
	}
	return res
}

func PrependPrefix(list []string, prefix string) []string {
	var res []string
	for _, n := range list {
		res = append(res, prefix+n)
	}
	return res

}
