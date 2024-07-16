package daos

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"gm/common/constant"
	"gm/global"
	"gm/model/pay"
	"gm/request"
	"gm/response"
	"gm/utils/ecode"
	"gm/utils/errUtil"
	req2 "gm/utils/request"
	"strconv"
	"strings"
)

// 查询余额
func InPayBalance(merchantId, privateKey string) (balance float64, err error) {
	params := request.InPayBalance{
		MerchantId: merchantId,
		Sign:       GetMd5(merchantId + privateKey), //test_sign_key
	}

	var jsonTxt string

	url := constant.InPayBalance

	jsonTxt, err = JsonString(params)
	if err != nil {
		global.Logger["err"].Errorf("InPayBalance JsonString failed,err:[%v]", err.Error())
		return
	}

	res := response.InPayBalanceRsp{}

	err = req2.Call(url, jsonTxt, &res)
	if err != nil {
		global.Logger["err"].Errorf("InPayBalance req2.Call failed,err:[%v]", err.Error())
		return
	}

	if res.Code != 100 {
		err = ecode.New(res.Code, res.Msg)
		return
	}

	balance, err = strconv.ParseFloat(res.Data.Balance, 64)

	return
}

// 代付
func InPayPayment(cfg pay.PayConfig, orderNo string, amount float64, bank pay.BankInfo) (err error) {
	url := constant.InPayWithdrawTest

	if global.ServerConfig.Mode == "pre" || global.ServerConfig.Mode == "release" {
		url = constant.InPayWithdraw
	}

	params := request.InPayWithdraw{
		MerchantId:  cfg.Merchant,
		OrderNumber: orderNo,
		OrderAmount: fmt.Sprintf("%v", amount),
		Type:        "BANK", //支付类型填： BANK/UPI（任选其一）
		Vpa:         bank.Vpa,
		Email:       bank.Email,
		Account:     bank.AccountNo,
		Name:        bank.Name,
		Ifsc:        bank.Ifsc,
		Phone:       bank.Phone,
		NotifyUrl:   cfg.PaymentBackUrl,
	}

	var p []byte

	p, err = json.Marshal(params)
	if err != nil {
		global.Logger["err"].Errorf("InPayPayment json.Marshal failed,err:[%v]", err.Error())
		return
	}

	var (
		privateKey *rsa.PrivateKey
		signByte   []byte
	)

	key := `-----BEGIN RSA PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAI9uLqgB+aEfkJP+vkti8L8jFa8KIbNNW/qYdDm46I36L/iefS2vXQcB3YQxn2TCg9Bj4LImYNC5vBG01kMm22asd1Xhg/pyiQksY6cRkm2I3H0MAwBWsLSKSyTvdco+cBWWBwco0rEassL6cmzqhmN/IcHB4j5E2WYJjmzD2FpdAgMBAAECgYAtn1B75FTw6Udlq8v/0rLdOV22VbSugdbV+RRPH//o2UHVBSSwGW2vwuohGF+o/y5KZNqeEBMPkWS/BRR6O6vhm/qbJjWKkcUsTMX6UKEAbUJkF/oz47ihjSb0TgY5ynf2+vLyTHvztXpX7d+E0IuPGmjEgqdC+7QbzZaUu7l0SQJBANi17hYJayxD18exVOgj99WsU8Ra0yi4J3ckU1YAPZDzifEy+T2TGgzATTUvXS9JcqhY33SRpIYPYNR3Ov6N9C8CQQCpbyL298l31nI8N222P5vdWXTTQibFI0vX4SuURTDotMwOatHGn/DVCdl2FXe44CdER48FxXEolbYZk9Z/fFszAkAvvs2Iz2RsaToWRgyl5J7K1d/SyAvz0bboOfmeXgkycWW33XoqRcmce5XHHPtT2sPHMWVyAlCUNLkptmcqBE6DAkB4OkxtQtbLGnhlEk/firNnFhs37TDlom4m+biatZ5HAkPp1xKUBto10Y9lo0YZAbXbVYu/ZKMvUUyuSaFjRTLRAkAUZ9nKeT+6tkx/X+QjqGokbOsv/ZMlSAjmGZg9QTIKyyzq8jfrUMV4Bb4aQi+rWTmj3HfnlljoXeYq0FxyFG/4
-----END RSA PRIVATE KEY-----`

	privateKey, err = parsePrivateKey(key)
	if err != nil {
		global.Logger["err"].Errorf("InPayPayment parsePrivateKey failed,err:[%v]", err.Error())
		return
	}

	signByte, err = PrivateEncrypt(privateKey, p)
	if err != nil {
		global.Logger["err"].Errorf("InPayPayment PrivateEncrypt failed,err:[%v]", err.Error())
		return
	}

	sign := base64.StdEncoding.EncodeToString(signByte)

	params.Sign = sign

	var (
		res response.InPayWithdrawRsp
	)

	err = req2.Call(url, params, &res)

	if err != nil {
		global.Logger["err"].Infof("InPayPayment req2.Call failed,orderNo:[%s], err: [%s]", orderNo, err.Error())

		return err
	}

	if res.Code == nil {
		global.Logger["err"].Infof("InPayPayment req2.Call res.Code == nil")
		err = errUtil.NewUnmarshalErr()
		return
	}

	if *res.Code != constant.InPayPaymentSuccessCode {
		global.Logger["err"].Errorf("InPayPayment failed, res.Code:[%v],res.Msg:[%v]", res.Code, res.Msg)
		err = errors.New(res.Msg)

	} else {
		//成功
		global.Logger["info"].Infof("InPayPayment success， giveMoney order_no:%v,pay_cfg_id: %v", orderNo, cfg.ID)
	}

	return
}

