package util

func SplitList(list []string) [][]string {
	var res [][]string
	var tempSlice []string
	for i := 0; i < len(list); i++ {
		tempSlice = append(tempSlice, list[i])
		if (i+1)%6 == 0 {
			res = append(res, tempSlice)
			tempSlice = nil
		}

	}
	return res
}
