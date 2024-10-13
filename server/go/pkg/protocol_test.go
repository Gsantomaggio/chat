package pkg

import (
	"bufio"
	"bytes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Protocol", func() {
	Context("CommandLogin", func() {
		It("has the correct attributes", func() {
			login := NewCommandLogin("user", 1)
			Expect(login.Username()).To(Equal("user"))
		})

		It("can encode itself into a binary sequence", func() {

			login := &CommandLogin{}
			byteSequence := []byte{
				0x00, 0x00, 0x00, 0x01, // uint32 correlation id
				0x00, 0x04, // uint 16 username len
			}
			byteSequence = append(byteSequence, []byte("user")...)

			Expect(login.UnmarshalBinary(byteSequence)).To(Succeed())
			Expect(login.Username()).To(Equal("user"))
			Expect(login.GetCorrelationId()).To(BeNumerically("==", 1))
		})

		It("can encode itself into a binary sequence", func() {
			login := NewCommandLogin("user", 1)
			byteSequence := []byte{
				0x00, 0x00, 0x00, 0x01, // uint32 correlation id
				0x00, 0x04, // uint 16 username len
			}
			byteSequence = append(byteSequence, []byte("user")...)
			Expect(login.SizeNeeded()).To(Equal(10 + chatProtocolHeaderSize))
		})

		It("can return the size needed to encode the frame", func() {
			login := NewCommandLogin("user", 1)
			Expect(login.SizeNeeded()).To(Equal(10 + chatProtocolHeaderSize))

			buff := &bytes.Buffer{}
			wr := bufio.NewWriter(buff)
			Expect(login.Write(wr)).To(BeNumerically("==", login.SizeNeeded()-chatProtocolHeaderSize))
			Expect(wr.Flush()).To(Succeed())

			Expect(buff.Bytes()).To(Equal([]byte{
				0x00, 0x00, 0x00, 0x01, // uint32 correlation id
				0x00, 0x04, // uint 16 username len
				0x75, 0x73, 0x65, 0x72, // user
			}))

		})

	})
})
