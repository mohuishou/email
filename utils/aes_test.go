package utils

import (
	"testing"
)

func TestAES(t *testing.T) {
	a := NewAES("1234567890123456")
	str := "hello world"
	plantext := a.Encrypt(str)
	b, err := a.Decrypt(plantext)
	if err != nil {
		t.Fatal("解密失败", "原始字符: ", str, "加密结果: ", plantext)
	}
	detext := string(b)
	if detext != str {
		t.Fatal("解密之后的字符串和原始字符串不一致！", "原始字符: ", str, "加密结果: ", plantext, "解密结果：", detext)
	}

}