// Encrypt 使用公钥对数据进行RSA加密
func Encrypt(publicKey *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	label := []byte("") // OAEP padding的label参数，可以为空
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, label)
}

// Decrypt 使用私钥对数据进行RSA解密
func Decrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
}

// Encrypt 使用私钥对数据进行RSA加密
func PrivateEncrypt(privateKey *rsa.PrivateKey, plaintext []byte) (signature []byte, err error) {
	// 签名
	hash := sha256.New()
	hash.Write(plaintext)
	hashed := hash.Sum(nil)

	signature, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	return
}

func pubDecrypt(cryptedString string) (err error) {
	publicKeyString := global.ServerConfig.InPay.PublicKey
	// 将公钥字符串转换为PEM格式
	pemBlock, _ := pem.Decode([]byte(addPEMHeaderFooter(publicKeyString)))
	if pemBlock == nil {
		global.Logger["err"].Infof("failed to decode PEM block containing the public key")
		err = errors.New("failed to decode PEM block containing the public key")
		return
	}

	// 解析公钥
	_, err = x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		global.Logger["err"].Infof("failed to parse DER encoded public key err:", err.Error())
		err = errors.New("failed to parse DER encoded public key: " + err.Error())
		//fmt.Errorf("failed to parse DER encoded public key: %w", err)
		return
	}

	crypted := []byte(cryptedString)

	// cryptedString 签名，请商户使用平台公钥解密，如果能直接解出参数来就是成功

	// 模拟的分块处理（这不是真正的公钥操作）
	chunkSize := 128 // 假设的块大小
	var decrypt []byte
	for len(crypted) > 0 {
		if len(crypted) < chunkSize {
			chunkSize = len(crypted)
		}
		chunk := crypted[:chunkSize]
		crypted = crypted[chunkSize:]

		// 在这里，我们将只是简单地返回数据，因为公钥不直接用于“解密”
		// 你可能需要根据实际情况执行某种验证或操作
		decrypt = append(decrypt, chunk...)
	}

	// 返回模拟“解密”后的数据（实际上只是原始数据的副本）
	return nil
}

// 辅助函数：为公钥字符串添加PEM头部和尾部
func addPEMHeaderFooter(publicKey string) string {
	return "-----BEGIN PUBLIC KEY-----\n" +
		strings.ReplaceAll(chunkSplit(publicKey, 64), "\n", "\n"+strings.Repeat(" ", 4)) +
		"\n-----END PUBLIC KEY-----"
}

// 辅助函数：将字符串按指定长度分块，并添加换行符（与PHP的chunk_split类似）
func chunkSplit(str string, chunkSize int) string {
	var chunks []string
	runes := []rune(str)
	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return strings.Join(chunks, "\n")
}
