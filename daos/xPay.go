package daos

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gm/common/constant"
	"gm/global"
	"gm/model/pay"
	"gm/request"
	"gm/response"
	"gm/utils/errUtil"
	req2 "gm/utils/request"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// md5 加密
func GetMd5(str string) (code string) {
	md := md5.New()
	md.Write([]byte(str))
	code = hex.EncodeToString(md.Sum(nil))
	return
}

// 字符串拼接
func StringJoin(str []string) string {
	var st bytes.Buffer
	for _, val := range str {
		st.WriteString(val)
	}
	return st.String()
}

// 格式化json字符串
func JsonString(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// 生成订单号
func GetOrderNo() (orderNo string) {
	now := time.Now().Local()
	ymdHis := now.Format("20060102150405")

	// 创建自定义的随机数生成器源
	source := rand.NewSource(time.Now().UnixNano())

	// 创建自定义的随机数生成器实例
	r := rand.New(source)

	num := r.Intn(1000) + 1000
	orderNo = StringJoin([]string{"NO", ymdHis, strconv.Itoa(num)})
	return orderNo
}

// 查询余额接口签名
func GetBalanceSign(req response.BalanceRequest, secret string) (sign string) {
	signStr := StringJoin([]string{"appId=", req.AppId, "&mchNo=", req.MchNo, "&reqTime=", req.ReqTime, "&key=", secret})
	return strings.ToUpper(GetMd5(signStr))
}

// 查询账户余额
func GetXPayBalance(payConfig pay.PayConfig) (data *response.XPayBalanceReturn, err error) {

	now := time.Now().Local()
	r := response.BalanceRequest{MchNo: payConfig.Merchant, AppId: payConfig.AppId, ReqTime: strconv.Itoa(int(now.UnixMilli()))}
	r.Sign = GetBalanceSign(r, payConfig.Secret)

	urlPath := payConfig.Url + constant.XPayBalanceUrl

	res := response.BalanceReturn{}

	err = req2.Call(urlPath, r, &res)
	if err != nil {
		global.Logger["err"].Errorf("GetXPayBalance req2.Call failed,err:[%v]", err.Error())
		return
	}

	if res.Code == 0 {
		data.Balance = res.Data.MchBalance
	}

	return
}

// 代付下单签名
func GetPaymentSign(req response.PaymentOrderRequest, secret string) (sign string) {
	signStr := ""
	if len(req.AccountAddress) > 0 {
		signStr = StringJoin([]string{"accountAddress=", req.AccountAddress, "&accountEmail=", req.AccountEmail, "&accountMobileNo=", req.AccountMobileNo, "&accountName=", req.AccountName})
	} else {
		signStr = StringJoin([]string{"accountEmail=", req.AccountEmail, "&accountMobileNo=", req.AccountMobileNo, "&accountName=", req.AccountName})
	}

	if len(req.AccountNo) > 0 {
		signStr = StringJoin([]string{signStr, "&accountNo=", req.AccountNo})
	}
	signStr = StringJoin([]string{signStr, "&appId=", req.AppId, "&bankCode=", req.BankCode})

	if len(req.BankName) > 0 {
		signStr = StringJoin([]string{signStr, "&bankName=", req.BankName})
	}
	signStr = StringJoin([]string{signStr, "&currency=INR"})
	if len(req.Ifsc) > 0 {
		signStr = StringJoin([]string{signStr, "&ifsc=", req.Ifsc})
	}

	signStr = StringJoin([]string{signStr, "&mchNo=", req.MchNo, "&mchOrderNo=", req.MchOrderNo, "&notifyUrl=", req.NotifyUrl, "&orderAmount=", req.OrderAmount, "&reqTime=", req.ReqTime})

	if len(req.TransferDesc) > 0 {
		signStr = StringJoin([]string{signStr, "&transferDesc=", req.TransferDesc})
	}
	signStr = StringJoin([]string{signStr, "&key=", secret})

	return strings.ToUpper(GetMd5(signStr))
}

func CheckFunPaySign(params url.Values, secret string) bool {
	// 对参数进行排序
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 创建一个新的缓冲区来存储排序后的参数
	var sortedParams bytes.Buffer

	thirdSign := ""

	// 打印排序后的参数
	for _, k := range keys {
		if k == "sign" {
			thirdSign = params[k][0]
			continue
		}
		values := params[k]
		for _, v := range values {
			if sortedParams.Len() > 0 {
				sortedParams.WriteByte('&')
			}
			sortedParams.WriteString(k)
			sortedParams.WriteByte('=')
			sortedParams.WriteString(v)
		}
	}

	sortedParams.WriteString("&key=")
	sortedParams.WriteString(secret)

	str := sortedParams.String()

	md := md5.New()
	md.Write([]byte(str))
	sign := strings.ToUpper(hex.EncodeToString(md.Sum(nil)))

	return sign == thirdSign
}

// 代付回调签名
func CheckPaymentCallbackSign(req request.XPaymentCallback, secret string) bool {
	signStr := StringJoin([]string{"accountName=", req.AccountName})

	if len(req.AccountNo) > 0 {
		signStr = StringJoin([]string{signStr, "&accountNo=", req.AccountNo})
	}

	signStr = StringJoin([]string{signStr, "&appId=", req.AppId, "&createdAt=", strconv.Itoa(req.CreatedAt), "&currency=", req.Currency})

	if len(req.ErrMsg) > 0 {
		signStr = StringJoin([]string{signStr, "&errMsg=", req.ErrMsg})
	}

	signStr = StringJoin([]string{
		signStr, "&mchNo=", req.MchNo,
		"&mchOrderNo=", req.MchOrderNo,
		"&orderAmount=", fmt.Sprintf("%v", req.OrderAmount),
		"&reqTime=", strconv.Itoa(req.ReqTime),
		"&state=", strconv.Itoa(req.State),
	})

	if req.SuccessTime > 0 {
		signStr = StringJoin([]string{signStr, "&successTime=", strconv.Itoa(req.SuccessTime)})
	}

	signStr = StringJoin([]string{signStr, "&transferId=", req.TransferId})

	if len(req.Utr) > 0 {
		signStr = StringJoin([]string{signStr, "&utr=", req.Utr})
	}

	if len(req.Vpa) > 0 {
		signStr = StringJoin([]string{signStr, "&vpa=", req.Vpa})
	}

	signStr = StringJoin([]string{signStr, "&key=", secret})

	sign := strings.ToUpper(GetMd5(signStr))

	if sign != req.Sign {
		global.Logger["err"].Infof("CheckPaymentCallbackSign failed sign:", sign)
	}

	return sign == req.Sign
}

// 代付下单
func TransferOrder(bank pay.BankInfo, payConfig pay.PayConfig, orderNo string, amount float64) (err error) {
	payOrderReq := response.PaymentOrderRequest{
		MchNo:           payConfig.Merchant,
		AppId:           payConfig.AppId,
		MchOrderNo:      orderNo,
		Currency:        "INR",
		OrderAmount:     fmt.Sprintf("%v", amount),
		BankCode:        bank.BankCode,
		AccountName:     bank.Name,
		AccountEmail:    bank.Email,
		AccountMobileNo: bank.Phone,
		TransferDesc:    "",
		NotifyUrl:       payConfig.PaymentBackUrl,
		ReqTime:         strconv.Itoa(int(time.Now().UnixMilli())),
	}
	if len(bank.BankName) > 1 {
		payOrderReq.BankName = bank.BankName
	}
	if len(bank.AccountNo) > 1 {
		payOrderReq.AccountNo = bank.AccountNo
	}
	if len(bank.Ifsc) > 1 {
		payOrderReq.Ifsc = bank.Ifsc
	}
	if len(bank.Address) > 1 {
		payOrderReq.AccountAddress = bank.Address
	}
	if len(bank.Vpa) > 1 {
		payOrderReq.Vpa = bank.Vpa
	}
	payOrderReq.Sign = GetPaymentSign(payOrderReq, payConfig.Secret)

	url := payConfig.Url + constant.XPayPaymentUrl
	res := response.PaymentOrderReturn{}
	err = req2.Call(url, payOrderReq, &res)

	if err != nil {
		global.Logger["err"].Infof("xPay TransferOrder req2.Call failed,orderNo:[%s], err: [%s]", orderNo, err.Error())
		return
	}

	//避免三方修改字段，导致原本成功发起的单被卡单，导致损失
	if res.Code == nil {
		global.Logger["err"].Infof("xPay TransferOrder req2.Call res.Code == nil")
		err = errUtil.NewUnmarshalErr()
		return
	}

	//订单状态 0-订单异常，1-代付中(回调)，2-代付成功（回调），3-代付下单失败(注意不回调，不回调，不回调) ，4-撤销
	if *res.Code == constant.XPayPaymentSuccessCode && (res.Data.State == constant.XPayPaymentStateProcessing || res.Data.State == constant.XPayPaymentStateSuccess) {
		//成功
		global.Logger["info"].Infof("xPay TransferOrder giveMoney order_no:%v,pay_cfg_id: %v", orderNo, payConfig.ID)
	} else {
		//失败
		global.Logger["err"].Infof("giveMoney order_no:%v,err: %s", orderNo, res.Data.ErrMsg)
		err = errors.New(res.Data.ErrMsg)
	}

	return
}
