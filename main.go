package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/nlopes/slack"
)

type Service interface {
	Has2FA(userID string) (bool, error)
}

type ServiceConfig struct {
	Type        string
	Description string
	Config      string
	UserMap     map[string]string
}

type Config struct {
	Services map[string]ServiceConfig
}

func getEnv(env []string, key string) string {
	for _, envEntry := range env {
		if strings.HasPrefix(envEntry, key+"=") {
			return strings.TrimPrefix(envEntry, key+"=")
		}
	}

	return ""
}

func ParseConfig(env []string) Config {
	serviceNamesAndTypes := strings.Split(getEnv(env, "SERVICES"), ",")

	config := Config{Services: make(map[string]ServiceConfig, len(serviceNamesAndTypes))}

	for _, serviceNameAndTypeString := range serviceNamesAndTypes {
		serviceNameAndType := strings.SplitN(serviceNameAndTypeString, ":", 2)
		serviceName := serviceNameAndType[0]
		serviceType := serviceNameAndType[1]

		userMapEntries := strings.Split(getEnv(env, strings.ToUpper(serviceName)+"_USER_MAP"), ",")
		userMap := make(map[string]string, len(userMapEntries))
		for _, entry := range userMapEntries {
			fromAndTo := strings.SplitN(entry, ":", 2)
			userMap[fromAndTo[0]] = fromAndTo[1]
		}

		config.Services[serviceName] = ServiceConfig{
			Type:        serviceType,
			Description: getEnv(env, strings.ToUpper(serviceName)+"_DESCRIPTION"),
			Config:      getEnv(env, strings.ToUpper(serviceName)+"_CONFIG"),
			UserMap:     userMap,
		}
	}

	return config
}

func createServices(config Config) map[string]Service {
	services := make(map[string]Service, len(config.Services))
	for alias, serviceConfig := range config.Services {
		service, err := NewService(serviceConfig.Type, serviceConfig.Config)
		if err != nil {
			fmt.Printf("could create %s service: %s\n", alias, err)
			os.Exit(1)
		}
		services[alias] = service
	}

	return services
}

func findStatuses(config Config, services map[string]Service) map[string][]string {
	statuses := make(map[string][]string)

	for alias, serviceConfig := range config.Services {
		for serviceUserID, chatUserID := range serviceConfig.UserMap {
			ok, err := services[alias].Has2FA(serviceUserID)
			if err != nil {
				fmt.Printf("couldn't get 2fa information for %s (chat ID: %s) on %s: %s\n", serviceUserID, chatUserID, alias, err)
				continue
			}

			if !ok {
				if _, ok := statuses[chatUserID]; !ok {
					statuses[chatUserID] = make([]string, 0)
				}

				statuses[chatUserID] = append(statuses[chatUserID], alias)
			}
		}
	}

	return statuses
}

type msgServiceInfo struct {
	Description string
	URL         string
}

func notifyUsers(config Config, status map[string][]string) {
	client := slack.New(os.Getenv("SLACK_TOKEN"))

	t := template.Must(template.New("message").Parse(os.Getenv("MESSAGE_TEMPLATE")))

	for chatUserID, services := range status {
		_, _, imID, err := client.OpenIMChannel(chatUserID)
		if err != nil {
			fmt.Printf("Couldn't open IM channel for user %s: %s\n", chatUserID, err)
			continue
		}

		data := make([]msgServiceInfo, 0, len(services))
		for _, service := range services {
			data = append(data, msgServiceInfo{
				Description: config.Services[service].Description,
				URL:         docURLForService(config.Services[service].Type),
			})
		}

		buf := new(bytes.Buffer)
		err = t.Execute(buf, data)
		if err != nil {
			fmt.Printf("Couldn't create message for %s: %s\n", chatUserID, err)
			continue
		}

		client.PostMessage(imID, buf.String(), slack.PostMessageParameters{
			UnfurlLinks: false,
			AsUser:      false,
		})
	}
}

func main() {
	config := ParseConfig(os.Environ())
	services := createServices(config)
	statuses := findStatuses(config, services)
	notifyUsers(config, statuses)
}
