package conf

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/config"

	"github.com/aide-family/rabbit/pkg/merr"
)

var (
	_ config.Source  = (*bytesSource)(nil)
	_ config.Watcher = (*noOpWatcher)(nil)
)

func NewBytesSource(data []byte) config.Source {
	d := bytesSource(data)
	return &d
}

type bytesSource []byte

// Load implements config.Source.
func (b *bytesSource) Load() ([]*config.KeyValue, error) {
	// Make a copy of the data to avoid external modifications
	data := make([]byte, len(*b))
	copy(data, *b)
	return []*config.KeyValue{
		{
			Key:    "server",
			Value:  data,
			Format: format(*b),
		},
	}, nil
}

// format detects the format from the data content.
func format(data []byte) string {
	content := strings.TrimSpace(string(data))
	if strings.HasPrefix(content, "{") || strings.HasPrefix(content, "[") {
		return "json"
	}
	return "yaml"
}

// Watch implements config.Source.
func (b *bytesSource) Watch() (config.Watcher, error) {
	return newNoOpWatcher(), nil
}

type noOpWatcher struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func newNoOpWatcher() config.Watcher {
	ctx, cancel := context.WithCancel(context.Background())
	return &noOpWatcher{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Next implements config.Watcher.
func (w *noOpWatcher) Next() ([]*config.KeyValue, error) {
	<-w.ctx.Done()
	return nil, w.ctx.Err()
}

// Stop implements config.Watcher.
func (w *noOpWatcher) Stop() error {
	w.cancel()
	return nil
}

func Load(bc any, sources ...config.Source) error {
	c := config.New(config.WithSource(sources...))
	if err := c.Load(); err != nil {
		return merr.ErrorInternal("load config failed").WithCause(err)
	}
	if err := c.Scan(bc); err != nil {
		return merr.ErrorInternal("scan config failed").WithCause(err)
	}
	return nil
}
