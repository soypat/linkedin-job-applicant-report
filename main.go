package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/exp/slices"
)

func main() {
	var allResults []result
	for i := 1; i <= 8; i++ {
		filename := fmt.Sprintf("p%d.html", i)
		fp, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		results, err := parseFile(fp)
		fmt.Println("got results:", len(results), "in", filename)
		if err != nil {
			panic(err)
		}
		allResults = append(allResults, results...)
	}
	resultsByApplicants := append([]result{}, allResults...)
	slices.SortStableFunc(allResults, func(a, b result) bool {
		return a.jobname < b.jobname
	})
	slices.SortStableFunc(resultsByApplicants, func(a, b result) bool {
		return a.applicants > b.applicants
	})

	locationBins := make(map[string]struct{})
	jobBins := make(map[string]struct{})
	for _, result := range allResults {
		locationBins[result.location] = struct{}{}
		jobBins[result.jobname] = struct{}{}
	}

	fpRes, _ := os.Create("results.csv")
	defer fpRes.Close()
	w := csv.NewWriter(fpRes)
	defer w.Flush()
	w.Write(result{}.CSVHeader())
	for _, result := range allResults {
		w.Write(result.CSVRecord())
	}

	// Location Analyses
	fploc, _ := os.Create("report_by_location.md")
	defer fploc.Close()
	fploc.WriteString("# Revelo Job Market Analysis (by Location)\n\n")
	for location := range locationBins {
		fploc.WriteString(fmt.Sprintf("### By Applicants in %s\n", location))
		fploc.WriteString("| Job Name | Applicants | Pay USD/mo |\n| --- | --- | --- |\n")
		for _, result := range resultsByApplicants {
			if result.location != location {
				continue
			}
			fploc.WriteString(fmt.Sprintf("| %s | %d | %d |\n", result.jobname, result.applicants, result.paygrade))
		}
	}

	// Technology Analyses
	fpjob, _ := os.Create("report_by_job.md")
	defer fpjob.Close()
	fpjob.WriteString("# Revelo Job Market Analysis (by Job description)\n\n")
	for job := range jobBins {
		fpjob.WriteString(fmt.Sprintf("### Job description: %s\n", job))
		fpjob.WriteString("| Location | Applicants | Pay USD/mo |\n| --- | --- | --- |\n")
		for _, result := range resultsByApplicants {
			if result.jobname != job {
				continue
			}
			fpjob.WriteString(fmt.Sprintf("| %s | %d | %d |\n", result.location, result.applicants, result.paygrade))
		}
	}

}

type result struct {
	jobname    string
	location   string
	paygrade   int
	applicants int
}

func (r result) String() string {
	s := fmt.Sprintf("%s @ %s %d/mo", r.jobname, r.location, r.paygrade)
	if r.applicants > 0 {
		s += fmt.Sprintf(" (%d applicants)", r.applicants)
	}
	return s
}
func (r result) CSVHeader() []string {
	return []string{"jobname", "location", "pay (USD)", "applicants"}
}
func (r result) CSVRecord() []string {
	return []string{r.jobname, r.location, strconv.Itoa(r.paygrade), strconv.Itoa(r.applicants)}
}

func parseFile(r io.Reader) ([]result, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	selection := doc.Find(`#main > div > div.scaffold-layout__list > div > ul > li.ember-view > div > div`)

	var skipped []string
	var results []result
	// fmt.Println(selection.Children().Text())
	selection.Each(func(i int, s *goquery.Selection) {
		got := strings.Fields(s.Text())
		filtered := strings.Join(got, " ")
		original := filtered
		jobname, filtered, ok := strings.Cut(filtered, "-")
		if !ok {
			skipped = append(skipped, original)
			return
		}
		paygrade, err := extractIntBefore(filtered, "/month")
		if err != nil {
			skipped = append(skipped, original)
			return
		}

		jobname = strings.TrimSpace(jobname)
		_, filtered, ok = strings.Cut(filtered, "Revelo ")
		if !ok {
			skipped = append(skipped, original)
			return
		}
		location, filtered, ok := strings.Cut(filtered, " Hide job")
		if !ok {
			skipped = append(skipped, original)
			return
		}
		location = strings.TrimSpace(strings.TrimSuffix(location, "Remote"))
		applicants, _ := extractIntBefore(filtered, " applicants")

		results = append(results, result{
			jobname:    jobname,
			location:   location,
			applicants: int(applicants),
			paygrade:   paygrade,
		})
	})
	return results, nil
}

func extractIntBefore(s, suffix string) (int, error) {
	got := strings.Index(s, suffix)
	if got < 1 {
		return 0, fmt.Errorf("no %q found in %q", suffix, s)
	}
	slice := strings.Fields(s[:got])
	if len(slice) < 1 {
		return 0, fmt.Errorf("no fields found in %q", s)
	}
	Int, err := strconv.ParseInt(slice[len(slice)-1], 10, 32)
	return int(Int), err
}
