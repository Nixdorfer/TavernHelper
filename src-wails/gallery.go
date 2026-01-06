package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetGalleryImages() ([]GalleryImage, error) {
	images, err := DBGetAllGalleryImages()
	if err != nil {
		return nil, err
	}

	for i := range images {
		if _, err := os.Stat(images[i].LocalPath); os.IsNotExist(err) {
			images[i].IsValid = false
		} else {
			images[i].IsValid = true
		}
	}

	return images, nil
}

func (a *App) GetGalleryFolders() ([]GalleryFolder, error) {
	execPath, _ := os.Executable()
	imageDir := filepath.Join(filepath.Dir(execPath), "images")

	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(imageDir)
	if err != nil {
		return nil, err
	}

	var folders []GalleryFolder
	for _, entry := range entries {
		if entry.IsDir() {
			info, _ := entry.Info()
			folders = append(folders, GalleryFolder{
				Name:      entry.Name(),
				Path:      entry.Name(),
				CreatedAt: info.ModTime().Format(time.RFC3339),
			})
		}
	}

	return folders, nil
}

func (a *App) CreateGalleryFolder(name string) (*GalleryFolder, error) {
	execPath, _ := os.Executable()
	imageDir := filepath.Join(filepath.Dir(execPath), "images")
	folderPath := filepath.Join(imageDir, name)

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return nil, fmt.Errorf("创建文件夹失败: %w", err)
	}

	return &GalleryFolder{
		Name:      name,
		Path:      name,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func (a *App) DeleteGalleryFolder(name string) error {
	execPath, _ := os.Executable()
	imageDir := filepath.Join(filepath.Dir(execPath), "images")
	folderPath := filepath.Join(imageDir, name)

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("读取文件夹失败: %w", err)
	}

	if len(entries) > 0 {
		return fmt.Errorf("文件夹不为空，无法删除")
	}

	return os.Remove(folderPath)
}

func (a *App) MoveGalleryImageToFolder(imageId, folderPath string) error {
	images, err := DBGetAllGalleryImages()
	if err != nil {
		return err
	}

	var targetImage *GalleryImage
	for _, img := range images {
		if img.ID == imageId {
			targetImage = &img
			break
		}
	}

	if targetImage == nil {
		return fmt.Errorf("图片不存在")
	}

	execPath, _ := os.Executable()
	imageDir := filepath.Join(filepath.Dir(execPath), "images")

	var newPath string
	if folderPath == "" {
		newPath = filepath.Join(imageDir, filepath.Base(targetImage.LocalPath))
	} else {
		destDir := filepath.Join(imageDir, folderPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("创建目标文件夹失败: %w", err)
		}
		newPath = filepath.Join(destDir, filepath.Base(targetImage.LocalPath))
	}

	if newPath != targetImage.LocalPath {
		if err := os.Rename(targetImage.LocalPath, newPath); err != nil {
			return fmt.Errorf("移动文件失败: %w", err)
		}

		_, err = db.Exec(`UPDATE wt_image SET local_path = ?, folder_path = ? WHERE id = ?`,
			newPath, folderPath, imageId)
		if err != nil {
			os.Rename(newPath, targetImage.LocalPath)
			return fmt.Errorf("更新数据库失败: %w", err)
		}
	} else {
		return DBUpdateGalleryImageFolder(imageId, folderPath)
	}

	return nil
}

func (a *App) SelectAndAddGalleryImage() (*GalleryImage, error) {
	return a.SelectAndAddGalleryImageToFolder("")
}

func (a *App) SelectAndAddGalleryImageToFolder(folderPath string) (*GalleryImage, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择图片",
		Filters: []runtime.FileFilter{
			{DisplayName: "图片文件", Pattern: "*.png;*.jpg;*.jpeg;*.gif;*.webp;*.bmp"},
		},
	})
	if err != nil {
		return nil, err
	}
	if filePath == "" {
		return nil, nil
	}

	return a.AddGalleryImageFromPathToFolder(filePath, folderPath)
}

func (a *App) AddGalleryImageFromPath(filePath string) (*GalleryImage, error) {
	return a.AddGalleryImageFromPathToFolder(filePath, "")
}

