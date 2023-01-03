package create_db

import (
	"archive/tar"
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	textrank "github.com/DavidBelicza/TextRank/v2"
	"github.com/DavidBelicza/TextRank/v2/convert"
	"github.com/DavidBelicza/TextRank/v2/parse"
	"github.com/DavidBelicza/TextRank/v2/rank"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"regexp"
	"runtime"

	"github.com/LINQQ1212/common/global"
	pb "github.com/LINQQ1212/common/grpc/grpc_server"
	"github.com/LINQQ1212/common/models"
	"github.com/LINQQ1212/common/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zip"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func New(req models.NewVersionReq) *Create {
	return &Create{
		Info:       req,
		domainInfo: &sync.Map{},
		CC:         make(chan [2][]byte, 30),
		S:          make(chan [3][]byte, 30),
		done:       make(chan struct{}, 1),
		wg:         &sync.WaitGroup{},
		cateInfo:   &sync.Map{},
		cates:      &sync.Map{},
	}
}

type Create struct {
	Info       models.NewVersionReq
	domainInfo *sync.Map //map[string]uint64{}
	cateInfo   *sync.Map //map[string]*models.Cate
	cates      *sync.Map //map[string]*models.Cate

	db *bbolt.Tx
	pb *bbolt.Bucket
	//btx          *bbolt.Tx
	wg       *sync.WaitGroup
	pId      uint64
	domainId uint64
	cateId   uint64
	conn2    pb.GreeterClient
	//txn          *badger.Txn
	googleImgZip  *zip.ReadCloser
	yahooDscZip   *zip.ReadCloser
	bingDscZip    *zip.ReadCloser
	youtobeDscZip *zip.ReadCloser
	CC            chan [2][]byte
	S             chan [3][]byte
	done          chan struct{}
}

func (c *Create) CreateBucketIfNotExists(db *bbolt.DB) error {
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(models.BInfo)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(models.BCate)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(models.BDomain)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(models.BProduct)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(models.BProductInfo)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(models.BCateProducts)
		if err != nil {
			return err
		}
		return nil
	})
}

