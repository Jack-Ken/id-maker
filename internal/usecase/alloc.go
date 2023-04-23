package usecase

import (
	"context"
	"errors"
	"id-maker/internal/entity"
	"id-maker/pkg/snowflake"
	"sync"
	"time"
)

type Alloc struct {
	Mu        sync.RWMutex
	BizTagMap map[string]*BizAlloc
}

type BizAlloc struct {
	Mu      sync.Mutex
	BizTag  string
	IdArray []*IdArray // 双缓存避免
	GetDb   bool
}
type IdArray struct {
	Cur   int64 // 当前发送到哪个位置
	Start int64 // 最小值
	End   int64 // 最大值
}

func (uc *SegmentUseCase) NewAllocId() (a *Alloc, err error) {
	var res []entity.Segments
	if res, err = uc.repo.GetList(); err != nil {
		return
	}
	a = &Alloc{
		BizTagMap: make(map[string]*BizAlloc),
	}
	for _, v := range res {
		a.BizTagMap[v.BizTag] = &BizAlloc{
			BizTag:  v.BizTag,
			GetDb:   false,
			IdArray: make([]*IdArray, 0),
		}
	}
	return
}

func (uc *SegmentUseCase) NewSnowFlake() (*snowflake.Worker, error) {
	return snowflake.NewWorker(1)
}

func (b *BizAlloc) GetId(uc *SegmentUseCase) (id int64, err error) {
	var (
		canGetId    bool
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	)
	b.Mu.Lock()
	if b.LeftIdCount() > 0 { // 查询分配数组中是否还有剩余的未分配ID
		id = b.PopId()
		canGetId = true
	}
	// 分配ID数组不足，开始新的goruntime去申请新的数组(双缓存)
	if len(b.IdArray) <= 1 && !b.GetDb {
		b.GetDb = true
		b.Mu.Unlock()
		go b.GetIdArray(cancel, uc)
	} else {
		b.Mu.Unlock()
		defer cancel()
	}
	if canGetId { // 已经分配过了ID直接返回
		return
	}
	select { // select的作用是阻塞进程直到获取新的数组进来
	case <-ctx.Done():
	}
	b.Mu.Lock()
	if b.LeftIdCount() > 0 {
		id = b.PopId()
	} else {
		err = errors.New("no get id")
	}
	b.Mu.Unlock()
	return
}

// 申请新的分配数组
func (b *BizAlloc) GetIdArray(cancel context.CancelFunc, uc *SegmentUseCase) {
	var (
		tryNum int
		ids    *entity.Segments
		err    error
	)
	defer cancel()
	for {
		if tryNum >= 3 {
			b.GetDb = false
			break
		}
		b.Mu.Lock()
		if len(b.IdArray) <= 1 {
			b.Mu.Unlock()
			ids, err = uc.repo.GetNextId(b.BizTag)
			if err != nil {
				tryNum++
			} else {
				tryNum = 0
				b.Mu.Lock()
				b.IdArray = append(b.IdArray, &IdArray{Start: ids.MaxId, End: ids.MaxId + ids.Step})
				if len(b.IdArray) > 1 {
					b.GetDb = false
					b.Mu.Unlock()
					break
				} else {
					b.Mu.Unlock()
				}

			}
		} else {
			b.Mu.Unlock()
		}
	}
}

// 查询剩余可分配的ID数量

func (b *BizAlloc) LeftIdCount() (count int64) {
	for _, v := range b.IdArray {
		arr := v
		count += arr.End - arr.Start - arr.Cur
	}
	return count
}

// 分配ID

func (b *BizAlloc) PopId() (id int64) {
	id = b.IdArray[0].Start + b.IdArray[0].Cur // 开始位置加上分配的次数
	b.IdArray[0].Cur++                         //分配次数 +1
	if id+1 >= b.IdArray[0].End {              //该数组里面没有ID了
		b.IdArray = append(b.IdArray[:0], b.IdArray[1:]...) //把分配完的数组移除
		//b.IdArray = b.IdArray[1:] //把分配完的数组移除
		// 两个数组可能不是相连的，当一个数组分配完之后，从新的分配数组重新分配id
		id = b.IdArray[0].Start + b.IdArray[0].Cur // 开始位置加上分配的次数
		b.IdArray[0].Cur++
	}
	return
}
