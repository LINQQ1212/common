package apis

import (
	"github.com/LINQQ1212/common/create_db_v2"
	"github.com/LINQQ1212/common/global"
	"github.com/LINQQ1212/common/models"
	"github.com/LINQQ1212/common/response"
	"github.com/LINQQ1212/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
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
		sftpClient, err := sftp.NewClient(sshClient)
		if err != nil {
			sendTgMessage(req.Domain + "\nIP:" + req.RemoteHost + "\n" + "远程服务器链接失败：" + err.Error())
			return
		}
		list, err := download2(sftpClient, pdir, req.TopDir)
		if err != nil {
			sendTgMessage(req.Domain + "\nIP:" + req.RemoteHost + "\n" + "复制远程服务器的文件失败：" + err.Error())
			return
		}

		sshClient.Close()
		if len(list) > 0 && req.AutoFilePath {
			for _, s := range list {
				if strings.HasSuffix(s, ".tar.gz") {
					req.ProductTarLink = s
					continue
				}
				if strings.Contains(s, "/gok/") {
					req.GoogleImg = s
					continue
				}
				if strings.Contains(s, "/yok/") {
					req.YahooDsc = s
					continue
				}
				if strings.Contains(s, "/bok/") {
					req.BingDsc = s
					continue
				}
				if strings.Contains(s, "/ytok/") {
					req.YoutubeDsc = s
					continue
				}
			}
			req.AutoFilePath = false
		}
	}
	if strings.HasPrefix(req.ProductTarLink, "http") {
		req.ProductTarLink = path.Join(pdir, dir+".tar.gz")
	}

	if !utils.FileExist(req.ProductTarLink) {
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

func download2(sourceClient *sftp.Client, sourcePath string, destPath string) (fs []string, err error) {
	var sourceFile *sftp.File
	sourceFile, err = sourceClient.Open(sourcePath)
	if err != nil {
		return
	}
	defer sourceFile.Close()
	var info os.FileInfo
	info, err = sourceFile.Stat()
	if info.Name() == "txt" || info.Name() == "log" {
		return
	}
	if err != nil {
		return
	}
	if info.IsDir() {
		nextDestPath := path.Join(destPath, info.Name())
		os.Mkdir(nextDestPath, info.Mode())
		var childInfos []os.FileInfo
		childInfos, err = sourceClient.ReadDir(sourcePath)
		if err != nil {
			return
		}
		for _, child := range childInfos {
			nextSourcePath := path.Join(sourcePath, child.Name())
			var list []string
			list, err = download2(sourceClient, nextSourcePath, nextDestPath)
			if err != nil {
				return
			}
			fs = append(fs, list...)
		}
	} else {
		var sourceFile *sftp.File
		sourceFile, err = sourceClient.Open(sourcePath)
		if err != nil {
			return
		}
		defer sourceFile.Close()
		destFileName := path.Join(destPath, info.Name())
		fs = append(fs, destFileName)
		var destFile *os.File
		destFile, err = os.Create(destFileName)
		if err != nil {
			return
		}
		defer destFile.Close()

		if _, err = sourceFile.WriteTo(destFile); err != nil {
			return
		}
	}
	return
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
