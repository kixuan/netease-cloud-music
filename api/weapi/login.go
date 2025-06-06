// MIT License
//
// Copyright (c) 2024 chaunsin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package weapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chaunsin/netease-cloud-music/api"
	"github.com/chaunsin/netease-cloud-music/api/types"
	"github.com/chaunsin/netease-cloud-music/pkg/crypto"

	"github.com/skip2/go-qrcode"
)

type QrcodeCreateKeyReq struct {
	types.ReqCommon
	Type int64 `json:"type"`
}

type QrcodeCreateKeyResp struct {
	types.RespCommon[any]
	UniKey string `json:"unikey"`
}

// QrcodeCreateKey 生成二维码需要得key
// 常见问题
// 1. 请求成功了,但是body为空值什么也没有,原因还是参数加密出现了问题。
// 2. crsftoken 可传可不传个人猜测前端写得通用框架传了
func (a *Api) QrcodeCreateKey(ctx context.Context, req *QrcodeCreateKeyReq) (*QrcodeCreateKeyResp, error) {
	var (
		url   = "https://music.163.com/weapi/login/qrcode/unikey"
		reply QrcodeCreateKeyResp
		opts  = api.NewOptions()
	)

	resp, err := a.client.Request(ctx, url, req, &reply, opts)
	if err != nil {
		return nil, fmt.Errorf("Request: %w", err)
	}
	_ = resp
	return &reply, nil
}

type QrcodeGenerateReq struct {
	CodeKey string
	Level   qrcode.RecoveryLevel // 二维码恢复率
}

type QrcodeGenerateResp struct {
	types.RespCommon[any]
	Qrcode      []byte
	QrcodePrint string
}

// QrcodeGenerate 根据 QrcodeCreateKey 接口生成得key生成生成二维码,注意此处不是调用服务接口。
func (a *Api) QrcodeGenerate(ctx context.Context, req *QrcodeGenerateReq) (*QrcodeGenerateResp, error) {
	var (
		content = fmt.Sprintf("https://music.163.com/login?codekey=%s", req.CodeKey)
		reply   QrcodeGenerateResp
	)

	qr, err := qrcode.New(content, req.Level)
	if err != nil {
		return nil, err
	}
	reply.Qrcode, err = qr.PNG(256)
	if err != nil {
		return nil, fmt.Errorf("PNG: %w", err)
	}
	reply.QrcodePrint = qr.ToSmallString(false)
	// if err := qr.WriteFile(256, "./qrcode.png"); err != nil {
	// 	return nil, fmt.Errorf("WriteFile: %w", err)
	// }

	// if err := qrcode.WriteFile(content, qrcode.Medium, 256, "./qrcode.png"); err != nil {
	// 	return nil, fmt.Errorf("WriteFile: %w", err)
	// }
	return &reply, nil
}

type QrcodeCheckReq struct {
	Key  string `json:"key"`  // QrcodeCreateKey()返回值codekey
	Type int64  `json:"type"` // 目前传1
}

type QrcodeCheckResp struct {
	types.RespCommon[any]
}

// QrcodeCheck 查询扫码状态
// 返回值:
// 800-二维码不存在或已过期
// 801-等待扫码
// 802-正在扫码授权中
// 803-授权登录成功
func (a *Api) QrcodeCheck(ctx context.Context, req *QrcodeCheckReq) (*QrcodeCheckResp, error) {
	var (
		url   = "https://music.163.com/weapi/login/qrcode/client/login"
		reply QrcodeCheckResp
		opts  = api.NewOptions()
	)

	resp, err := a.client.Request(ctx, url, req, &reply, opts)
	if err != nil {
		return nil, fmt.Errorf("Request: %w", err)
	}
	_ = resp
	return &reply, nil
}

type GetUserInfoReq struct {
	types.ReqCommon
}

type GetUserInfoResp struct {
	types.RespCommon[any]
	Account *GetUserInfoRespAccount `json:"account"`
	Profile *GetUserInfoRespProfile `json:"profile"`
}

type GetUserInfoRespAccount struct {
	Id                 int64  `json:"id"`
	UserName           string `json:"userName"`
	Type               int64  `json:"type"`
	Status             int64  `json:"status"`
	WhitelistAuthority int64  `json:"whitelistAuthority"`
	CreateTime         int64  `json:"createTime"`
	TokenVersion       int64  `json:"tokenVersion"`
	Ban                int64  `json:"ban"`
	BaoyueVersion      int64  `json:"baoyueVersion"`
	DonateVersion      int64  `json:"donateVersion"`
	VipType            int64  `json:"vipType"`
	AnonimousUser      bool   `json:"anonimousUser"`
	PaidFee            bool   `json:"paidFee"`
}

