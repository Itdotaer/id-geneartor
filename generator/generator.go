package generator

import (
	"fmt"
	"sync"
)

// Generator 生成器
type Generator struct {
	Mutex    sync.Mutex
	Business string          // 业务标识
	Segments []*Segment      // 双Buffer, 最少0个, 最多2个号段在内存
	IsAlloc  bool            // 是否在分配号段
	Map      map[int64]int64 // 这个为本地统计是否生成ID是否冲突，正式版本可以不要
}

// Segment 号段
type Segment struct {
	CurrentId int64 // 当前号码
	Offset    int64 // 消费偏移
	Step      int64 // 步长
}

func (gen *Generator) GenerateNextId() int64 {
	segment := gen.Segments[0]
	nextId := segment.CurrentId - segment.Step + segment.Offset

	gen.Segments[0].Offset++
	if nextId+1 >= segment.CurrentId {
		gen.Segments = append(gen.Segments[:0], gen.Segments[1:]...) // 弹出第一个seg, 后续seg向前移动
	}

	// 这个为本地统计是否生成ID是否冲突，正式版本可以不要
	if value, ok := gen.Map[nextId]; ok {
		println(fmt.Sprintf("业务%s冲突:%d", gen.Business, value))
	}

	return nextId
}

func (gen *Generator) Left() int64 {
	count := int64(0)
	for i := 0; i < len(gen.Segments); i++ {
		count += gen.Segments[i].Step - gen.Segments[i].Offset
	}
	return count
}

func (gen *Generator) AppendSegment() error {
	var (
		segment *Segment
		err     error
	)

	gen.Mutex.Lock()
	if len(gen.Segments) <= 1 {
		if segment, err = GMysql.NextSegment(gen.Business); err != nil {
			return err
		}

		gen.Segments = append(gen.Segments, segment)
	}
	gen.Mutex.Unlock()
	gen.IsAlloc = false

	return nil
}
