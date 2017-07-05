package utils

import (
	"testing"
)

func TestApp(t *testing.T) {
	aes := AES{key: []byte("1234567890123456")}
	appkey := NewAppToken("test", "mohuishou", "127.0.0.1", aes)
	app, err := DecryptToken(appkey, aes)
	if err != nil {
		t.Fatal("appkey 解密失败", "key: ", appkey)
	}
	if app.Author != "mohuishou" || app.Name != "test" || app.IP != "127.0.0.1" {
		t.Fatal("解密信息错误！", "期望得到的解密结果：[test mohuishou 127.0.0.1]")
	}
}