type GetUserInfoRespProfile struct {
	UserId              int64       `json:"userId"`
	UserType            int64       `json:"userType"`
	Nickname            string      `json:"nickname"`
	AvatarImgId         int64       `json:"avatarImgId"`
	AvatarUrl           string      `json:"avatarUrl"`
	BackgroundImgId     int64       `json:"backgroundImgId"`
	BackgroundUrl       string      `json:"backgroundUrl"`
	Signature           string      `json:"signature"`
	CreateTime          int64       `json:"createTime"`
	UserName            string      `json:"userName"`
	AccountType         int64       `json:"accountType"`
	ShortUserName       string      `json:"shortUserName"`
	Birthday            int64       `json:"birthday"`
	Authority           int64       `json:"authority"`
	Gender              int64       `json:"gender"`
	AccountStatus       int64       `json:"accountStatus"`
	Province            int64       `json:"province"`
	City                int64       `json:"city"`
	AuthStatus          int64       `json:"authStatus"`
	Description         interface{} `json:"description"`
	DetailDescription   interface{} `json:"detailDescription"`
	DefaultAvatar       bool        `json:"defaultAvatar"`
	ExpertTags          interface{} `json:"expertTags"`
	Experts             interface{} `json:"experts"`
	DjStatus            int64       `json:"djStatus"`
	LocationStatus      int64       `json:"locationStatus"`
	VipType             int64       `json:"vipType"`
	Followed            bool        `json:"followed"`
	Mutual              bool        `json:"mutual"`
	Authenticated       bool        `json:"authenticated"`
	LastLoginTime       int64       `json:"lastLoginTime"`
	LastLoginIP         string      `json:"lastLoginIP"`
	RemarkName          interface{} `json:"remarkName"`
	ViptypeVersion      int64       `json:"viptypeVersion"`
	AuthenticationTypes int64       `json:"authenticationTypes"`
	AvatarDetail        interface{} `json:"avatarDetail"`
	Anchor              bool        `json:"anchor"`
}

// GetUserInfo 获取用户信息
func (a *Api) GetUserInfo(ctx context.Context, req *GetUserInfoReq) (*GetUserInfoResp, error) {
	var (
		url   = "https://music.163.com/weapi/w/nuser/account/get"
		reply GetUserInfoResp
		opts  = api.NewOptions()
	)

	resp, err := a.client.Request(ctx, url, req, &reply, opts)
	if err != nil {
		return nil, fmt.Errorf("Request: %w", err)
	}
	_ = resp
	return &reply, nil
}

type TokenRefreshReq struct {
	types.ReqCommon
}

type TokenRefreshResp struct {
	types.RespCommon[any]
	BizCode string `json:"bizCode"` // 201:貌似刷新成功 400:貌似刷新不成功 504:貌似token已经过期了或者无效了
}

// TokenRefresh 登录token刷新
// har:
func (a *Api) TokenRefresh(ctx context.Context, req *TokenRefreshReq) (*TokenRefreshResp, error) {
	var (
		url   = "https://music.163.com/weapi/login/token/refresh"
		reply TokenRefreshResp
		opts  = api.NewOptions()
	)
	if req.CSRFToken == "" {
		csrf, _ := a.client.GetCSRF(url)
		req.CSRFToken = csrf
	}

	// 以下参数分析从eapi中分析得来
	// 请求头重需要传,此外此token也在v6/playlist中也有使用:
	// x-anticheattoken=9ca17ae2e6ffcda170e2e6ee88fb7db79eaf96f0409ab48aa3c54b929e9ab0d670b1ee8891d55fed93fd85b52af0feaec3b92af8f1e1a2e65293eb8c91c45b869a9fa6d45e948997daec44ad9b98a6cc70b59dee9e
	// MUSIC_R_U=00C572559E9EC4370FB21EB2CDFC28BA79632C61958228B75DA68C65488B3719DE982C68ED14E9026C527B9896FC29CF399F86469F18716A44AAC30F6FEF8A40BCD5575D6D311B95ACE21C05E94AF988B7
	// 参数中要传：
	// "checkToken":"9ca17ae2e6ffcda170e2e6ee88fb7db79eaf96f0409ab48aa3c54b929e9ab0d670b1ee8891d55fed93fd85b52af0feaec3b92af8f1e1a2e65293eb8c91c45b869a9fa6d45e948997daec44ad9b98a6cc70b59dee9e"
	// 其中header结构体中得字段X-antiCheatToken也传和checkToken同样之

	// 经测试MUSIC_R_U需要传参,否则会返回bizCode返回400错误
	// opts.SetHeader("x-anticheattoken", "9ca17ae2e6ffcda170e2e6ee88fb7db79eaf96f0409ab48aa3c54b929e9ab0d670b1ee8891d55fed93fd85b52af0feaec3b92af8f1e1a2e65293eb8c91c45b869a9fa6d45e948997daec44ad9b98a6cc70b59dee9e")
	// opts.SetCookies(&http.Cookie{Name: "MUSIC_R_U", Value: "00C572559E9EC4370FB21EB2CDFC28BA79632C61958228B75DA68C65488B3719DE982C68ED14E9026C527B9896FC29CF399F86469F18716A44AAC30F6FEF8A40BCD5575D6D311B95ACE21C05E94AF988B7"})
	opts.SetCookies(&http.Cookie{Name: "os", Value: "pc"}) // 解决400问题
	resp, err := a.client.Request(ctx, url, req, &reply, opts)
	if err != nil {
		return nil, fmt.Errorf("Request: %w", err)
	}
	_ = resp
	return &reply, nil
}

type RegisterAnonymousReq struct {
	types.ReqCommon
	Username string `json:"username"` // 设备id如果为空则设备id为ncmctl
}

type RegisterAnonymousResp struct {
	types.RespCommon[any]
}

// RegisterAnonymous 匿名用户注册
// har: 33.har
func (a *Api) RegisterAnonymous(ctx context.Context, req *RegisterAnonymousReq) (*RegisterAnonymousResp, error) {
	var (
		url   = "https://interface.music.163.com/weapi/register/anonimous"
		reply RegisterAnonymousResp
		opts  = api.NewOptions()
	)
	if req.Username == "" {
		req.Username = "ncmctl" // 默认用户名
	}
	username, err := crypto.Anonymous(req.Username)
	if err != nil {
		return nil, fmt.Errorf("Anonymous: %w", err)
	}
	req.Username = username

	resp, err := a.client.Request(ctx, url, req, &reply, opts)
	if err != nil {
		return nil, fmt.Errorf("Request: %w", err)
	}
	_ = resp
	return &reply, nil
}
