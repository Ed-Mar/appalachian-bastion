package saga

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Saga struct {
	Name         string
	Transactions []*Transaction
}

const retriable = "retriable"
const compensatable = "compensatable"
const pivot = "pivot"

type Transaction struct {
	ServiceName          string
	TransactionType      string
	sagaTransactionProxy *TransactionRESTProxy
}

type TransactionRESTProxy struct {
	name         string
	httpFunction string
	url          string
	body         *io.Reader
	Response     *http.Response
}
type ExecuteTransactionViaRestProxy interface {
}

// NewTransactionRestProxy is a constructor for the Transaction REST Proxy creation just to make sure values make sense before making that call to the other services, or it own service
func (trp *TransactionRESTProxy) NewTransactionRestProxy(name string, httpFunction string, url string, body *io.Reader) error {
	trp.name = name
	groomedHTTPFunction := strings.ToUpper(strings.TrimSpace(httpFunction))
	switch groomedHTTPFunction {
	case "GET":
		{
			if body != nil {
				ErrMismatchHttpFunctionWithBody := fmt.Errorf("\"[ERROR] [SAGA] [TRANSACTION] [REST PROXY] [CREATION] | Cannot have body with this %+v | GET or DELETE", httpFunction)
				return ErrMismatchHttpFunctionWithBody
			} else {
				trp.httpFunction = groomedHTTPFunction
				trp.body = body
			}

		}
	case "DELETE":
		{
			if body != nil {
				ErrMismatchHttpFunctionWithBody := fmt.Errorf("\"[ERROR] [SAGA] [TRANSACTION] [REST PROXY] [CREATION] | Cannot have body with this %+v | GET or DELETE", httpFunction)
				return ErrMismatchHttpFunctionWithBody
			} else {
				trp.httpFunction = groomedHTTPFunction
				trp.body = body
			}
		}
	case "UPDATE":
		{
			if body == nil {
				ErrMissingBody := fmt.Errorf("\"[ERROR] [SAGA] [TRANSACTION] [REST PROXY] [CREATION] | Cannot have missingbody for this fucntion %+v | POST or UPDATE", httpFunction)
				return ErrMissingBody
			} else {
				trp.httpFunction = groomedHTTPFunction
				trp.body = body
			}
		}
	case "POST":
		{
			if body == nil {
				ErrMissingBody := fmt.Errorf("\"[ERROR] [SAGA] [TRANSACTION] [REST PROXY] [CREATION] | Cannot have missingbody for this fucntion %+v | POST or UPDATE", httpFunction)
				return ErrMissingBody
			} else {
				trp.httpFunction = groomedHTTPFunction
				trp.body = body
			}
		}
	default:
		{
			ErrNonSupportedHTTPFunction := fmt.Errorf("[ERROR] [SAGA] [TRANSACTION] [REST PROXY] [CREATION] | %+v is not a supported http Fuction: Supported Functions: POST, GET, UPDATE, and DELETE", httpFunction)
			return ErrNonSupportedHTTPFunction
		}
	}
	//TODO Figure out if should do any check on the url
	//I am guessing I should check if it in the correct domain like internal, but later
	//That means I need a file where I keep the domain at for later
	trp.url = url

	return nil
}
func (trp *TransactionRESTProxy) executeTransactionViaRestProxy() {

	client := &http.Client{}
	req, err := http.NewRequest(trp.httpFunction, trp.url, *trp.body)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	trp.Response = resp
}
