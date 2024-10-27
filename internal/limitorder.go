package internal

import (
	"context"
	"fmt"
)

// Limit price orders combined as a FIFO queue
type LimitOrder struct {
	Price float64

	orders          *ordersQueue
	totalVolume     float64
	cacheRepository Cache
}

func NewLimitOrder(price float64, cr Cache) LimitOrder {
	q := NewOrdersQueue()
	return LimitOrder{
		Price:           price,
		orders:          &q,
		cacheRepository: cr,
	}
}

func (this *LimitOrder) TotalVolume() float64 {
	return this.totalVolume
}

func (this *LimitOrder) Size() int {
	res, _ := this.cacheRepository.Size(context.Background(), fmt.Sprintf("%f", this.Price))
	return res
}

func (this *LimitOrder) Enqueue(o *Order) {
	o.Limit = this
	if err := this.cacheRepository.Enqueue(context.Background(), fmt.Sprintf("%f", this.Price), o); err != nil {
		panic(fmt.Sprintf("error from redis: %v", err.Error()))
	}

	this.totalVolume += o.Volume
}

func (this *LimitOrder) Dequeue() *Order {
	ok, err := this.cacheRepository.IsEmpty(context.Background(), fmt.Sprintf("%f", this.Price))
	if err != nil {
		panic(fmt.Sprintf("error from redis: %v", err.Error()))
	}
	if ok {
		return nil
	}

	o, err := this.cacheRepository.Dequeue(context.Background(), fmt.Sprintf("%f", this.Price))
	if err != nil {
		panic(fmt.Sprintf("error from redis: %v", err.Error()))
	}

	this.totalVolume -= o.Volume
	return o
}

func (this *LimitOrder) Delete(o *Order) {
	if o.Limit != this {
		panic("order does not belong to the limit")
	}

	if err := this.cacheRepository.Delete(context.Background(), fmt.Sprintf("%f", this.Price), o); err != nil {
		panic(fmt.Sprintf("error from redis: %v", err.Error()))
	}

	o.Limit = nil
	this.totalVolume -= o.Volume
}

func (this *LimitOrder) Clear() {
	if err := this.cacheRepository.DeleteAll(context.Background(), fmt.Sprintf("%f", this.Price)); err != nil {
		panic(fmt.Sprintf("error from redis: %v", err.Error()))
	}
	this.totalVolume = 0
}
