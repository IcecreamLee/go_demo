# Typora图床

写好配置文件config.json，然后放入程序同级文件夹：
```
{
    "AccessKey": "xxx",
    "SecretKey": "xxxx",
    "Bucket": "xxxxx",
    "BucketPath": "images/",
    "Domain": "https://xxx.com/"
}
```

Typor进入偏好设置-图像上传服务设定 \
上传服务选择框选择Custom Command
自定义命令填入go build 后 可执行文件的路径