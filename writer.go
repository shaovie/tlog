package tlog

type FileStoreModeT int

const (
	DailySplit    FileStoreModeT = 1 // Outputting everything to a single file
	AppendOneFile FileStoreModeT = 2 // Outputting to specific files based on severity level.
)

type Writer interface {
	Write(e Encoder, p []byte) (n int, err error)
}
