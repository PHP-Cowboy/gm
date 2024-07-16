package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gm/utils/errUtil"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"gm/global"
)

type HttpRsp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func PostTest() {
	url := "http://121.196.60.92:19090/api/v1/remote/get/goods/by/id"
	method := "POST"

	payload := strings.NewReader(`{"order_id":[255]}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func Post(url string, responseData interface{}, timeout time.Duration) ([]byte, error) {

	client := &http.Client{}

	jData, err := json.Marshal(responseData)
	if err != nil {
		global.Logger["err"].Infof("json.Marshal failed,err: [%s]", err.Error())
		return nil, err
	}

	global.Logger["info"].Infof("url:%s,params:%s", url, string(jData))

	// 创建请求
	rq, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(jData))
	if err != nil {
		global.Logger["err"].Infof("http.NewRequest failed,err: %s", err.Error())
		return nil, err
	}

	// 设置请求头
	rq.Header.Set("Content-Type", "application/json")

	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // 确保在函数结束时取消上下文

	// 发送请求
	res, err := client.Do(rq.WithContext(ctx))

	if err != nil {
		global.Logger["err"].Infof("client.Do failed,err: %s", err.Error())
		return nil, err
	}

	defer res.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(res.Body)
	if err != nil {
		global.Logger["err"].Infof("io.ReadAll failed,err: [%s]", err.Error())
		return nil, err
	}

	global.Logger["info"].Infof(string(body))

	return body, nil
}

func SrvPost(url string, responseData interface{}) ([]byte, error) {

	client := &http.Client{}

	jData, err := json.Marshal(responseData)
	if err != nil {
		return nil, err
	}

	global.Logger["info"].Infof("url:%s", url)
	global.Logger["info"].Infof("params:%s", string(jData))

	rq, err := http.NewRequest("POST", url, bytes.NewReader(jData))

	if err != nil {
		return nil, err
	}

	rq.Header.Add("Content-Type", "application/json")
	rq.Header.Add("token", global.ServerConfig.Srv.Token)

	res, err := client.Do(rq)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	global.Logger["info"].Infof(string(body))

	return body, nil
}

func Get(url string) ([]byte, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func TestGet() ([]byte, error) {
	url := "http://121.196.60.92:19090/api/v1/remote/pick/shop/list"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// res 传指针
func Call(uri string, params interface{}, res interface{}) (err error) {
	var (
		body []byte
	)

	body, err = Post(uri, params, 10*time.Second)

	if err != nil {
		global.Logger["err"].Infof("Call Post failed, params:[%s], err: [%s]", params, err.Error())
		return
	}

	err = json.Unmarshal(body, &res)

	if err != nil {
		global.Logger["err"].Infof("Call json.Unmarshal failed, params:[%s], err: [%s]", params, err.Error())
		err = errUtil.NewUnmarshalErr()
		return
	}

	//rspCode, rspMsg := getRspCode(res)
	//
	//if rspCode != 200 {
	//	return errors.New(rspMsg)
	//}

	return
}

func getRspCode(rsp interface{}) (code int, msg string) {
	code = -1
	msg = "未知错误"

	if rsp == nil {
		return
	}
	v := reflect.ValueOf(rsp)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		//最多取两层
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}
	//kind := v.Kind()
	// 判断是否是结构体
	if v.Kind() != reflect.Struct {
		return
	}
	codeValue := v.FieldByName("Code")
	if !codeValue.IsValid() {
		return
	}

	msgValue := v.FieldByName("Msg")

	if !codeValue.IsValid() {
		return
	}

	code = int(codeValue.Int())

	msg = msgValue.String()

	return
}
