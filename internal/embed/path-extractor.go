package pathextractor

import (
	"embed"
	"log"
	"regexp"
)

func ExtractAndTransformRoutes(fileFS embed.FS, file string) ([]string, error) {
	routes, err := extractRoutes(fileFS, file)
	if err != nil {
		log.Fatal("couldn't extract routes from ", file, "\n", err)
	}
	newRoutes := transformRoutes(routes)
	return newRoutes, nil
}

func extractRoutes(fileFS embed.FS, file string) ([]string, error) {
	routes := []string{}

	dat, err := fileFS.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// regexp '/route':
	// r := regexp.MustCompile(`"([^{}]*)":`)
	r := regexp.MustCompile(`'([^{}]*)':`) // TODO check the double quote after save
	matches := r.FindAllStringSubmatch(string(dat), -1)

	for _, r := range matches {
		route := r[1]
		routes = append(routes, route)
	}

	return routes, nil
}

func transformRoutes(routes []string) []string {
	newRoutes := []string{}
	for _, route := range routes {
		routeWithNoParams := removeParams(route)
		// removed the trailing stash from /pokemons/
		// or else page refresh won't work on refresh
		newRoute := removeTrailingSlash(routeWithNoParams)
		newRoutes = append(newRoutes, newRoute)
	}
	return newRoutes
}

func removeParams(route string) string {
	reg_dollar := regexp.MustCompile(`\$`)
	found_dollar := reg_dollar.FindStringIndex(route)

	newRoute := route

	if len(found_dollar) > 0 {
		start := found_dollar[0]
		url_after_dolar := route[start:]

		reg_slash := regexp.MustCompile(`/`)
		found_slash := reg_slash.FindStringIndex(url_after_dolar)

		// end := len(url_after_dolar)
		if len(found_slash) > 0 {
			end := found_slash[0]
			newRoute = route[:start] + "*" + route[start+end:]
		} else {
			newRoute = route[:start] + "*"
		}

		// if it contains $ run it again
		found_another_dollar := reg_dollar.FindStringIndex(newRoute)
		if len(found_another_dollar) > 0 {
			newRoute = removeParams(newRoute)
		}
	}

	return newRoute
}

func removeTrailingSlash(route string) string {
	if len(route) > 1 {
		if route[len(route)-1:] == "/" {
			return route[:len(route)-1]
		}
	}
	return route
}
