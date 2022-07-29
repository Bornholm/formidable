package yaml

import (
	"bytes"
	"io"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"forge.cadoles.com/wpetit/formidable/internal/data/format"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const YAMLTagAnsibleVaultValuesQueryParam = "ansible_vault"

type EncoderHandler struct{}

func (d *EncoderHandler) Match(url *url.URL) bool {
	ext := filepath.Ext(path.Join(url.Host, url.Path))

	return ExtensionYAML.MatchString(ext) ||
		format.MatchURLQueryFormat(url, FormatYAML)
}

func (d *EncoderHandler) Encode(url *url.URL, data interface{}) (io.Reader, error) {
	var output bytes.Buffer

	encoder := yaml.NewEncoder(&output)

	if err := encoder.Encode(data); err != nil {
		return nil, errors.WithStack(err)
	}

	if shouldTransformAnsibleVault(url) {
		if err := tagAnsibleVaultValues(&output); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &output, nil
}

func NewEncoderHandler() *EncoderHandler {
	return &EncoderHandler{}
}

func shouldTransformAnsibleVault(url *url.URL) bool {
	return !url.Query().Has(YAMLTagAnsibleVaultValuesQueryParam) || url.Query().Get(YAMLTagAnsibleVaultValuesQueryParam) == "yes"
}

func tagAnsibleVaultValues(buf *bytes.Buffer) error {
	decoder := yaml.NewDecoder(buf)

	var node yaml.Node

	if err := decoder.Decode(&node); err != nil {
		return errors.WithStack(err)
	}

	walkNodeTree(&node, func(node *yaml.Node) {
		isAnsibleVaultNode := node.Tag == "!!str" && strings.HasPrefix(strings.TrimSpace(node.Value), "$ANSIBLE_VAULT")
		if !isAnsibleVaultNode {
			return
		}

		node.Tag = "!vault"
		node.Style = yaml.LiteralStyle | yaml.TaggedStyle
	})

	buf.Reset()

	encoder := yaml.NewEncoder(buf)

	if err := encoder.Encode(node.Content[0]); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func walkNodeTree(node *yaml.Node, fn func(node *yaml.Node)) {
	fn(node)

	if node.Content == nil {
		return
	}

	for _, sub := range node.Content {
		walkNodeTree(sub, fn)
	}
}
