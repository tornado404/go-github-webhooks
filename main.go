package main

import (
	"fmt"
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
	// 此处是windows版本
	// c := exec.Command("cmd", "/C", cmd)
	output, err := c.CombinedOutput()
	fmt.Println(string(output))
	return err
}

func updateHandler(c *gin.Context) {
	log.Infof("get a request")

	hook, _ := github.New(github.Options.Secret("paddleflow123"))
	payload, err := hook.Parse(c.Request, github.WorkflowJobEvent, github.WorkflowRunEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
		log.Errorf("err: %v", err)
	}
	log.Infof("parse request")
	switch payload.(type) {
	case github.WorkflowJobPayload:
		jobBody := payload.(github.WorkflowJobPayload)
		// Do whatever you want from here...
		log.Infof("WorkflowJobPayload: %+v", jobBody)
		Command("ls")
	case github.WorkflowRunPayload:
		jobBody := payload.(github.WorkflowRunPayload)
		// Do whatever you want from here...
		log.Infof("WorkflowRunPayload: %+v", jobBody)
		Command("ls")
	default:
		log.Infof("bad request body")
	}
	c.String(http.StatusOK, "pong")
}