func (c *Create) Start() error {
	if c.Info.DownMainPic {
		conn, err := grpc.Dial(global.CONFIG.System.ImageGrcp, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return err
		}
		defer conn.Close()
		c.conn2 = pb.NewGreeterClient(conn)
	}
	//c.txn = global.IMGDB.NewTransaction(true)
	global.LOG.Info(c.Info.Domain + " start")
	global.LOG.Info(c.Info.Domain + " start 22")
	if c.Info.GoogleImg != "" {
		var err error
		c.googleImgZip, err = zip.OpenReader(c.Info.GoogleImg)
		if err != nil {
			global.LOG.Error("GoogleImg", zap.Error(err))
			return err
		}
		defer c.googleImgZip.Close()
	}

	if c.Info.YahooDsc != "" {
		var err error
		c.yahooDscZip, err = zip.OpenReader(c.Info.YahooDsc)
		if err != nil {
			global.LOG.Error("YahooDsc", zap.Error(err))
			return err
		}
		defer c.yahooDscZip.Close()
	}

	if c.Info.BingDsc != "" {
		var err error
		c.bingDscZip, err = zip.OpenReader(c.Info.BingDsc)
		if err != nil {
			global.LOG.Error("BingDsc", zap.Error(err))
			return err
		}
		defer c.bingDscZip.Close()
	}
	if c.Info.YoutubeDsc != "" {
		var err error
		c.youtobeDscZip, err = zip.OpenReader(c.Info.YoutubeDsc)
		if err != nil {
			global.LOG.Error("YoutubeDsc", zap.Error(err))
			return err
		}
		defer c.youtobeDscZip.Close()
	}

	db, err := bbolt.Open(path.Join(global.VersionDir, "~"+c.Info.Domain+".db"), 0666, bbolt.DefaultOptions)
	//db, err := badger.Open(badger.DefaultOptions(path.Join(global.VersionDir, "~"+c.domain)).WithLoggingLevel(badger.ERROR))
	//db, err := storm.Open(path.Join(global.VersionDir, "~"+c.domain+".db"), storm.Codec(protobuf.Codec))
	if err != nil {
		global.LOG.Error("Open bbolt db", zap.Error(err))
		return err
	}
	err = c.CreateBucketIfNotExists(db)
	if err != nil {
		global.LOG.Error("CreateBucketIfNotExists", zap.Error(err))
		return err
	}

	go func() {
		var tx *bbolt.Tx
		tx, err = db.Begin(true)
		if err != nil {
			global.LOG.Error("SaveProduct", zap.Error(err))
			return
		}
		index := 0
		for bs := range c.S {
			if index >= 10000 {
				index = 0
				tx.Commit()
				tx, err = db.Begin(true)
				if err != nil {
					global.LOG.Error("SaveProduct", zap.Error(err))
				}
			}
			if err = tx.Bucket(models.BProduct).Put(bs[0], bs[1]); err != nil {
				global.LOG.Error("SaveProduct", zap.Error(err))
			}
			if err = tx.Bucket(models.BProductInfo).Put(bs[0], bs[2]); err != nil {
				global.LOG.Error("SaveProductInfo", zap.Error(err))
			}
			index++
		}
		if index > 0 {
			err = tx.Commit()
			if err != nil {
				global.LOG.Error("SaveProduct Commit", zap.Error(err))
				return
			}
		}

		c.done <- struct{}{}
	}()

	f, err := os.Open(c.Info.ProductTarLink)
	if err != nil {
		global.LOG.Error("Open："+c.Info.ProductTarLink, zap.Error(err))
		return err
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		global.LOG.Error("gzip.NewReader", zap.Error(err))
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	var h *tar.Header
	c.RunCC()
	for {
		h, err = tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			global.LOG.Error("tr.Next()", zap.Error(err))
			break
		}
		c.one(h.Name, tr)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	close(c.CC)
	c.wg.Wait()
	close(c.S)
	<-c.done
	/*err = tx.Commit()
	if err != nil {
		return err
	}*/
	global.LOG.Info(c.Info.Domain + "product end")
	vinfo := models.VersionInfo{
		Name:     c.Info.Domain,
		FileName: strings.TrimSuffix(filepath.Base(c.Info.ProductTarLink), ".tar.gz"),
		Count:    c.pId,
		DownPic:  c.Info.DownMainPic,
	}

	c.db, err = db.Begin(true)
	if err != nil {
		global.LOG.Error("Begin", zap.Error(err))
		return err
	}
	vb, err := proto.Marshal(&vinfo)
	if err != nil {
		global.LOG.Error("vinfo", zap.Error(err))
		return err
	}
	err = c.db.Bucket(models.BInfo).Put(models.BInfo, vb)
	if err != nil {
		global.LOG.Error("BInfo", zap.Error(err))
		return err
	}
	err = c.db.Bucket(models.BCate).Put(models.BCate, c.GetCateByte())
	if err != nil {
		global.LOG.Error("BCate", zap.Error(err))
		return err
	}
	err = c.db.Bucket(models.BDomain).Put(models.BDomain, c.GetDomainByte())
	if err != nil {
		global.LOG.Error("BDomain", zap.Error(err))
		return err
	}
	err = c.db.Commit()
	if err != nil {
		global.LOG.Error("Commit", zap.Error(err))
		return err
	}
	err = db.Close()
	if err != nil {
		global.LOG.Error("Close", zap.Error(err))
		return err
	}
	global.LOG.Info(c.Info.Domain + " end")
	runtime.GC()
	os.Rename(path.Join(global.VersionDir, "~"+c.Info.Domain+".db"), path.Join(global.VersionDir, c.Info.Domain+".db"))
	return nil
}

func (c *Create) GetDomain() string {
	return c.Info.Domain
}

func (c *Create) GetDomainByte() []byte {
	data := map[uint64]string{}
	c.domainInfo.Range(func(s, id any) bool {
		data[id.(uint64)] = s.(string)
		return true
	})
	b, _ := json.Marshal(&data)
	return b
}

func (c *Create) getDomainId(domain string) (uint64, error) {
	if v, ok := c.domainInfo.Load(domain); ok {
		return v.(uint64), nil
	}
	c.domainInfo.Store(domain, atomic.AddUint64(&c.domainId, 1))
	return c.domainId, nil
}

func (c *Create) GetCateByte() []byte {
	var data []*models.Cate
	//tx := c.db.Bucket(models.BCateProducts)
	c.cates.Range(func(key, v any) bool {
		cate := v.(*models.Cate)
		cate.Count = uint64(len(cate.Products))
		/*if cate.Count > 0 {
			b, err := tx.CreateBucketIfNotExists(utils.Itob(cate.ID))
			if err != nil {
				global.LOG.Error("Cate Save", zap.Error(err))
				return false
			}
			for i := 0; i < len(cate.Products); i++ {

				err = b.Put(utils.Itob(uint64(i+1)), utils.Itob(cate.Products[i]))
				if err != nil {
					global.LOG.Error("Cate id Save", zap.Error(err))
					return false
				}
			}
		}
		cate.Products = nil
		*/
		data = append(data, cate)
		return true
	})
	b, _ := json.Marshal(&data)
	return b
}

func (c *Create) getCateId(cates string, DomainID uint64) (*models.Cate, error) {
	if v, ok := c.cateInfo.Load(cates + ","); ok {
		return v.(*models.Cate), nil
	}
	arr := strings.Split(cates, ",")
	var pid uint64
	var cate *models.Cate
	names := ""
	for _, s := range arr {
		names += s + ","
		var cate2 *models.Cate
		cate3, ok := c.cates.Load(names)
		if !ok {

			cate2 = &models.Cate{}

			arr2 := strings.Split(s, ":")
			name := strings.TrimSpace(arr2[0])
			if len(arr2) > 1 {
				cate2.CategoryId = arr2[1]
			}
			cate2.Name = name
			cate2.ID = atomic.AddUint64(&c.cateId, 1)
			cate2.DomainID = DomainID
			cate2.ParentId = pid
			c.cates.Store(names, cate2)
		} else {
			cate2 = cate3.(*models.Cate)
		}
		pid = cate2.ID
		cate = cate2
	}
	c.cateInfo.Store(names, cate)
	return cate, nil
}

func (c *Create) RunCC() {
	for i := 0; i < 60; i++ {
		c.wg.Add(1)
		go func() {
			for i2 := range c.CC {
				if err1 := c.handleOneRow(i2[0], i2[1]); err1 != nil {
					fmt.Println(err1)
				}
			}
			c.wg.Done()
		}()
	}

}

func (c *Create) one(fname string, r io.Reader) {
	fname = path.Base(filepath.Base(fname))
	a := strings.Split(fname, "/")
	fname = a[len(a)-1]
	arr := strings.SplitN(fname, "_", 2)
	domain := []byte(arr[0])
	bfRd := bufio.NewReader(r)
	for {
		line, err := bfRd.ReadBytes('\n')
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			break
		}
		c.CC <- [2][]byte{domain, bytes.TrimSpace(line)}
	}
}

