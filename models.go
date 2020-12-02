package workwx

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

type reqAccessToken struct {
	CorpID     string
	CorpSecret string
}

var _ urlValuer = reqAccessToken{}

func (x reqAccessToken) intoURLValues() url.Values {
	return url.Values{
		"corpid":     {x.CorpID},
		"corpsecret": {x.CorpSecret},
	}
}

type respCommon struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// IsOK 响应体是否为一次成功请求的响应
//
// 实现依据: https://work.weixin.qq.com/api/doc#10013
//
// > 企业微信所有接口，返回包里都有errcode、errmsg。
// > 开发者需根据errcode是否为0判断是否调用成功(errcode意义请见全局错误码)。
// > 而errmsg仅作参考，后续可能会有变动，因此不可作为是否调用成功的判据。
func (x *respCommon) IsOK() bool {
	return x.ErrCode == 0
}

func (x *respCommon) TryIntoErr() error {
	if x.IsOK() {
		return nil
	}

	return &WorkwxClientError{
		Code: x.ErrCode,
		Msg:  x.ErrMsg,
	}
}

type respAccessToken struct {
	respCommon

	AccessToken   string `json:"access_token"`
	ExpiresInSecs int64  `json:"expires_in"`
}

type reqJSAPITicketAgentConfig struct{}

var _ urlValuer = reqJSAPITicketAgentConfig{}

func (x reqJSAPITicketAgentConfig) intoURLValues() url.Values {
	return url.Values{
		"type": {"agent_config"},
	}
}

type reqJSAPITicket struct{}

var _ urlValuer = reqJSAPITicket{}

func (x reqJSAPITicket) intoURLValues() url.Values {
	return url.Values{}
}

type respJSAPITicket struct {
	respCommon

	Ticket        string `json:"ticket"`
	ExpiresInSecs int64  `json:"expires_in"`
}

// reqMessage 消息发送请求
type reqMessage struct {
	ToUser  []string
	ToParty []string
	ToTag   []string
	ChatID  string
	AgentID int64
	MsgType string
	Content map[string]interface{}
	IsSafe  bool
}

var _ bodyer = reqMessage{}

