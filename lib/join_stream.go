package honeybadger

type joinStream struct {
	hb                         *HoneyBadger
	triple                     Triple
	matcher                    matcherFn
	maskUpdater                maskUpdaterFn
	limit, limitCounter, index uint
	hasEnded                   bool
	tripleCh                   chan Triple
}

func newJoinStream(hb *HoneyBadger, qp QueryPattern, limit, index uint) *joinStream {
	return &joinStream{
		triple:  qp.Triple,
		matcher: matcher(qp),
		// mask:        queryMask(t),
		maskUpdater: maskUpdater(qp),
		limit:       limit,
		hb:          hb,
		index:       index,
		tripleCh:    make(chan Triple),
	}
}
