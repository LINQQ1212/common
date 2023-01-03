package words

import (
	"math/rand"
	"os"
	"strings"
	"time"
)

func NewWords(f string) *Words {
	w := &Words{}
	//
	b, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}
	arr := strings.Split(string(b), "\n")
	for _, s := range arr {
		s = strings.TrimSpace(s)
		if s != "" {
			w.list = append(w.list, s)
			w.count++
		}
	}
	return w
}

type Words struct {
	list  []string
	count int
}

func (w *Words) GetRand() string {
	return w.list[rand.Int()%w.count]
}

func (w *Words) GetByIndex(id uint64, index int) string {
	return w.list[(int(id)+index)%w.count]
}

func (w *Words) GetById(id uint64) string {
	return w.list[int(id)%w.count]
}

func (w *Words) GetArrayByNum(index int, num int) []string {
	var strs []string
	for i := 1; i <= num; i++ {
		n := index * i
		if n < 0 {
			n = ^n
		}
		strs = append(strs, w.list[n%w.count])
	}
	return strs
}

func (w *Words) GetByTime(id uint64) string {
	return w.list[(int(id)+time.Now().Day())%w.count]
}

func (w *Words) GetByTimeDay(id uint64) string {
	return w.list[(int(id)+time.Now().Day())%w.count]
}

func (w *Words) GetByTimeHour(id uint64) string {
	return w.list[(int(id)+time.Now().Hour())%w.count]
}
func (w *Words) GetByTimeMinute(id uint64) string {
	return w.list[(int(id)+time.Now().Minute())%w.count]
}

func (w *Words) GetByTimeWeek(id uint64) string {
	_, Week := time.Now().ISOWeek()
	return w.list[(int(id)+Week)%w.count]
}

func (w *Words) GetName(pid uint64, index int, num int) string {
	if num > 1000 {
		num = num + index
	}
	id := (pid % uint64(num)) * uint64(index)
	if id >= uint64(w.count) {
		id = id % uint64(w.count)
	}
	if id < 0 {
		id = 0
	}
	return w.list[id]
}
