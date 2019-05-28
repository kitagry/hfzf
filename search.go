package main

import "log"

func fuzzyFind(keyword string, data map[interface{}]interface{}) map[interface{}]interface{} {
	threshold := len([]rune(text))

	result := make(map[interface{}]interface{})
	for k, value := range data {
		if smithWaterman(k.(string), keyword) >= threshold {
			result[k] = value
			continue
		}

		switch v := value.(type) {
		case string:
			if smithWaterman(v, keyword) >= threshold {
				result[k] = v
			}
		case []interface{}:
			var tmpData []interface{}
			for _, el := range v {
				if smithWaterman(el.(string), keyword) >= threshold {
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

func smithWaterman(s1, s2 string) int {
	s1Rune := []rune(s1)
	s2Rune := []rune(s2)
	gap := 0
	match := 1
	mismatch := 1

	matrix := make([][]int, len(s1Rune)+1)
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
			maxScore = max(maxScore, matrix[i][j])
		}
	}

	return maxScore
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
