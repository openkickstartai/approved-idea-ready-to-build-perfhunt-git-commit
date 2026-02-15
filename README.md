# 🔍 PerfHunt

**Find the exact commit that killed your performance — with statistical proof.**

PerfHunt uses binary search over Git history + the Mann-Whitney U test to pinpoint performance regressions with statistical significance. Like `git bisect`, but automated, repeatable, and backed by math.

## 🚀 Quick Start

```bash
# Install
go install github.com/perfhunt/perfhunt@latest

# Find which commit made your app slow
perfhunt --from v1.0.0 --to HEAD --cmd "go test -bench=BenchmarkAPI -benchtime=1x" --n 10

# Use with any benchmark command
perfhunt --from abc123 --cmd "./run_benchmark.sh" --alpha 0.01
```

### Output

```
📊 Baseline abc12345: 0.1042s
  ✅ def45678: 0.1038s (p=0.8721)
  ⚠️  ghi99887: 0.3012s (p=0.0003)

🔴 Performance regression found!
   Commit:   ghi9988776655aabb
   Baseline: 0.1042s → Regressed: 0.3012s
   Slowdown: +189.1%
   P-value:  0.000312 (statistically significant)
```

## How It Works

1. **Binary search** over your commit range (tests O(log n) commits)
2. **Runs your benchmark** N times per commit to collect samples
3. **Mann-Whitney U test** compares each sample against the baseline
4. **Reports** the first commit where performance is statistically worse

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--from` | (required) | Good commit (baseline) |
| `--to` | `HEAD` | Bad commit (latest) |
| `--cmd` | (required) | Benchmark command |
| `--n` | `10` | Iterations per commit |
| `--alpha` | `0.05` | Significance level |

## 📊 Why Statistical Testing Matters

Benchmarks are noisy. A single run can vary ±20%. PerfHunt runs your benchmark multiple times and uses the **Mann-Whitney U test** (non-parametric, no normality assumption) to determine if the difference is real or just noise.

---

## 💰 Pricing

| Feature | Free (OSS) | Pro ($19/mo) | Team ($49/mo) | Enterprise ($199/mo) |
|---|---|---|---|---|
| Binary search bisect | ✅ | ✅ | ✅ | ✅ |
| Mann-Whitney U test | ✅ | ✅ | ✅ | ✅ |
| Single metric | ✅ | ✅ | ✅ | ✅ |
| Multi-metric (CPU, mem, latency) | ❌ | ✅ | ✅ | ✅ |
| JSON / JUnit output | ❌ | ✅ | ✅ | ✅ |
| GitHub Actions integration | ❌ | ✅ | ✅ | ✅ |
| Historical trend tracking | ❌ | ❌ | ✅ | ✅ |
| Slack / webhook alerts | ❌ | ❌ | ✅ | ✅ |
| Web dashboard | ❌ | ❌ | ✅ | ✅ |
| SSO & audit logs | ❌ | ❌ | ❌ | ✅ |
| Priority support | ❌ | ❌ | ❌ | ✅ |

### Why Pay?

- **Pro**: Save hours per regression. JSON output plugs into CI gates. Multi-metric catches memory leaks too.
- **Team**: Track performance over weeks. Get Slacked before your users notice.
- **Enterprise**: Compliance, SSO, and a direct line to our engineers.

**ROI**: One engineer spending 4 hours hunting a regression = $200+. PerfHunt finds it in minutes.

## License

MIT — Free core forever. Pro features require a license key.
