package saga

import (
	"backend/channel-service/models"
	saga_pattern "backend/internal/saga-pattern"
	serverMeta "backend/server-service/meta"
)

var errCreateChannelSagaPrefixError = "[ERROR] [SAGA] [CHANNEL] [CREATE_CHANNEL]"

func CreateChannelSaga(channel models.Channel) (*saga_pattern.SagaOrchestrator, error) {

	var verifyServerExist = saga_pattern.NewPivotTransaction("VerifyServerExist", saga_pattern.NewInternalHttpRESTRequest("GetServerViaServerID", serverMeta.GetGETServerViaHTTPServerIDRequest(channel.ServerID.String())))
	var addChannel = saga_pattern.NewRetryableTransaction("AddChannel", NewHostSagaChannelCommandAddChannel(channel))

	t := []saga_pattern.Transaction{verifyServerExist, addChannel}
	orchestrator, err := saga_pattern.NewSagaOrchestrator("create-channel-saga", t)
	if err != nil {

		return nil, err
	}
	return orchestrator, nil

}

type hostSagaChannelCommandAddChannel struct {
	TransactionCommandName string
	channel                models.Channel
}

func NewHostSagaChannelCommandAddChannel(channel models.Channel) *hostSagaChannelCommandAddChannel {
	hh := new(hostSagaChannelCommandAddChannel)
	hh.channel = channel
	return hh
}
func (h hostSagaChannelCommandAddChannel) GetTransactionCommandType() string {
	return "Host Service Create"
}
func (h *hostSagaChannelCommandAddChannel) GetTransactionCommandName() string {
	return h.TransactionCommandName
}

func (h hostSagaChannelCommandAddChannel) Execute() (success bool, err error) {
	err = models.AddChannel(h.channel)
	if err != nil {
		return false, err
	}
	return true, err

}

//type hostChannelServiceCommand struct {
//	CommandType string
//	SearchID    uuid.UUID
//}
//
//var listOfChannelModelFunctions = []string{"GetEveryChannel", "GetChannelsViaServerID", "GetChannelViaChannelID", "AddChannel", "UpdateChannel", "DeleteChannel", "doesChannelExistWithMatchingID"}
//
//func NewHostChannelServiceCommand(commandType string, searchID uuid.UUID) (hostChannelServiceCommand, error) {
//	if contains(listOfChannelModelFunctions, commandType) {
//		return hostChannelServiceCommand{
//			CommandType: commandType,
//			SearchID:    searchID,
//		}, nil
//	}
//	var errNonSupportChannelMethodName = errors.New("[ERROR] [SAGA] [CHANNEL]: Unsupported models.channel method name passed for HOST sage service Transaction Command")
//	return hostChannelServiceCommand{}, errNonSupportChannelMethodName
//}
//
//func (h hostChannelServiceCommand) Execute(i interface{}) (response interface{}, err error) {
//
//}
//func contains(slice []string, searchFor string) bool {
//	for _, a := range slice {
//		if a == searchFor {
//			return true
//		}
//	}
//	return false
//}
