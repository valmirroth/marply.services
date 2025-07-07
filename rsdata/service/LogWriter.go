package service

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger é um logger customizado que escreve logs em arquivo e console.
type Logger struct {
	mu       sync.Mutex
	fileDir  string
	basename string
	level    LogLevel
	writer   io.Writer
	logger   *log.Logger
}

// LogLevel define o nível do log.
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// NewLogger cria um novo Logger.
// fileDir: diretório onde os arquivos de log serão salvos.
// basename: nome base do arquivo de log (sem extensão).
func NewLogger(fileDir, basename string, level LogLevel) (*Logger, error) {
	// Garantir que o diretório exista
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório de log: %w", err)
	}

	l := &Logger{
		fileDir:  fileDir,
		basename: basename,
		level:    level,
	}
	// Inicializa writer e logger
	if err := l.rotateWriter(); err != nil {
		return nil, err
	}
	return l, nil
}

// rotateWriter abre o arquivo de log do dia atual e configura o escritor.
func (l *Logger) rotateWriter() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Formatar nome de arquivo com data
	now := time.Now().Format("2006-01-02")
	filename := filepath.Join(l.fileDir, fmt.Sprintf("%s_%s.log", l.basename, now))

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("falha ao abrir arquivo de log: %w", err)
	}

	// Escritor multiplo: arquivo + console
	l.writer = io.MultiWriter(f, os.Stdout)
	l.logger = log.New(l.writer, "", log.LstdFlags|log.Lmicroseconds)
	return nil
}

// logf escreve a mensagem formatada se o nível for adequado.
func (l *Logger) logf(level LogLevel, prefix string, format string, v ...interface{}) {
	if level < l.level {
		return
	}
	// Rotaciona a cada dia
	if time.Now().Day() != time.Now().Day() {
		l.rotateWriter()
	}
	l.logger.Printf(prefix+" "+format, v...)
}

// Debug registra mensagem de debug.
func (l *Logger) Debug(format string, v ...interface{}) {
	l.logf(DEBUG, "DEBUG:", format, v...)
}

// Info registra mensagem de informação.
func (l *Logger) Info(format string, v ...interface{}) {
	l.logf(INFO, "INFO:", format, v...)
}

// Warn registra mensagem de aviso.
func (l *Logger) Warn(format string, v ...interface{}) {
	l.logf(WARN, "WARN:", format, v...)
}

// Error registra mensagem de erro.
func (l *Logger) Error(format string, v ...interface{}) {
	l.logf(ERROR, "ERROR:", format, v...)
}
