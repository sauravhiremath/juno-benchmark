package src

import uuid "github.com/satori/go.uuid"

/*

GetInitializeModuleMessage constructs the Juno Initialize module message.

@optional params -> requestID

*/
func GetInitializeModuleMessage(moduleID string, params ...string) InitializeModuleRequest {
	var requestID string
	if len(params) > 0 {
		requestID = params[0]
	} else {
		requestID = uuid.NewV4().String()
	}
	return InitializeModuleRequest{
		RequestID: requestID,
		Type:      1,
		ModuleID:  "Module-" + requestID,
		Version:   "1.0.0",
	}
}

/*

GetFunctionCallMessage constructs the Juno Call Function message.

@optional params -> requestID

*/
func GetFunctionCallMessage(moduleID string, functionName string, params ...string) FunctionCallRequest {
	var requestID string
	if len(params) > 0 {
		requestID = params[0]
	} else {
		requestID = uuid.NewV4().String()
	}
	return FunctionCallRequest{
		RequestID: requestID,
		Type:      3,
		Function:  moduleID + "." + functionName,
	}
}

/*

GetDeclareFunctionMessage constructs the Juno Declare Function message.

@optional params -> requestID

*/
func GetDeclareFunctionMessage(function string, params ...string) DeclareFunctionRequest {
	var requestID string
	if len(params) > 0 {
		requestID = params[0]
	} else {
		requestID = uuid.NewV4().String()
	}
	return DeclareFunctionRequest{
		RequestID: requestID,
		Type:      9,
		Function:  function,
	}
}
