package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	var files []string
	for idx, args := range os.Args {
		if idx <= 0 {
			continue
		}
		file, err := uploadImg(args)
		if err != nil {
			panic(err)
		}
		files = append(files, file)
	}
	if len(files) > 0 {
		fmt.Println("Upload Success:")
		for _, file := range files {
			fmt.Println(config.Domain + file)
		}
	}
}

func uploadImg(localFile string) (string, error) {
	bucket := config.Bucket
	imgSuffix := path.Ext(localFile)
	key := config.BucketPath + genFilename() + imgSuffix

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(config.AccessKey, config.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	// 可选配置
	putExtra := storage.PutExtra{}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		return "", err
	}
	return ret.Key, nil
}

func genFilename() string {
	tm := strconv.FormatInt(time.Now().UnixNano(), 10)
	md5String := md5.Sum([]byte(tm))
	return hex.EncodeToString(md5String[:])
}

type Config struct {
	AccessKey  string
	SecretKey  string
	Bucket     string
	BucketPath string
	Domain     string
}

var config *Config

// 加载配置文件
func init() {
	config = &Config{}
	file, _ := os.Open(getCurrentPath() + "config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(config)
	if err != nil {
		fmt.Println("Load Config Error:", err)
	}
}

// GetCurrentPath 返回当前程序运行的路径
func getCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) + "/"
}
