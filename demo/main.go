package main

import (
	"fmt"
	"github.com/Transsion-mios/paynicorn-sdk-golang"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	var appkey = "PUT_YOUR_APPKEY_HERE"
	var merchantkey = "PUT_YOUR_MERCHANT_SECRET_HERE"


	//raise a payment request to PAYNICORN
	request := paynicorn.InitPaymentRequest{}
	request.OrderId="PUT_YOUR_ORDER_ID_HERE"
	request.CountryCode="NG"
	request.Currency="NGN"
	request.Amount="10"
	request.CpFrontPage="PUT_YOUR_WEB_REDIRECT_URL_HERE"
	request.OrderDescription="TEST GOODS NAME"
	response := paynicorn.InitPayment(appkey,merchantkey,request)
	if response != nil{
		fmt.Println(response)
	}



	//query a payment status from PAYNICORN
	request1 := paynicorn.QueryPaymentRequest{}
	request1.OrderId=request.OrderId
	request1.TxnType=paynicorn.PAYMENT
	response1 := paynicorn.QueryPayment(appkey,merchantkey,request1)
	if response1 != nil{
		fmt.Println(response1)
	}


	//query support payment method from PAYNICORN
	request2 := paynicorn.QueryMethodRequest{}
	request2.TxnType = paynicorn.PAYMENT
	request2.CountryCode = "NG"
	request2.Currency = "NGN"
	response2 := paynicorn.QueryMethod(appkey,merchantkey,request2)
	if response2 != nil{
		fmt.Println(response2)
	}


	//receive a payment status postback from PAYNICORN
	r := gin.Default()
	r.POST("/postback", func(context *gin.Context) {

		var req paynicorn.PostbackRequest
		if err := context.BindJSON(&req); err != nil{
			context.String(http.StatusInternalServerError,"")
		}else{
			postbackInfo := paynicorn.ReceivePostback(merchantkey,req)
			if postbackInfo != nil && postbackInfo.Verified{
				fmt.Println(postbackInfo)
				context.String(http.StatusOK,"success_"+postbackInfo.TxnId)
			}else{
				context.String(http.StatusInternalServerError,"")
			}
		}

	})
	r.Run(":8080")


}
