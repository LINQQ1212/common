package vhost

import (
	"errors"
	pb "github.com/LINQQ1212/common/grpc/grpc_server"
	"github.com/LINQQ1212/common/utils"
	"github.com/cornelk/hashmap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func New(dir, grpc string) *V {
	v := &V{
		dir:  dir,
		grpc: grpc,
	}
	dir2 := path.Join(v.dir, "vhost")
	if ok, _ := utils.PathExists(dir2); !ok {
		os.Mkdir(dir2, 0766)
	}
	v.Reload()
	return v
}

type V struct {
	dir  string
	grpc string
	list *hashmap.Map[string, *hashmap.Map[string, *Site]]
}

func (v *V) Reload() {
	v.list = hashmap.New[string, *hashmap.Map[string, *Site]]()
	dirs, err := os.ReadDir(path.Join(v.dir, "vhost"))
	if err != nil {
		panic(err)
		return
	}

	for i := 0; i < len(dirs); i++ {
		version := dirs[i].Name()
		fs, err := filepath.Glob(path.Join(v.dir, "vhost", version, "/*.yaml"))
		if err != nil {
			panic(err)
			continue
		}

		v.list.Set(version, hashmap.New[string, *Site]())
		for _, f := range fs {
			s := &Site{}
			err = utils.ReadYamlConfig(f, s)
			if err != nil {
				panic(err)
				continue
			}
			s.InitTableMap()
			domain := strings.TrimSuffix(filepath.Base(f), ".yaml")
			v.Set(version, domain, s)
		}
	}
}

func (v *V) NewVersion(domain string) error {
	_, ok := v.list.Get(domain)
	if !ok {
		v.list.Set(domain, hashmap.New[string, *Site]())
		return os.MkdirAll(path.Join(v.dir, "vhost", domain), os.ModePerm)
	}
	return nil
}

func (v *V) Set(version, domain string, s *Site) {
	l, ok := v.list.Get(version)
	if !ok {
		l = hashmap.New[string, *Site]()
		v.list.Set(version, l)
	}
	s.Domain = domain
	if s.Index == -1 {
		s.Index = l.Len()
	}
	l.Set(domain, s)
}

func (v *V) Get(version, domain string) (*Site, bool) {
	l, ok := v.list.Get(version)
	if !ok {
		return nil, false
	}
	return l.Get(domain)
}

func (v *V) GetOrNew(version, domain, f string) (*Site, error) {
	l, ok := v.list.Get(version)
	if !ok {
		return nil, errors.New("404")
	}
	s, ok := l.Get(domain)
	if !ok {
		var err error
		s, err = v.NewSite(version, domain, f)
		if err != nil {
			return nil, err
		}
		s.Index = l.Len() + 1
		l.Set(domain, s)
		err = utils.SaveYamlConfig(path.Join(v.dir, "vhost", version, domain+".yaml"), s)
		return s, err
	}
	return s, nil
}

var ErrTemplateArticleNil = errors.New("not has template")
var ErrTemplateListNil = errors.New("not has list template")

func (v *V) NewSite(version, domain, f string) (*Site, error) {
	s := &Site{}
	s.Domain = domain
	s.Index = -1
	s.Num = rand.Intn(3000)
	s.Num2 = uint64(rand.Int63n(8000))
	s.Table = utils.RandChars()
	s.Key = s.Table[0:32]
	s.InitTableMap()
	s.Path = utils.GetPath()
	s.CPath = utils.GetCategoryPath()
	s.Suffix = GetSuffix()
	s.IsParam = strings.Contains(f, ".php")
	s.SearchPath = GetSearch()
	s.Type = 1
	s.InitHtml()
	s.MinLength = utils.RandInt(3, 12)
	if rand.Float64() > 0.5 {
		s.Type = 2
	}
	tempDir := path.Join(v.dir, "views", version, s.Domain)
	err := utils.CreateDir(tempDir)
	if err != nil {
		return s, err
	}
	err = v.getViews(s, tempDir)
	return s, err
}

func (v *V) getViews(s *Site, dir string) error {
	conn, err := grpc.Dial(v.grpc, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()
	conn2 := pb.NewGreeterClient(conn)

	res, err := conn2.GetListView(context.Background(), &pb.NullRequest{})
	if err != nil {
		return err
	}
	if len(res.Data) == 0 {
		return ErrTemplateListNil
	}

	err = os.WriteFile(path.Join(dir, "list.jet"), res.Data, 0644)
	if err != nil {
		return err
	}
	s.ListTemp = "list.jet"

	ares, err := conn2.GetArticleView(context.Background(), &pb.ArticleReq{Num: 1})
	if err != nil {
		return err
	}
	al := len(ares.Datas)
	if al == 0 {
		return ErrTemplateArticleNil
	}
	for i := 0; i < al; i++ {
		name := "a.jet"
		err = os.WriteFile(path.Join(dir, name), ares.Datas[i], 0644)
		if err != nil {
			return err
		}
		s.Temps = append(s.Temps, name)
	}
	s.TempsLen = len(s.Temps)
	return nil
}
