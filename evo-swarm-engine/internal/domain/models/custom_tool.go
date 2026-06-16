package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// CustomTool representa a tabela "evo_core_custom_tools" do banco de dados PostgreSQL.
type CustomTool struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name          string         `gorm:"type:varchar;not null" json:"name"`
	Description   *string        `gorm:"type:text" json:"description"`
	Method        string         `gorm:"type:varchar;not null;check:method IN ('GET', 'POST', 'PUT', 'DELETE', 'PATCH')" json:"method"`
	Endpoint      string         `gorm:"type:varchar;not null" json:"endpoint"`
	Headers       datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"headers"`
	PathParams    datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"path_params"`
	QueryParams   datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"query_params"`
	BodyParams    datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"body_params"`
	ErrorHandling datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"error_handling"`
	Values        datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"values"`
	Tags          datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"tags"`
	Examples      datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"examples"`
	InputModes    datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"input_modes"`
	OutputModes   datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"output_modes"`
	CreatedAt     time.Time      `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at"`
}

// TableName define o nome exato da tabela no PostgreSQL.
func (CustomTool) TableName() string {
	return "evo_core_custom_tools"
}
