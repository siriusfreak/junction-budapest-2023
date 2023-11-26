package interfaces

type VideoStorage interface {
	SaveFile(data []byte) (string, error)
	GetFile(uid string) ([]byte, error)
}