func (a *App) AddGalleryImageFromPathToFolder(filePath string, folderPath string) (*GalleryImage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, fmt.Errorf("计算哈希失败: %w", err)
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	existing, err := DBGetGalleryImageByHash(hash)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		existing.IsValid = true
		return existing, nil
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	execPath, _ := os.Executable()
	imageDir := filepath.Join(filepath.Dir(execPath), "images")
	if folderPath != "" {
		imageDir = filepath.Join(imageDir, folderPath)
	}
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return nil, fmt.Errorf("创建图片目录失败: %w", err)
	}

	ext := filepath.Ext(filePath)
	newFileName := hash[:16] + ext
	destPath := filepath.Join(imageDir, newFileName)

	file.Seek(0, 0)
	destFile, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		return nil, fmt.Errorf("复制文件失败: %w", err)
	}

	img := &GalleryImage{
		ID:         fmt.Sprintf("img_%d", time.Now().UnixNano()),
		Hash:       hash,
		LocalPath:  destPath,
		FileName:   filepath.Base(filePath),
		FileSize:   fileInfo.Size(),
		CreatedAt:  time.Now().Format(time.RFC3339),
		FolderPath: folderPath,
		IsValid:    true,
	}

	if err := DBAddGalleryImage(img); err != nil {
		os.Remove(destPath)
		return nil, fmt.Errorf("保存图片记录失败: %w", err)
	}

	return img, nil
}

func (a *App) AddGalleryImageFromBase64(base64Data string, fileName string) (*GalleryImage, error) {
	return a.AddGalleryImageFromBase64ToFolder(base64Data, fileName, "")
}

func (a *App) AddGalleryImageFromBase64ToFolder(base64Data string, fileName string, folderPath string) (*GalleryImage, error) {
	idx := 0
	for i := 0; i < len(base64Data); i++ {
		if base64Data[i] == ',' {
			idx = i + 1
			break
		}
	}
	if idx == 0 {
		idx = 0
	}
	data, err := base64Decode(base64Data[idx:])
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %w", err)
	}
	return a.AddGalleryImageFromBytesToFolder(data, fileName, folderPath)
}

func (a *App) AddGalleryImageFromBytes(data []byte, fileName string) (*GalleryImage, error) {
	return a.AddGalleryImageFromBytesToFolder(data, fileName, "")
}

func (a *App) AddGalleryImageFromBytesToFolder(data []byte, fileName string, folderPath string) (*GalleryImage, error) {
	hasher := sha256.New()
	hasher.Write(data)
	hash := hex.EncodeToString(hasher.Sum(nil))

	existing, err := DBGetGalleryImageByHash(hash)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		existing.IsValid = true
		return existing, nil
	}

	execPath, _ := os.Executable()
	imageDir := filepath.Join(filepath.Dir(execPath), "images")
	if folderPath != "" {
		imageDir = filepath.Join(imageDir, folderPath)
	}
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return nil, fmt.Errorf("创建图片目录失败: %w", err)
	}

	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".png"
	}
	newFileName := hash[:16] + ext
	destPath := filepath.Join(imageDir, newFileName)

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return nil, fmt.Errorf("保存图片失败: %w", err)
	}

	img := &GalleryImage{
		ID:         fmt.Sprintf("img_%d", time.Now().UnixNano()),
		Hash:       hash,
		LocalPath:  destPath,
		FileName:   fileName,
		FileSize:   int64(len(data)),
		CreatedAt:  time.Now().Format(time.RFC3339),
		FolderPath: folderPath,
		IsValid:    true,
	}

	if err := DBAddGalleryImage(img); err != nil {
		os.Remove(destPath)
		return nil, fmt.Errorf("保存图片记录失败: %w", err)
	}

	return img, nil
}

func (a *App) UpdateGalleryImageURL(id, remoteURL string) error {
	return DBUpdateGalleryImageURL(id, remoteURL)
}

func (a *App) RenameGalleryImage(id, newName string) error {
	return DBUpdateGalleryImageName(id, newName)
}

func (a *App) RenameGalleryFolder(oldName, newName string) error {
	execPath, _ := os.Executable()
	imageDir := filepath.Join(filepath.Dir(execPath), "images")
	oldPath := filepath.Join(imageDir, oldName)
	newPath := filepath.Join(imageDir, newName)

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("重命名文件夹失败: %w", err)
	}

	if err := DBRenameFolderImages(oldName, newName); err != nil {
		os.Rename(newPath, oldPath)
		return fmt.Errorf("更新数据库失败: %w", err)
	}

	return nil
}

