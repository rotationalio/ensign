package mock

import (
	"encoding/json"
	"fmt"
	"os"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
)

var jsonpb = &protojson.UnmarshalOptions{
	AllowPartial:   true,
	DiscardUnknown: true,
}

func EventListFixture(path string) (_ []*api.EventWrapper, err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, err
	}
	return UnmarshalEventList(data)
}

func UnmarshalEventList(data []byte) (events []*api.EventWrapper, err error) {
	items := make([]interface{}, 0)
	if err = json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("could not json unmarshal fixture: %w", err)
	}

	events = make([]*api.EventWrapper, 0, len(items))
	for _, item := range items {
		var buf []byte
		if buf, err = json.Marshal(item); err != nil {
			return nil, err
		}

		event := &api.EventWrapper{}
		if err = jsonpb.Unmarshal(buf, event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func TopicListFixture(path string) (_ []*api.Topic, err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, err
	}
	return UnmarshalTopicList(data)
}

func UnmarshalTopicList(data []byte) (topics []*api.Topic, err error) {
	items := make([]interface{}, 0)
	if err = json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("could not json unmarshal fixture: %w", err)
	}

	topics = make([]*api.Topic, 0, len(items))
	for _, item := range items {
		var buf []byte
		if buf, err = json.Marshal(item); err != nil {
			return nil, err
		}

		topic := &api.Topic{}
		if err = jsonpb.Unmarshal(buf, topic); err != nil {
			return nil, err
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

func TopicNamesListFixture(path string) (_ []*api.TopicName, err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, err
	}
	return UnmarshalTopicNamesList(data)
}

func UnmarshalTopicNamesList(data []byte) (names []*api.TopicName, err error) {
	items := make([]interface{}, 0)
	if err = json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("could not json unmarshal fixture: %w", err)
	}

	names = make([]*api.TopicName, 0, len(items))
	for _, item := range items {
		var buf []byte
		if buf, err = json.Marshal(item); err != nil {
			return nil, err
		}

		name := &api.TopicName{}
		if err = jsonpb.Unmarshal(buf, name); err != nil {
			return nil, err
		}

		names = append(names, name)
	}

	return names, nil
}

func TopicInfoListFixture(path string) (_ map[string]*api.TopicInfo, err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, err
	}
	return UnmarshalTopicInfoList(data)
}

func UnmarshalTopicInfoList(data []byte) (infos map[string]*api.TopicInfo, err error) {
	items := make(map[string]interface{}, 0)
	if err = json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("could not json unmarshal fixture: %w", err)
	}

	infos = make(map[string]*api.TopicInfo, len(items))
	for key, item := range items {
		var buf []byte
		if buf, err = json.Marshal(item); err != nil {
			return nil, err
		}

		info := &api.TopicInfo{}
		if err = jsonpb.Unmarshal(buf, info); err != nil {
			return nil, err
		}

		infos[key] = info
	}

	return infos, nil
}
