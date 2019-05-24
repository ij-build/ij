package ssh

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type KeysSuite struct {
	path     string
	agent    agent.Agent
	listener net.Listener
}

func (s *KeysSuite) SetUpTest(t sweet.T) {
	s.path = filepath.Join(
		os.TempDir(),
		fmt.Sprintf("ij-%d.sock", time.Now().UnixNano()%1000),
	)

	os.Setenv("SSH_AUTH_SOCK", s.path)

	s.agent = agent.NewKeyring()
	s.listener, _ = net.Listen("unix", s.path)

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				break
			}

			go agent.ServeAgent(s.agent, conn)
		}
	}()
}

func (s *KeysSuite) TearDownTest(t sweet.T) {
	s.listener.Close()
	os.Remove(s.path)
}

func (s *KeysSuite) TestEnsureKeysAvailableMatching(t sweet.T) {
	key := makeKey()
	pub, _ := ssh.NewPublicKey(key.Public())

	fingerprints := []string{
		ssh.FingerprintSHA256(pub),
		ssh.FingerprintLegacyMD5(pub),
	}

	// Add same key
	s.agent.Add(agent.AddedKey{PrivateKey: key})

	for _, fingerprint := range fingerprints {
		ok, err := EnsureKeysAvailable([]string{fingerprint})
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	}
}

func (s *KeysSuite) TestEnsureKeysAvailableNonMatching(t sweet.T) {
	key := makeKey()
	pub, _ := ssh.NewPublicKey(key.Public())

	fingerprints := []string{
		ssh.FingerprintSHA256(pub),
		ssh.FingerprintLegacyMD5(pub),
	}

	// Add a different key
	s.agent.Add(agent.AddedKey{PrivateKey: makeKey()})

	for _, fingerprint := range fingerprints {
		_, err := EnsureKeysAvailable([]string{fingerprint})
		Expect(err).To(MatchError("available ssh keys do not match expected identities"))
	}
}

func (s *KeysSuite) TestEnsureKeysAvailableWildcard(t sweet.T) {
	s.agent.Add(agent.AddedKey{PrivateKey: makeKey()})

	ok, err := EnsureKeysAvailable([]string{"*"})
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())
}

func (s *KeysSuite) TestEnsureKeysAvailableNoKeys(t sweet.T) {
	_, err := EnsureKeysAvailable([]string{"*"})
	Expect(err).To(MatchError("no ssh keys available"))
}

func (s *KeysSuite) TestEnsureKeysAvailableNoIdentities(t sweet.T) {
	ok, err := EnsureKeysAvailable(nil)
	Expect(err).To(BeNil())
	Expect(ok).To(BeFalse())
}

//
// Helpers

func makeKey() *ecdsa.PrivateKey {
	key, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	return key
}
