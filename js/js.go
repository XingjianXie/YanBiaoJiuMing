package js

import (
	"github.com/mark07x/YanBiaoJiuMing/shuffle"
	"honnef.co/go/js/dom"
	"math/rand"
	"strconv"
	"time"
)

func Main() {
	d := dom.GetWindow().Document()
	d.QuerySelector(".shuffle").AddEventListener("click", false, func(event dom.Event) {
		content := []rune(d.QuerySelector(".content").(*dom.HTMLInputElement).Value)
		seed, err := strconv.ParseInt(d.QuerySelector(".seed").(*dom.HTMLInputElement).Value, 10, 64)
		size, err := strconv.Atoi(d.QuerySelector(".size").(*dom.HTMLInputElement).Value)
		if seed == 0 || err != nil {
			seed = time.Now().UnixNano()
		}
		if size == 0 || err != nil {
			size = 3
		}
		r := rand.New(rand.NewSource(seed))
		shuffle.Shuffle(content, r, size) // 执行打乱
		d.QuerySelector(".result").(*dom.HTMLInputElement).Value = string(content)
	})
}
