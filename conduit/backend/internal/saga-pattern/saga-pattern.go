package saga_pattern

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type SagaOrchestrator struct {
	Name         string
	Transactions []Transaction
}

func (s *SagaOrchestrator) Log() {
	log.Println("SagaOrchestrator Name: " + s.Name)
	for _, trans := range s.Transactions {
		log.Println("Transaction Name: " + trans.GetTransactionName())
		log.Printf("Transaction Type: %+v\n", trans.GetTransactionType())
		log.Printf("Transaction Command : %+v\n", trans.GetTransactionCommands()[0])
		if trans.GetTransactionType() == GetCompensatableTransactionType() {
			log.Printf("Transaction Command : %+v\n", trans.GetTransactionCommands()[1])
		}
		log.Println("------------------------------")
	}
}

func NewSagaOrchestrator(name string, transactions []Transaction) (*SagaOrchestrator, error) {
	var sagaCreationErrorPrefix = "[ERROR] [SAGA_PATTERN] "

	so := new(SagaOrchestrator)
	so.Name = name

	var countOfCompensatableTransactions = 0
	var countOfPivotTransactions = 0
	var countOfRetryableTransactions = 0
	lastIndex := len(transactions) - 1
	// Loops through the list of transactions to check for acceptable order has been passed
	for index, transaction := range transactions {

		//Checks for Compensatable Transactions
		//Ruling:
		// There cannot be a Pivot Transaction prior to any Compensatable Transaction
		// There cannot be Retryable Transactions prior to any Retryable Transactions
		if transaction.GetTransactionType() == GetCompensatableTransactionType() {
			countOfCompensatableTransactions = countOfCompensatableTransactions + 1
			//Checking if it's before any Pivot Transactions
			if countOfPivotTransactions > 0 {
				return nil, errors.New(sagaCreationErrorPrefix + "[SAGA_CREATION] [Compensatable]| Cannot put Compensatable Transactions after a Pivot Transaction")
			}
			//Checking if it's before any Retryable Transactions
			if countOfRetryableTransactions > 0 {
				return nil, errors.New(sagaCreationErrorPrefix + "[SAGA_CREATION] [Compensatable]| Cannot put Compensatable Transactions before any Retryable Transaction(s)")
			}
		}

		//Checks for Pivot Transaction
		//Ruling:
		// There cannot be more than on pivot transaction
		// There are no Retryable prior to the Pivot
		if transaction.GetTransactionType() == GetPivotTransactionType() {
			countOfPivotTransactions = countOfPivotTransactions + 1
			//Checks that there is not already a Pivot Transaction in the List of Transactions
			if countOfPivotTransactions > 1 {
				return nil, errors.New(sagaCreationErrorPrefix + "[SAGA_CREATION] [Pivot]| There is more than one Pivot Transaction")
			}
			//Checks that there is not any Retryable Transactions prior in the list of Transactions
			if countOfRetryableTransactions > 0 {
				return nil, errors.New(sagaCreationErrorPrefix + "[SAGA_CREATION] [Pivot]| Cannot put Pivot Transaction After Retryable Transactions")
			}

		}
		//Checking Retryable
		//Note: currently 220801 that I do not think I need to check for anything for Retryable, cause the other two cover it.
		//but I will leave this code here just in case
		if transaction.GetTransactionType() == GetPivotTransactionType() {
			countOfRetryableTransactions = countOfRetryableTransactions + 1
		}
		//Checks at last Transaction in the list
		if index == lastIndex {
			//Checks that there is a pivot if there is any Compensatable Transactions
			if countOfPivotTransactions < 1 && countOfCompensatableTransactions > 0 {
				return nil, errors.New(sagaCreationErrorPrefix + "[SAGA_CREATION] | Compensatable Transactions have not Pivot Transactions")
			}
		}
	}
	// So it passes the checks above its good to go.
	so.Transactions = transactions
	return so, nil
}

type Transaction interface {
	GetTransactionType() sagaTransactionType
	GetTransactionCommands() []TransactionCommand
	GetTransactionName() string
}

//compensatableTransaction s are Transactions that have reverse Transaction encase a subsequent Transaction fails
type compensatableTransaction struct {
	Name string
	TransactionCommand
	CompensationCommand TransactionCommand
}

func (c *compensatableTransaction) GetTransactionCommands() []TransactionCommand {
	tc := []TransactionCommand{c.TransactionCommand, c.CompensationCommand}
	return tc

}
func (c *compensatableTransaction) GetTransactionType() sagaTransactionType {
	return GetCompensatableTransactionType()
}
func (c *compensatableTransaction) GetTransactionName() string {
	return c.Name
}

//NewCompensatableTransaction Constructor for compensatableTransaction
func NewCompensatableTransaction(name string, transactionCommand TransactionCommand, compensationCommand TransactionCommand) *compensatableTransaction {
	ct := &compensatableTransaction{
		Name:                name,
		TransactionCommand:  transactionCommand,
		CompensationCommand: compensationCommand,
	}
	return ct
}

type pivotTransaction struct {
	Name string
	TransactionCommand
}

func (p *pivotTransaction) GetTransactionCommands() []TransactionCommand {
	tc := []TransactionCommand{p.TransactionCommand}
	return tc
}

