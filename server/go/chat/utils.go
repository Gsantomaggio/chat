package chat

import (
	"bufio"
	"bytes"
	"io"
	"time"
)

func FormResponseCodeToString(responseCode uint16) string {

	fromCodeToString := "Success"
	switch responseCode {
	case ResponseCodeErrorUserAlreadyLogged:
		fromCodeToString = "ErrorUserAlreadyLogged"
	case ResponseCodeErrorUserNotFound:
		fromCodeToString = "ErrorUserNotFound"
	}
	return fromCodeToString
}

// TODO: Explain the REST problem and how this function solves it

func ReadFullBufferFromSource(sourceStream io.Reader) (*bufio.Reader, error) {
	dataLength, _ := readUInt(sourceStream)
	var bytesBuffer = make([]byte, int(dataLength))
	_, err := io.ReadFull(sourceStream, bytesBuffer)
	bufferReader := bytes.NewReader(bytesBuffer)
	dataReader := bufio.NewReader(bufferReader)
	return dataReader, err
}

func ConvertTimeToUint64(t time.Time) uint64 {
	return uint64(t.UnixNano())
}

// ConvertUint64ToTime converts a Unix timestamp in nanoseconds (uint64) to time.Time
func ConvertUint64ToTime(nanoseconds uint64) time.Time {
	return time.Unix(0, int64(nanoseconds))
}

// Convert and format

func ConvertUint64ToTimeFormatted(nanoseconds uint64) string {
	return ConvertUint64ToTime(nanoseconds).Format(time.RFC822)
}