var brandReg = regexp.MustCompile("<b>ブランド</b></th><th>(.+?)</th>")
var htmlReg = regexp.MustCompile("<.+?>|</.+?>")

var NameReg = regexp.MustCompile(`(^((★|❤️|❤|☆|『|【|「|\(|✨|★|■|❣|♥|●)(.{1,8})(★|❤️|❤|☆|』|】|」|\)|✨|★|■|❣|♥|●))+)|(新品未使用|新品)(　| |★|❤️|❤|☆|✨|★|■|❣|●)`)
var NameReg2 = regexp.MustCompile(`(★|❤️|❤|☆|✨|★|■|❣|♥|●|\d+$)`)

func (c *Create) handleOneRow(domain []byte, line []byte) error {

	domainId, err := c.getDomainId(string(domain))
	if err != nil {
		return err
	}

	arrb := bytes.Split(line, []byte("|"))
	if len(arrb) < 12 {
		return errors.New("len not 12")
	}
	if len(arrb[1]) == 0 {
		println(string(line))
		return errors.New("len2 error")
	}
	arr := make([]string, len(arrb))
	for i, i2 := range arrb {
		d := make([]byte, base64.StdEncoding.DecodedLen(len(i2)))
		base64.StdEncoding.Decode(d, i2)
		arr[i] = string(bytes.Trim(d, "\x00"))
	}
	arr[4] = strings.ReplaceAll(arr[4], ",http", "|||http")
	arr[4] = strings.ReplaceAll(arr[4], ",//", "|||//")
	// /分类:1,分类2:2|cPath|产品ID|产品型号|产品图片^产品图片2|产品价格|优惠价格|产品名称|产品详情|标题|关键词|描述|

	imgArr := strings.Split(arr[4], "|||")
	for i := 0; i < len(imgArr); i++ {
		if v := strings.Index(imgArr[i], "?"); v > 0 {
			imgArr[i] = imgArr[i][:v]
		}
	}
	mainImg := imgArr[0]
	if len(imgArr) == 1 {
		imgArr = []string{}
	} else {
		imgArr = imgArr[1:]
	}
	/*if c.Info.DownMainPic {
		img2 := c.downPic(mainImg)
		if mainImg != img2 {
			mainImg = "/images/" + img2
		}
	}*/

	p := &models.Product{
		ID:       atomic.AddUint64(&c.pId, 1),
		DomainID: domainId,
		Image:    mainImg,
		Name:     arr[7],
	}
	pi := &models.ProductInfo{
		Pid:          arr[2], //
		Model:        arr[3],
		Images:       imgArr,
		Price:        arr[5],
		Specials:     arr[6],
		Description:  arr[8],
		MTitle:       arr[9],
		MKeywords:    arr[10],
		MDescription: arr[11],
	}

	p.Name = strings.TrimSpace(NameReg.ReplaceAllString(p.Name, ""))
	p.Name = strings.TrimSpace(NameReg2.ReplaceAllString(p.Name, " "))
	p.Name = strings.TrimSpace(strings.ReplaceAll(p.Name, "　", " "))
	p.Name = strings.TrimSpace(strings.ReplaceAll(p.Name, "【送料無料】", ""))
	p.Name = strings.TrimSpace(strings.ReplaceAll(p.Name, "送料無料", ""))

	brandArr := brandReg.FindStringSubmatch(pi.Description)
	if len(brandArr) > 1 {
		p.Brand = brandArr[1]
	}
	pi.Description = strings.ReplaceAll(pi.Description, "<h2>商品の情報</h2>", "")
	pi.Description = strings.ReplaceAll(pi.Description, "<h2>商品情報</h2>", "")
	pi.Description = strings.ReplaceAll(pi.Description, "<h2 class=\"Heading Heading-f\">商品情報</h2>", "")

	pppp := ""
	ttti := strings.Index(pi.Description, "<table border=\"1\">")
	if ttti > 0 {
		pi.Description = pi.Description[0:ttti]
		pppp = htmlReg.ReplaceAllString(pi.Description, " ")
	}

	pi.Description = strings.Replace(pi.Description, "<style>table th{border: 1px solid #ccc !important;}</style>", "", 1)
	pi.Description = strings.TrimSuffix(pi.Description, "</p>")
	pi.Description = strings.TrimPrefix(pi.Description, "<p>")
	pi.Description = htmlReg.ReplaceAllString(pi.Description, "")

	arrkws := GetKeys(p.Name + " " + pi.MDescription + " " + pppp)

	for _, arrkw := range arrkws {
		if len(arrkw) < 12 {
			p.Keywords = append(p.Keywords, arrkw)
		}
	}
	cate, err := c.getCateId(strings.TrimSpace(arr[0]), domainId)
	if err == nil {
		p.CateId = cate.ID
		cate.Products = append(cate.Products, p.ID)
	}

	/*cateArr := strings.Split(arr[0], ",")
	for i := 0; i < len(cateArr); i++ {
		cateArr2 := strings.Split(cateArr[i], ":")
		p.Categories = append(p.Categories, strings.TrimSpace(cateArr2[0]))
	}*/
	fname := string(domain) + "/" + pi.Pid + ".json"
	if c.Info.GoogleImg != "" {
		if fp, err := c.googleImgZip.Open(fname); err == nil {
			if b, err := io.ReadAll(fp); err == nil {
				json.Unmarshal(b, &pi.GoogleImgs)
				pi.GoogleImgs = lo.Shuffle(pi.GoogleImgs)
			}
			fp.Close()
		}
	}

	if c.Info.YoutubeDsc != "" {
		if fp, err := c.youtobeDscZip.Open(fname); err == nil {
			if b, err := io.ReadAll(fp); err == nil {
				json.Unmarshal(b, &pi.Youtube)
				pi.Youtube = lo.Shuffle(pi.Youtube)
			}
			fp.Close()
		}
	}

	if c.Info.YahooDsc != "" {
		if fp, err := c.yahooDscZip.Open(fname); err == nil {
			if b, err := io.ReadAll(fp); err == nil {
				var arr3 []models.YahooDsc
				err = json.Unmarshal(b, &arr3)
				arr3l := len(arr3)
				if err == nil && arr3l > 0 {
					for i := 0; i < arr3l; i++ {
						if arr3[i].Des != "" || arr3[i].Title != "" {
							pi.YahooDesc = append(pi.YahooDesc, &arr3[i])
						}
					}
					pi.YahooDesc = lo.Shuffle(pi.YahooDesc)
				}
			}
			fp.Close()
		}
	}

	if c.Info.BingDsc != "" {
		if fp, err := c.bingDscZip.Open(fname); err == nil {
			if b, err := io.ReadAll(fp); err == nil {
				var arr3 []*models.YahooDsc
				err = json.Unmarshal(b, &arr3)
				if err == nil {
					pi.BingDesc = lo.Shuffle(arr3)
				}
			}
			fp.Close()
		}
	}

	if err != nil {
		fmt.Println(err, "---")
		return err
	}
	//c.db.Save(p)
	//c.S <- p
	return c.SaveProduct(p, pi)
}