func (p *pivotTransaction) GetTransactionType() sagaTransactionType {
	return GetPivotTransactionType()
}
func (p *pivotTransaction) GetTransactionName() string {
	return p.Name
}

// NewPivotTransaction Constructor for pivotTransaction
// pivotTransaction The go/no-go point in a saga. If the pivot Transaction commits, the saga will run until completion.
// A pivot Transaction can be a Transaction that’s neither compensatable nor retriable. Alternatively, it can be the last compensatable
// Transaction or the first retriable Transaction.
func NewPivotTransaction(name string, transactionCommand TransactionCommand) *pivotTransaction {
	p := &pivotTransaction{
		Name:               name,
		TransactionCommand: transactionCommand,
	}
	return p
}

// retryableTransaction a Transaction that should be retried if it fails  retryableTransaction s should only
type retryableTransaction struct {
	Name string
	TransactionCommand
}

func (r *retryableTransaction) GetTransactionCommands() []TransactionCommand {
	tc := []TransactionCommand{r.TransactionCommand}
	return tc
}

func (r *retryableTransaction) GetTransactionType() sagaTransactionType {
	return GetRetriableTransactionType()
}
func (r *retryableTransaction) GetTransactionName() string {
	return r.Name
}

// NewRetryableTransaction Constructor for retryableTransaction
func NewRetryableTransaction(name string, transactionCommand TransactionCommand) *retryableTransaction {
	r := &retryableTransaction{
		Name:               name,
		TransactionCommand: transactionCommand,
	}
	return r
}

// TransactionCommand this action part of the Transaction in which Execute would do it and give the response if any.
// this was made to allow both REST and Kafka thus functioning taking an anonymous interface. Doing this allows for both
// of the possible solutions Actions to be in the same collection
type TransactionCommand interface {
	Execute() (success bool, err error)
	GetTransactionCommandType() string
	GetTransactionCommandName() string
}

type internalHttpRESTRequest struct {
	TransactionCommandName string
	request                *http.Request
}

func (i *internalHttpRESTRequest) setInternalHTTPRequest(req *http.Request) {
	i.request = req
}
func NewInternalHttpRESTRequest(transactionCommandName string, req *http.Request) *internalHttpRESTRequest {
	i := new(internalHttpRESTRequest)
	i.TransactionCommandName = transactionCommandName
	i.setInternalHTTPRequest(req)
	return i

}
func (i *internalHttpRESTRequest) GetTransactionCommandName() string {
	return i.TransactionCommandName
}

func (i *internalHttpRESTRequest) GetTransactionCommandType() string {
	return "Internal HTTP Request"
}

func (i *internalHttpRESTRequest) Execute() (success bool, err error) {

	var errNotAnInternalURLRequest = errors.New("[ERROR] [SAGA-PATTERN] [internalHttpRESTRequest]: URL is not an internal domain")
	const correctAPIDomain = "localhost"

	//TODO change this to actually check if the URL in the http.Request is with in the correct domain
	// Cause doing a regex or tokenizer to pick off the port info just for this work in the dev env
	// is a waste of time and compute JUST MAKE SURE YOU COME BACK AND FIX THIS
	//if httpReq.URL.Host != correctAPIDomain{}
	if i.request.URL.Host == "" {
		err = fmt.Errorf(i.request.URL.Host, errNotAnInternalURLRequest)
		log.Println(err)
		return false, err
	}
	// I think I need to not use the http.DefaultClient cuase it limits the HTTP fuctions I can do HEAD,GET,POST and my call use PUT/DELETE Aswell
	simpleClient := http.Client{
		// If nil, DefaultTransport is used.
		Transport: nil,
		//If CheckRedirect is nil, the Client uses its default policy,
		//which is to stop after 10 consecutive requests.
		CheckRedirect: nil,
		// If Jar is nil, cookies are only sent if they are explicitly
		// set on the Request.
		Jar: nil,
		// So the default time out is zero, but not for Get,Post,Head Its odd thing
		Timeout: 60 * time.Millisecond,
	}

	reqDump, err := httputil.DumpRequestOut(i.request, true)
	if err != nil {
		log.Println(err)
	}

	log.Printf("REQUEST:\n%s", string(reqDump))

	// Do the actual transaction
	res, httpErr := simpleClient.Do(i.request)

	respDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("RESPONSE:\n%s", string(respDump))

	//TODO Keep this In mind building out
	// Okay So I know at some point I am going to need to come back make this better to handle more than just the status
	// Codes, but I think I am okay for now for internal HTTP Calls
	// Checking if the Status Code are the ones that indicate if the saga should move forward.

	switch res.StatusCode {
	case http.StatusOK: //200
	case http.StatusCreated: //201
	case http.StatusAccepted: //202
	case http.StatusNoContent: //204
	default:
		return false, errors.New("Non Supported Status Code: " + res.Status)
	}
	// So if one of the expected response codes then I am going to return a true to continue forward
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	bodyString := string(bodyBytes)
	log.Println("Response Body: " + bodyString)

	return true, httpErr

}

type HostSagaServiceCommand interface {
}
