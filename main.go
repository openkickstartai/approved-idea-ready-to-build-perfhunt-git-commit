package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	from := flag.String("from", "", "good commit (baseline)")
	to := flag.String("to", "HEAD", "bad commit (latest)")
	cmd := flag.String("cmd", "", "benchmark command to run")
	n := flag.Int("n", 10, "iterations per commit")
	alpha := flag.Float64("alpha", 0.05, "significance level")
	flag.Parse()

	if *from == "" || *cmd == "" {
		fmt.Fprintln(os.Stderr, "PerfHunt — find the commit that killed your performance")
		fmt.Fprintln(os.Stderr, "\nUsage: perfhunt --from <commit> --cmd <benchmark>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	commits, err := gitLog(*from, *to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	commits = append([]string{*from}, commits...)

	iter, benchCmd := *n, *cmd
	runner := func(hash string) ([]float64, error) {
		if err := exec.Command("git", "checkout", hash, "--quiet").Run(); err != nil {
			return nil, fmt.Errorf("checkout %s: %w", hash, err)
		}
		samples := make([]float64, iter)
		for i := range samples {
			start := time.Now()
			if err := exec.Command("sh", "-c", benchCmd).Run(); err != nil {
				return nil, fmt.Errorf("benchmark failed: %w", err)
			}
			samples[i] = time.Since(start).Seconds()
		}
		return samples, nil
	}

	result, err := Hunt(commits, runner, *alpha)
	exec.Command("git", "checkout", "-", "--quiet").Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	printResult(result)
}

func gitLog(from, to string) ([]string, error) {
	out, err := exec.Command("git", "log", "--format=%H", "--reverse", from+".."+to).Output()
	if err != nil {
		return nil, fmt.Errorf("git log: %w", err)
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return nil, fmt.Errorf("no commits in range %s..%s", from, to)
	}
	return strings.Split(s, "\n"), nil
}

func printResult(r *HuntResult) {
	if !r.Found {
		fmt.Println("\n✅ No statistically significant regression detected.")
		return
	}
	fmt.Println("\n🔴 Performance regression found!")
	fmt.Printf("   Commit:   %s\n", r.Culprit.Hash)
	fmt.Printf("   Baseline: %.4fs → Regressed: %.4fs\n", r.Baseline.Mean, r.Culprit.Mean)
	fmt.Printf("   Slowdown: +%.1f%%\n", r.Slowdown)
	fmt.Printf("   P-value:  %.6f (statistically significant)\n", r.PValue)
}
