package controller

var (
	APIs []ReminderApi

	NoAuthPath map[string][]string
)

func init() {
	NoAuthPath = map[string][]string{}
}

func registerApi(api ReminderApi) {
	APIs = append(APIs, api)
}

func registerNoAuthPath(method, path string) {
	if len(NoAuthPath[method]) == 0 {
		NoAuthPath[method] = []string{path}
		return
	}
	NoAuthPath[method] = append(NoAuthPath[method], path)
}
