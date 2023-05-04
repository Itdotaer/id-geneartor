package service

import (
	"errors"
	"github.com/itdotaer/id-generator/store"
	"sync"
)

type GeneratorServiceImpl struct {
	Mutex       sync.Mutex
	BusinessMap map[string]*Generator
}

func NewGeneratorService() GeneratorService {
	return &GeneratorServiceImpl{
		BusinessMap: make(map[string]*Generator),
	}
}

func (service *GeneratorServiceImpl) NextId(business string) (int64, error) {
	var (
		gen   *Generator
		exist bool
	)

	service.Mutex.Lock()
	if gen, exist = service.BusinessMap[business]; !exist {
		gen = &Generator{
			Business: business,
			Segments: make([]*Segment, 0),
			IsAlloc:  false,
			Map:      make(map[int64]int64),
			Store:    store.GRedis,
		}

		service.BusinessMap[business] = gen
	}
	service.Mutex.Unlock()

	if len(gen.Segments) <= 1 && !gen.IsAlloc {
		gen.IsAlloc = true
		go gen.AppendSegment()
	}

	if gen.Left() == 0 {
		return 0, errors.New("no left id")
	}
	return gen.GenerateNextId(), nil
}
