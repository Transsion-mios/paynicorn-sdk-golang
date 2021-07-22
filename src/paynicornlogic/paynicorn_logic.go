package paynicornlogic

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

var(
	SUCCESS_CODE = "000000"
)

type TxnTypeEnum string

const (
	PAYMENT TxnTypeEnum = "payment"

	PAYOUT TxnTypeEnum = "payout"

	AUTHPAY TxnTypeEnum = "authpay"

	REFUND TxnTypeEnum = "refund"
)

/**
request body struct
content is the base64 string of request json data
sign is the md5 value of content

 */
type RequestBody struct{
	Content string `json:"content"`
	Sign string `json:"sign"`
	AppKey string `json:"appKey"`
}

/**
paynicorn response body struct
 */
type ResponseBody struct{
	ResponseCode string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Content string `json:"content"`
	Sign string `json:"sign"`
}


type PostbackResponse struct{
	TxnId string `json:"txnId"`
	OrderId string `json:"orderId"`
	Amount string `json:"amount"`
	Currency string `json:"currency"`
	CountryCode string `json:"countryCode"`
	Status string `json:"status"`
	Code string `json:"code"`
	Message string `json:"message"`
}

type PostbackRequest struct{
	Content string `json:"content"`
	Sign string `json:"sign"`
}

/**
query  transaction status request
 */
type QueryTransactionRequest struct{
	OrderId string `json:"orderId"`//MANDATORY unique transaction id generate by paynicorn.you can use it to query your transaction status or wait for a postback
	TxnType  TxnTypeEnum `json:"txnType"`
}

/**
query  transaction status response
 */
type QueryTransactionResponse struct{

	Code	string    `json:"code"` // MANDATORY response code represent current request is success response or not refer to https://www.paynicorn.com/#/docs 6.5

	Message	string `json:"message"` //	MANDATORY response code refer to https://www.paynicorn.com/#/docs 6.5

	TxnId	string `json:"txnId"`//	MANDATORY unique transaction id generate by paynicorn.you can use it to query your transaction status or wait for a postback

	Status	string `json:"status"`//	OPTIONAL transaction status (1:for success；-1：processing；0：failure)

	CompleteTime	string `json:"completeTime"` //OPTIONAL transaction complete time
}

/**
cashier payment request model
*/
type RaiseCashierRequest struct {
	Amount string `json:"amount"`//mandatory ,local currency for a country.the range of amount is depend on currency.refer to develop docs https://www.paynicorn.com/#/docs 6.1

	CountryCode string `json:"countryCode"` //mandatory,country code define in iso 3166 alpha-2 code ，refer to develop docs https://www.paynicorn.com/#/docs 6.1

	Currency string `json:"currency"` //mandatory, currency short cod define in iso4217，refer to https://www.paynicorn.com/#/docs 6.1

	OrderId string `json:"orderId"`//mandatory,unique id for a transaction，if a request has the same orderId as an old request.you will get the same response as old request.

	OrderDescription string `json:"orderDescription"` //mandatory,this field will show on paynicorn cashier

	PayMethod string `json:"payMethod"`//optional, payment method refer to https://www.paynicorn.com/#/docs 6.3 if this filed is set paynicorn will use the payment method you set,otherwise paynicorn will show all available payment method on cashier

	Language string `json:"language"`//optional,default paynicorn will set language as user device local language.if you set we use it

	CpFrontPage string `json:"cpFrontPage"`//optional,if a payment is done ,paynicorn will redirect to this uri.

	UserId string `json:"userId"`//optional,deprecated

	Email string `json:"email"`//optional,user email kyc required

	Phone string `json:"phone"`//optional,user phone kyc required

	PayByLocalCurrency string `json:"payByLocalCurrency"`//optional ,if you use just one currency to price all your service or goods set it true paynicorn will change your currency to local currency

	Memo string `json:"memo"`//optional,you send it to paynicorn,paynicorn will return it back.
}

/**
raise cashier response
*/
type RaiseCashierResponse struct {

	Code	string    `json:"code"` // MANDATORY response code represent current request is success response or not refer to https://www.paynicorn.com/#/docs 6.5

	Message	string `json:"message"` //	OPTIONAL response code refer to https://www.paynicorn.com/#/docs 6.5

	TxnId	string `json:"txnId"`//	OPTIONAL unique transaction id generate by paynicorn.you can use it to query your transaction status or wait for a postback

	Status	string `json:"status"`//	OPTIONAL transaction status (1:for success；-1：processing；0：failure)

	WebUrl	string `json:"webUrl"`//	OPTIONAL paynicorn cashier uri

}




/**
raise a payment cashier，most time you will get a web url and you need to open it in webview or browse
appKey ：merchant creat a app will get a appKey,refer to https://console.paynicorn.com/#/app/apply
merchantSecret ：merchant's secret use it to sign your data ,refer to https://console.paynicorn.com/#/developer
RaiseCashierRequest ：raise an online payment cashier parameters
 */
