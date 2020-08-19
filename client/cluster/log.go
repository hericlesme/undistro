/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package cluster

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	undistrov1 "github.com/getupcloud/undistro/api/v1alpha1"
	"github.com/getupcloud/undistro/internal/scheme"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	undistroNamespace = "undistro-system"
)

type LogStreamer interface {
	Stream(context.Context, *rest.Config, chan<- string, types.NamespacedName) error
}

type logStreamer struct {
	client     *rest.RESTClient
	kubeClient client.Client
}

func (l *logStreamer) Stream(ctx context.Context, c *rest.Config, logChan chan<- string, nm types.NamespacedName) error {
	var err error
	l.client, err = rest.UnversionedRESTClientFor(c)
	if err != nil {
		return err
	}
	l.kubeClient, err = client.New(c, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return err
	}
	var podList corev1.PodList
	err = l.kubeClient.List(ctx, &podList, client.InNamespace(undistroNamespace))
	if err != nil {
		return err
	}
	logs := l.getLogs(podList.Items)
	return l.streamLogs(ctx, logs, nm, logChan)
}

func (l *logStreamer) streamLogs(ctx context.Context, reqs []logRequest, nm types.NamespacedName, logChan chan<- string) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(reqs))
	for _, req := range reqs {
		go l.readLog(ctx, wg, nm, req, logChan)
	}
	wg.Wait()
	return nil
}

func (l *logStreamer) readLog(ctx context.Context, wg *sync.WaitGroup, nm types.NamespacedName, r logRequest, logChan chan<- string) {
	reader, err := r.req.Stream()
	if err != nil {
		logChan <- fmt.Sprintf("%s - %v", r.podName, err)
		return
	}
	defer reader.Close()
	ready, err := l.clusterIsReady(ctx, nm)
	if err != nil {
		logChan <- fmt.Sprintf("%s - %v", r.podName, err)
		return
	}
	bReader := bufio.NewReader(reader)
	for !ready {
		var str string
		str, err = bReader.ReadString('\n')
		if err != nil && err != io.EOF {
			logChan <- fmt.Sprintf("%s - %v", r.podName, err)
		}
		if strings.Contains(str, nm.Name) {
			logChan <- fmt.Sprintf("%s - %s", r.podName, str)
		}
		ready, err = l.clusterIsReady(ctx, nm)
		if err != nil {
			logChan <- fmt.Sprintf("%s - %v", r.podName, err)
		}
	}
	close(logChan)
}

func (l *logStreamer) clusterIsReady(ctx context.Context, nm types.NamespacedName) (bool, error) {
	var cl undistrov1.Cluster
	err := l.kubeClient.Get(ctx, nm, &cl)
	if err != nil {
		return false, err
	}
	return cl.Status.Ready, nil
}

func (l *logStreamer) getLogs(pods []corev1.Pod) []logRequest {
	reqs := make([]logRequest, len(pods))
	for i, pod := range pods {
		req := l.client.Get().
			Name(pod.Name).
			Namespace(undistroNamespace).
			Resource("pods").
			SubResource("log").
			Param("follow", "true").
			Param("timestamps", "true").
			Param("previous", "false").
			Param("container", "manager")
		reqs[i] = logRequest{
			podName: pod.Name,
			req:     req,
		}
	}
	return reqs
}

func NewLogStreamer() *logStreamer {
	return &logStreamer{}
}

type logRequest struct {
	podName string
	req     *rest.Request
}
