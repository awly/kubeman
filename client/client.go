package client

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/clientcmd"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/kubectl/cmd/config"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

type Client struct {
	c *client.Client
}

func Connect() (*Client, error) {
	opts := config.NewDefaultPathOptions()
	tconf, err := opts.GetStartingConfig()
	if err != nil {
		return nil, err
	}
	conf, err := clientcmd.NewDefaultClientConfig(*tconf, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}

	c, err := client.New(conf)
	if err != nil {
		return nil, err
	}

	return &Client{c: c}, nil
}

func (c *Client) Services() ([]api.Service, error) {
	sl, err := c.c.Services("default").List(labels.Everything())
	if err != nil {
		return nil, err
	}
	return sl.Items, nil
}

func (c *Client) WatchServices() (<-chan watch.Event, error) {
	w, err := c.c.Services("default").Watch(labels.Everything(), fields.Everything(), "")
	if err != nil {
		return nil, err
	}
	return w.ResultChan(), nil
}

func (c *Client) Pods() ([]api.Pod, error) {
	pl, err := c.c.Pods("default").List(labels.Everything(), fields.Everything())
	if err != nil {
		return nil, err
	}
	return pl.Items, nil
}

func (c *Client) StopPod(name string) error {
	return c.c.Pods("default").Delete(name, nil)
}

func (c *Client) WatchPods() (<-chan watch.Event, error) {
	w, err := c.c.Pods("default").Watch(labels.Everything(), fields.Everything(), "")
	if err != nil {
		return nil, err
	}
	return w.ResultChan(), nil
}

func (c *Client) RCs() ([]api.ReplicationController, error) {
	pl, err := c.c.ReplicationControllers("default").List(labels.Everything())
	if err != nil {
		return nil, err
	}
	return pl.Items, nil
}

func (c *Client) WatchRCs() (<-chan watch.Event, error) {
	w, err := c.c.ReplicationControllers("default").Watch(labels.Everything(), fields.Everything(), "")
	if err != nil {
		return nil, err
	}
	return w.ResultChan(), nil
}

func (c *Client) Nodes() ([]api.Node, error) {
	pl, err := c.c.Nodes().List(labels.Everything(), fields.Everything())
	if err != nil {
		return nil, err
	}
	return pl.Items, nil
}

func (c *Client) WatchNodes() (<-chan watch.Event, error) {
	w, err := c.c.Nodes().Watch(labels.Everything(), fields.Everything(), "")
	if err != nil {
		return nil, err
	}
	return w.ResultChan(), nil
}
