package ssh

import (
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func EnureKeysAvailable(identities []string) (bool, error) {
	if len(identities) == 0 {
		return false, nil
	}

	rw, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return false, err
	}

	keys, err := agent.NewClient(rw).List()
	if err != nil {
		return false, err
	}

	for _, identity := range identities {
		if matchIdentity(identity, keys) {
			return true, nil
		}
	}

	return false, fmt.Errorf("no ssh keys available")
}

func matchIdentity(identity string, keys []*agent.Key) bool {
	if identity == "*" && len(keys) > 0 {
		return true
	}

	for _, key := range keys {
		var (
			fpSHA256     = ssh.FingerprintSHA256(key)
			fpMD5        = ssh.FingerprintLegacyMD5(key)
			fingerprints = []string{
				fpSHA256,
				fpMD5,
				fmt.Sprintf("sha256:%s", fpSHA256),
				fmt.Sprintf("md5:%s", fpMD5),
			}
		)

		for _, fp := range fingerprints {
			if strings.ToUpper(identity) == strings.ToUpper(fp) {
				return true
			}
		}
	}

	return false
}
