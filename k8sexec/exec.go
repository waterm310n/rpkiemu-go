package k8sexec

/*
目标：
	创建一个结构体方便得让我操作pod，在里面随意添加容器，以及在容器中执行命令
从kubectl exec模块中删改后使用
*/
import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"time"
)

type ExecOptions struct {
	config        *restclient.Config
	clientset     *kubernetes.Clientset
	Namespace     string
	PodName       string
	ContainerName string
	Command       [3]string
}

func NewExecOptions(namespace, podName, containerName string) (*ExecOptions, error) {
	kubeConfig := "/home/master/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &ExecOptions{
		config:        config,
		clientset:     clientset,
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		Command:       [3]string{"/bin/sh", "-c", ""},
	}, nil
}

func (p *ExecOptions) Exec(cmd string) ([]byte, error) {
	p.Command[2] = cmd
	req := p.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(p.PodName).
		Namespace(p.Namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command:   p.Command[:],
			Container: p.ContainerName,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)
	// 执行命令
	executor, err := remotecommand.NewSPDYExecutor(p.config, "POST", req.URL())
	if err != nil {
		return nil, err
	}
	ctx,cancel := context.WithTimeout(context.Background(),time.Second*2)
	defer cancel()
	var stdout, stderr bytes.Buffer
	if err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	}); err != nil {
		if err.Error() == "context deadline exceeded"{
			slog.Warn(err.Error(),"cmd",cmd)
		}else{
			return stdout.Bytes(), fmt.Errorf(stderr.String())
		}
	}
	return stdout.Bytes(),nil
}

// 文件上传到pod中，要求容器中有tar命令。
func (p *ExecOptions) Upload(srcFile string, dstFile string) error {
	/*
		使用tar进行文件或者文件夹复制
		tar cf - <文件> | tar xf - -C <目的地址>/
	*/
	src := fileSpec{
		File: newLocalPath(srcFile),
	}
	dest := fileSpec{
		PodName:      p.PodName,
		PodNamespace: p.Namespace,
		File:         newRemotePath(dstFile),
	}
	o := NewCopyOptions(genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stdout})
	if err := o.copyToPod(src, dest,p); err != nil {
		return err
	}
	return nil
}

func (p *ExecOptions) Download(srcFile string, dstFile string) error {
	slog.Info("Download func run")
	// TODO
	return nil
}

// 获取当前namespace下pod的容器的日志，limitLine限制日志行数数量获取
func (p *ExecOptions) GetLog(limitLine int) ([]string, error) {
	req := p.clientset.CoreV1().Pods(p.Namespace).GetLogs(p.PodName, &v1.PodLogOptions{Container: p.ContainerName})
	readCloser, err := req.Stream(context.TODO())
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()
	logs := make([]string, limitLine)
	scanner := bufio.NewScanner(readCloser)
	for i := 0; i < limitLine && scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())
		if len(line) != 0 {
			logs[i] = line
		}
	}
	return logs, nil
}
