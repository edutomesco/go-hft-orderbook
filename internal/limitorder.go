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
	//return this.orders.Size()
}

func (this *LimitOrder) Enqueue(o *Order) {
	//this.orders.Enqueue(o)
	o.Limit = this
	_ = this.cacheRepository.Enqueue(context.Background(), fmt.Sprintf("%f", this.Price), o)

	this.totalVolume += o.Volume
}

func (this *LimitOrder) Dequeue() *Order {
	/*if this.orders.IsEmpty() {
		return nil
	}*/

	ok, _ := this.cacheRepository.IsEmpty(context.Background(), fmt.Sprintf("%f", this.Price))
	if ok {
		return nil
	}

	//o := this.orders.Dequeue()

	o, _ := this.cacheRepository.Dequeue(context.Background(), fmt.Sprintf("%f", this.Price))

	this.totalVolume -= o.Volume
	return o
}

func (this *LimitOrder) Delete(o *Order) {
	if o.Limit != this {
		panic("order does not belong to the limit")
	}

	//this.orders.Delete(o)
	_ = this.cacheRepository.Delete(context.Background(), fmt.Sprintf("%f", this.Price), o)

	o.Limit = nil
	this.totalVolume -= o.Volume
}

func (this *LimitOrder) Clear() {
	/*q := NewOrdersQueue()
	this.orders = &q*/

	_ = this.cacheRepository.DeleteAll(context.Background(), fmt.Sprintf("%f", this.Price))
	this.totalVolume = 0
}
