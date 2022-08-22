package arbitrage

import (
	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/transactionlogger"
	"github.com/gstdl/crypto-arbitrage/internal/pkg/luno"
)

func (q *queue) Enqueue(element interface{}) {
	q.elements = append(q.elements, element)
}

func (q *queue) Dequeue() (firstElement interface{}, ok bool) {
	if len(q.elements) <= 0 {
		return
	}
	firstElement, q.elements, ok = q.elements[0], q.elements[1:], true
	return
}

// Implements BFS on each arbitrage target to find all unique paths
func FindPathsFromLunoPairSetup(ps map[string][]luno.PreOrder) (new ArbitrageService) {
	p := make(PathMap)

	for target, options := range ps {
		var paths []Path
		var q queue
		for _, option := range options {
			q.Enqueue(Path{option})
		}

		for len(q.elements) > 0 {
			currentPathInterface, ok := q.Dequeue()
			if !ok {
				break
			}
			if currentPath, ok := currentPathInterface.(Path); ok {
				if len(currentPath) < len(ps) {
					lastPreOrder := currentPath[len(currentPath)-1]
					if options, ok := ps[lastPreOrder.OrderResult()]; ok {
						for _, option := range options {
							if newPath, ok := currentPath.Append(option); ok {
								if option.OrderResult() == target {
									if len(newPath) > 2 {
										paths = append(paths, newPath)
									}
								} else {
									q.Enqueue(newPath)
								}
							}
						}
					}
				}
			}

			p[target] = paths

		}
	}

	logger := transactionlogger.New("pkg/arbitrage", "LUNO")

	new = ArbitrageService{paths: p, logger: logger}

	return
}

func (p Path) Append(toAppend PreOrder) (newPath Path, ok bool) {
	for _, cur := range p {
		if cur.OrderResult() == toAppend.OrderResult() {
			return
		}
	}
	newPath, ok = append(append(Path{}, p...), toAppend), true
	return
}

func (p Path) GetInitialCurrency() (curr string) {
	curr = p[0].OrderBase()
	return
}