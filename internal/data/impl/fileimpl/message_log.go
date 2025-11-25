package fileimpl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/encoding"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/merr"
	"github.com/aide-family/rabbit/pkg/middler"
)

const (
	logFilePrefix = "message_"
	logFileSuffix = ".log"
	dateFormat    = "20060102"
)

func NewMessageLogRepository(d *data.Data, helper *klog.Helper) repository.MessageLog {
	repo := &messageLogRepositoryImpl{
		helper:        klog.NewHelper(klog.With(helper.Logger(), "data", "fileimpl.messageLogRepository")),
		d:             d,
		cache:         safety.NewMap(make(map[string]struct{})),
		uidToLocation: safety.NewSyncMap(make(map[snowflake.ID]*fileLocation)),
		fileMutex:     &sync.Mutex{},
		lastIDByDate:  safety.NewSyncMap(make(map[string]uint32)),
		codec:         encoding.GetCodec("json"),
	}

	// 获取当前工作目录
	baseDir, err := os.Getwd()
	if err != nil {
		repo.helper.Errorf("failed to get current directory: %v", err)
		baseDir = "."
	}
	repo.baseDir = baseDir

	// 启动时加载数据
	repo.loadMessageLogs()

	return repo
}

// fileLocation 记录消息日志在文件中的位置
type fileLocation struct {
	filePath string
	lineNum  int // 行号从1开始
}

type messageLogRepositoryImpl struct {
	helper        *klog.Helper
	d             *data.Data
	cache         *safety.Map[string, struct{}]
	uidToLocation *safety.SyncMap[snowflake.ID, *fileLocation] // UID 到文件位置的映射
	fileMutex     *sync.Mutex
	baseDir       string
	lastIDByDate  *safety.SyncMap[string, uint32] // 按日期存储的lastID，格式：YYYYMMDD -> lastID
	codec         encoding.Codec
}

// loadMessageLogs 加载消息日志数据
// 按时间倒序加载所有 message_YYYYMMDD.log 格式的文件
func (m *messageLogRepositoryImpl) loadMessageLogs() {
	// 查找所有日志文件并按时间倒序加载
	files, err := m.findHistoryFiles()
	if err != nil {
		m.helper.Warnf("failed to find log files: %v", err)
		return
	}

	// 按时间倒序排序（最新的在前）
	sort.Slice(files, func(i, j int) bool {
		return files[i].date.After(files[j].date)
	})

	// 加载所有文件，只建立 UID 到文件位置的映射
	for _, file := range files {
		if err := m.loadFile(file.path); err != nil {
			m.helper.Warnf("failed to load log file %s: %v", file.path, err)
		}
	}

	m.helper.Infof("loaded %d message log locations from files", m.uidToLocation.Len())
}

type historyFile struct {
	path string
	date time.Time
}

// findHistoryFiles 查找所有日志文件（message_YYYYMMDD.log 格式）
func (m *messageLogRepositoryImpl) findHistoryFiles() ([]historyFile, error) {
	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		return nil, err
	}

	var files []historyFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// 匹配 message_YYYYMMDD.log 格式
		if strings.HasPrefix(name, logFilePrefix) && strings.HasSuffix(name, logFileSuffix) {
			// 提取日期部分
			dateStr := strings.TrimPrefix(name, logFilePrefix)
			dateStr = strings.TrimSuffix(dateStr, logFileSuffix)

			date, err := time.Parse(dateFormat, dateStr)
			if err != nil {
				m.helper.Warnf("invalid date format in filename %s: %v", name, err)
				continue
			}

			files = append(files, historyFile{
				path: filepath.Join(m.baseDir, name),
				date: date,
			})
		}
	}

	return files, nil
}

