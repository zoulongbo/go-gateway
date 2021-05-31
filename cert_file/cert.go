package cert_file

import (
	"path/filepath"
	"runtime"
)

/*
证书签名生成方式:
//CA私钥
openssl genrsa -out ca.key 2048
//CA数据证书
openssl req -x509 -new -nodes -key ca.key -subj "/CN=example1.com" -days 5000 -out ca.crt
//服务器私钥（默认由CA签发）
openssl genrsa -out server.key 2048
//服务器证书签名请求：Certificate Sign Request，简称csr（example1.com代表你的域名）
openssl req -new -key server.key -subj "/CN=example1.com" -out server.csr
//上面2个文件生成服务器证书（days代表有效期）
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 5000
*/

var basePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath = filepath.Dir(currentFile)
}

func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basePath, rel)
}

