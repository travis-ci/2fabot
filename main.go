package main

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/Sirupsen/logrus"
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
			logrus.WithFields(logrus.Fields{
				"service_name": alias,
				"err":          err,
			}).Fatal("couldn't create service")
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
				logrus.WithFields(logrus.Fields{
					"service_user_id": serviceUserID,
					"chat_user_id":    chatUserID,
					"service_name":    alias,
					"err":             err,
				}).Error("couldn't get 2fa information")
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
			logrus.WithFields(logrus.Fields{
				"chat_user_id": chatUserID,
				"err":          err,
			}).Error("couldn't open IM channel")
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
			logrus.WithFields(logrus.Fields{
				"chat_user_id": chatUserID,
				"err":          err,
			}).Error("couldn't create message from template")
			continue
		}

		client.PostMessage(imID, buf.String(), slack.PostMessageParameters{
			UnfurlLinks: false,
			AsUser:      false,
			Username:    "2fabot",
			IconEmoji:   ":lock:",
		})

		logrus.WithFields(logrus.Fields{
			"chat_user_id": chatUserID,
		}).Info("notified user")
	}
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})

	config := ParseConfig(os.Environ())
	services := createServices(config)
	statuses := findStatuses(config, services)
	notifyUsers(config, statuses)
}
