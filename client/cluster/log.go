/*
Copyright 2020 Getup Cloud. All rights reserved.
*/

package cluster

import (
	"context"
	"io"
	"sync"
	"time"

	undistrov1 "github.com/getupcloud/undistro/api/v1alpha1"
	"github.com/getupcloud/undistro/internal/scheme"
	"github.com/getupcloud/undistro/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	undistroNamespace = "undistro-system"
)

type LogStreamer interface {
	Stream(context.Context, *rest.Config, io.Writer, types.NamespacedName) error
}

type logStreamer struct {
	client     *kubernetes.Clientset
	kubeClient client.Client
}

func (l *logStreamer) Stream(ctx context.Context, c *rest.Config, w io.Writer, nm types.NamespacedName) error {
	var (
		err    error
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(ctx)
	l.kubeClient, err = client.New(c, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		cancel()
		return err
	}
	l.client, err = kubernetes.NewForConfig(c)
	if err != nil {
		cancel()
		return err
	}
	var podList corev1.PodList
	err = l.kubeClient.List(ctx, &podList, client.InNamespace(undistroNamespace))
	if err != nil {
		cancel()
		return err
	}
	logs := l.getLogs(podList.Items)
	cerr := make(chan error, 1)
	go func(ctx context.Context, cancel context.CancelFunc) {
		ready, err := l.clusterIsReady(ctx, nm)
		if err != nil {
			cerr <- err
			return
		}
		for !ready {
			ready, err = l.clusterIsReady(ctx, nm)
			if err != nil {
				cerr <- err
				return
			}
			<-time.After(10 * time.Second)
		}
		cancel()
	}(ctx, cancel)
	go func(ctx context.Context) {
		cerr <- l.streamLogs(ctx, logs, w)
	}(ctx)
	select {
	case err = <-cerr:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (l *logStreamer) streamLogs(ctx context.Context, reqs []logRequest, writer io.Writer) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(reqs))
	for _, req := range reqs {
		reader, err := req.req.Stream()
		if err != nil {
			return err
		}
		go l.readLog(ctx, reader, writer)
	}
	wg.Wait()
	return nil
}

func (l *logStreamer) readLog(ctx context.Context, reader io.ReadCloser, writer io.Writer) {
	_, err := io.Copy(writer, reader)
	if err != nil && err != context.DeadlineExceeded {
		log.Log.Error(err, "fail to copy")
		return
	}
	defer reader.Close()
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
		req := l.client.
			CoreV1().
			Pods(undistroNamespace).
			GetLogs(pod.Name, &corev1.PodLogOptions{
				Container: "manager",
				Follow:    true,
			})
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
