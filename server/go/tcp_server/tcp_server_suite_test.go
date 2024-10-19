package tcp_server_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTcpServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TcpServer Suite")
}
