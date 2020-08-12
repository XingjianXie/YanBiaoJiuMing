package cmd

import (
	"flag"
	"github.com/mark07x/YanBiaoJiuMing/shuffle"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	seed          int64
	size          int
	inputPath     string
	outputPath    string
)

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
		f, err = os.OpenFile(path, flag, os.ModePerm) // 打开文件
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

func init() {
	// 给flag绑定对应的变量来自动处理程序的命令行参数
	flag.Int64Var(&seed, "s", 0, "The seed for shuffle, 0 represents a time related seed")
	flag.IntVar(&size, "z", 3, "The maximum length of characters to shuffle at once")
	flag.StringVar(&inputPath, "i", "", "Input file path, empty represents stdin")
	flag.StringVar(&outputPath, "o", "", "Input file path, empty represents stdout")
}

func Main() {
	flag.Parse() // 处理程序的命令行参数
	r := newRand() // 获取一个rand.Rand实例
	inputFile := openFile(inputPath, os.Stdin, os.O_RDONLY, "input") // O_RDONLY代表只读
	outputFile := openFile(outputPath, os.Stdout, os.O_CREATE | os.O_WRONLY, "output") // 只写，不存在则创建新文件
	content := readFile(inputFile) // 用刚刚的文件进行读取（如果是标准输入，也以同样方式读取）
	shuffle.Shuffle(content, r, size) // 执行打乱
	writeFile(outputFile, content) // 写入文件（如果是标准输出，也以同样方式写入）
}
