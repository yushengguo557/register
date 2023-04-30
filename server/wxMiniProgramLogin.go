package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yushengguo557/register/common"
	"github.com/yushengguo557/register/util"
)

// Code2SessionResponse 微信小程序登陆 code2Session 接口响应
type Code2SessionResponse struct {
	SessionKey string `json:"session_key"` // 会话密钥
	Unionid    string `json:"unionid"`     // 用户在开放平台的唯一标识符 若当前小程序已绑定到微信开放平台帐号下会返回
	Errmsg     string `json:"errmsg"`      // 错误信息
	Openid     string `json:"openid"`      // 用户唯一标识
	Errcode    int32  `json:"errcode"`     // 错误码 (=0标识没有错误🙅‍♂️)
}

// WXMiniProgramLogin 微信小程序登录
func (s *Server) WXMiniProgramLogin(code string) (common.Response, error) {
	AppID := "wx0f79c8e2fa2000cd"                   // AppID(小程序ID)
	AppSecret := "4837f4bdc5b3fdf2af7416ea224a1ed4" // AppSecret(小程序密钥)

	// 1.请求微信接口，获取 openid 和 session_key
	// GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
	// var url bytes.Buffer
	// url.WriteString("https://api.weixin.qq.com/sns/jscode2session?appid=")
	// url.WriteString(AppID)
	// url.WriteString("&secret=")
	// url.WriteString(AppSecret)
	// url.WriteString("&js_code=")
	// url.WriteString(code)
	// url.WriteString("&grant_type=authorization_code")
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		AppID, AppSecret, code,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取 openid 和 session_key 失败: %w", err)
	}
	defer resp.Body.Close()

	// 2.读取 并 解析返回的 JSON 数据
	// body, err := io.ReadAll(resp.Body)
	reader := bufio.NewReader(resp.Request.Body)
	var p []byte
	n, err := reader.Read(p)
	if err != nil {
		return nil, fmt.Errorf("读取 openid 和 session_key 失败: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("读取 openid 和 session_key 失败: %w", err)
	}
	var rsp Code2SessionResponse
	err = json.Unmarshal(p[:n], &rsp)
	if err != nil {
		return nil, fmt.Errorf("解析 openid 和 session_key 失败: %w", err)
	}

	// 3.处理错误码
	if rsp.Errcode != 0 {
		return nil, fmt.Errorf("错误码: %d 错误信息: %s", rsp.Errcode, rsp.Errmsg)
	}

	// 4.生成一个用户唯一标识 token 并返回给客户端
	token, err := util.GenerateToken(rsp.Openid, rsp.SessionKey)
	if err != nil {
		return nil, fmt.Errorf("token 生成失败: %w", err)
	}

	// 5.将 openid 和 session_key token 存储到数据库中
	// ...
	fmt.Println(rsp.Openid, rsp.SessionKey)

	return common.Response{"token": token}, nil
}

// GetPhoneNumberRequest 获取电话号码请求
type GetPhoneNumberRequest struct {
	AccessToken string `json:"access_token"`
	Code        string `json:"code"`
}

// PhoneInfo 电话信息☎️
type PhoneInfo struct {
	PhoneNumber     string             `json:"phoneNumber"`     // 用户绑定的手机号（国外手机号会有区号）
	PurePhoneNumber string             `json:"purePhoneNumber"` // 没有区号的手机号
	CountryCode     string             `json:"countryCode"`     // 区号
	Watermark       `json:"watermark"` //  数据水印
}

// Watermark 水印
type Watermark struct {
	Timestamp int    `json:"timestamp"` // 用户获取手机号操作的时间戳
	Appid     string `json:"appid"`     // 小程序appid
}

// GetPhoneNumberResponse 获取手机号码响应
type GetPhoneNumberResponse struct {
	Errmsg    string              `json:"errmsg"`  // 错误信息
	Errcode   int32               `json:"errcode"` // 错误码 (=0标识没有错误🙅‍♂️)
	PhoneInfo `json:"phone_info"` // 用户手机号信息
}

// GetPhoneNumber 微信小程序登陆成功后获取手机号
func (s *Server) GetPhoneNumber(code, token string) error {
	// 1. 向微信接口服务发送 POST 请
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", token)
	jsonStr, err := json.Marshal(GetPhoneNumberRequest{AccessToken: token, Code: code})
	if err != nil { // 序列化失败
		return fmt.Errorf("序列化 GetPhoneNumberRequest 结构体失败: %w", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("向微信服务器发送POST请求获取手机号码失败: %w", err)
	}
	defer resp.Body.Close()

	// 2.读取 并 解析响应体json数据
	reader := bufio.NewReader(resp.Body)
	var p []byte
	if _, err = reader.Read(p); err != nil {
		return fmt.Errorf("从响应体中读取json数据失败: %w", err)
	}
	var ret GetPhoneNumberResponse
	err = json.Unmarshal(p, &ret)
	if err != nil {
		return fmt.Errorf("响应体反序列化失败: %w", err)
	}

	// 3.处理错误码
	if ret.Errcode != 0 {
		return fmt.Errorf("错误码: %d 错误信息: %s", ret.Errcode, ret.Errmsg)
	}
	// 4.保存手机号码  ret.CountryCode(区号) ret.PurePhoneNumber(无区号手机号)
	fmt.Println(ret.CountryCode, ret.PurePhoneNumber)

	return nil
}
