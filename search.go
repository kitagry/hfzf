package main

import (
	"log"
	"strconv"
)

// FuzzyFind except data which is not near by keyword.
// If parent data is near, the child data is not excepted.
func FuzzyFind(keyword string, data interface{}) interface{} {
	threshold := len([]rune(keyword))

	switch d := data.(type) {
	case map[interface{}]interface{}:
		result := make(map[interface{}]interface{})
		for k, value := range d {
			if wordsScore(k.(string), keyword) >= threshold {
				result[k] = value
				continue
			}

			tmp := FuzzyFind(keyword, value)
			if tmp != nil {
				result[k] = tmp
			}
		}

		if len(result) > 0 {
			return result
		}
	case string:
		if wordsScore(d, keyword) >= threshold {
			return d
		}
	case int:
		if wordsScore(string(d), keyword) >= threshold {
			return string(d)
		}
	case bool:
		if wordsScore(strconv.FormatBool(d), keyword) >= threshold {
			return strconv.FormatBool(d)
		}
	case []interface{}:
		var tmpData []interface{}
		for _, el := range d {
			switch eld := el.(type) {
			case string:
				if wordsScore(eld, keyword) >= threshold {
					tmpData = append(tmpData, eld)
				}
			case map[interface{}]interface{}:
				tmp := FuzzyFind(keyword, eld)
				if tmp != nil {
					tmpData = append(tmpData, tmp)
				}
			}
		}

		if len(tmpData) > 0 {
			return tmpData
		}
	default:
		log.Printf("FuzzyFind method is not supported for %s \n", d)
	}
	return nil
}

func wordsScore(s1, s2 string) int {
	matrix, maxI, maxJ := smithWaterman(s1, s2)
	return matrix[maxI][maxJ]
}

// PointPlace returns places of s1 which is same as s2.
// For example s1 is "abcdef" and s2 is "ace",
// then PointPlace returns [0, 2, 4].
func PointPlace(s1, s2 string) []int {
	place := make([]int, 0)
	matrix, i, j := smithWaterman(s1, s2)
	if matrix[i][j] < len([]rune(s2)) {
		return place
	}

	for i > 0 {
		if matrix[i][j] == matrix[i-1][j] {
			i--
		} else if matrix[i][j] == matrix[i][j-1] {
			j--
		} else if matrix[i][j] == matrix[i-1][j-1]+1 {
			i--
			j--
			place = append([]int{i}, place...)
		} else {
			i--
			j--
		}

		if matrix[i][j] == 0 {
			break
		}
	}

	return place
}

func smithWaterman(s1, s2 string) (matrix [][]int, maxI, maxJ int) {
	s1Rune := []rune(s1)
	s2Rune := []rune(s2)
	gap := 0
	match := 1
	mismatch := 1

	matrix = make([][]int, len(s1Rune)+1)
	for i := 0; i < len(s1Rune)+1; i++ {
		matrix[i] = make([]int, len(s2Rune)+1)
	}

	maxScore := 0
	for i := 1; i < len(s1Rune)+1; i++ {
		for j := 1; j < len(s2Rune)+1; j++ {
			s1Gap := matrix[i-1][j] - gap
			s2Gap := matrix[i][j-1] - gap

			match := matrix[i-1][j-1] + match
			if s1Rune[i-1] != s2Rune[j-1] {
				match = matrix[i-1][j-1] - mismatch
			}

			matrix[i][j] = max(s1Gap, s2Gap, match, 0)
			if maxScore < matrix[i][j] {
				maxI = i
				maxJ = j
				maxScore = matrix[i][j]
			}
		}
	}
	return
}

func max(s ...int) int {
	maxInt := s[0]
	for _, el := range s {
		if el > maxInt {
			maxInt = el
		}
	}
	return maxInt
}
