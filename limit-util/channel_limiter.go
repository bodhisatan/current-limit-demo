package limit_util

type ChannelLimiter struct {
	bufferChannel chan int
}

func NewChannelLimiter(limit int) *ChannelLimiter {
	return &ChannelLimiter{bufferChannel: make(chan int, limit)}
}

func (c *ChannelLimiter) Allow() bool {
	select {
	case c.bufferChannel <- 1:
		return true
	default:
		return false
	}
}

func (c *ChannelLimiter) Release() {
	<-c.bufferChannel
}
