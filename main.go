package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	var mailConfig MailConfig
	config := ReadConfigFromEnv()
	log.Println("check", config.ResourceType, config.ResourceName, "with namespace", config.Namespace, "error log content", config.ErrContent, "(ignore \""+config.IgnoreContent+"\")")
	if config.SendEmail {
		mailConfig = ReadMailConfigFromEnv()
	}
	if exist, podName, errorLog := HasErrorLogInResource(config); exist {
		log.Println("send email to ", mailConfig.ToEmailAddress)
		subject := "[alert] error log in " + config.Namespace + " " + podName
		body := "error log \"" + config.ErrContent + "(ignore \"" + config.IgnoreContent + "\")\" exist in pod " + podName + " with namespace " + config.Namespace + "\nerror logs: \n  " + errorLog
		if config.SendEmail {
			SendMail(mailConfig, subject, body)
		} else {
			log.Println(subject)
			log.Println(body)
		}
	} else {
		log.Println("no error log")
	}
}

func HasErrorLogInPod(config K8sLogCheckerConfig, podName string) (bool, string, string) {
	hasError := false
	cmd := exec.Command("kubectl", "-n", config.Namespace, "logs", "--tail", strconv.Itoa(config.ScanLogCount), podName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("execute " + cmd.String() + " fail")
	}
	for _, m := range bytes.Split(out, []byte("\n")) {
		for _, errContent := range strings.Split(config.ErrContent, "|") {
			if len(strings.TrimSpace(config.IgnoreContent)) > 0 {
				for _, ignoreContent := range strings.Split(config.IgnoreContent, "|") {
					if strings.Contains(string(m), errContent) && !strings.Contains(string(m), ignoreContent) {
						hasError = true
					}
				}
			} else {
				if strings.Contains(string(m), errContent) {
					hasError = true
				}
			}
		}

	}
	if hasError {
		return true, podName, string(out)

	} else {
		return false, "", ""
	}
}

func GetResourceReplicas(config K8sLogCheckerConfig) int {
	var rt string
	switch config.ResourceType {
	case StatfulSets:
		rt = "sts"
	case Deployment:
		rt = "deploy"
	default:
		panic("resource type \"" + config.ResourceType + "\" invalid")
	}
	cmd := exec.Command("kubectl", "-n", config.Namespace, "get", rt, config.ResourceName, "-o=jsonpath={.status.replicas}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("get resource rep fail in namespace ", config.Namespace, " ", string(out))
	}
	replica, err := strconv.Atoi(string(out))
	if err != nil {
		log.Fatal("get resource rep fail in namespace ", config.Namespace, " ", string(out))
	}
	return replica
}

func HasErrorLogInResource(config K8sLogCheckerConfig) (bool, string, string) {
	if config.ResourceType == Pod {
		return HasErrorLogInPod(config, config.ResourceName)
	} else {
		resourceReplica := GetResourceReplicas(config)
		for i := 0; i < resourceReplica; i++ {
			podName := config.ResourceName + "-" + strconv.Itoa(i)
			return HasErrorLogInPod(config, podName)
		}
		return false, "", ""
	}
}

// send email

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func SendMail(mailConfig MailConfig, subject string, body string) {
	header := make(map[string]string)
	header["From"] = mailConfig.From
	header["To"] = mailConfig.ToEmailAddress
	header["Subject"] = subject + " " + mailConfig.AddSubject
	header["Content-Type"] = "text/html; charset=UTF-8"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + strings.ReplaceAll(body, "\n", "</br>")

	to := strings.Split(mailConfig.ToEmailAddress, ";")
	address := mailConfig.Host + ":" + strconv.Itoa(mailConfig.Port)

	auth := LoginAuth(mailConfig.From, mailConfig.Password)
	err := smtp.SendMail(address, auth, mailConfig.From, to, []byte(message))
	if err != nil {
		panic(err)
	}
}

// config

type RP string

const (
	Pod         RP = "pod"
	StatfulSets RP = "statfulset"
	Deployment  RP = "deployment"
)

type K8sLogCheckerConfig struct {
	Namespace     string
	ResourceType  RP
	ResourceName  string
	ScanLogCount  int
	ErrContent    string
	IgnoreContent string
	SendEmail     bool
}

type MailConfig struct {
	From           string
	Password       string
	Host           string
	Port           int
	ToEmailAddress string
	AddSubject     string
}

func ReadMailConfigFromEnv() MailConfig {
	mc := MailConfig{}
	mc.From = GetRequiredEnv("MAIL_FROM")
	mc.Password = GetRequiredEnv("MAIL_PASSWORD")
	mc.Host = GetRequiredEnv("MAIL_HOST")
	mc.Port = GetIntEnv("MAIL_PORT", 587)
	mc.ToEmailAddress = GetRequiredEnv("MAIL_TO")
	mc.AddSubject = GetRequiredEnv("MAIL_ADD_SUBJECT")
	return mc
}
func ReadConfigFromEnv() K8sLogCheckerConfig {
	c := K8sLogCheckerConfig{}

	c.Namespace = GetEnv("KLC_NAMESPACE", "default")
	c.ResourceType = GetcResourceTypeEnv("KLC_RESOURCE_TYPE")
	c.ScanLogCount = GetIntEnv("KLC_SCAN_LOG_COUNT", 20)
	c.ResourceName = GetRequiredEnv("KLC_RESOURCE_NAME")
	c.ErrContent = GetRequiredEnv("KLC_ERR_CONTENT")
	c.IgnoreContent = GetEnv("KLC_IGNORE_CONTENT", "")
	c.SendEmail = GetBoolEnv("KLC_SEND_MAIL")

	return c
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetIntEnv(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
		panic("env \"" + key + "\" is not number")
	}
	return fallback
}

func GetRequiredEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("env \"" + key + "\" is required")
}
func GetcResourceTypeEnv(key string) RP {
	if value, ok := os.LookupEnv(key); ok {
		switch value {
		case "pod":
			return Pod
		case "statfulset":
			return StatfulSets
		case "deployment":
			return Deployment
		default:
			panic("env \"" + key + "\" is required")
		}
	}
	panic("env \"" + key + "\" is required")
}

func GetBoolEnv(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}
