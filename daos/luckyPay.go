package daos

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"gm/common/constant"
	"gm/global"
	"gm/model/pay"
	"gm/request"
	"gm/response"
	"gm/utils/errUtil"
	req2 "gm/utils/request"
	"regexp"
	"strconv"
	"time"
)

var luckyPubKey = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjOrF8H78wLlVUd0Xe3qe+mBvkzchxccuBqtPZGbK9ebhXbqOaIgRQByjWSuC1/qJcT6eBumSX/iHUVaTdf3v2XO1QCbUZZm74vcSsO5f64bmzTtk/cLa1S4068Rfde0YoX/avNfiX2nRz1SBzy3zaiwA8pmhkbpU+L3k+a8qd3WH3Hq2MpX/dD1haqzTHYiA8+xJANixvNtKWq5AMyrs9rSkwGTEHV33ETLUZgpHjtgI+L/PAqzIaW89dgh4EIga4sNGr/3iUkZU1vYKj2PHY6/pfgs5wCS0BhxHF35LYCHgyZV7UH4oxu0kOtGGmdOAeIylCWRXdPkaXGGnOlmYmQIDAQAB
-----END RSA PUBLIC KEY-----`

var luckyPriKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCK9LOBPXl1JblCNEytEmhlxPTWMnRtJ8C7+3dOiBUomTBDC1ifSD01bnMbsXuuhhmYv/1HWYiGDvr5vWDp293lhsGIg/iJ8WIw8xUT02fBjGCiJ7EWkLp0UPX0qP6xOBWsJeaZFwYkBYpxr9mrnjj3sAM0x6XB8n7MJeotF5it1oAyoQMbfzXcAxZH5ZRz6FSjwKj70cb9x3rO9VrTgfmlDk+SgmDHokF2Up7jLAOTwthpXBWC9GGTzi7IcLcDjgGPAzlj/fL5G2Ya1i71XdMYoqU7OOy85TczqueTX1FAAWpAu9kkKPq9iUr3rtlySwQNFaQXX2ys4rAGS6xs5i0JAgMBAAECggEAVilHSQMSVeZvSjLdvI9mdnw7CLo+YRS+OJHDt8k/vW3HdDsL17BWa9QxweE6fMwgAQxDA8PrmJyfWM3p0viHTfRR+tFCAl7PcOS5lBKJCoV9we50qzZQtyEgcXV8f+zz20uhAho5r4pT/wx4Hyc70bGtO0P0bx9rKjD0VMOF43vs3+UN9xeLZi7XQOlbES3rFNQUyL2onqoyApmOEM146WuaYoas26fWbTMHKi8WiXfWzRxBt2rGMW80VDQsUmjcniU18eJq2QhjpSevicsNbnyEBPU5m80lgkavZ8aCAda4H1tnPf2CWo6wynIejPC2VniLUeept8VeFMxm4Gx30QKBgQDFEec5YeYj9bAyyCceaZ+yv5xQvgusA4zseqOuwraj4cZ8+pS78aP3xv401W06E6tr0SxAFsNPviO6q0bQiMaqCTG4rZuVCsWpWYcfrK9s+aSZLYXdxQCDHO921+LxkFo06SHjYTbtv15vk8dulLQIhxznpre+9dkIIIJ95WvffQKBgQC0ggzGRzSGIwNr+g4APtRVF+llDlQEFRd7Ph2LXjbD5PsH4bJ39hK8butB/IU1tnSRzPpQDeBI+XIMOS51PT4FdoyNX8r7JjLvArDru1ghNJtdoaB+/TFkz50o6cOJSxfXLjC4WYp2hk2cMOcAbqHEDBYRlwhhiVu0KxeFMenRfQKBgHJHM4K/Fqn7qzWHg8fLEGSjYI4h5rrVRD9NzuVk3GykXGnVFbL3KVJh/r+8lB59kYZwQezYRmJTrHxvHh2Hc0zfEAo8AmCnYoAV/pmLlh7nlV+OWnaD6wwdF3AfOOdEAkt0dcZZdXTg6G5jj968NLpNP6HFt//wqO5hi8pd4/QxAoGAM0x2Vu8rhSd0NH6G9hjk/R3jjX2p8NMRrkuz18S5qahskwvYTXFYV8bqAwHB1cb5j+oCBTg8UmDZwZGKm3UBKEpNnWvo4sEzXmuUXhoK7Lznno9tbkmEfRLnphXxJRZ7OwL8g5em2xJGAip/q8bFIFMS/oK+tgF0V1qbv5W0zn0CgYACRyPgNPse6cIZA0bM7YEVZL3fwrR7gD0ydVSPvTvTQTRfbqIRhsRKK9d8BmNFeFVc9QvLtl3waYbOlog0AcwrC3NbmB1MiMqPmAsm45q6Q8I2ENf2nDBS288/5iRrH9Mhdi435c1sobd+jNtzpdYSL1kvvH0yyG+meGMFEVlIOQ==
-----END RSA PRIVATE KEY-----`

// isIFSC 检查字符串是否可能是IFSC代码
func isIFSC(code string) bool {
	// 假设IFSC是11位，并且只包含字母数字字符
	if len(code) != 11 {
		return false
	}
	// 这里可以添加更复杂的验证逻辑，但出于简化考虑，我们只检查长度和字符类型
	for _, char := range code {
		if (char < '0' || char > '9') && (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') {
			return false
		}
	}
	return true
}

