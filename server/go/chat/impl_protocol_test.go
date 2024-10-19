package chat

import (
	"bufio"
	"bytes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Protocol", func() {
	Context("ChatHeader", func() {
		It("can encode itself into a binary sequence", func() {
			header := NewChatHeader(5, 2, 1)

			buff := &bytes.Buffer{}
			wr := bufio.NewWriter(buff)
			Expect(header.Write(wr)).To(BeNumerically("==", chatProtocolKeySizeInt+chatProtocolKeySizeUint16+chatProtocolKeySizeUint16))
			Expect(wr.Flush()).To(Succeed())

			Expect(buff.Bytes()).To(Equal([]byte{
				0x00, 0x00, 0x00, 0x005, // length
				0x00, 0x02, // version
				0x00, 0x01, // command
			}))
		})

		It("can decode a binary sequence", func() {
			header := &ChatHeader{}
			byteSequence := []byte{
				0x00, 0x00, 0x00, 0x008, // length
				0x00, 0x02, // version
				0x00, 0x01, // command
			}

			Expect(header.Read(bufio.NewReader(bytes.NewReader(byteSequence)))).To(Succeed())
			Expect(header.Version()).To(BeNumerically("==", 0x0002))
			Expect(header.Key()).To(BeNumerically("==", 0x001))
			Expect(header.Length()).To(BeNumerically("==", 0x0008))

			headerB := &ChatHeader{}
			buff := bytes.NewReader(byteSequence)
			Expect(headerB.Read(bufio.NewReader(buff))).To(Succeed())
			Expect(headerB.Version()).To(BeNumerically("==", 0x0002))
			Expect(headerB.Key()).To(BeNumerically("==", 0x001))
			Expect(headerB.Length()).To(BeNumerically("==", 0x0008))
		})

	})

	Context("CommandLogin", func() {
		It("has the correct attributes", func() {
			login := NewCommandLoginWithCorrelation("user", 1)
			Expect(login.Username()).To(Equal("user"))
		})

		It("can encode itself into a binary sequence", func() {
			login := &CommandLogin{}
			byteSequence := []byte{
				0x00, 0x00, 0x00, 0x01, // uint32 correlation id
				0x00, 0x04, // uint 16 username len
			}
			byteSequence = append(byteSequence, []byte("user")...)
			buff := bytes.NewReader(byteSequence)
			Expect(login.Read(bufio.NewReader(buff))).To(Succeed())
			Expect(login.Username()).To(Equal("user"))
			Expect(login.GetCorrelationId()).To(BeNumerically("==", 1))

			loginB := &CommandLogin{}
			buffB := bytes.NewReader(byteSequence)
			Expect(loginB.Read(bufio.NewReader(buffB))).To(Succeed())
			Expect(loginB.Username()).To(Equal("user"))
			Expect(loginB.GetCorrelationId()).To(BeNumerically("==", 1))

		})

		It("can return the size needed to encode the frame", func() {
			login := NewCommandLoginWithCorrelation("user", 1)
			expectedSize := 2 + 2 + 4 + // key ID + version + correlation ID {header}
				2 + 4 // uint16 for the username string  + uint32 username string length

			Expect(login.SizeNeeded()).To(Equal(expectedSize))

			buff := &bytes.Buffer{}
			wr := bufio.NewWriter(buff)
			Expect(login.Write(wr)).To(BeNumerically("==", login.SizeNeeded()-chatProtocolHeaderSizeBytes))
			Expect(wr.Flush()).To(Succeed())

			Expect(buff.Bytes()).To(Equal([]byte{
				0x00, 0x00, 0x00, 0x01, // uint32 correlation id
				0x00, 0x04, // uint 16 username len
				0x75, 0x73, 0x65, 0x72, // user
			}))

		})
	})

	Context("CommandMessage", func() {
		It("has the correct attributes", func() {
			msg := NewCommandMessageWithCorrelationId("hello", "from", "to", 55)
			Expect(msg.Message).To(Equal("hello"))
			Expect(msg.From).To(Equal("from"))
			Expect(msg.To).To(Equal("to"))
			Expect(msg.CorrelationId()).To(BeNumerically("==", 55))
		})

		It("can encode itself into a binary sequence", func() {

			msg := &CommandMessage{}
			byteSequence := []byte{
				0x00, 0x00, 0x00, 0x01, // uint32 correlation id
				0x00, 0x05, // uint 16 message len
			}
			byteSequence = append(byteSequence, []byte("hello")...)
			byteSequence = append(byteSequence, 0x00, 0x04) // uint 16 to len
			byteSequence = append(byteSequence, []byte("from")...)
			byteSequence = append(byteSequence, 0x00, 0x02) // uint 16 to len
			byteSequence = append(byteSequence, []byte("to")...)

			buff := bytes.NewReader(byteSequence)
			Expect(msg.Read(bufio.NewReader(buff))).To(Succeed())
			Expect(msg.Message).To(Equal("hello"))
			Expect(msg.To).To(Equal("to"))
			Expect(msg.From).To(Equal("from"))
		})

		It("can return the size needed to encode the frame", func() {
			msg := NewCommandMessageWithCorrelationId("hello", "from", "to", 1)
			expectedSize := 2 + 2 + 4 + // key ID + version + correlation ID {header}
				2 + 5 + // uint16 for the message string  + uint32 message string length
				2 + 4 + // from uint16 for the to string  + uint32 to string length
				2 + 2 // to uint16 for the to string  + uint32 to string length

			Expect(msg.SizeNeeded()).To(Equal(expectedSize))

			buff := &bytes.Buffer{}
			wr := bufio.NewWriter(buff)
			Expect(msg.Write(wr)).To(BeNumerically("==", msg.SizeNeeded()-chatProtocolHeaderSizeBytes))
			Expect(wr.Flush()).To(Succeed())

			Expect(buff.Bytes()).To(Equal([]byte{
				0x00, 0x00, 0x00, 0x01, // uint32 correlation id
				0x00, 0x05, // uint 16 message len
				0x68, 0x65, 0x6c, 0x6c, 0x6f, // hello
				0x00, 0x04, // uint 16 from len
				0x66, 0x72, 0x6f, 0x6d, // from
				0x00, 0x02, // uint 16 to len
				0x74, 0x6f, // to
			}))
		})
	})
	Context("Header + Commands", func() {
		It("Header + CommandLogin should encode and decode ", func() {

			login := NewCommandLoginWithCorrelation("user", 1)

			buff := &bytes.Buffer{}
			writer := bufio.NewWriter(buff)

			err := WriteCommandWithHeader(login, writer)
			Expect(err).To(Succeed())

			reader := bufio.NewReader(buff)
			chatHeaderRead := &ChatHeader{}
			loginRead := &CommandLogin{}

			err = chatHeaderRead.Read(reader)
			Expect(err).To(Succeed())
			Expect(chatHeaderRead.Version()).To(BeNumerically("==", 0x0001))
			Expect(chatHeaderRead.Key()).To(BeNumerically("==", 0x01))
			Expect(chatHeaderRead.Length()).To(BeNumerically("==", 14))

			err = loginRead.Read(reader)
			Expect(err).To(Succeed())
			Expect(login.Username()).To(Equal("user"))
			Expect(login.GetCorrelationId()).To(BeNumerically("==", 1))
		})

		It("Header + CommandMessage should encode and decode ", func() {

			msg := NewCommandMessageWithCorrelationId("hello", "user_from", "user_to", 1)

			buff := &bytes.Buffer{}
			writer := bufio.NewWriter(buff)

			err := WriteCommandWithHeader(msg, writer)
			Expect(err).To(Succeed())

			reader := bufio.NewReader(buff)
			chatHeaderRead := &ChatHeader{}
			msgRead := &CommandMessage{}

			err = chatHeaderRead.Read(reader)
			Expect(err).To(Succeed())
			Expect(chatHeaderRead.Version()).To(BeNumerically("==", 0x0001))
			Expect(chatHeaderRead.Key()).To(BeNumerically("==", 0x02))
			Expect(chatHeaderRead.Length()).To(BeNumerically("==", 35))
			err = msgRead.Read(reader)
			Expect(err).To(Succeed())
			Expect(msgRead.Message).To(Equal("hello"))
			Expect(msgRead.From).To(Equal("user_from"))
			Expect(msgRead.To).To(Equal("user_to"))
			Expect(msgRead.CorrelationId()).To(BeNumerically("==", 1))
		})
	})

})