// loadFile 从文件加载消息日志的 UID 到文件位置的映射，不加载完整数据
func (m *messageLogRepositoryImpl) loadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在不算错误
		}
		return err
	}
	defer file.Close()

	// 从文件路径提取日期（message_YYYYMMDD.log）
	fileName := filepath.Base(filePath)
	dateStr := strings.TrimPrefix(fileName, logFilePrefix)
	dateStr = strings.TrimSuffix(dateStr, logFileSuffix)

	var maxID uint32
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// 只解析 UID 和 ID，不加载完整数据
		var msgLog do.MessageLog
		if err := m.codec.Unmarshal([]byte(line), &msgLog); err != nil {
			m.helper.Warnf("failed to unmarshal line %d in %s: %v", lineNum, filePath, err)
			continue
		}

		// 更新该日期的最大ID
		if msgLog.ID > maxID {
			maxID = msgLog.ID
		}

		// 只建立 UID 到文件位置的映射（如果已存在则更新，保留最新的文件位置）
		if existing, ok := m.uidToLocation.Get(msgLog.UID); ok {
			// 如果新文件的时间更晚，则更新位置映射
			// 通过比较文件路径中的日期来判断
			existingDateStr := strings.TrimPrefix(filepath.Base(existing.filePath), logFilePrefix)
			existingDateStr = strings.TrimSuffix(existingDateStr, logFileSuffix)
			if dateStr > existingDateStr {
				m.uidToLocation.Set(msgLog.UID, &fileLocation{
					filePath: filePath,
					lineNum:  lineNum,
				})
			}
		} else {
			// 建立 UID 到文件位置的映射
			m.uidToLocation.Set(msgLog.UID, &fileLocation{
				filePath: filePath,
				lineNum:  lineNum,
			})
		}
	}

	// 更新该日期的lastID（如果文件中有数据）
	if maxID > 0 {
		// 如果已经存在该日期的lastID，取最大值
		if existingLastID, ok := m.lastIDByDate.Get(dateStr); ok {
			if maxID > existingLastID {
				m.lastIDByDate.Set(dateStr, maxID)
			}
		} else {
			m.lastIDByDate.Set(dateStr, maxID)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return nil
}

// getDateLogFile 根据日期获取日志文件路径（message_YYYYMMDD.log 格式）
func (m *messageLogRepositoryImpl) getDateLogFile(date time.Time) string {
	dateStr := date.Format(dateFormat)
	filename := logFilePrefix + dateStr + logFileSuffix
	return filepath.Join(m.baseDir, filename)
}

// writeMessageLog 写入消息日志到文件（追加模式，如果文件不存在则创建）
// 返回写入的行号
func (m *messageLogRepositoryImpl) writeMessageLog(msgLog *do.MessageLog) (int, error) {
	m.fileMutex.Lock()
	defer m.fileMutex.Unlock()

	// 根据 SendAt 确定文件路径（统一使用 message_YYYYMMDD.log 格式）
	filePath := m.getDateLogFile(msgLog.SendAt)

	// 计算当前文件的行数（用于确定新行的行号）
	lineCount, err := m.countLines(filePath)
	if err != nil && !os.IsNotExist(err) {
		return 0, fmt.Errorf("failed to count lines in file %s: %w", filePath, err)
	}

	// 打开文件（追加模式，如果不存在则创建）
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return 0, fmt.Errorf("failed to open log file %s: %w", filePath, err)
	}
	defer file.Close()

	// 序列化为 JSON
	dataBytes, err := m.codec.Marshal(msgLog)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal message log: %w", err)
	}

	// 写入文件（一行一条）
	if _, err := file.Write(append(dataBytes, '\n')); err != nil {
		return 0, fmt.Errorf("failed to write to log file %s: %w", filePath, err)
	}

	// 新行的行号 = 原行数 + 1
	newLineNum := lineCount + 1

	// 更新位置映射
	m.uidToLocation.Set(msgLog.UID, &fileLocation{
		filePath: filePath,
		lineNum:  newLineNum,
	})

	return newLineNum, nil
}

// countLines 统计文件的行数（不包括空行）
func (m *messageLogRepositoryImpl) countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

