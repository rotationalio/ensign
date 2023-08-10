package mock

import (
	"encoding/json"
	"fmt"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
)

func UnmarshalTopicList(data []byte, jsonpb *protojson.UnmarshalOptions) (topics []*api.Topic, err error) {
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

func UnmarshalTopicInfoList(data []byte, jsonpb *protojson.UnmarshalOptions) (infos map[string]*api.TopicInfo, err error) {
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
