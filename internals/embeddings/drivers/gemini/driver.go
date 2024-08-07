package gemini

import (
	"context"
	"fmt"
	"github.com/alash3al/vecdb/internals/embeddings"
	"github.com/alash3al/vecdb/internals/vector"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var _ embeddings.Driver = (*Driver)(nil)

type Driver struct {
	client             *genai.Client
	textEmbeddingModel string
}

func (d *Driver) Open(args map[string]any) error {
	apiKey, ok := args["api_key"].(string)
	if !ok {
		return fmt.Errorf("unable to find gemini `args.api_key`")
	}

	textEmbeddingModel, ok := args["text_embedding_model"].(string)
	if !ok {
		return fmt.Errorf("unable to find gemini `args.text_embedding_model`")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return err
	}

	d.client = client
	d.textEmbeddingModel = textEmbeddingModel

	return nil
}

func (d *Driver) TextEmbedding(ctx context.Context, content string) (vector.Vec, error) {
	res, err := d.client.
		EmbeddingModel(d.textEmbeddingModel).
		EmbedContent(ctx, genai.Text(content))

	if err != nil {
		panic(err)
	}

	var vec vector.Vec

	for _, f := range res.Embedding.Values {
		vec = append(vec, float64(f))
	}

	return vec, nil
}

func (d *Driver) Close() error {
	return d.client.Close()
}
