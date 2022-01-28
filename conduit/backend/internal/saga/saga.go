package saga

type Saga struct {
	name         string
	transactions []*Transaction
}

const compensatable = "compensatable"
const pivot = "pivot"
const retriable = "retriable"

type Transaction struct {
	serviceName          string
	transactionType      string
	sagaTransactionProxy *TransactionProxy
}

type TransactionProxy struct {
}
