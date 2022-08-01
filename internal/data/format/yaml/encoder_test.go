package yaml

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

// YAML string containing an ansible-vault encrypted variable
const ansibleVaultYAML = `
unencrypted: foo
encrypted: !vault |
  $ANSIBLE_VAULT;1.1;AES256
  63393636613562663937383964323839376239663230366130386566393131313963386265303632
  3133356532346437653338343032303732646530303431660a383862353766326334306138613734
  36313438626564623435373365616531353533663765663335616134656430323134323537336661
  3437653863343331370a393136653735643333373962633631663539653664313936303964303866
  3933
`

func TestEncoderAnsibleVault(t *testing.T) {
	_, err := exec.LookPath("ansible")
	if err != nil {
		t.Skip("The 'ansible' command seems not to be available on this system. Skipping.")

		return
	}

	var data interface{}

	if err := yaml.Unmarshal([]byte(ansibleVaultYAML), &data); err != nil {
		t.Fatal(errors.WithStack(err))
	}

	encoder := NewEncoderHandler()

	url, err := url.Parse("stdout://local.yml?ansible_vault=yes")
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	reader, err := encoder.Encode(url, data)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	temp, err := os.CreateTemp(os.TempDir(), "formidable_test_*.yml")
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	defer func() {
		if err := os.Remove(temp.Name()); err != nil {
			panic(errors.WithStack(err))
		}
	}()

	t.Logf("Writing encoded YAML content in file '%s'...", temp.Name())

	if _, err := io.Copy(temp, reader); err != nil {
		t.Fatal(errors.WithStack(err))
	}

	args := []string{
		"localhost",
		"-m", "debug",
		"--vault-password-file", "./testdata/vault.txt",
		"-e", fmt.Sprintf("@%s", temp.Name()),
		"-a", "var=encrypted",
	}

	t.Logf("Running command 'ansible %s'", strings.Join(args, " "))

	cmd := exec.Command("ansible", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatal(errors.WithStack(err))
	}
}
