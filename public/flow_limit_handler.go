package public

import (
	"golang.org/x/time/rate"
	"sync"
)

var FlowLimiterHandler *FlowLimiter


type FlowLimiter struct {
	FlowLimiterMap map[string]*FlowLimiterItem
	FlowLimiterSlice []*FlowLimiterItem
	Locker sync.RWMutex
}

type FlowLimiterItem struct {
	ServiceName string
	Limiter *rate.Limiter
}


func init()  {
	FlowLimiterHandler = NewFlowLimiter()
}

func NewFlowLimiter()  *FlowLimiter {
	return &FlowLimiter{
		FlowLimiterMap:   map[string]*FlowLimiterItem{},
		FlowLimiterSlice: []*FlowLimiterItem{},
		Locker:           sync.RWMutex{},
	}
}

func (f *FlowLimiter) GetLimiter (serviceName string, qps float64) *rate.Limiter {
	for _, flowLimiter := range f.FlowLimiterSlice {
		if flowLimiter.ServiceName == serviceName {
			return flowLimiter.Limiter
		}
	}
	limiter := rate.NewLimiter(rate.Limit(qps), int(3 * qps))
	newFlowLimiter := &FlowLimiterItem{
		ServiceName: serviceName,
		Limiter:    limiter ,
	}
	f.FlowLimiterSlice = append(f.FlowLimiterSlice, newFlowLimiter)
	f.Locker.Lock()
	defer f.Locker.Unlock()
	f.FlowLimiterMap[serviceName] = newFlowLimiter

	return limiter
}