/*
	SegmentUseCase业务逻辑对外实现的类，使用时通过使用创建SegmentUseCase实例来使用去功能
*/
package usecase

import (
	"id-maker/internal/entity"
	"id-maker/pkg/snowflake"
)

type SegmentUseCase struct {
	repo      SegmentRepo
	alloc     *Alloc
	snowFlake *snowflake.Worker
}

func New(r SegmentRepo) *SegmentUseCase {
	var err error

	s := &SegmentUseCase{}
	s.repo = r
	if s.alloc, err = s.NewAllocId(); err != nil {
		panic(err)
	}
	if s.snowFlake, err = s.NewSnowFlake(); err != nil {
		panic(err)
	}
	return s
}

func (uc *SegmentUseCase) CreateTag(e *entity.Segments) (err error) {
	if err = uc.repo.Add(e); err != nil {
		return
	}
	b := &BizAlloc{
		BizTag:  e.BizTag,
		GetDb:   false,
		IdArray: make([]*IdArray, 0),
	}
	b.IdArray = append(b.IdArray, &IdArray{
		Cur:   1,
		Start: 0,
		End:   e.Step,
	})
	uc.alloc.BizTagMap[e.BizTag] = b
	return
}

func (uc *SegmentUseCase) SnowFlakeGetId() int64 {
	return uc.snowFlake.GetId()
}

func (uc *SegmentUseCase) GetId(tag string) (id int64, err error) {
	uc.alloc.Mu.Lock()
	defer uc.alloc.Mu.Unlock()

	val, ok := uc.alloc.BizTagMap[tag]
	if !ok {
		// tag不存在就创建新的tag
		if err = uc.CreateTag(&entity.Segments{
			BizTag: tag,
			MaxId:  1,
			Step:   1000,
		}); err != nil {
			return 0, err
		}
		val, _ = uc.alloc.BizTagMap[tag]
	}
	return val.GetId(uc)
}
