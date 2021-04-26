package shorturl

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/busyfree/shorturl-go/util/crypto/md5"
)

func ShortUrl(params ...string) string {
	if len(params) == 0 {
		return ""
	}
	var key = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var longURL = params[0]
	if len(params) > 1 {
		key = params[1]
	}
	var text = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	md5text := md5.EncryptString(key + longURL)
	var shortURLs = make([]string, 4, 4)
	for i := 0; i < 4; i++ {
		str := md5text[i*8 : i*8+8]
		//取出MD5中第i个字节，并忽略超过30位的部分
		var num int
		num64, _ := strconv.ParseInt(fmt.Sprintf("%x", str), 10, 64) //把str里的十六进制表示转化成int
		num = int(num64)
		num &= 0x3FFFFFFF //选择低30位
		//取30位的后6位与0x0000003D进行逻辑与操作，结果范围是0~61，作为text的下标选择字符
		//把num右移5位重复进行，得到6个字符组成短URL
		var shortURL string
		for j := 0; j < 6; j++ {
			str2 := text[num&0x0000003D]
			shortURL += string(str2)
			num >>= 5
		}
		shortURLs[i] = shortURL
	}
	//随机返回四个短URL中的一个
	//fmt.Println(shortURLs)
	return shortURLs[rand.Intn(len(shortURLs))]
}
