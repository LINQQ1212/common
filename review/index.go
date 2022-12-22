package review

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/LINQQ1212/common/utils"
	"go.etcd.io/bbolt"
	"io"
	"os"
	"time"
)

func NewReview(path, namePath string) (*Review, error) {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	fp, err := os.Open(namePath)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin(false)
	r := &Review{
		db:    db,
		Count: 0,
		B:     tx.Bucket([]byte("review")),
	}
	b := r.B.Get([]byte("count"))
	if b == nil {
		panic(errors.New("review count is null"))
	}
	r.Count = utils.Btoi(b)
	rd := bufio.NewReader(fp)
	for {
		lineBytes, err := rd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		s := string(bytes.TrimSpace(lineBytes))
		if s != "" {
			r.Names = append(r.Names, s)
			r.NCount++
		}

	}

	return r, err
}

type Review struct {
	db     *bbolt.DB
	Count  uint64
	B      *bbolt.Bucket
	Names  []string
	NCount uint64
}

type Info struct {
	Name   string
	Detail string
	Score  float64
	Time   time.Time
}

func (r *Review) Close() error {
	return r.db.Close()
}

func (r *Review) GetReviews(id uint64, index int, size int) []string {
	var data []string
	n1 := uint64(index) + id
	num := n1 % r.Count
	for i := 0; i < size; i++ {
		if num >= r.Count {
			num = 0
		}
		d := string(r.B.Get(utils.Itob(num)))
		if d != "" {
			data = append(data, d)
		}
		num++
	}
	return data
}

func (r *Review) GetName(pid uint64, index int, num int) string {
	if num > 1000 {
		num = num + index
	}
	id := (pid % uint64(num)) * uint64(index)
	if id >= r.NCount {
		id = id - r.NCount
	}
	if id < 0 {
		id = 0
	}
	return r.Names[id]
}

func (r *Review) Get(id uint64, index int, size int) []Info {
	tt, _ := time.Now().ISOWeek()
	n1 := uint64(tt+index) + id
	num := n1 % r.Count
	num2 := n1 % r.NCount
	var data []Info

	for i := 0; i < size; i++ {
		if num >= r.Count {
			num = 0
		}
		if num2 >= r.NCount {
			num2 = 0
		}
		n2 := tt + index
		d := string(r.B.Get(utils.Itob(num)))
		if d != "" {
			l := n2 + len(d)
			t := (time.Hour * time.Duration(n2%7)) + (time.Minute * time.Duration(l%60)) + (time.Second * time.Duration(l%30))
			data = append(data, Info{
				Name:   r.Names[num2],
				Detail: d,
				Time:   time.Now().Add(-t),
				Score:  4 + 0.3 + float64((l+n2)%7)/10,
			})
		}
		num++
		num2++
	}
	return data
}
