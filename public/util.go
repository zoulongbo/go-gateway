package public

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
)

func GenSaltPassword(salt, password string) string {
	ps := sha256.New()
	ps.Write([]byte(password))
	psStr := fmt.Sprintf("%x", ps.Sum(nil))

	ss := sha256.New()
	ss.Write([]byte(psStr + salt))
	return fmt.Sprintf("%x", ss.Sum(nil))
}

//MD5 md5加密
func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Obj2Json(s interface{}) string  {
	str, _ := json.Marshal(s)
	return string(str)
}