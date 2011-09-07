package main

import (
	"fmt"
	"http"
	"url"
)

// support for running "redwood -test http://example.com"

// runURLTest prints debugging information about how the URL and its content would be rated.
func runURLTest(u string) {
	url_, err := url.Parse(u)
	if err != nil {
		fmt.Println("Could not parse the URL.")
		return
	}

	if url_.Scheme == "" {
		url2, err := url.Parse("http://" + u)
		if err == nil {
			url_ = url2
		}
	}

	fmt.Println("URL:", url_)
	fmt.Println()

	urlTally := URLRules.MatchingRules(url_)
	if len(urlTally) == 0 {
		fmt.Println("No URL rules match.")
	} else {
		fmt.Println("The following URL rules match:")
		for s, _ := range urlTally {
			fmt.Println(s)
		}
	}

	urlScores := categoryScores(urlTally)
	if len(urlScores) > 0 {
		fmt.Println()
		fmt.Println("The request has the following category scores:")
		printSortedTally(urlScores)

		blocked := blockedCategories(urlScores)
		if len(blocked) > 0 {
			fmt.Println()
			fmt.Println("The request is blocked by the following categories:")
			for _, c := range blocked {
				fmt.Println(c)
			}
			return
		}
	}

	fmt.Println()
	fmt.Println("Downloading content...")
	res, err := http.Get(url_.String())
	if err != nil {
		fmt.Println(err)
		return
	}

	phraseTally := phrasesInResponse(res)

	fmt.Println()

	if len(phraseTally) == 0 {
		fmt.Println("No content phrases match.")
	} else {
		fmt.Println("The following content phrases match:")
		printSortedTally(phraseTally)
	}

	pageScores := categoryScores(phraseTally)
	if len(pageScores) > 0 {
		// Add the URL scores to the page scores.
		for c, s := range urlScores {
			pageScores[c] += s
		}
		fmt.Println()
		fmt.Println("The response has the following category scores:")
		printSortedTally(pageScores)

		blocked := blockedCategories(pageScores)
		if len(blocked) > 0 {
			fmt.Println()
			fmt.Println("The page is blocked by the following categories:")
			for _, c := range blocked {
				fmt.Println(c)
			}
		}
	}
}

// printSortedTally prints tally's keys and values in descending order by value.
func printSortedTally(tally map[string]int) {
	for _, rule := range sortedKeys(tally) {
		fmt.Println(rule, tally[rule])
	}
}
