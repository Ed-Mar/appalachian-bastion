package holdthis

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Saga struct {
	SageName     string
	Transactions []*Transaction
}

type Transaction struct {
	TransactionName string
	ServiceName     string
	TransactionAction
}

func (t *Transaction) GetTransaction() *Transaction {
	return t
}

//TODO find a better name for this. "TransactionAction"
// I have no idea wtf to call this. I have spent four+ hours looking for what to call this.
// It's the action part of the Transaction. I am floored on what to call this.
// If anyone ever read this please let me.
// Cause a Transaction has data and it's a function at the same time, technically its a noun but idk

//TransactionAction generic for the action part of the Transaction
type TransactionAction interface {
	Execute() error
}

//These Definitions are straight form the book
//Microservices Patterns By Chris Richardson

//RetryableTransaction Transactions that follow the pivot Transaction and are guaranteed to succeed.
type RetryableTransaction struct {
	Transaction
}

//NewRetriableTransaction Constructor for RetryableTransaction
func NewRetriableTransaction(transactionName string, serviceName string, transactionAction TransactionAction) *RetryableTransaction {
	return &RetryableTransaction{Transaction{
		TransactionName:   transactionName,
		ServiceName:       serviceName,
		TransactionAction: transactionAction,
	}}
}

//CompensatableTransaction Transaction that can potentially be rolled back using a compensating Transaction.
type CompensatableTransaction struct {
	Transaction
	//The is the action that un does this Transaction if it fails
	CompensatableAction TransactionAction
}

//NewCompensatableTransaction Constructor for CompensatableTransaction
func NewCompensatableTransaction(transactionName string, serviceName string, transactionAction TransactionAction, compensatableAction TransactionAction) *CompensatableTransaction {
	return &CompensatableTransaction{
		Transaction: Transaction{
			TransactionName:   transactionName,
			ServiceName:       serviceName,
			TransactionAction: transactionAction,
		},
		CompensatableAction: compensatableAction,
	}
}

//PivotTransaction The go/no-go point in a saga. If the pivot Transaction commits, the saga will run until completion. A pivot Transaction can be a Transaction that’s neither compensatable nor retriable. Alternatively, it can be the last compensatable Transaction or the first retriable Transaction.
type PivotTransaction struct {
	Transaction
}

//NewPivotTransaction Constructor for NewPivotTransaction
func NewPivotTransaction(transactionName string, serviceName string, transactionAction TransactionAction) *RetryableTransaction {
	return &RetryableTransaction{Transaction{
		TransactionName:   transactionName,
		ServiceName:       serviceName,
		TransactionAction: transactionAction,
	}}
}

// OKay I am looking back after 2 months of writing most of this, and I am not positive how the theis RestProxy connects
// to the code above

type TransactionRESTProxy struct {
	TransactionAction
	HttpFunction string
	URL          string
	Body         *io.Reader
	Response     *http.Response
}

// NewTransactionActionRestProxy is a constructor for the Transaction REST Proxy creation just to make sure values make sense before making that call to the other services, or it own service
func NewTransactionActionRestProxy(httpFunction string, url string, body *io.Reader) (*TransactionRESTProxy, error) {
	var ErrMismatchHttpFunctionWithBody = fmt.Errorf("[ERROR] [SAGA] [TRANSACTION] [REST PROXY] [CREATION] | Cannot have body with this %+v | GET or DELETE", httpFunction)
	var ErrMissingBodyForAssociatedHttpFunction = fmt.Errorf("[ERROR] [SAGA] [TRANSACTION] [REST PROXY] [CREATION] | Cannot have missingbody for this fucntion %+v | POST or UPDATE", httpFunction)
	var ErrNonSupportedHTTPFunction = fmt.Errorf("[ERROR] [SAGA] [TRANSACTION] [REST PROXY] [CREATION] | %+v is not a supported http Fuction: Supported Functions: POST, GET, UPDATE, and DELETE", httpFunction)
	groomedHTTPFunction := strings.ToUpper(strings.TrimSpace(httpFunction))
	switch groomedHTTPFunction {
	case "GET":
		{
			if body != nil {
				return nil, ErrMismatchHttpFunctionWithBody
			}

		}
	case "DELETE":
		{
			if body != nil {
				return nil, ErrMismatchHttpFunctionWithBody
			}
		}
	case "UPDATE":
		{
			if body == nil {
				return nil, ErrMissingBodyForAssociatedHttpFunction
			}
		}
	case "POST":
		{
			if body == nil {
				return nil, ErrMissingBodyForAssociatedHttpFunction
			}
		}
	default:
		{
			return nil, ErrNonSupportedHTTPFunction
		}
	}
	//TODO Figure out if should do any check on the url
	//I am guessing I should check if it in the correct domain like internal, but later
	//That means I need a file where I keep the domain at for later
	//
	// I could also have it be dynamic like a list that updated when services in the "backend" so have script
	// run at the start of each service
	// I have a lot idea about this, but none seem to pressing to getting this operational

	return &TransactionRESTProxy{HttpFunction: groomedHTTPFunction, URL: url, Body: body, Response: nil}, nil
}

func (trp *TransactionRESTProxy) Execute() error {

	client := &http.Client{}
	req, err := http.NewRequest(trp.HttpFunction, trp.URL, *trp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// Need to find a way to make this return an error.
			//
		}
	}(resp.Body)
	trp.Response = resp
	return nil

}
