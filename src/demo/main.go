package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.shalltry.com/paynicorn/go-sdk/src/paynicornlogic"
)

func main_sample() {
	var appkey = "PUT_YOUR_APPKEY_HERE"
	var merchantkey = "PUT_YOUR_MERCHANT_SECRET_HERE"


	//raise a payment request to PAYNICORN
	cashierRequest := paynicornlogic.RaiseCashierRequest{}
	cashierRequest.OrderId="PUT_YOUR_ORDER_ID_HERE"
	cashierRequest.CountryCode="NG"
	cashierRequest.Currency="NGN"
	cashierRequest.Amount="10"
	cashierRequest.CpFrontPage="PUT_YOUR_WEB_REDIRECT_URL_HERE"
	cashierRequest.OrderDescription="TEST GOODS NAME"
	cashierresponse,_ := paynicornlogic.RaiseCashierPayment(appkey,merchantkey,cashierRequest)
	fmt.Println(cashierresponse)


	//query a payment status from PAYNICORN
	request := paynicornlogic.QueryTransactionRequest{}
	request.OrderId=cashierRequest.OrderId
	request.TxnType=paynicornlogic.PAYMENT
	response := paynicornlogic.QueryPaymentStatus(appkey,merchantkey,request)
	fmt.Println(response)


	//receive a payment status postback from PAYNICORN
	r := gin.Default()
	r.POST("/postback", func(context *gin.Context) {
		postbackresponse,_ := paynicornlogic.Postback(context,merchantkey)
		fmt.Println(postbackresponse)

	})
	r.Run(":80")


}
