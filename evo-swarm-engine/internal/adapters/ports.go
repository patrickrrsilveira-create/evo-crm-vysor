package adapters

import (
	"context"
)

// ChannelAdapter define a interface padrão que todos os canais de comunicação devem implementar
type ChannelAdapter interface {
	// Name retorna o nome do adaptador (ex: "whatsapp", "teams", "email")
	Name() string

	// Start inicializa as rotinas de escuta do adaptador (Webhooks ou Long Polling)
	Start(ctx context.Context) error

	// SendMessage envia uma mensagem padronizada para o destino através do canal específico
	SendMessage(ctx context.Context, to string, content string) error
}

// DriveAdapter define a interface para abstrair armazenamentos de arquivos na nuvem (Google Drive, OneDrive)
type DriveAdapter interface {
	Name() string
	UploadFile(ctx context.Context, folderID string, fileName string, fileData []byte) (string, error)
	DownloadFile(ctx context.Context, fileID string) ([]byte, error)
}
