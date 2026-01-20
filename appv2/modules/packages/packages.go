package packages

import (
	"github.com/qor5/admin/v3/presets"
	"gorm.io/gorm"

	"gohome.4gophers.ru/getapp/gohome/appv2/modules/packages/models"
)

type Packages struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Packages {
	return &Packages{
		db: db,
	}
}

func (m *Packages) Configure(b *presets.Builder) {
	ma := b.Model(&models.Package{}).
		MenuIcon("mdi-account-group").
		// Label("Пакеты").
		RightDrawerWidth("1000")

	ma.Listing("ID", "Title", "Info", "Repo", "Package", "Active")
}

func (m *Packages) Migrate() {
	err := m.db.AutoMigrate(
		&models.Package{},
	)
	if err != nil {
		panic(err)
	}
}
