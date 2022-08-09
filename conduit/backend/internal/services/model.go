package services

const baseDomain = "http://localhost:"

type Service struct {
	BaseDomain string `default:"http://localhost:"`
}

func GetService() Service {
	return Service{BaseDomain: baseDomain}
}
