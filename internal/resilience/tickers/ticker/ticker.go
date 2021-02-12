package ticker

// Ticker delivers `ticks' of a clock at intervals.
type Ticker interface {

	// Stop turns off a ticker. After Stop, no more ticks will be sent.
	Stop()
}
