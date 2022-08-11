package saga_pattern

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var prefixSagaExecutionError = "[ERROR] [SAGA] [EXECUTION] "

const RetryableTransactionRetryLimit = 3

type sagaRunner struct {
	logger             *log.Logger
	generalErrorPrefix error
}

//ExecuteSaga
/*

	Order of Transaction is import in passed ListOfTransactions
	- Compensatable Transactions:
		-- NOTE: Compensatable Transactions should have two Transaction Commands. The first one for the
				intended Transaction and the second one to Roll back that Transaction encase of failure/issues
				in any of the following Transaction
		-- Can only be in place prior in order to the Pivot Transaction and never after
		-- If there is no Pivot Transaction Action then there cannot be any Compensatable Transactions
	- Pivot Transaction:
		-- There can only be one per list of Transactions or
	- Retryable Transactions:
		-- If the Pivot Transaction exist in the List then Retryable Transactions can only be in after that Pivot Transaction
		-- Retryable Transactions should never be in the list of Transactions prior to the Pivot Transaction if it exists

	- NOTE Other Use Cases:
		-- You can have a Pivot Transaction without Compensatable Transaction(s)
		-- You can just have list of Transactions without Pivot nor Compensatable Transactions(So just Retryable Transactions (Granted idk why you would tho))
*/
func ExecuteSaga(orchestrator SagaOrchestrator) error {
	prefixSagaExecutionError := errors.New(prefixSagaExecutionError + "[" + orchestrator.Name + "]")
	sagaRunnerLogger := sagaRunner{
		logger:             log.New(os.Stdout, orchestrator.Name+" | ", log.LstdFlags),
		generalErrorPrefix: prefixSagaExecutionError,
	}
	//This is the list of the undo transaction in the event of a failure or
	var compensatableTransactionCommands []TransactionCommand
	var hasCompensatableTransaction = false

	sagaRunnerLogger.logger.Println("------- Saga Execution Start --------")
	for _, transaction := range orchestrator.Transactions {
		//Pull the List of Transactions Commands
		// [1] is the intended transaction Command
		// [2] would be a compensatableTransaction Command
		transactionCommands := transaction.GetTransactionCommands()

		sagaRunnerLogger.logger.Println("Transaction Name: " + transaction.GetTransactionName())
		sagaRunnerLogger.logger.Printf("Transaction Type: %+v\n", transaction.GetTransactionType())

		//Execute the Transaction
		// To handle result a lil later
		canProceed, err := transactionCommands[0].Execute()

		sagaRunnerLogger.logger.Println("--- Transaction Executed Results ---")
		sagaRunnerLogger.logger.Printf("Can Proceed to Next Transaction: %t\n", canProceed)
		sagaRunnerLogger.logger.Printf("Was There an Error In that Transaction: %v\n", err)
		//COMPENSATABLE TRANSACTION
		//If it's a Compensatable Transaction Type add the second transaction to the list undo Transaction Commands
		if transaction.GetTransactionType() == GetCompensatableTransactionType() {
			hasCompensatableTransaction = true
			compensatableTransactionCommands = append(compensatableTransactionCommands, transactionCommands[2])
		}

		//FAILURE IN COMPENSATABLE OR PIVOT TRANSACTION
		// If get the do not proceed flag or an error (the canProceed is for Transaction Commands that get responses)
		if !canProceed || err != nil {
			// This is just in case you get a Pivot Transaction without any CompensatableTransaction(s) in an effort to try and save some compute
			if !hasCompensatableTransaction {
				if transaction.GetTransactionType() == GetPivotTransactionType() {
					if err == nil {
						err = errors.New("got a Do Not Proceed, but No Error(So now you get an Error)")
					}
					log.Println(err)
					return err
				}
			}
			//Check if the Transaction Type is Pivot/Compensatable then undo the Transactions providing the list of undo Transaction Commands
			if transaction.GetTransactionType() == GetCompensatableTransactionType() || transaction.GetTransactionType() == GetPivotTransactionType() {
				undoErr := undoCompensatableTransactions(compensatableTransactionCommands, &sagaRunnerLogger)
				if undoErr != nil {
					// This process fails I really don't want to tell the end user. So I guess I'll just log it.
					sagaRunnerLogger.logger.Println(undoErr)
				}
				if err == nil {
					err = errors.New("got a Do Not Proceed, but No Error(So now you get an Error)")
				}
				return err

			}
		}
		//Check if the Transaction Type is Retryable then try to reattempt them to preset retry limit
		if transaction.GetTransactionType() == GetRetriableTransactionType() {
			if !canProceed || err != nil {
				didItFinallyWork, err := retryRetractableTransactionToPresetLimit(transactionCommands[0], &sagaRunnerLogger)
				if err != nil {
					return err
				}
				if !didItFinallyWork {
					retriesFailed := fmt.Errorf("%v [UNDO] | After the Preset Number (#%d) of Retries this RetryTabiable Transaction(%v) has still Failed. Look at contruction of this Transaction and its Transaction Command(%v)", sagaRunnerLogger.generalErrorPrefix, RetryableTransactionRetryLimit, transaction.GetTransactionName(), transaction.GetTransactionCommands()[0].GetTransactionCommandName())
					sagaRunnerLogger.logger.Println(retriesFailed)
					// I am not positive If I am supposed to return this error I will for now
					return fmt.Errorf("%v  | %v ", retriesFailed, err)
				}
			}

		}
		sagaRunnerLogger.logger.Println("------- Next Transaction --------")
	} //End of loop
	sagaRunnerLogger.logger.Println("------- End of Saga Execution --------")
	return nil
}
func retryRetractableTransactionToPresetLimit(command TransactionCommand, runner *sagaRunner) (bool, error) {
	var didItFinallyWork = false
	var err error
	for i := 1; i <= RetryableTransactionRetryLimit; i++ {
		runner.logger.Println(fmt.Sprintf("%v [RETRYABLE] Retrying Retryable Transaction %d of %d", runner.generalErrorPrefix, i, RetryableTransactionRetryLimit))
		var err1 error
		didItFinallyWork, err1 = command.Execute()
		//TODO Learn proper error handling
		err = fmt.Errorf("retry #: %d | error: %v | %v", i, err1, err)

	}
	return didItFinallyWork, err
}

func undoCompensatableTransactions(listOfCompensatableTransactions []TransactionCommand, runner *sagaRunner) error {
	//TODO Need to learn to Wrap errors
	undo := errors.New("[UNDO] |")
	errCollector := errors.New("")
	for _, transactionCommand := range listOfCompensatableTransactions {
		workAsExpected, err := transactionCommand.Execute()
		if err != nil {
			runner.logger.Printf(" %v %v | TransactionCommandName: %v, Error: %v", runner.generalErrorPrefix, undo, transactionCommand.GetTransactionCommandName(), err)
			errCollector = fmt.Errorf("%v | %v | %v", runner.generalErrorPrefix, undo, err)
			return errCollector
		}
		if !workAsExpected {
			runner.logger.Printf("%v %v | | TransactionCommandName: %v . Got an do not proceed on the undo that's not good so here is your sign", runner.generalErrorPrefix, undo, transactionCommand.GetTransactionCommandName())
		}
	}
	return nil
}
