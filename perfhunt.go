package main

import (
	"fmt"
	"math"
	"sort"
)

type BenchRunner func(string) ([]float64, error)

type CommitResult struct {
	Hash    string
	Samples []float64
	Mean    float64
}

type HuntResult struct {
	Found    bool
	Baseline *CommitResult
	Culprit  *CommitResult
	PValue   float64
	Slowdown float64
}

func bench(hash string, run BenchRunner) (*CommitResult, error) {
	samples, err := run(hash)
	if err != nil {
		return nil, fmt.Errorf("bench %s: %w", hash, err)
	}
	return &CommitResult{Hash: hash, Samples: samples, Mean: avg(samples)}, nil
}

func Hunt(commits []string, run BenchRunner, alpha float64) (*HuntResult, error) {
	if len(commits) < 2 {
		return nil, fmt.Errorf("need at least 2 commits, got %d", len(commits))
	}
	base, err := bench(commits[0], run)
	if err != nil {
		return nil, err
	}
	h := commits[0]
	if len(h) > 8 {
		h = h[:8]
	}
	fmt.Printf("📊 Baseline %s: %.4fs\n", h, base.Mean)
	lo, hi := 0, len(commits)-1
	var culprit *CommitResult
	for lo < hi-1 {
		mid := (lo + hi) / 2
		r, err := bench(commits[mid], run)
		if err != nil {
			return nil, err
		}
		_, p := MannWhitneyU(base.Samples, r.Samples)
		reg := p < alpha && r.Mean > base.Mean
		icon := "✅"
		if reg {
			icon = "⚠️"
		}
		ch := commits[mid]
		if len(ch) > 8 {
			ch = ch[:8]
		}
		fmt.Printf("  %s %s: %.4fs (p=%.4f)\n", icon, ch, r.Mean, p)
		if reg {
			hi, culprit = mid, r
		} else {
			lo = mid
		}
	}
	if culprit == nil && hi > lo {
		r, err := bench(commits[hi], run)
		if err != nil {
			return nil, err
		}
		_, p := MannWhitneyU(base.Samples, r.Samples)
		if p < alpha && r.Mean > base.Mean {
			culprit = r
		}
	}
	if culprit == nil {
		return &HuntResult{Found: false, Baseline: base}, nil
	}
	_, p := MannWhitneyU(base.Samples, culprit.Samples)
	return &HuntResult{Found: true, Baseline: base, Culprit: culprit,
		PValue: p, Slowdown: (culprit.Mean - base.Mean) / base.Mean * 100}, nil
}

func MannWhitneyU(x, y []float64) (float64, float64) {
	type obs struct {
		val float64
		grp int
	}
	n1, n2 := float64(len(x)), float64(len(y))
	all := make([]obs, 0, len(x)+len(y))
	for _, v := range x {
		all = append(all, obs{v, 0})
	}
	for _, v := range y {
		all = append(all, obs{v, 1})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].val < all[j].val })
	ranks := make([]float64, len(all))
	for i := 0; i < len(all); {
		j := i + 1
		for j < len(all) && all[j].val == all[i].val {
			j++
		}
		r := float64(i+j+1) / 2.0
		for k := i; k < j; k++ {
			ranks[k] = r
		}
		i = j
	}
	r1 := 0.0
	for i, o := range all {
		if o.grp == 0 {
			r1 += ranks[i]
		}
	}
	u := n1*n2 + n1*(n1+1)/2 - r1
	mu := n1 * n2 / 2
	sig := math.Sqrt(n1 * n2 * (n1 + n2 + 1) / 12)
	if sig == 0 {
		return u, 1.0
	}
	z := math.Abs((u - mu) / sig)
	return u, 2 * (1 - 0.5*(1+math.Erf(z/math.Sqrt2)))
}

func avg(s []float64) float64 {
	sum := 0.0
	for _, v := range s {
		sum += v
	}
	return sum / float64(len(s))
}