// readMessageLogFromFile 根据文件位置从文件中读取消息日志
func (m *messageLogRepositoryImpl) readMessageLogFromFile(location *fileLocation) (*do.MessageLog, error) {
	file, err := os.Open(location.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, merr.ErrorNotFound("message log file not found: %s", location.filePath)
		}
		return nil, fmt.Errorf("failed to open log file %s: %w", location.filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	for scanner.Scan() {
		currentLine++
		if currentLine == location.lineNum {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				return nil, merr.ErrorNotFound("message log at line %d is empty", location.lineNum)
			}

			var msgLog do.MessageLog
			if err := m.codec.Unmarshal([]byte(line), &msgLog); err != nil {
				return nil, fmt.Errorf("failed to unmarshal message log at line %d: %w", location.lineNum, err)
			}
			return &msgLog, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", location.filePath, err)
	}

	return nil, merr.ErrorNotFound("message log at line %d not found in file %s", location.lineNum, location.filePath)
}

// getNextID 获取指定日期的下一个ID
func (m *messageLogRepositoryImpl) getNextID(date time.Time) uint32 {
	m.fileMutex.Lock()
	defer m.fileMutex.Unlock()

	dateStr := date.Format(dateFormat)
	lastID, ok := m.lastIDByDate.Get(dateStr)
	if !ok {
		// 第一天ID从1开始
		lastID = 0
	}

	// ID自增
	newID := lastID + 1
	m.lastIDByDate.Set(dateStr, newID)

	return newID
}

// CreateMessageLog implements repository.MessageLog.
func (m *messageLogRepositoryImpl) CreateMessageLog(ctx context.Context, messageLog *do.MessageLog) error {
	node, err := snowflake.NewNode(hello.NodeID())
	if err != nil {
		return err
	}
	messageLog.CreatedAt = time.Now()
	messageLog.UpdatedAt = messageLog.CreatedAt
	messageLog.WithCreator(ctx)
	messageLog.WithUID(node.Generate())
	if strutil.IsEmpty(messageLog.Namespace) {
		messageLog.WithNamespace(middler.GetNamespace(ctx))
	}

	// 根据SendAt确定日期，生成当天的ID
	if messageLog.SendAt.IsZero() {
		messageLog.SendAt = messageLog.CreatedAt
	}
	messageLog.ID = m.getNextID(messageLog.SendAt)

	// 写入文件（会自动更新位置映射）
	if _, err = m.writeMessageLog(messageLog); err != nil {
		return err
	}

	return nil
}

// GetMessageLog implements repository.MessageLog.
func (m *messageLogRepositoryImpl) GetMessageLog(ctx context.Context, uid snowflake.ID) (*do.MessageLog, error) {
	location, ok := m.uidToLocation.Get(uid)
	if !ok {
		return nil, merr.ErrorNotFound("message log %d not found", uid.Int64())
	}

	// 从文件读取数据
	return m.readMessageLogFromFile(location)
}

// GetMessageLogWithLock implements repository.MessageLog.
// 文件实现不需要真正的锁，直接返回结果
func (m *messageLogRepositoryImpl) GetMessageLogWithLock(ctx context.Context, uid snowflake.ID) (*do.MessageLog, error) {
	return m.GetMessageLog(ctx, uid)
}

// ListMessageLog implements repository.MessageLog.
func (m *messageLogRepositoryImpl) ListMessageLog(ctx context.Context, req *bo.ListMessageLogBo) (*bo.PageResponseBo[*do.MessageLog], error) {
	// 查找所有日志文件
	files, err := m.findHistoryFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to find log files: %w", err)
	}

	// 按时间倒序排序（最新的在前）
	sort.Slice(files, func(i, j int) bool {
		return files[i].date.After(files[j].date)
	})

	// 收集所有符合条件的消息日志
	allLogs := make([]*do.MessageLog, 0)

	// 遍历所有文件
	for _, file := range files {
		// 如果指定了时间范围，可以提前过滤文件（按日期比较）
		fileDateStr := file.date.Format(dateFormat)
		if !req.StartAt.IsZero() {
			startDateStr := req.StartAt.Format(dateFormat)
			if fileDateStr < startDateStr {
				continue
			}
		}
		if !req.EndAt.IsZero() {
			endDateStr := req.EndAt.Format(dateFormat)
			if fileDateStr > endDateStr {
				continue
			}
		}

		// 读取文件中的所有消息日志
		fileLogs, err := m.readAllLogsFromFile(file.path, req)
		if err != nil {
			m.helper.Warnf("failed to read logs from file %s: %v", file.path, err)
			continue
		}
		allLogs = append(allLogs, fileLogs...)
	}

	// 按创建时间倒序排序
	sort.Slice(allLogs, func(i, j int) bool {
		return allLogs[i].CreatedAt.After(allLogs[j].CreatedAt)
	})

	total := int64(len(allLogs))

	// 分页
	if req.PageRequestBo != nil {
		req.WithTotal(total)
		offset := req.Offset()
		limit := req.Limit()

		if offset >= len(allLogs) {
			return bo.NewPageResponseBo(req.PageRequestBo, []*do.MessageLog{}), nil
		}

		end := offset + limit
		if end > len(allLogs) {
			end = len(allLogs)
		}

		allLogs = allLogs[offset:end]
	}

	return bo.NewPageResponseBo(req.PageRequestBo, allLogs), nil
}

// readAllLogsFromFile 从文件中读取所有符合条件的消息日志
func (m *messageLogRepositoryImpl) readAllLogsFromFile(filePath string, req *bo.ListMessageLogBo) ([]*do.MessageLog, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var logs []*do.MessageLog
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var msgLog do.MessageLog
		if err := m.codec.Unmarshal([]byte(line), &msgLog); err != nil {
			m.helper.Warnf("failed to unmarshal line in %s: %v", filePath, err)
			continue
		}

		// 过滤条件
		if !req.StartAt.IsZero() && msgLog.SendAt.Before(req.StartAt) {
			continue
		}
		if !req.EndAt.IsZero() && msgLog.SendAt.After(req.EndAt) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && msgLog.Status != req.Status {
			continue
		}
		if req.Type.Exist() && !req.Type.IsUnknown() && msgLog.Type != req.Type {
			continue
		}

		logs = append(logs, &msgLog)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return logs, nil
}

