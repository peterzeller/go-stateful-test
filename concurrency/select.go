package concurrency

import "github.com/peterzeller/go-fun/zero"

type ChannelReceiver[In, Out any] struct {
	channel   <-chan In
	onReceive func(In) Out
}

func (c ChannelReceiver[In, Out]) Await() Out {
	return zero.Value[Out]()
}

func Receive[In, Out any](channel <-chan In, onReceive func(In) Out) ChannelReceiver[In, Out] {
	return ChannelReceiver[In, Out]{
		channel:   channel,
		onReceive: onReceive,
	}
}

// Select is a generic implementation of select on multiple go-routines
func Select[A, B any]() {

}