func (x reqMessage) intoBody() ([]byte, error) {
	// fuck
	safeInt := 0
	if x.IsSafe {
		safeInt = 1
	}

	obj := map[string]interface{}{
		"msgtype": x.MsgType,
		"agentid": x.AgentID,
		"safe":    safeInt,
	}

	// msgtype polymorphism
	obj[x.MsgType] = x.Content

	// 复用这个结构体，因为是 package-private 的所以这么做没风险
	if x.ChatID != "" {
		obj["chatid"] = x.ChatID
	} else {
		obj["touser"] = strings.Join(x.ToUser, "|")
		obj["toparty"] = strings.Join(x.ToParty, "|")
		obj["totag"] = strings.Join(x.ToTag, "|")
	}

	result, err := json.Marshal(obj)
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

// respMessageSend 消息发送响应
type respMessageSend struct {
	respCommon

	InvalidUsers   string `json:"invaliduser"`
	InvalidParties string `json:"invalidparty"`
	InvalidTags    string `json:"invalidtag"`
}

type reqUserGet struct {
	UserID string
}

var _ urlValuer = reqUserGet{}

func (x reqUserGet) intoURLValues() url.Values {
	return url.Values{
		"userid": {x.UserID},
	}
}

// respUserDetail 成员详细信息的公共字段
type respUserDetail struct {
	UserID         string   `json:"userid"`
	Name           string   `json:"name"`
	DeptIDs        []int64  `json:"department"`
	DeptOrder      []uint32 `json:"order"`
	Position       string   `json:"position"`
	Mobile         string   `json:"mobile"`
	Gender         string   `json:"gender"`
	Email          string   `json:"email"`
	IsLeaderInDept []int    `json:"is_leader_in_dept"`
	AvatarURL      string   `json:"avatar"`
	Telephone      string   `json:"telephone"`
	IsEnabled      int      `json:"enable"`
	Alias          string   `json:"alias"`
	Status         int      `json:"status"`
	QRCodeURL      string   `json:"qr_code"`
	// TODO: extattr external_profile external_position
}

// respUserGet 读取成员响应
type respUserGet struct {
	respCommon

	respUserDetail
}

// reqUserList 部门成员请求
type reqUserList struct {
	DeptID     int64
	FetchChild bool
}

var _ urlValuer = reqUserList{}

func (x reqUserList) intoURLValues() url.Values {
	var fetchChild int64
	if x.FetchChild {
		fetchChild = 1
	}

	return url.Values{
		"department_id": {strconv.FormatInt(x.DeptID, 10)},
		"fetch_child":   {strconv.FormatInt(fetchChild, 10)},
	}
}

// respUsersByDeptID 部门成员详情响应
type respUserList struct {
	respCommon

	Users []*respUserDetail `json:"userlist"`
}

// reqUserIDByMobile 手机号获取 userid 请求
type reqUserIDByMobile struct {
	Mobile string `json:"mobile"`
}

var _ bodyer = reqUserIDByMobile{}

func (x reqUserIDByMobile) intoBody() ([]byte, error) {
	body, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// respUserIDByMobile 手机号获取 userid 响应
type respUserIDByMobile struct {
	respCommon

	UserID string `json:"userid"`
}

type reqDeptList struct {
	HaveID bool
	ID     int64
}

var _ urlValuer = reqDeptList{}

func (x reqDeptList) intoURLValues() url.Values {
	if !x.HaveID {
		return url.Values{}
	}

	return url.Values{
		"id": {strconv.FormatInt(x.ID, 10)},
	}
}

// respDeptList 部门列表响应
type respDeptList struct {
	respCommon

	// TODO: 不要懒惰，把 API 层的类型写好
	Department []*DeptInfo `json:"department"`
}

// reqAppchatGet 获取群聊会话请求
type reqAppchatGet struct {
	ChatID string
}

var _ urlValuer = reqAppchatGet{}

func (x reqAppchatGet) intoURLValues() url.Values {
	return url.Values{
		"chatid": {x.ChatID},
	}
}

// respAppchatGet 获取群聊会话响应
type respAppchatGet struct {
	respCommon

	ChatInfo *ChatInfo `json:"chat_info"`
}

// reqAppchatCreate 创建群聊会话请求
type reqAppchatCreate struct {
	ChatInfo *ChatInfo
}

var _ bodyer = reqAppchatCreate{}

func (x reqAppchatCreate) intoBody() ([]byte, error) {
	result, err := json.Marshal(x.ChatInfo)
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

// respAppchatCreate 创建群聊会话响应
type respAppchatCreate struct {
	respCommon

	ChatID string `json:"chatid"`
}

// reqMediaUpload 临时素材上传请求
type reqMediaUpload struct {
	Type  string
	Media *Media
}

var _ urlValuer = reqMediaUpload{}
var _ mediaUploader = reqMediaUpload{}

func (x reqMediaUpload) intoURLValues() url.Values {
	return url.Values{
		"type": {x.Type},
	}
}

func (x reqMediaUpload) getMedia() *Media {
	return x.Media
}

// respMediaUpload 临时素材上传响应
type respMediaUpload struct {
	respCommon

	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

// reqMediaUploadImg 永久图片素材上传请求
type reqMediaUploadImg struct {
	Media *Media
}

var _ urlValuer = reqMediaUploadImg{}
var _ mediaUploader = reqMediaUploadImg{}

func (x reqMediaUploadImg) intoURLValues() url.Values {
	return url.Values{}
}

func (x reqMediaUploadImg) getMedia() *Media {
	return x.Media
}

// respMediaUploadImg 永久图片素材上传响应
type respMediaUploadImg struct {
	respCommon

	URL string `json:"url"`
}

// reqExternalContactList 获取客户列表
type reqExternalContactList struct {
	UserID string `json:"userid"`
}

var _ urlValuer = reqExternalContactList{}

func (x reqExternalContactList) intoURLValues() url.Values {
	return url.Values{
		"userid": {x.UserID},
	}
}

// respExternalContactList 获取客户列表
type respExternalContactList struct {
	respCommon

	ExternalUserID []string `json:"external_userid"`
}

// reqExternalContactGet 获取客户详情
type reqExternalContactGet struct {
	ExternalUserID string `json:"external_userid"`
}

var _ urlValuer = reqExternalContactGet{}

func (x reqExternalContactGet) intoURLValues() url.Values {
	return url.Values{
		"external_userid": {x.ExternalUserID},
	}
}

// respExternalContactGet 获取客户详情
type respExternalContactGet struct {
	respCommon
	ExternalContactInfo
}

type ExternalContactInfo struct {
	ExternalContact ExternalContact `json:"external_contact"`
	FollowUser      []FollowUser    `json:"follow_user"`
}

// reqExternalContactRemark 获取客户详情
type reqExternalContactRemark struct {
	Remark *ExternalContactRemark
}

var _ bodyer = reqExternalContactRemark{}

func (x reqExternalContactRemark) intoBody() ([]byte, error) {
	result, err := json.Marshal(x.Remark)
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

// respExternalContactRemark 获取客户详情
type respExternalContactRemark struct {
	respCommon
}

// reqUserInfoGet 获取访问用户身份
type reqUserInfoGet struct {
	// 通过成员授权获取到的code，最大为512字节。每次成员授权带上的code将不一样，code只能使用一次，5分钟未被使用自动过期。
	Code string
}

var _ urlValuer = reqUserInfoGet{}

func (x reqUserInfoGet) intoURLValues() url.Values {
	return url.Values{
		"code": {x.Code},
	}
}

// respUserInfoGet 部门列表响应
type respUserInfoGet struct {
	respCommon
	UserIdentityInfo
}

// reqExternalContactListCorpTags 获取企业标签库
type reqExternalContactListCorpTags struct {
	// 要查询的标签id，如果不填则获取该企业的所有客户标签，目前暂不支持标签组id
	TagIDs []string `json:"tag_id"`
}

var _ bodyer = reqExternalContactListCorpTags{}

func (x reqExternalContactListCorpTags) intoBody() ([]byte, error) {
	result, err := json.Marshal(x)
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

// respExternalContactListCorpTags 获取企业标签库
type respExternalContactListCorpTags struct {
	respCommon
	// 标签组列表
	TagGroup []ExternalContactCorpTagGroup `json:"tag_group"`
}

// reqExternalContactAddCorpTag 添加企业客户标签
type reqExternalContactAddCorpTag struct {
	ExternalContactCorpTagGroup
}

var _ bodyer = reqExternalContactAddCorpTag{}

func (x reqExternalContactAddCorpTag) intoBody() ([]byte, error) {
	result, err := json.Marshal(x.ExternalContactCorpTagGroup)
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

// respExternalContactAddCorpTag 添加企业客户标签
type respExternalContactAddCorpTag struct {
	respCommon
	// 标签组列表
	TagGroup []ExternalContactCorpTagGroup `json:"tag_group"`
}

// reqExternalContactEditCorpTag 编辑企业客户标签
type reqExternalContactEditCorpTag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Order uint32 `json:"order"`
}

var _ bodyer = reqExternalContactEditCorpTag{}

func (x reqExternalContactEditCorpTag) intoBody() ([]byte, error) {
	result, err := json.Marshal(x)
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

// respExternalContactEditCorpTag 编辑企业客户标签
type respExternalContactEditCorpTag struct {
	respCommon
}

// reqExternalContactDelCorpTag 删除企业客户标签
type reqExternalContactDelCorpTag struct {
	TagID   []string `json:"tag_id"`
	GroupID []string `json:"group_id"`
}

var _ bodyer = reqExternalContactDelCorpTag{}

func (x reqExternalContactDelCorpTag) intoBody() ([]byte, error) {
	result, err := json.Marshal(x)
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

// reqExternalContactDelCorpTag 删除企业客户标签
type respExternalContactDelCorpTag struct {
	respCommon
}

// reqExternalContactMarkTag 编辑企业客户标签
type reqExternalContactMarkTag struct {
	UserID         string   `json:"userid"`
	ExternalUserID string   `json:"external_userid"`
	AddTag         []string `json:"add_tag"`
	RemoveTag      []string `json:"remove_tag"`
}

var _ bodyer = reqExternalContactMarkTag{}

func (x reqExternalContactMarkTag) intoBody() ([]byte, error) {
	result, err := json.Marshal(x)
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

// reqExternalContactMarkTag 编辑企业客户标签
type respExternalContactMarkTag struct {
	respCommon
}

// +----------------------------------------------------------------------
// | 标签管理
// +----------------------------------------------------------------------

// reqTag 添加标签请求
type reqTag struct {
	Tag
}

var _ bodyer = reqTag{}

func (x reqTag) intoBody() ([]byte, error) {
	result, err := json.Marshal(x.Tag)
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

// respTag 添加标签响应
type respTag struct {
	respCommon
	// 标签组列表
	TagID int `json:"tagid"`
}

// reqEditTag 编辑标签请求
type reqEditTag struct {
	// tagname 标签名称
	TagName string `json:"tagname"`
	// tagid 标签id
	TagID int `json:"tagid"`
}

var _ bodyer = reqEditTag{}

func (x reqEditTag) intoBody() ([]byte, error) {
	result, err := json.Marshal(x)
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

// respEditTag 编辑标签响应
type respEditTag struct {
	respCommon
}
