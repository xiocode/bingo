/**
 * Author:        Tony.Shao
 * Email:         xiocode@gmail.com
 * Github:        github.com/xiocode
 * File:          login.go
 * Description:   WebQQ Login
 */

package bingo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

type webclient struct {
	Id           string // QQ號
	Password     string // 密碼
	Nickname     string // 昵稱
	aid          string // aid 固定
	client_id    int32  // 客户端id 随机固定
	message_id   int32  // 消息id, 随机初始化
	requireCheck bool   // 是否需要验证码
	polling      bool   // 开始拉取消息和心跳
	http_client  *http.Client

	// WebQQ登录期间需要保存的数据
	checkCode string
	skey      string
	ptwebqq   string

	check_data         string        // 检查时返回的数据
	blogin_data        string        // 登录前返回的数据
	friends_list       []interface{} // 好友列表
	groups_list        []interface{} // QQ 群信息
	group_members_list []interface{} // QQ 群成員信息

	hb_time             int // 心跳間隔時常
	login_time          int // 登錄時間
	last_group_msg_time int // 最近群信息
	last_msg_content    int // 最近發送信息
	last_msg_numbers    int // 剩餘消息數量
}

func GetQQClient() *webclient {

	client := &webclient{
		aid:         "1003903",
		http_client: new(http.Client),
	}
	return client
}

func (this *webclient) check() []string {
	fmt.Println("開始嘗試登錄！")
	url := "https://ssl.ptlogin2.qq.com/check"
	params := map[string]interface{}{
		"uin":   this.Id,
		"appid": this.aid,
		"r":     rand.Float64(),
	}
	url_params, err := encodeParams(params)
	if err != nil {
		panic(err)
	}
	request_url := fmt.Sprintf("%v?%v", url, url_params)
	check_resp, err := this.http_client.Get(request_url)
	if err != nil {
		panic(err)
	}
	defer check_resp.Body.Close()
	body, err := read_body(check_resp)
	if err != nil {
		panic(err)
	}
	body = body[strings.Index(body, "(")+1 : strings.Index(body, ")")]
	rest := strings.Split(body, ",")
	return rest
}

func (this *webclient) encodePassword() string {

}

func (this *webclient) Login() {
	url := "https://ssl.ptlogin2.qq.com/login"
	check_result := this.check()
	if check_result != nil && len(check_result) == 3 {
		checkCode := check_result[1]
		var captcha string
		if strings.HasPrefix(checkCode, "!") {
			// TODO 不需要驗證碼，直接登錄
			captcha = checkCode
		} else {
			this.getCaptcha()
			fmt.Println("請輸入驗證碼！")
			fmt.Scan(&captcha)
			fmt.Println(captcha)
		}
	}

}

func (this *webclient) getCaptcha() error {
	url := "https://ssl.captcha.qq.com/getimage"
	params := map[string]interface{}{
		"uin":   this.Id,
		"appid": this.aid,
		"r":     rand.Float64(),
	}
	url_params, err := encodeParams(params)
	if err != nil {
		panic(err)
	}
	request_url := fmt.Sprintf("%v?%v", url, url_params)
	resp, err := this.http_client.Get(request_url)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("captcha.jpg", data, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (this *webclient) get_login_sig() (string, error) {
	resp, err := this.http_client.Get("https://ui.ptlogin2.qq.com/cgi-bin/login?target=self&style=5&mibao_css=m_webqq&appid=1003903&enable_qlogin=0&no_verifyimg=1&s_url=http%3A%2F%2Fweb.qq.com%2Floginproxy.html&f_url=loginerroralert&strong_login=1&login_state=10&t=20130516001")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := read_body(resp)
	if err != nil {
		return "", err
	}
	pattern := regexp.MustCompile(`g_login_sig="([^"]+)";`)
	result := pattern.FindStringSubmatch(body)
	if len(result) == 0 {
		return "", errors.New("安全参数")
	}
	return result[1], nil
}