func (a *App) GetFolderImages(folderPath string) ([]GalleryImage, error) {
	images, err := DBGetAllGalleryImages()
	if err != nil {
		return nil, err
	}

	var result []GalleryImage
	for _, img := range images {
		if img.FolderPath == folderPath {
			if _, err := os.Stat(img.LocalPath); os.IsNotExist(err) {
				img.IsValid = false
			} else {
				img.IsValid = true
			}
			result = append(result, img)
		}
	}

	return result, nil
}

func (a *App) DeleteGalleryImage(id string) error {
	images, err := DBGetAllGalleryImages()
	if err != nil {
		return err
	}

	var targetImage *GalleryImage
	for _, img := range images {
		if img.ID == id {
			targetImage = &img
			break
		}
	}

	if targetImage == nil {
		return fmt.Errorf("图片不存在")
	}

	if err := DBDeleteGalleryImage(id); err != nil {
		return err
	}

	if targetImage.LocalPath != "" {
		os.Remove(targetImage.LocalPath)
	}

	return nil
}

func (a *App) DeleteGalleryImages(ids []string) error {
	for _, id := range ids {
		if err := a.DeleteGalleryImage(id); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) GetGalleryImageLocalURL(id string) (string, error) {
	images, err := DBGetAllGalleryImages()
	if err != nil {
		return "", err
	}

	for _, img := range images {
		if img.ID == id {
			return img.LocalPath, nil
		}
	}
	return "", fmt.Errorf("图片不存在")
}

func (a *App) ReadGalleryImageAsBase64(id string) (string, error) {
	images, err := DBGetAllGalleryImages()
	if err != nil {
		return "", err
	}

	var targetImage *GalleryImage
	for _, img := range images {
		if img.ID == id {
			targetImage = &img
			break
		}
	}

	if targetImage == nil {
		return "", fmt.Errorf("图片不存在")
	}

	data, err := os.ReadFile(targetImage.LocalPath)
	if err != nil {
		return "", fmt.Errorf("读取图片失败: %w", err)
	}

	ext := filepath.Ext(targetImage.LocalPath)
	mimeType := "image/png"
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	case ".bmp":
		mimeType = "image/bmp"
	}

	base64Data := "data:" + mimeType + ";base64," + base64Encode(data)
	return base64Data, nil
}

func base64Encode(data []byte) string {
	const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	result := make([]byte, 0, (len(data)+2)/3*4)
	for i := 0; i < len(data); i += 3 {
		var n uint32
		remaining := len(data) - i
		if remaining >= 3 {
			n = uint32(data[i])<<16 | uint32(data[i+1])<<8 | uint32(data[i+2])
			result = append(result, base64Chars[n>>18&0x3F], base64Chars[n>>12&0x3F], base64Chars[n>>6&0x3F], base64Chars[n&0x3F])
		} else if remaining == 2 {
			n = uint32(data[i])<<16 | uint32(data[i+1])<<8
			result = append(result, base64Chars[n>>18&0x3F], base64Chars[n>>12&0x3F], base64Chars[n>>6&0x3F], '=')
		} else {
			n = uint32(data[i]) << 16
			result = append(result, base64Chars[n>>18&0x3F], base64Chars[n>>12&0x3F], '=', '=')
		}
	}
	return string(result)
}

func base64Decode(s string) ([]byte, error) {
	const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	decodeMap := make(map[byte]byte)
	for i := 0; i < len(base64Chars); i++ {
		decodeMap[base64Chars[i]] = byte(i)
	}
	padCount := 0
	for i := len(s) - 1; i >= 0 && s[i] == '='; i-- {
		padCount++
	}
	s = s[:len(s)-padCount]
	result := make([]byte, 0, len(s)*3/4)
	for i := 0; i < len(s); i += 4 {
		var n uint32
		count := 4
		if i+4 > len(s) {
			count = len(s) - i
		}
		for j := 0; j < count; j++ {
			n = n<<6 | uint32(decodeMap[s[i+j]])
		}
		for j := count; j < 4; j++ {
			n = n << 6
		}
		if count >= 2 {
			result = append(result, byte(n>>16))
		}
		if count >= 3 {
			result = append(result, byte(n>>8))
		}
		if count >= 4 {
			result = append(result, byte(n))
		}
	}
	return result, nil
}
