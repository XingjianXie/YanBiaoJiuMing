package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
	"unicode"
)

var (
	seed          int64
	size          int
	inputPath     string
	outputPath    string
)

func init() {
	// 给flag绑定对应的变量来自动处理程序的命令行参数
	flag.Int64Var(&seed, "s", 0, "The seed for shuffle, 0 represents a time related seed")
	flag.IntVar(&size, "z", 3, "The maximum length of characters to shuffle at once")
	flag.StringVar(&inputPath, "i", "", "Input file path, empty represents stdin")
	flag.StringVar(&outputPath, "o", "", "Input file path, empty represents stdout")
}

func newRand() *rand.Rand {
	var r *rand.Rand
	if seed == 0 {
		r = rand.New(rand.NewSource(time.Now().UnixNano())) // 以当前时间作为随机数种子
	} else {
		r = rand.New(rand.NewSource(seed)) // 以命令行参数作为随机数种子
	}
	return r
}

func openFile(path string, empty *os.File, flag int, name string) *os.File {
	var f *os.File
	if path == "" {
		f = empty
	} else {
		var err error
		f, err = os.OpenFile(path, flag, 0) // 打开文件
		if err != nil {
			log.Fatal("failed to open " + name + " file") // 错误处理
		}
	}
	return f
}

func readFile(file *os.File) []rune {
	bytes, err := ioutil.ReadAll(file) // 标准输入和通过OpenFile()打开的文件同样是*os.File，可以统一通过该方式读取
	if err != nil {
		log.Fatal("failed to read file") // 错误处理
	}
	return []rune(string(bytes))
	// 通过string作为中间类型分两次把[]bytes转换成[]rune
	// 目的是确保在处理unicode字符时每次取出一个完整字符
}

func writeFile(file *os.File, content []rune) {
	_, err := file.Write([]byte(string(content))) // 标准输出和通过OpenFile()打开的文件同样是*os.File，可以统一通过该方式写入
	if err != nil {
		log.Fatal("failed to write file") // 错误处理
	}
}

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
	case '是', '和', '与', '我', '你', '她', '会', '不', '好', '用', '有', '的', '在', '这', '那', '还', '又', '而', '之':
		// 排除句子的结构词，防止读不懂
		return true
	}
	return false
}

func shuffle(content []rune, r *rand.Rand) { // 因为数组切片以引用传递，不需要通过另外返回数组，直接在原数组上操作
	lowerBound := 0
	for upperBound, char := range content {
		if !unicode.Is(unicode.Han, char) || // 通过unicode包判断是否在汉字区域
			!unicode.IsLetter(char) || // 是否是字母（也就是说不是标点）
			structureChar(char) { // 结构词
			shuffleSlice(content[lowerBound:upperBound],  r) // 数组的slice是左开右闭，因此不会算到当前这个不该打乱的字符
			lowerBound = upperBound + 1 // 防止下次把不该打乱的字符打乱
		} else if upperBound - lowerBound >= size { // 最长打乱长度
			shuffleSlice(content[lowerBound:upperBound],  r) // 同上
			lowerBound = upperBound // 与上面不同，因为长度限制造成的终止打乱并不需要考虑当前字符的会不会在下次被打乱
	    }
	}
}

func main() {
	flag.Parse() // 处理程序的命令行参数
	r := newRand() // 获取一个rand.Rand实例
	inputFile := openFile(inputPath, os.Stdin, os.O_RDONLY, "input") // O_RDONLY代表只读
	outputFile := openFile(outputPath, os.Stdout, os.O_WRONLY, "output") // O_WRONLY代表只写
	content := readFile(inputFile) // 用刚刚的文件进行读取（如果是标准输入，也以同样方式读取）
	shuffle(content, r) // 执行打乱
	writeFile(outputFile, content) // 写入文件（如果是标准输出，也以同样方式写入）
}