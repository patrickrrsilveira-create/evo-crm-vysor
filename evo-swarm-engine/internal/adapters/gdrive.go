package adapters

import (
	"context"
	"log"
)

// GoogleDriveAdapter implementa DriveAdapter para manipulação de arquivos
type GoogleDriveAdapter struct {
	ServiceAccountJSON string
}

func NewGoogleDriveAdapter(serviceAccount string) *GoogleDriveAdapter {
	return &GoogleDriveAdapter{
		ServiceAccountJSON: serviceAccount,
	}
}

func (a *GoogleDriveAdapter) Name() string {
	return "google_drive"
}

func (a *GoogleDriveAdapter) UploadFile(ctx context.Context, folderID string, fileName string, fileData []byte) (string, error) {
	log.Printf("📁 [GoogleDriveAdapter] Fazendo upload de '%s' (%d bytes) para a pasta %s", fileName, len(fileData), folderID)
	// Chamada real da API de Google Drive
	return "mock_file_id_12345", nil
}

func (a *GoogleDriveAdapter) DownloadFile(ctx context.Context, fileID string) ([]byte, error) {
	log.Printf("📁 [GoogleDriveAdapter] Baixando arquivo %s", fileID)
	// Download binário via API do Google Drive
	return []byte("mock data"), nil
}
