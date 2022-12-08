/*
Controller package is the interface of the module.
It acts as the input receiver for other services or for external users.
*/
package controller

import (
	"database/sql"
	"fmt"

	"github.com/blocklords/gosds/account"
	"github.com/blocklords/gosds/env"
	"github.com/blocklords/gosds/message"

	zmq "github.com/pebbe/zmq4"
)

type CommandHandlers map[string]interface{}

// Creates a new Reply controller using ZeroMQ
// The requesters is the list of curve public keys that are allowed to connect to the socket.
func ReplyController(db *sql.DB, commands CommandHandlers, e *env.Env, public_keys []string) {
	if !e.PortExist() {
		panic(fmt.Errorf("missing .env variable: Please set '" + e.ServiceName() + "' port"))
	}

	zmq.AuthSetVerbose(true)
	err := zmq.AuthStart()
	if err != nil {
		panic(err)
	}

	// allow income from any ip address
	// for any domain name where this controller is running.
	zmq.AuthAllow("*")
	// only whitelisted users are allowed
	zmq.AuthCurveAdd("*", public_keys...)

	handler := func(version string, request_id string, domain string, address string, identity string, mechanism string, credentials ...string) (metadata map[string]string) {
		metadata = map[string]string{
			"request_id": request_id,
			"Identity":   zmq.Z85encode(credentials[0]),
			"address":    address,
			"pub_key":    zmq.Z85encode(credentials[0]), // if mechanism is not curve, it will fail
		}
		return metadata
	}
	zmq.AuthSetMetadataHandler(handler)

	// Socket to talk to clients
	socket, _ := zmq.NewSocket(zmq.REP)
	socket.ServerAuthCurve(e.DomainName(), e.SecretKey())
	defer socket.Close()
	defer zmq.AuthStop()
	if err := socket.Bind("tcp://*:" + e.Port()); err != nil {
		println("error to bind socket for '"+e.ServiceName()+" - "+e.Url()+"' : ", err.Error())
		panic(err)
	}

	println("'" + e.ServiceName() + "' request-reply server runs on port " + e.Port())

	for {
		msg_raw, err := socket.RecvMessage(0)
		if err != nil {
			println(fmt.Errorf("receiving: %w", err))
			fail := message.Fail("invalid command " + err.Error())
			reply := fail.ToString()
			if _, err := socket.SendMessage(reply); err != nil {
				println(fmt.Errorf(" reply: %w", err))
			}
			continue
		}
		request, err := message.ParseRequest(msg_raw)
		if err != nil {
			fail := message.Fail("invalid request " + err.Error())
			reply := fail.ToString()
			if _, err := socket.SendMessage(reply); err != nil {
				println(fmt.Errorf("sending reply: %w", err))
			}
			continue
		}

		// Any request types is compatible with the Request.
		if commands[request.Command] == nil {
			fail := message.Fail("invalid command " + request.Command)
			reply := fail.ToString()
			if _, err := socket.SendMessage(reply); err != nil {
				println(fmt.Errorf(" reply: %w", err))
			}
			continue
		}

		var reply message.Reply

		command_handler, ok := commands[request.Command].(func(*sql.DB, message.SmartcontractDeveloperRequest, *account.SmartcontractDeveloper) message.Reply)
		if ok {
			smartcontract_developer_request, err := message.ParseSmartcontractDeveloperRequest(msg_raw)
			if err != nil {
				fail := message.Fail("invalid smartcontract developer request " + err.Error())
				reply := fail.ToString()
				if _, err := socket.SendMessage(reply); err != nil {
					println(fmt.Errorf("sending reply: %w", err))
				}
				continue
			}

			smartcontract_developer, err := smartcontract_developer_request.GetAccount()
			if err != nil {
				fail := message.Fail("invalid smartcontract developer request " + err.Error())
				reply := fail.ToString()
				if _, err := socket.SendMessage(reply); err != nil {
					println(fmt.Errorf("sending reply: %w", err))
				}
				continue
			}

			reply = command_handler(db, smartcontract_developer_request, smartcontract_developer)
		} else {
			reply = commands[request.Command].(func(*sql.DB, message.Request) message.Reply)(db, request)
		}

		if _, err := socket.SendMessage(reply.ToString()); err != nil {
			println(fmt.Errorf("error sending controller reply: %w", err))
		}
	}
}
