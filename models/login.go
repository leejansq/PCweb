package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	domain_login string = "http://115.29.106.70/v2/users/login?"
	//base64Table         = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
	domain_list string = "http://115.29.106.70/v2/devices/list?"
)

type LData struct {
	Userid       int
	Name         string
	Token        string
	Token_secret string
	Flag         bool
}

type Login_obj struct {
	Data LData
	Code string
	Msg  string
}

type VData struct {
	Uid      string
	Name     string
	Password string
	Message  string
	Flag     string
	Share    string
	Online   string
}

type List_obj struct {
	Code string
	Msg  string
	Data []VData
}

func YILogin(user, key string) (*Login_obj, error) {
	suffix := "seq=1&account=" + user
	m5 := md5.New()
	m5.Write([]byte(key))
	hmacv := m5.Sum(nil)
	suffix = suffix + "&hmac=" + url_base64_hmac_sha1(suffix, hex.EncodeToString(hmacv))
	log.Println(domain_login + suffix)
	resp, err := http.Get(domain_login + suffix)
	if err != nil {
		return nil, errors.New("Network Error!")
	}
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New("Return Body Error!")
		}
		ret := &Login_obj{}
		json.Unmarshal(body, ret)
		return ret, nil
	}
	return nil, errors.New("Not 200 Error!")
}

func GetList(token string, token_secret string, userid int) (*List_obj, error) {
	suffix := "seq=1&userid=" + strconv.Itoa(userid) + "&token=" + token
	log.Println(suffix)
	key := token + "&" + token_secret
	suffix = suffix + "&hmac=" + url_base64_hmac_sha1(suffix, key)
	resp, err := http.Get(domain_list + suffix)

	if err != nil {
		return nil, errors.New("Network Error!")
	}
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New("Return Body Error!")
		}
		ret := &List_obj{}
		json.Unmarshal(body, ret)
		return ret, nil
	}
	return nil, errors.New("Not 200 Error!")
}