func RaiseCashierPayment(appKey string,merchantSecret string,request RaiseCashierRequest)(raiseCashierResponse RaiseCashierResponse ,err error){
	raiseCashierResponse = RaiseCashierResponse{}
	logger,_ := zap.NewProduction()
	defer logger.Sync()
	paynicornUrl := "https://api.paynicorn.com/trade/v3/transaction/pay"

	jsonStr,err :=json.Marshal(request)
	var base64Str string//
	if err == nil{
		base64Str = base64.StdEncoding.EncodeToString(jsonStr)
	}


	b := RequestBody{}

	b.Content = base64Str
	b.Sign = fmt.Sprintf("%x",md5.Sum([]byte(b.Content+merchantSecret)))
	b.AppKey = appKey



	client := &http.Client{}

	jsonBytes, err := json.Marshal(b)
	if err != nil {
		logger.Error("json marshal meet error %v",zap.Error(err))
	}

	postRequest, err := http.NewRequest("POST", paynicornUrl, strings.NewReader(string(jsonBytes)))
	if err != nil {
		logger.Error("NewRequest meet error %v",zap.Error(err))
	}

	postRequest.Header.Add("Content-Type", "application/json")

	var buffer []byte
	if response, err := client.Do(postRequest); err != nil {

	} else {
		if buffer, err = ioutil.ReadAll(response.Body); err != nil {
			logger.Error("post to  meet error %v",zap.Error(err))
		}
	}

	rsp := ResponseBody{}
	err = json.Unmarshal(buffer, &rsp)


	if rsp.ResponseCode == SUCCESS_CODE{

		if sign := fmt.Sprintf("%x",md5.Sum([]byte(rsp.Content+merchantSecret))); sign == rsp.Sign{

			content, err := base64.StdEncoding.DecodeString(rsp.Content)
			if err == nil {
				err = json.Unmarshal([]byte(content),&raiseCashierResponse)
				return raiseCashierResponse,err

			}
		}
	}else{
		raiseCashierResponse.Code = rsp.ResponseCode
		raiseCashierResponse.Message = rsp.ResponseMessage
	}
	return raiseCashierResponse,err
}








func Postback(context *gin.Context,merchantSecret string ) (postbackResponse PostbackResponse,err error){

	var req PostbackRequest
	postbackResponse = PostbackResponse{}
	logger,_ := zap.NewProduction()
	defer logger.Sync()
	if err := context.BindJSON(&req); err != nil{
		fmt.Println("bind json failed")
		postbackResponse.Code="204005"
		postbackResponse.Message="bind json failed"
		return postbackResponse,err
	}else{
		if s := fmt.Sprintf("%x",md5.Sum([]byte(req.Content+merchantSecret))); req.Sign == s{

			c, err := base64.StdEncoding.DecodeString(req.Content)

			if err == nil {
				err = json.Unmarshal(c,postbackResponse)
				if err == nil{
					return postbackResponse,err
				}
			}
		}else{
			fmt.Println("sign verify failed")
			postbackResponse.Code = "204005"
			postbackResponse.Message = "sign verify failed"
		}
		return postbackResponse,err
	}
}


func QueryPaymentStatus(appKey string,merchantSecret string,queryTransactionRequest QueryTransactionRequest) QueryTransactionResponse{
	queryTransactionResponse := QueryTransactionResponse{}
	logger,_ := zap.NewProduction()
	defer logger.Sync()

	jsonStr,err := json.Marshal(queryTransactionRequest)


	requestBody := RequestBody{}

	requestBody.Content =  base64.StdEncoding.EncodeToString(jsonStr)
	requestBody.Sign = fmt.Sprintf("%x",md5.Sum([]byte(requestBody.Content+merchantSecret)))
	requestBody.AppKey = appKey



	client := &http.Client{}

	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {

		queryTransactionResponse.Code = "204005"
		queryTransactionResponse.Message = "query fail,retry later"
		return queryTransactionResponse
	}

	request, err := http.NewRequest("POST", "https://api.paynicorn.com/trade/v3/transaction/query", strings.NewReader(string(jsonBytes)))
	if err != nil {
		queryTransactionResponse.Code = "204005"
		queryTransactionResponse.Message = "query fail,retry later"
		return queryTransactionResponse
	}

	request.Header.Add("Content-Type", "application/json")

	var buffer []byte
	if response, err := client.Do(request); err != nil {

	} else {
		if buffer, err = ioutil.ReadAll(response.Body); err != nil {
			queryTransactionResponse.Code = "204005"
			queryTransactionResponse.Message = "query fail,retry later"
			return queryTransactionResponse
		}
	}

	rsp := ResponseBody{}
	err = json.Unmarshal(buffer, &rsp)


	if rsp.ResponseCode == "000000"{

		if sign := fmt.Sprintf("%x",md5.Sum([]byte(rsp.Content+merchantSecret))); sign == rsp.Sign{

			content, err := base64.StdEncoding.DecodeString(rsp.Content)
			if err == nil {
				err = json.Unmarshal(content, &queryTransactionResponse)
				return queryTransactionResponse

			}
		}
	}else{
		queryTransactionResponse.Code = rsp.ResponseCode
		queryTransactionResponse.Message = rsp.ResponseMessage
		return queryTransactionResponse
	}

	return queryTransactionResponse

}

