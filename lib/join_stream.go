package honeybadger

type joinStream struct {
	hb                         *HoneyBadger
	triple                     Triple
	mask                       map[string][]byte
	matcher                    matcherFn
	maskUpdater                maskUpdaterFn
	limit, limitCounter, index uint
	hasEnded                   bool
	tripleCh                   chan Triple
	lastSolution               Solutions
}

func newJoinStream(hb *HoneyBadger, qp QueryPattern, limit, index uint) *joinStream {
	return &joinStream{
		triple:      qp.Triple,
		matcher:     matcher(qp),
		mask:        queryMask(qp),
		maskUpdater: maskUpdater(qp),
		limit:       limit,
		hb:          hb,
		index:       index,
		tripleCh:    make(chan Triple),
	}
}

func (js *joinStream) transform(solutions Solutions) {
	newMask := js.maskUpdater(solutions, js.mask)
	js.lastSolution = solutions

	qp := maskToQueryPattern(newMask)
	js.tripleCh = js.hb.SearchCh(qp)
}

func maskToQueryPattern(mask Solutions) QueryPattern {
	qp := QueryPattern{}

	s, ok := mask["subject"]
	if ok {
		qp.Triple.Subject = s
	}

	p, ok := mask["predicate"]
	if ok {
		qp.Triple.Predicate = p
	}

	o, ok := mask["object"]
	if ok {
		qp.Triple.Object = o
	}

	return qp
}
