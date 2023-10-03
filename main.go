package main

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jstemmer/go-junit-report/v2/junit"
)

//go:embed style.css
var styles string

func printTest(testSuite junit.Testsuite, testCase junit.Testcase) {
	id := fmt.Sprintf("%s.%s.%s", testSuite.Name, testCase.Classname, testCase.Name)
	class, text := "passed", "Pass"
	failure := testCase.Failure
	if failure != nil {
		class, text = "failed", "Fail"
	}
	skipped := testCase.Skipped
	if skipped != nil {
		class, text = "skipped", "Skip"
	}

	fmt.Printf("<div class='%s' id='%s'>\n", class, "div-"+id)

	fmt.Printf("<label for='%s' class='toggle'>%s<span class='badge'>%s</span></a></label>\n", id, testCase.Name, text)
	fmt.Printf("<input type='checkbox' name='one' id='%s' class='hide-input'>", id)
	fmt.Printf("<div class='toggle-el'>\n")
	if failure != nil {
		failure.Data = strings.ReplaceAll(failure.Data, `<bool>`, `"bool"`)
		testCase.SystemErr.Data = strings.ReplaceAll(testCase.SystemErr.Data, `<bool>`, `"bool"`)
		fmt.Printf("<div class='content'>%s</div>\n", failure.Data)
		fmt.Printf("<div class='content'>%s</div>\n", testCase.SystemErr.Data)
	} else if skipped != nil {
		fmt.Printf("<div class='content'>%s</div>\n", skipped.Message)
	}
	d, _ := time.ParseDuration(testCase.Time)
	fmt.Printf("<p class='duration' title='Test duration'>%v</p>\n", d)
	fmt.Printf("</div>\n")
	fmt.Printf("</div>\n")

}

func main() {
	suites := &junit.Testsuites{}

	err := xml.NewDecoder(os.Stdin).Decode(suites)
	if err != nil {
		panic(err)
	}

	fmt.Println("<html>")
	fmt.Println("<head>")
	fmt.Println("<meta charset=\"UTF-8\">")
	fmt.Println("<style>")
	fmt.Println(styles)
	fmt.Println("</style>")
	fmt.Println("</head>")
	fmt.Println("<body>")
	failures, total := 0, 0
	for _, s := range suites.Suites {
		failures += s.Failures
		total += len(s.Testcases)
	}
	fmt.Printf("<p>%d of %d tests failed</p>\n", failures, total)
	for _, s := range suites.Suites {
		if s.Failures > 0 {
			printSuiteHeader(s)
			for _, c := range s.Testcases {
				if f := c.Failure; f != nil {
					printTest(s, c)
				}
			}
		}
	}
	for _, s := range suites.Suites {
		printSuiteHeader(s)
		for _, c := range s.Testcases {
			if c.Failure == nil {
				printTest(s, c)
			}
		}
	}
	fmt.Println("</body>")
	fmt.Println("</html>")
}

func printSuiteHeader(s junit.Testsuite) {
	fmt.Println("<h4>")
	fmt.Println(s.Name)
	for _, p := range *s.Properties {
		if strings.HasPrefix(p.Name, "coverage.") {
			v, _ := strconv.ParseFloat(p.Value, 32)
			fmt.Printf("<span class='coverage' title='%s'>%.0f%%</span>\n", p.Name, v)
		}
	}
	fmt.Println("</h4>")
}
