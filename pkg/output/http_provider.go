package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/ncarlier/feedpushr/autogen/app"
	"github.com/ncarlier/feedpushr/pkg/model"
)

// HTTPOutputProvider HTTP output provider
type HTTPOutputProvider struct {
	name      string
	desc      string
	tags      []string
	nbError   uint64
	nbSuccess uint64
	uri       string
}

func newHTTPOutputProvider(output *app.Output) *HTTPOutputProvider {
	u, ok := output.Props["url"]
	if !ok {
		return nil
	}
	return &HTTPOutputProvider{
		name: "http",
		desc: "New articles are sent as JSON document to an HTTP endpoint (POST).\n\n" + jsonFormatDesc,
		tags: output.Tags,
		uri:  fmt.Sprintf("%v", u),
	}
}

// Send article to HTTP endpoint.
func (op *HTTPOutputProvider) Send(article *model.Article) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(article)
	resp, err := http.Post(op.uri, "application/json; charset=utf-8", b)
	if err != nil {
		atomic.AddUint64(&op.nbError, 1)
		return err
	} else if resp.StatusCode >= 300 {
		atomic.AddUint64(&op.nbError, 1)
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	atomic.AddUint64(&op.nbSuccess, 1)
	return nil
}

// GetSpec return output provider specifications
func (op *HTTPOutputProvider) GetSpec() model.OutputSpec {
	result := model.OutputSpec{
		Name: op.name,
		Desc: op.desc,
		Tags: op.tags,
	}
	result.Props = map[string]interface{}{
		"uri":       op.uri,
		"nbError":   op.nbError,
		"nbSuccess": op.nbSuccess,
	}
	return result
}
