package k8sexec

/*
目标：
	创建一个结构体方便得让我操作pod，在里面随意添加容器，以及在容器中执行命令
从kubectl exec模块中删改后使用
*/
import (
	"bytes"
	"context"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type ExecOptions struct {
	StreamOptions
	Command   [3]string
}

type StreamOptions struct {
	Namespace     string
	PodName       string
	ContainerName string
}
func NewExecOptions(namespace,podName,containerName string) *ExecOptions{
	return &ExecOptions{
		StreamOptions: StreamOptions{
			namespace,
			podName,
			containerName,
		},
		Command: [3]string{"/bin/sh","-c",""},
	}
}

func Helper(){
	kubeConfig := viper.GetString("kubeConfig")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil{
		slog.Error(err.Error())
	}
	ret, err := cmdExecuter(config,"r5","bgp","r5-frr","touch ~/.hello")
	if err != nil{
		slog.Error(err.Error())
	}
	for k,v:= range ret{
		slog.Error(k,"msg",v)
	}
}

func cmdExecuter(config *restclient.Config, podName, namespace, containerName,cmd string) (map[string]string, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
    }
    // 构造执行命令请求
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: []string{"/bin/sh", "-c", cmd},
			Container: containerName,
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
        }, scheme.ParameterCodec)
    // 执行命令
	executor, err := remotecommand.NewSPDYExecutor(config, "POST",req.URL())
	if err != nil {
		return nil, err
    }
    // 使用bytes.Buffer变量接收标准输出和标准错误
	var stdout, stderr bytes.Buffer
	if err = executor.StreamWithContext(context.TODO(),remotecommand.StreamOptions{
		Stdin: strings.NewReader(""),
		Stdout: &stdout,
		Stderr: &stderr,
	}); err != nil {
		return nil, err
    }
    // 返回数据
	ret := map[string]string{"stdout":stdout.String(), "stderr":stderr.String(), "pod_name": podName}
	return ret, nil
}



// // RemoteExecutor defines the interface accepted by the Exec command - provided for test stubbing
// type RemoteExecutor interface {
// 	Execute(method string, url *url.URL, config *restclient.Config, stdin io.Reader, stdout, stderr io.Writer) error
// }

// // DefaultRemoteExecutor is the standard implementation of remote command execution
// type DefaultRemoteExecutor struct{}

// func (*DefaultRemoteExecutor) Execute(method string, url *url.URL, config *restclient.Config, stdin io.Reader, stdout, stderr io.Writer) error {
// 	exec, err := remotecommand.NewSPDYExecutor(config, method, url)
// 	if err != nil {
// 		return err
// 	}
// 	return exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
// 		Stdin:  stdin,
// 		Stdout: stdout,
// 		Stderr: stderr,
// 		Tty:    false,
// 	})
// }

// func NewExecOptions(namespace, podName, containerName string) *ExecOptions {
// 	kubeConfig := viper.GetString("kubeConfig")
// 	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		os.Exit(1)
// 	}
// 	options := &ExecOptions{
// 		StreamOptions: StreamOptions{
// 			Namespace:     namespace,
// 			PodName:       podName,
// 			ContainerName: containerName,
// 		},
// 		Command: []string{"/bin/sh", "-c", "ip addr"},
// 		PodClient: coreclient.NewForConfigOrDie(config),
// 		Config: config,
// 		Executor: &DefaultRemoteExecutor{},
// 		In:       os.Stdin,
// 		Out:      os.Stdout,
// 		ErrOut:   os.Stderr,
// 	}
// 	return options
// }

// func (p *ExecOptions) Run() error {
// 	var err error
// 	if len(p.PodName) != 0 {
// 		p.Pod, err = p.PodClient.Pods(p.Namespace).Get(context.TODO(), p.PodName, metav1.GetOptions{})
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	pod := p.Pod
// 	if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
// 		slog.Error(fmt.Sprintf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase))
// 		return fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase)
// 	}

// 	containerName := p.ContainerName
// 	if len(containerName) == 0 {
// 		slog.Error("len(containerName) == 0")
// 		return fmt.Errorf("len(containerName) == 0")
// 	}

// 	fn := func() error {
// 		restClient, err := restclient.RESTClientFor(p.Config)
// 		if err != nil {
// 			return err
// 		}
// 		req := restClient.Post().
// 			Resource("pods").
// 			Name(pod.Name).
// 			Namespace(pod.Namespace).
// 			SubResource("exec")
// 		req.VersionedParams(&corev1.PodExecOptions{
// 			Container: containerName,
// 			Command:   p.Command,
// 		}, scheme.ParameterCodec)

// 		return p.Executor.Execute("POST", req.URL(), p.Config, p.In, p.Out, p.ErrOut)
// 	}
// 	if err := fn(); err != nil {
// 		slog.Error(err.Error())
// 		return err
// 	}
// 	return err
// }