func (c *Create) SaveProduct(p *models.Product, pi *models.ProductInfo) error {
	var (
		b   []byte
		b2  []byte
		err error
	)

	//b, err := json.Marshal(p)
	if b, err = proto.Marshal(p); err != nil {
		global.LOG.Error("protobuf.Codec Product", zap.Error(err), zap.Any("id", p.ID))
		return err
	}
	if b2, err = proto.Marshal(pi); err != nil {
		global.LOG.Error("protobuf.Codec ProductInfo", zap.Error(err), zap.Any("id", p.ID))
		return err
	}

	c.S <- [3][]byte{utils.Itob(p.ID), b, b2}
	p = nil
	pi = nil
	return nil
}

var ww = []string{"あそこ", "あっ", "あの", "あのかた", "あの人", "あり", "あります", "ある", "あれ", "い", "いう", "います", "いる", "う", "うち", "え", "お", "および", "おり", "おります", "か", "かつて", "から", "が", "き", "ここ", "こちら", "こと", "この", "これ", "これら", "さ", "さらに", "し", "しかし", "する", "ず", "せ", "せる", "そこ", "そして", "その", "その他", "その後", "それ", "それぞれ", "それで", "た", "ただし", "たち", "ため", "たり", "だ", "だっ", "だれ", "つ", "て", "で", "でき", "できる", "です", "では", "でも", "と", "という", "といった", "とき", "ところ", "として", "とともに", "とも", "と共に", "どこ", "どの", "な", "ない", "なお", "なかっ", "ながら", "なく", "なっ", "など", "なに", "なら", "なり", "なる", "なん", "に", "において", "における", "について", "にて", "によって", "により", "による", "に対して", "に対する", "に関する", "の", "ので", "のみ", "は", "ば", "へ", "ほか", "ほとんど", "ほど", "ます", "また", "または", "まで", "も", "もの", "ものの", "や", "よう", "より", "ら", "られ", "られる", "れ", "れる", "を", "ん", "何", "及び", "彼", "彼女", "我々", "特に", "私", "私達", "貴方", "貴方方"}

