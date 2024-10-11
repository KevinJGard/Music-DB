package controller

import "strings"

// addValues appends values from the given section to the results map under the specified key.
func addValues(results map[string][]string, key, seccion string) {
	values := strings.Split(seccion, "&&")
	results[key] = append(results[key], values...)
}

// splitString processes the search string and separates its content into titles, artists, 
// albums, years and genres, according to the set search language.
func splitString(search string) map[string][]string {
	results := map[string][]string{
		"titles": {},
		"artists": {},
		"albums": {},
		"years": {},
		"genres": {},
	}

	sections := strings.Split(search, "||")

	for _, seccion := range sections {
		if strings.HasPrefix(seccion, "ti:") {
			addValues(results, "titles", strings.TrimPrefix(seccion, "ti:"))
		} else if strings.HasPrefix(seccion, "ar:") {
			addValues(results, "artists", strings.TrimPrefix(seccion, "ar:"))
		} else if strings.HasPrefix(seccion, "al:") {
			addValues(results, "albums", strings.TrimPrefix(seccion, "al:"))
		} else if strings.HasPrefix(seccion, "ye:") {
			addValues(results, "years", strings.TrimPrefix(seccion, "ye:"))
		} else if strings.HasPrefix(seccion, "ge:") {
			addValues(results, "genres", strings.TrimPrefix(seccion, "ge:"))
		}
	}
	return results
}