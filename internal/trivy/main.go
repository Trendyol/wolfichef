package trivy

import (
	"context"
	"encoding/json"
	dbtypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/commands/artifact"
	"github.com/aquasecurity/trivy/pkg/flag"
	"github.com/aquasecurity/trivy/pkg/types"
	"io"
	_ "modernc.org/sqlite"
	"os"
	"time"
)

type Image struct {
	Name   string
	Output string
}

type Result struct {
	*types.Report
}

var defaultOptions flag.Options

func init() {
	defaultOptions = flag.Options{
		GlobalOptions: flag.GlobalOptions{
			CacheDir: "/tmp/trivy",
			Quiet:    true,
			Timeout:  1 * time.Minute,
		},
		ScanOptions: flag.ScanOptions{
			Scanners: types.AllScanners,
		},
		VulnerabilityOptions: flag.VulnerabilityOptions{
			VulnType: types.VulnTypes,
		},
		ReportOptions: flag.ReportOptions{
			Format:     "json",
			Severities: []dbtypes.Severity{dbtypes.SeverityHigh, dbtypes.SeverityCritical},
		},
		DBOptions: flag.DBOptions{
			DBRepository:     "ghcr.io/aquasecurity/trivy-db",
			JavaDBRepository: "ghcr.io/aquasecurity/trivy-java-db",
			NoProgress:       true,
		},
	}
}

func (img *Image) Scan() error {
	f, err := os.CreateTemp("/tmp", "trivy-*")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img.Output = f.Name()

	var writer io.Writer = f

	opts := defaultOptions
	opts.ImageOptions = flag.ImageOptions{Input: img.Name}
	opts.ReportOptions.Output = writer

	err = artifact.Run(context.TODO(), opts, artifact.TargetImageArchive)
	if err != nil {
		return err
	}
	return nil
}

func (img *Image) Result() (Result, error) {
	var report Result
	data, err := os.ReadFile(img.Output)
	defer os.Remove(img.Output)
	defer os.Remove(img.Name)

	if err != nil {
		return report, err
	}

	err = json.Unmarshal(data, &report)

	if err != nil {
		return report, err
	}

	return report, nil
}

func (r *Result) Categorize() (criticalCount, highCount int) {
	for _, result := range r.Results {
		for _, vuln := range result.Vulnerabilities {
			switch vuln.Severity {
			case dbtypes.SeverityCritical.String():
				criticalCount++
			case dbtypes.SeverityHigh.String():
				highCount++
			}
		}
	}

	return criticalCount, highCount
}
