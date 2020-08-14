package util

//SplitList reduced list into smaller chunks.
func SplitList(list []string, size int) [][]string {
	var res [][]string
	var tempSlice []string
	for i := 0; i < len(list); i++ {
		tempSlice = append(tempSlice, list[i])
		if size%6 == 0 {
			res = append(res, tempSlice)
			tempSlice = nil
		}

	}
	return res
}
