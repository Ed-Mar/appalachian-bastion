package saga_pattern

type SagaTransactionType interface {
	GetSagaTransactionType() sagaTransactionType
}

type sagaTransactionType struct {
	transactionType string
}

func (s sagaTransactionType) GetSagaTransactionType() sagaTransactionType {
	return sagaTransactionType{s.transactionType}
}

//These Definitions are straight form the book
//Microservices Patterns By Chris Richardson

// CompensatableTransactionType
type CompensatableTransactionType struct {
	sagaTransactionType
}

func (c *CompensatableTransactionType) GetSagaTransactionType() sagaTransactionType {
	c.transactionType = "Compensatable"
	return sagaTransactionType{transactionType: "Compensatable"}
}

//GetCompensatableTransactionType This is here due to this one be special, and I can't think a better to do right now
func GetCompensatableTransactionType() sagaTransactionType {
	return sagaTransactionType{transactionType: "Compensatable"}

}

type PivotTransactionType struct {
	sagaTransactionType
}

func (p *PivotTransactionType) GetSagaTransactionType() sagaTransactionType {
	p.transactionType = "Pivot"
	return sagaTransactionType{transactionType: "Pivot"}
}
func GetPivotTransactionType() sagaTransactionType {
	return sagaTransactionType{transactionType: "Pivot"}
}

//RetriableTransactionType Transactions that follow the pivot Transaction and are guaranteed to succeed.
type RetriableTransactionType struct {
	sagaTransactionType
}

func GetRetriableTransactionType() sagaTransactionType {
	return sagaTransactionType{transactionType: "Retriable"}
}

func (r *RetriableTransactionType) GetSagaTransactionType() sagaTransactionType {
	r.transactionType = "Retriable"
	return sagaTransactionType{transactionType: "Retriable"}
}
