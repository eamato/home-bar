package database

import (
	"home-bar/configs"
	"home-bar/internal"
	"os"
	"os/exec"
)

const backupFileName = "backup.sql"

type homeBarDatabaseBackup struct {
	config *configs.Config
}

func NewHomeBarDatabaseBackup(config *configs.Config) Backup {
	return &homeBarDatabaseBackup{
		config: config,
	}
}

func (b *homeBarDatabaseBackup) CreateBackup() error {
	if b.config == nil {
		internal.PrintFatal("config file nil in backup", nil)
	}

	file, err := os.Create(backupFileName)
	if err != nil {
		return err
	}

	username := b.config.DatabaseConfig.User
	password := b.config.DatabaseConfig.Password
	hostname := b.config.DatabaseConfig.Host
	port := b.config.DatabaseConfig.Port
	databaseName := b.config.DatabaseConfig.Name

	cmdArgs := []string{
		"-u",
		username,
		"-p" + password,
		"-h",
		hostname,
		"-P",
		port,
		databaseName,
	}
	name := b.config.DatabaseConfig.DumpToolPath
	if name == "" {
		name = "mysqldump"
	}
	cmd := exec.Command(name, cmdArgs...)
	cmd.Stdout = file

	err = cmd.Run()
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	internal.PrintMessage("Backup created")

	return nil
}

func (b *homeBarDatabaseBackup) DeleteBackup() error {
	if b.config == nil {
		internal.PrintFatal("config file nil in backup", nil)
	}

	err := os.Remove(backupFileName)
	if err != nil {
		return err
	}

	internal.PrintMessage("Backup deleted")

	return nil
}

func (b *homeBarDatabaseBackup) Backup() error {
	if b.config == nil {
		internal.PrintFatal("config file nil in backup", nil)
	}

	file, err := os.Open(backupFileName)
	if err != nil {
		return err
	}

	username := b.config.DatabaseConfig.User
	password := b.config.DatabaseConfig.Password
	hostname := b.config.DatabaseConfig.Host
	port := b.config.DatabaseConfig.Port
	databaseName := b.config.DatabaseConfig.Name

	cmdArgs := []string{
		"-u",
		username,
		"-p" + password,
		"-h",
		hostname,
		"-P",
		port,
		databaseName,
	}
	name := b.config.DatabaseConfig.BackupToolPath
	if name == "" {
		name = "mysql"
	}
	backupCmd := exec.Command(name, cmdArgs...)
	backupCmd.Stdin = file

	err = backupCmd.Run()
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	internal.PrintMessage("Backup done")

	return nil
}

func (b *homeBarDatabaseBackup) BackupAndDelete() error {
	err := b.Backup()
	_ = b.DeleteBackup()
	return err
}
