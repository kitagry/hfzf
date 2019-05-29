package main

import "log"

func fuzzyFind(keyword string, data map[interface{}]interface{}) map[interface{}]interface{} {
	threshold := len([]rune(text))

	result := make(map[interface{}]interface{})
	for k, value := range data {
		if wordsScore(k.(string), keyword) >= threshold {
			result[k] = value
			continue
		}

		switch v := value.(type) {
		case string:
			if wordsScore(v, keyword) >= threshold {
				result[k] = v
			}
		case []interface{}:
			var tmpData []interface{}
			for _, el := range v {
				if wordsScore(el.(string), keyword) >= threshold {
					tmpData = append(tmpData, el)
				}
			}

			if len(tmpData) != 0 {
				result[k] = tmpData
			}
		case map[interface{}]interface{}:
			tmp := fuzzyFind(keyword, v)
			if len(tmp) != 0 {
				result[k] = tmp
			}
		default:
			log.Println(v)
		}
	}
	return result
}

func wordsScore(s1, s2 string) int {
	matrix, maxI, maxJ := smithWaterman(s1, s2)
	return matrix[maxI][maxJ]
}

func pointPlace(s1, s2 string) []int {
	place := make([]int, 0)
	matrix, i, j := smithWaterman(s1, s2)
	if matrix[i][j] < len([]rune(s2)) {
		return place
	}

	for i > 0 {
		if matrix[i][j] == matrix[i-1][j] {
			i -= 1
		} else if matrix[i][j] == matrix[i][j-1] {
			j -= 1
		} else if matrix[i][j] == matrix[i-1][j-1]+1 {
			i -= 1
			j -= 1
			place = append([]int{i}, place...)
		} else {
			i -= 1
			j -= 1
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

func in(el int, array []int) bool {
	for _, ar := range array {
		if ar == el {
			return true
		}
	}
	return false
}
