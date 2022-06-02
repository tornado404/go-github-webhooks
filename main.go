package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/webhooks/v6/github"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	router.POST("/ping", updateHandler)
	router.GET("/ping", printPong)
	router.Run(":8081")

	// By default it serves on :8080 unless a
	// PORT environment
	log.Info("PaddleFlow webhook server exiting")

}

func printPong(c *gin.Context) {
	log.Info("get pong")
	c.String(http.StatusOK, "pong")
}

func Command(cmd string) error {
	c := exec.Command("bash", "-c", cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr

	// 此处是windows版本
	// c := exec.Command("cmd", "/C", cmd)
	if err := c.Run(); err != nil {
		log.Errorf("run command failed. err: %v, stderr: %s", err, stderr.String())
		return err
	} else {
		log.Infof("stdout: %s", out.String())
	}
	return nil
}

func updateHandler(c *gin.Context) {
	log.Infof("get a request from host[%s]", c.ClientIP())

	hook, _ := github.New(github.Options.Secret("YourSecret"))
	payload, err := hook.Parse(c.Request, github.WorkflowJobEvent, github.WorkflowRunEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
		log.Errorf("err: %v", err)
	}
	log.Infof("parse request")
	switch payload.(type) {
	case github.WorkflowRunPayload:
		jobBody := payload.(github.WorkflowRunPayload)
		if jobBody.Action != "completed" {
			log.Infof("uncompleted action, pass it")
			break
		}
		// Do whatever you want from here...
		jsonBytes, _ := json.Marshal(jobBody)
		log.Infof("WorkflowRunPayload: %s", string(jsonBytes))
		Command("./build.sh")
		// if build.sh absent, can use following
		//Command("rm -rf PaddleFlow && git clone -c http.proxy=\"https://agent.baidu.com:8118\" https://github.com/tornado404/go-github-webhooks.git ")
	// more case reference
	default:
		log.Infof("no relation request")
	}
	c.String(http.StatusOK, "pong")
}
