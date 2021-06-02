package public

import (
	"sync"
	"time"
)

var FlowCountHandler *FlowCounter


type FlowCounter struct {
	RedisFlowCountMap map[string]*RedisFlowCountService
	RedisFlowCountSlice []*RedisFlowCountService
	Locker sync.RWMutex
}

func init()  {
	FlowCountHandler = NewFlowCounter()
}

func NewFlowCounter()  *FlowCounter {
	return &FlowCounter{
		RedisFlowCountMap:   map[string]*RedisFlowCountService{},
		RedisFlowCountSlice: []*RedisFlowCountService{},
		Locker:              sync.RWMutex{},
	}
}

func (f *FlowCounter) GetFlowCounter (appId string) (*RedisFlowCountService, error) {
	for _, flowCounter := range f.RedisFlowCountSlice {
		if flowCounter.AppID == appId {
			return flowCounter, nil
		}
	}
	newFlowCounter, err := NewRedisFlowCountService(appId, time.Second)
	if err != nil {
		return nil, err
	}
	f.RedisFlowCountSlice = append(f.RedisFlowCountSlice, newFlowCounter)
	f.Locker.Lock()
	defer f.Locker.Unlock()
	f.RedisFlowCountMap[appId] = newFlowCounter

	return newFlowCounter, nil
}