// updateMessageLogInFile 更新文件中指定 UID 的消息日志
// 使用内存中的位置映射直接定位到行号，避免全文件扫描
func (m *messageLogRepositoryImpl) updateMessageLogInFile(msgLog *do.MessageLog) error {
	m.fileMutex.Lock()
	defer m.fileMutex.Unlock()

	// 从内存映射中获取文件位置
	location, ok := m.uidToLocation.Get(msgLog.UID)
	if !ok {
		// 如果没有位置信息，说明是新记录，直接写入
		_, err := m.writeMessageLogUnlocked(msgLog)
		return err
	}

	filePath := location.filePath
	lineNum := location.lineNum

	// 读取文件所有行
	file, err := os.OpenFile(filePath, os.O_RDWR, 0o644)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，创建新文件并写入
			_, err := m.writeMessageLogUnlocked(msgLog)
			return err
		}
		return fmt.Errorf("failed to open log file %s: %w", filePath, err)
	}
	defer file.Close()

	// 读取所有行
	scanner := bufio.NewScanner(file)
	var lines []string
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		line := scanner.Text()

		// 如果是目标行，更新它
		if currentLine == lineNum {
			// 序列化更新后的数据
			updatedBytes, err := m.codec.Marshal(msgLog)
			if err != nil {
				return fmt.Errorf("failed to marshal updated message log: %w", err)
			}
			lines = append(lines, string(updatedBytes))
		} else {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	// 如果目标行号超出文件范围，追加到末尾
	if lineNum > currentLine {
		updatedBytes, err := m.codec.Marshal(msgLog)
		if err != nil {
			return fmt.Errorf("failed to marshal message log: %w", err)
		}
		lines = append(lines, string(updatedBytes))
		// 更新位置映射
		m.uidToLocation.Set(msgLog.UID, &fileLocation{
			filePath: filePath,
			lineNum:  len(lines),
		})
	}

	// 重新写入文件
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file %s: %w", filePath, err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file %s: %w", filePath, err)
	}

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write line to file %s: %w", filePath, err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush file %s: %w", filePath, err)
	}

	return nil
}

// writeMessageLogUnlocked 写入消息日志到文件（不加锁版本，用于已持有锁的情况）
func (m *messageLogRepositoryImpl) writeMessageLogUnlocked(msgLog *do.MessageLog) (int, error) {
	// 根据 SendAt 确定文件路径
	filePath := m.getDateLogFile(msgLog.SendAt)

	// 计算当前文件的行数
	lineCount, err := m.countLines(filePath)
	if err != nil && !os.IsNotExist(err) {
		return 0, fmt.Errorf("failed to count lines in file %s: %w", filePath, err)
	}

	// 打开文件（追加模式，如果不存在则创建）
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return 0, fmt.Errorf("failed to open log file %s: %w", filePath, err)
	}
	defer file.Close()

	// 序列化为 JSON
	dataBytes, err := m.codec.Marshal(msgLog)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal message log: %w", err)
	}

	// 写入文件（一行一条）
	if _, err := file.Write(append(dataBytes, '\n')); err != nil {
		return 0, fmt.Errorf("failed to write to log file %s: %w", filePath, err)
	}

	// 新行的行号 = 原行数 + 1
	newLineNum := lineCount + 1

	// 更新位置映射
	m.uidToLocation.Set(msgLog.UID, &fileLocation{
		filePath: filePath,
		lineNum:  newLineNum,
	})

	return newLineNum, nil
}

// UpdateMessageLogStatusIf implements repository.MessageLog.
func (m *messageLogRepositoryImpl) UpdateMessageLogStatusIf(ctx context.Context, uid snowflake.ID, oldStatus vobj.MessageStatus, newStatus vobj.MessageStatus) (bool, error) {
	location, ok := m.uidToLocation.Get(uid)
	if !ok {
		return false, merr.ErrorNotFound("message log %d not found", uid.Int64())
	}

	// 从文件读取数据
	msgLog, err := m.readMessageLogFromFile(location)
	if err != nil {
		return false, err
	}

	// 检查当前状态是否匹配
	if msgLog.Status != oldStatus {
		return false, nil
	}

	// 更新状态
	msgLog.Status = newStatus
	msgLog.UpdatedAt = time.Now()

	// 更新文件中的对应行
	if err := m.updateMessageLogInFile(msgLog); err != nil {
		return false, fmt.Errorf("failed to update message log in file: %w", err)
	}

	return true, nil
}