var rule *parse.RuleDefault
var language *convert.LanguageDefault
var algorithmDef *rank.AlgorithmDefault

func init() {
	rule = textrank.NewDefaultRule()
	// Default Language for filtering stop words.
	language = textrank.NewDefaultLanguage()
	language.SetWords("ja", ww)
	// Active the Spanish.
	language.SetActiveLanguage("ja")

	// Default algorithm for ranking text.
	algorithmDef = textrank.NewDefaultAlgorithm()
}

func GetKeys(txt string) (data []string) {
	tr := textrank.NewTextRank()
	// Add text.
	tr.Populate(txt, language, rule)
	// Run the ranking.
	tr.Ranking(algorithmDef)
	// Get all phrases order by weight.
	tmp := map[string]struct{}{}

	rankedPhrases := textrank.FindPhrases(tr)
	for _, phrase := range rankedPhrases {
		if _, ok := tmp[phrase.Left]; !ok {
			tmp[phrase.Left] = struct{}{}
			data = append(data, phrase.Left)
		}
		if _, ok := tmp[phrase.Right]; !ok {
			tmp[phrase.Right] = struct{}{}
			data = append(data, phrase.Right)
		}
	}
	return
}

/*
func (c *Create) downPic(img string) string {
	req, err := http.NewRequest("GET", img, nil)
	if err != nil {
		return img
	}
	req.Header.Set("accept", "")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return img
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return img
	}
	res2, err := c.conn2.GetFiles(context.Background(), &pb.FilesReq{Image: b})
	if err != nil {
		global.LOG.Error("downpic grpc", zap.Error(err))
		return img
	}
	imgid := xid.New().String()
	err = c.txn.Set([]byte(imgid), res2.Image)
	if err != nil {
		return img
	}
	return imgid
}
*/
