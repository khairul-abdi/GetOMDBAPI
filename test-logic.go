package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	anagram := []string{"kita", "atik", "tika", "aku", "kia", "makan", "kua"}

	var (
		arr    []string
		temp   [][]string
		res    [][]string
		status = true
	)

	for _, v := range anagram {
		arr = append(arr, sortString(v))
	}

	for i, v := range arr {
		for j := range temp {
			if arr[i] == temp[j][0] {
				status = false
				temp[j] = append(temp[j], arr[i])
				res[j] = append(res[j], anagram[i])
			}
		}

		if status {
			temp = append(temp, []string{v})
			res = append(res, []string{anagram[i]})
		}
		status = true
	}
	fmt.Println("Hasil ==> ", res)
}

func sortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}