// isVPA 检查字符串是否可能是VPA
func isVPA(address string) bool {
	// 使用正则表达式检查字符串中是否包含@符号
	vpaRegex := regexp.MustCompile(`[^@]+@[^@]+\.[^@]+`)
	return vpaRegex.MatchString(address)
}

// identifyCodeType 尝试确定传入的字符串是IFSC还是VPA
func identifyCodeType(code string) string {
	if isIFSC(code) {
		return "IFSC"
	}
	if isVPA(code) {
		return "UPI" //VPA账号，类型填UPI
	}
	return "Unknown"
}

// 代付
func LuckyPayPayment(cfg pay.PayConfig, orderNo string, amount float64, bank pay.BankInfo) (err error) {

	param := request.LuckyPayment{
		AccountCode:  bank.Ifsc,
		AccountEmail: bank.Email,
		AccountName:  bank.Name,
		AccountNo:    bank.AccountNo,
		AccountPhone: bank.Phone,
		AccountType:  "IFSC",
		Currency:     "INR",
		CustomerIp:   "127.0.0.1",
		MchNo:        cfg.Merchant,
		MchOrderNo:   orderNo,
		NotifyUrl:    cfg.PaymentBackUrl,
		PayAmount:    fmt.Sprintf("%2.f", amount),
		ReqTime:      strconv.FormatInt(time.Now().UnixMilli(), 10),
		Summary:      "xxx",
	}

	var (
		privateKey *rsa.PrivateKey
		signByte   []byte
	)

	signStr := StringJoin([]string{
		"accountCode=", param.AccountCode,
		"&accountEmail=", param.AccountEmail,
		"&accountName=", param.AccountName,
		"&accountNo=", param.AccountNo,
		"&accountPhone=", param.AccountPhone,
		"&accountType=", param.AccountType,
		"&currency=", param.Currency,
		"&customerIp=", param.CustomerIp,
		"&mchNo=", param.MchNo,
		"&mchOrderNo=", param.MchOrderNo,
		"&notifyUrl=", param.NotifyUrl,
		"&payAmount=", param.PayAmount,
		"&reqTime=", param.ReqTime,
		"&summary=", param.Summary,
	})

	privateKey, err = parsePrivateKey(luckyPriKey)
	if err != nil {

		return
	}

	signByte, err = privateEncrypt(privateKey, []byte(signStr))
	if err != nil {

		return
	}

	param.Sign = base64.StdEncoding.EncodeToString(signByte)

	var (
		res response.LuckyPayPaymentRsp
	)

	url := constant.LuckyPayWithdraw

	err = req2.Call(url, param, &res)
	if err != nil {
		global.Logger["err"].Infof("LuckyPayPayment req2.Call failed,orderNo:[%s], err: [%s]", orderNo, err.Error())
		return err
	}

	if res.Code == nil {
		global.Logger["err"].Infof("LuckyPayPayment req2.Call res.Code == nil")
		err = errUtil.NewUnmarshalErr()
		return
	}

	//支付状态（1-代付中、2-代付成功 、3-代付失败 ）
	if *res.Code == constant.LuckyInPayPaymentSuccessCode && (res.Data.PayState == constant.LuckyInPayPaymentStateProcessing || res.Data.PayState == constant.LuckyInPayPaymentStateSuccess) {
		//成功
		global.Logger["info"].Infof("LuckyPayPayment giveMoney order_no:%v,pay_cfg_id: %v", orderNo, cfg.ID)
	} else {
		//失败
		global.Logger["err"].Infof("LuckyPayPayment failed,giveMoney order_no:%v,err: %s", orderNo, res.Message)
		err = errors.New(res.Message)
	}

	return
}

// 代付回调验签
func LuckyPayCallbackCheckSign(req request.LuckyPayCallback) (err error) {

	signStr := StringJoin([]string{"currency=", req.Currency})

	if req.ErrMsg != "" {
		signStr = StringJoin([]string{signStr, "&errMsg=", req.ErrMsg})
	}

	signStr = StringJoin([]string{
		signStr,
		"&mchNo=", req.MchNo,
		"&mchOrderNo=", req.MchOrderNo,
		"&payAmount=", req.PayAmount,
		"&payFinishTime=", req.PayFinishTime,
		"&payInitiateTime=", req.PayInitiateTime,
		"&payOrderNo=", req.PayOrderNo,
		"&payState=", strconv.Itoa(req.PayState),
	})

	// 假设这是你要验证的签名和原始数据
	originalData := []byte(signStr)

	// 使用SHA256哈希算法计算原始数据的哈希值
	hashed := sha256.Sum256(originalData)

	var (
		rsaPub *rsa.PublicKey
	)

	rsaPub, err = parsePublicKey(luckyPubKey)
	if err != nil {
		logger.Logger.Errorf("(server InitService) GetLuckyPayNotifyPay: parsePublicKey err:[%v]", err.Error())
		return err
	}

	// 将Base64编码的签名解码为字节切片
	signature, err := base64.StdEncoding.DecodeString(req.Sign)
	if err != nil {
		logger.Logger.Errorf("(server InitService) GetLuckyPayNotifyPay: CheckSign base64.StdEncoding.DecodeString err:[%v]", err.Error())
		return
	}

	// 验证签名
	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed[:], signature)
	if err != nil {
		logger.Logger.Errorf("(server InitService) GetLuckyPayNotifyPay: rsa.VerifyPKCS1v15 err:[%v]", err.Error())
		return
	}

	return
}
