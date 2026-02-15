package main

import (
	"math"
	"testing"
)

func TestMannWhitneyU_SameDistribution(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{1, 2, 3, 4, 5}
	_, p := MannWhitneyU(x, y)
	if p < 0.05 {
		t.Errorf("identical distributions should not be significant, got p=%.4f", p)
	}
}

func TestMannWhitneyU_DifferentDistributions(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{10, 11, 12, 13, 14}
	_, p := MannWhitneyU(x, y)
	if p >= 0.05 {
		t.Errorf("clearly different distributions should be significant, got p=%.4f", p)
	}
}

func TestMannWhitneyU_OverlappingReturnsValidP(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{3, 4, 5, 6, 7}
	_, p := MannWhitneyU(x, y)
	if p < 0 || p > 1 {
		t.Errorf("p-value out of [0,1] range: %.4f", p)
	}
}

func TestHunt_FindsRegression(t *testing.T) {
	commits := []string{"aaa11111", "bbb22222", "ccc33333", "ddd44444", "eee55555"}
	runner := func(hash string) ([]float64, error) {
		if hash == "ddd44444" || hash == "eee55555" {
			return []float64{0.50, 0.52, 0.48, 0.51, 0.49, 0.50, 0.53, 0.47, 0.51, 0.50}, nil
		}
		return []float64{0.10, 0.12, 0.09, 0.11, 0.10, 0.11, 0.09, 0.10, 0.12, 0.10}, nil
	}
	result, err := Hunt(commits, runner, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Found {
		t.Fatal("expected regression to be found")
	}
	if result.Culprit.Hash != "ddd44444" {
		t.Errorf("expected culprit ddd44444, got %s", result.Culprit.Hash)
	}
	if result.Slowdown < 100 {
		t.Errorf("expected >100%% slowdown, got %.1f%%", result.Slowdown)
	}
}

func TestHunt_NoRegression(t *testing.T) {
	commits := []string{"aaa11111", "bbb22222", "ccc33333"}
	runner := func(hash string) ([]float64, error) {
		return []float64{0.10, 0.11, 0.09, 0.10, 0.10, 0.11, 0.09, 0.10, 0.10, 0.11}, nil
	}
	result, err := Hunt(commits, runner, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Found {
		t.Error("expected no regression to be found")
	}
}

func TestAvg(t *testing.T) {
	m := avg([]float64{2, 4, 6})
	if math.Abs(m-4.0) > 1e-9 {
		t.Errorf("expected mean 4.0, got %f", m)
	}
}

func TestHunt_TwoCommits(t *testing.T) {
	commits := []string{"good1234", "bad56789"}
	runner := func(hash string) ([]float64, error) {
		if hash == "bad56789" {
			return []float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}, nil
		}
		return []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1}, nil
	}
	result, err := Hunt(commits, runner, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Found {
		t.Fatal("expected regression with 2 commits")
	}
	if result.Culprit.Hash != "bad56789" {
		t.Errorf("expected bad56789, got %s", result.Culprit.Hash)
	}
}
