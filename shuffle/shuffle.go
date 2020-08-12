package shuffle

import (
	"math/rand"
	"unicode"
)

func shuffleSlice(content []rune, r *rand.Rand) { // 打乱当前子串，因为数组切片以引用传递，此处的变更对整个数组有效
	if len(content) == 0 {
		return
	}
	r.Shuffle(len(content), func(i, j int) { // Shuffle函数需要一个交换函数作为参数，在这个交换函数中对当前子串进行操作
		content[i], content[j] = content[j], content[i]
	})
}

func structureChar(char rune) bool {
	switch char {
	case '是', '和', '与',
	     '我', '你', '她',
	     '会', '好', '有',
	     '用', '不', '的',
	     '在', '这', '那',
	     '还', '又', '而',
	     '之', '也', '为':
		// 排除句子的结构词，防止读不懂
		return true
	}
	return false
}

func Shuffle(content []rune, r *rand.Rand, size int) { // 因为数组切片以引用传递，不需要通过另外返回数组，直接在原数组上操作
	lowerBound := 0
	for upperBound, char := range content {
		if !unicode.Is(unicode.Han, char) || // 通过unicode包判断是否在汉字区域
			!unicode.IsLetter(char) || // 是否是字母（也就是说不是标点）
			structureChar(char) { // 结构词
			shuffleSlice(content[lowerBound:upperBound],  r) // 数组的slice是左开右闭，因此不会算到当前这个不该打乱的字符
			lowerBound = upperBound + 1                      // 防止下次把不该打乱的字符打乱
		} else if upperBound - lowerBound >= size { // 最长打乱长度
			shuffleSlice(content[lowerBound:upperBound],  r) // 同上
			lowerBound = upperBound                          // 与上面不同，因为长度限制造成的终止打乱并不需要考虑当前字符的会不会在下次被打乱
	    }
	}
}

