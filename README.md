# xcos
一个简易的cos上传工具。

## use
```text
Usage:
  xcos [flags]
  xcos [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  download    download
  help        Help about any command
  upload      upload

Flags:
  -f, --filename string   filename
  -h, --help              help for xcos
  -r, --regin string      cos regin

Use "xcos [command] --help" for more information about a command.
```
```text
# 上传到腾讯云
# 需要设置秘钥和存储桶名称
t.BaseBucket = BaseBucket{
		Secretid:  "a", 
		Secretkey: "a",
		Bucketurl: "https://xxxx-xxxx.cos.ap-nanjing.myqcloud.com",
	}
go run main.go upload -r tencent -f filename
go run main.go download -r tencent -f filename
```