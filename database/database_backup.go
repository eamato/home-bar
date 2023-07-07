package database

type Backup interface {
	CreateBackup() error
	DeleteBackup() error
	Backup() error
	BackupAndDelete() error
}
