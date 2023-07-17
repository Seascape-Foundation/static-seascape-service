package handler

import (
	"github.com/ahmetson/common-lib/data_type/database"
	"github.com/ahmetson/service-lib/communication/command"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/service-lib/remote"
	"github.com/ahmetson/static-service/abi"

	"github.com/ahmetson/service-lib/communication/message"
)

type GetAbiRequest struct {
	Id string `json:"abi_id"`
}
type GetAbiReply = abi.Abi

type SetAbiRequest struct {
	Body interface{} `json:"body"`
}
type SetAbiReply = abi.Abi

// AbiGet returns the abi
// Depends on the database extension
var AbiGet = func(request message.Request, _ log.Logger, extensions remote.Clients) message.Reply {
	if !remote.ClientExist(extensions, "database") {
		return message.Fail("missing extension")
	}

	var reqParameters GetAbiRequest
	err := request.Parameters.Interface(&reqParameters)
	if err != nil {
		return message.Fail("request.Parameters -> Command Parameter: " + err.Error())
	}
	if len(reqParameters.Id) == 0 {
		return message.Fail("missing abi id")
	}

	dbCon := remote.GetClient(extensions, "database")
	var selectedAbi = abi.Abi{Id: reqParameters.Id}
	saveErr := selectedAbi.Select(dbCon)
	if saveErr != nil {
		return message.Fail("database error:" + saveErr.Error())
	}

	replyMessage, err := command.Reply(selectedAbi)
	if err != nil {
		return message.Fail("failed to reply")
	}
	return replyMessage
}

func AbiRegister(request message.Request, _ log.Logger, extensions remote.Clients) message.Reply {
	if !remote.ClientExist(extensions, "database") {
		return message.Fail("missing extension")
	}

	var requestParameters SetAbiRequest
	err := request.Parameters.Interface(&requestParameters)
	if err != nil {
		return message.Fail("failed to parse data")
	}

	if requestParameters.Body == nil {
		return message.Fail("missing body")
	}

	newAbi, err := abi.NewFromInterface(requestParameters.Body)
	if err != nil {
		return message.Fail("abi.NewFromInterface: " + err.Error())
	}
	if len(newAbi.Body) == 0 {
		return message.Fail("body is empty")
	}

	replyMessage, err := command.Reply(newAbi)
	if err != nil {
		return message.Fail("failed to reply")
	}

	dbCon := remote.GetClient(extensions, "database")
	var crud database.Crud = newAbi
	saveErr := crud.Insert(dbCon)
	if saveErr != nil {
		return message.Fail("database error:" + saveErr.Error())
	}

	return replyMessage
}
