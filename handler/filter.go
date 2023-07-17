package handler

import (
	"fmt"

	"github.com/Seascape-Foundation/sds-service-lib/log"
	"github.com/Seascape-Foundation/static-seascape-service/configuration"
	"github.com/Seascape-Foundation/static-seascape-service/smartcontract"

	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
	"github.com/Seascape-Foundation/sds-common-lib/smartcontract_key"
	"github.com/Seascape-Foundation/sds-common-lib/topic"

	"github.com/Seascape-Foundation/sds-service-lib/communication/command"
	"github.com/Seascape-Foundation/sds-service-lib/communication/message"
)

type FilterSmartcontractsRequest = topic.Filter
type FilterSmartcontractsReply struct {
	Smartcontracts []*smartcontract.Smartcontract `json:"smartcontracts"`
	TopicStrings   []topic.Id                     `json:"topic_strings"`
}

type FilterSmartcontractKeysRequest = topic.Filter
type FilterSmartcontractKeysReply struct {
	SmartcontractKeys []smartcontract_key.Key `json:"smartcontract_keys"`
	TopicStrings      []topic.Id              `json:"topic_strings"`
}

func filterOrganization(configurations *key_value.List, paths []string) *key_value.List {
	if len(paths) == 0 {
		return configurations
	}

	filtered := key_value.NewList()
	if configurations == nil {
		return filtered
	}

	list := configurations.List()
	for key, value := range list {
		conf := value.(*configuration.Configuration)

		for _, path := range paths {
			if conf.Topic.Organization == path {
				_ = filtered.Add(key, value)
				break
			}
		}
	}

	return filtered
}

func filterProject(configurations *key_value.List, paths []string) *key_value.List {
	if len(paths) == 0 {
		return configurations
	}

	filtered := key_value.NewList()
	if configurations == nil {
		return filtered
	}

	list := configurations.List()
	for key, value := range list {
		conf := value.(*configuration.Configuration)

		for _, path := range paths {
			if conf.Topic.Project == path {
				_ = filtered.Add(key, value)
				break
			}
		}
	}

	return filtered
}

func filterNetworkId(configurations *key_value.List, paths []string) *key_value.List {
	if len(paths) == 0 {
		return configurations
	}

	filtered := key_value.NewList()
	if configurations == nil {
		return filtered
	}

	list := configurations.List()
	for key, value := range list {
		conf := value.(*configuration.Configuration)

		for _, path := range paths {
			if conf.Topic.NetworkId == path {
				_ = filtered.Add(key, value)
				break
			}
		}
	}

	return filtered
}

func filterGroup(configurations *key_value.List, paths []string) *key_value.List {
	if len(paths) == 0 {
		return configurations
	}

	filtered := key_value.NewList()
	if configurations == nil {
		return filtered
	}

	list := configurations.List()
	for key, value := range list {
		conf := value.(*configuration.Configuration)

		for _, path := range paths {
			if conf.Topic.Group == path {
				_ = filtered.Add(key, value)
				break
			}
		}
	}

	return filtered
}

func filterSmartcontractName(configurations *key_value.List, paths []string) *key_value.List {
	if len(paths) == 0 {
		return configurations
	}

	filtered := key_value.NewList()
	if configurations == nil {
		return filtered
	}

	list := configurations.List()
	for key, value := range list {
		conf := value.(*configuration.Configuration)

		for _, path := range paths {
			if conf.Topic.Name == path {
				_ = filtered.Add(key, value)
				break
			}
		}
	}

	return filtered
}

func filterConfiguration(configurationList *key_value.List, topicFilter *topic.Filter) []*configuration.Configuration {
	list := key_value.NewList()

	if len(topicFilter.Organizations) != 0 {
		list = filterOrganization(configurationList, topicFilter.Organizations)
	}

	if len(topicFilter.Projects) != 0 {
		list = filterProject(list, topicFilter.Projects)
	}

	if len(topicFilter.NetworkIds) != 0 {
		list = filterNetworkId(list, topicFilter.NetworkIds)
	}

	if len(topicFilter.Groups) != 0 {
		list = filterGroup(list, topicFilter.Groups)
	}

	if len(topicFilter.Smartcontracts) != 0 {
		list = filterSmartcontractName(list, topicFilter.Smartcontracts)
	}

	configs := make([]*configuration.Configuration, list.Len())

	i := 0
	for _, value := range list.List() {
		conf := value.(*configuration.Configuration)
		configs[i] = conf
		i++
	}

	return configs
}

