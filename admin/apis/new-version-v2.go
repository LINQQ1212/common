package apis

import (
	"github.com/LINQQ1212/common/create_db_v2"
	"github.com/LINQQ1212/common/global"
	"github.com/LINQQ1212/common/models"
	"github.com/LINQQ1212/common/response"
	"github.com/LINQQ1212/common/utils"
	"github.com/dtylman/scp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func NewVersionV2(c *gin.Context) {
	req := models.NewVersionReqV2{}
	if err := c.BindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	dir := strings.TrimSuffix(path.Base(req.ProductTarLink), ".tar.gz")
	if dir == "" {
		response.FailWithMessage("产品tar链接异常，为空", c)
		return
	}
	go NewVersionV2Start(req, dir)
	response.OkWithMessage("后台执行中", c)
}

func NewVersionV2Start(req models.NewVersionReqV2, dir string) {
	pdir := path.Join(req.TopDir, dir)
	if req.RemoteCopy {
		config := &ssh.ClientConfig{
			Timeout:         10 * time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
			User:            req.RemoteUser,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全
			Auth:            []ssh.AuthMethod{ssh.Password(req.RemotePwd)},
			//HostKeyCallback: hostKeyCallBackFunc(h.Host),
		}
		sshClient, err := ssh.Dial("tcp", req.RemoteHost+":"+req.RemotePort, config)
		if err != nil {
			sendTgMessage(req.Domain + "\nIP:" + req.RemoteHost + "\n" + "远程服务器链接失败：" + err.Error())
			return
		}
		if ok, _ := utils.PathExists(pdir); ok {
			os.Remove(pdir)
		}
		os.MkdirAll(pdir, os.ModeDir)
		_, err = scp.CopyFrom(sshClient, pdir+"/*.tar.gz", pdir)
		if err != nil {
			sendTgMessage(req.Domain + "\nIP:" + req.RemoteHost + "\n" + "复制远程 tar.gz 文件错误：" + err.Error())
			return
		}
		err = copyZip(sshClient, pdir, "gok", req.Domain+"\nIP:"+req.RemoteHost+"\n 复制远程 Google图片")
		if err != nil && !req.GErrorSkip {
			global.LOG.Error("Copy Google图片", zap.Error(err))
			return
		}

		err = copyZip(sshClient, pdir, "yok", req.Domain+"\nIP:"+req.RemoteHost+"\n 复制远程 雅虎描述")
		if err != nil && !req.YErrorSkip {
			global.LOG.Error("Copy 雅虎描述", zap.Error(err))
			return
		}

		err = copyZip(sshClient, pdir, "bok", req.Domain+"\nIP:"+req.RemoteHost+"\n 复制远程 Bing描述")
		if err != nil && !req.BErrorSkip {
			global.LOG.Error("Copy Bing描述", zap.Error(err))
			return
		}

		err = copyZip(sshClient, pdir, "ytok", req.Domain+"\nIP:"+req.RemoteHost+"\n 复制远程 Youtube")
		if err != nil && !req.YtErrorSkip {
			global.LOG.Error("Copy Youtube", zap.Error(err))
			return
		}

		sshClient.Close()
	}
	if strings.HasPrefix(req.ProductTarLink, "http") {
		req.ProductTarLink = path.Join(pdir, dir+".tar.gz")
	}

	if utils.FileExist(req.ProductTarLink) {
		sendTgMessage(req.Domain + "\n" + req.ProductTarLink + " 文件不存在")
		return
	}
	if req.AutoFilePath {
		reg := regexp.MustCompile(`_\d+$`)
		fname := reg.ReplaceAllString(dir, "")
		req.GoogleImg = getZipFile(pdir, "gok", fname, req.GErrorSkip)
		req.YahooDsc = getZipFile(pdir, "yok", fname, req.YErrorSkip)
		req.BingDsc = getZipFile(pdir, "bok", fname, req.BErrorSkip)
		req.YoutubeDsc = getZipFile(pdir, "ytok", fname, req.YtErrorSkip)
	}

	cdb := create_db_v2.New(req)
	err := cdb.Start()
	if err != nil {
		global.LOG.Error("Start", zap.Error(err))
		sendTgMessage(req.Domain + "\n" + err.Error())
		return
	}
	v, err := models.NewVersion(global.VersionDir, cdb.GetDomain())
	if err != nil {
		global.LOG.Error("NewVersion", zap.Error(err))
		sendTgMessage(req.Domain + "\n 读取数据库错误：" + err.Error())
		return
	}
	global.Versions.Set(cdb.GetDomain(), v)
	err = global.VHost.NewVersion(cdb.GetDomain())
	if err != nil {
		sendTgMessage(req.Domain + "\n 创建虚拟主机错误：" + err.Error())
		global.LOG.Error("NewVersion", zap.String("version", req.Domain), zap.Error(err))
		return
	}
	sendTgMessage(req.Domain + " #创建完成#")
	if req.EndRemove {
		os.Remove(pdir)
	}
	if req.RemoteEndRemove {
		config := &ssh.ClientConfig{
			Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
			User:            req.RemoteUser,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全
			Auth:            []ssh.AuthMethod{ssh.Password(req.RemotePwd)},
			//HostKeyCallback: hostKeyCallBackFunc(h.Host),
		}
		sshClient, err := ssh.Dial("tcp", req.RemoteHost+":"+req.RemotePort, config)
		if err != nil {
			return
		}
		defer sshClient.Close()
		session, err := sshClient.NewSession()
		if err != nil {
			return
		}
		defer session.Close()
		session.Run("rm " + pdir)
	}
}

func copyZip(sshClient *ssh.Client, dir, dir2, mess string) error {
	ndir := path.Join(dir, dir2)
	os.MkdirAll(ndir, os.ModeDir)
	_, err := scp.CopyFrom(sshClient, ndir+"/*.zip", ndir)
	if err != nil {
		sendTgMessage(mess + "\n错误：" + err.Error())
		return err
	}
	return nil
}

func getZipFile(topDir, dir, fname string, skip bool) string {
	mdir := path.Join(topDir, dir)
	file := path.Join(mdir, fname+".zip")
	if utils.FileExist(file) {
		return file
	}
	fs, err := filepath.Glob(mdir + "/*.zip")
	if err == nil && len(fs) > 0 {
		return fs[0]
	}
	if !skip {
		return file
	}
	return ""
}

func sendTgMessage(s string) {
	http.Post("https://api.telegram.org/"+global.CONFIG.System.TGToken+"/sendMessage", "application/json", strings.NewReader(`{"chat_id":`+global.CONFIG.System.TGChatId+`,"text":"`+s+`"}`))
}
