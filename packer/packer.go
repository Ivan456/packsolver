package packer

import (
	"errors"
	"fmt"
	"sort"
)

type Result struct {
	ShippedTotal int            `json:"shipped_total"`
	Packs        map[string]int `json:"packs"`
	PackCount    int            `json:"pack_count"`
}

// Solve calculates which packs to ship according to the rules:
// 1. Only whole packs can be sent.
// 2. Minimize total shipped items.
// 3. Minimize number of packs if multiple combinations tie on total shipped items.
func Solve(amount int, sizes []int) (*Result, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}
	if len(sizes) == 0 {
		return nil, errors.New("sizes must be provided")
	}

	// Remove duplicates and invalid sizes
	sizeSet := map[int]struct{}{}
	for _, s := range sizes {
		if s > 0 {
			sizeSet[s] = struct{}{}
		}
	}
	clean := make([]int, 0, len(sizeSet))
	for s := range sizeSet {
		clean = append(clean, s)
	}
	if len(clean) == 0 {
		return nil, errors.New("no valid sizes")
	}

	// Sort descending to try larger packs first
	sort.Sort(sort.Reverse(sort.IntSlice(clean)))

	type combo struct {
		total int
		count int
		packs map[int]int
	}

	// BFS queue
	queue := []combo{{total: 0, count: 0, packs: map[int]int{}}}

	best := combo{total: -1}

	// visited total -> minimal pack count
	visited := map[int]int{}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		// Check if valid candidate
		if cur.total >= amount {
			if best.total == -1 || cur.total < best.total ||
				(cur.total == best.total && cur.count < best.count) {
				best = cur
			}
			// No need to continue from here
			continue
		}

		for _, s := range clean {
			nextTotal := cur.total + s
			nextCount := cur.count + 1

			if minCount, ok := visited[nextTotal]; ok && nextCount >= minCount {
				continue
			}
			visited[nextTotal] = nextCount

			nextPacks := make(map[int]int)
			for k, v := range cur.packs {
				nextPacks[k] = v
			}
			nextPacks[s]++

			queue = append(queue, combo{total: nextTotal, count: nextCount, packs: nextPacks})
		}
	}

	if best.total == -1 {
		return nil, errors.New("no solution found")
	}

	// Convert pack keys to string
	packsStr := map[string]int{}
	for k, v := range best.packs {
		packsStr[fmt.Sprintf("%d", k)] = v
	}

	return &Result{
		ShippedTotal: best.total,
		Packs:        packsStr,
		PackCount:    best.count,
	}, nil
}