func filterSmartcontract(
	configurations []*configuration.Configuration,
	list *key_value.List) ([]*smartcontract.Smartcontract, []topic.Id, error) {

	smartcontracts := make([]*smartcontract.Smartcontract, 0)
	topicStrings := make([]topic.Id, 0)

	for _, conf := range configurations {
		key, err := smartcontract_key.New(conf.Topic.NetworkId, conf.Topics[0].Name)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create smartcontract key: %w", err)
		}

		value, err := list.Get(key)
		if err != nil {
			fmt.Println("not found")
			continue
		}
		sm := value.(*smartcontract.Smartcontract)

		smartcontracts = append(smartcontracts, sm)
		topicStrings = append(topicStrings, conf.Topic.Id())
	}

	return smartcontracts, topicStrings, nil
}

// SmartcontractFilter /*
func SmartcontractFilter(request message.Request, _ log.Logger, parameters ...interface{}) message.Reply {
	var topicFilter FilterSmartcontractKeysRequest
	err := request.Parameters.Interface(&topicFilter)
	if err != nil {
		return message.Fail("failed to parse data")
	}

	allConfigurations := parameters[3].(*key_value.List)
	configurations := filterConfiguration(allConfigurations, &topicFilter)
	if len(configurations) == 0 {
		reply := FilterSmartcontractsReply{
			Smartcontracts: []*smartcontract.Smartcontract{},
			TopicStrings:   []topic.Id{},
		}
		replyMessage, err := command.Reply(&reply)
		if err != nil {
			return message.Fail("failed to reply: " + err.Error())
		}
		return replyMessage
	}

	allSmartcontracts := parameters[2].(*key_value.List)
	smartcontracts, topicStrings, err := filterSmartcontract(configurations, allSmartcontracts)
	if err != nil {
		return message.Fail("failed to reply: " + err.Error())
	}

	reply := FilterSmartcontractsReply{
		Smartcontracts: smartcontracts,
		TopicStrings:   topicStrings,
	}
	replyMessage, err := command.Reply(&reply)
	if err != nil {
		return message.Fail("failed to reply: " + err.Error())
	}
	return replyMessage
}

// SmartcontractKeyFilter returns smartcontract keys and topic of the smartcontract
// by given topic filter
//
//	returns {
//			"smartcontract_keys" (where key is smartcontract key, value is a topic string)
//	}
func SmartcontractKeyFilter(request message.Request, _ log.Logger, parameters ...interface{}) message.Reply {
	var topicFilter FilterSmartcontractKeysRequest
	err := request.Parameters.Interface(&topicFilter)
	if err != nil {
		return message.Fail("failed to parse data")
	}

	allConfigurations := parameters[3].(*key_value.List)
	configurations := filterConfiguration(allConfigurations, &topicFilter)
	if len(configurations) == 0 {
		reply := FilterSmartcontractKeysReply{
			SmartcontractKeys: []smartcontract_key.Key{},
			TopicStrings:      []topic.Id{},
		}
		replyMessage, err := command.Reply(&reply)
		if err != nil {
			return message.Fail("failed to reply: " + err.Error())
		}
		return replyMessage
	}

	allSmartcontracts := parameters[2].(*key_value.List)
	smartcontracts, topicStrings, err := filterSmartcontract(configurations, allSmartcontracts)
	if err != nil {
		return message.Fail("failed to reply: " + err.Error())
	}

	keys := make([]smartcontract_key.Key, len(smartcontracts))
	//for i, sm := range smartcontracts {
	//keys[i] = sm.SmartcontractKey
	//}

	reply := FilterSmartcontractKeysReply{
		SmartcontractKeys: keys,
		TopicStrings:      topicStrings,
	}
	replyMessage, err := command.Reply(&reply)
	if err != nil {
		return message.Fail("failed to reply: " + err.Error())
	}
	return replyMessage
}
