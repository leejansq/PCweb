package models

import (
	"crypto/hmac"
	//"crypto/md5"
	"crypto/aes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"net/url"
)

func url_base64_hmac_sha1(val, key string) string {
	Hs1 := hmac.New(sha1.New, []byte(key))
	Hs1.Write([]byte(val))
	return url.QueryEscape(base64.StdEncoding.EncodeToString(Hs1.Sum(nil)))
}

func AesDecrpto(src, key string) string {
	srcby, _ := hex.DecodeString(src)
	if len(key) < 16 {
		return ""
	}
	block, err := aes.NewCipher([]byte(key[:16]))
	if err != nil {
		return ""
	}
	num_block := len(srcby) / 16
	//if num_block*block.BlockSize() != len(src) {
	//	num_block++
	//}
	//fmt.Println(block.BlockSize())
	var ret []byte
	dst := make([]byte, 16)
	for i := 0; i < num_block; i++ {
		block.Decrypt(dst, srcby[i*16:(i+1)*16])
		//fmt.Println(string(dst))
		ret = append(ret, dst...)
	}
	return string(ret)
}
