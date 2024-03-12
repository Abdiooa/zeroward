package decryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

func DecryptKey(filePath string, kekKey []byte) ([]byte, error) {

	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	offset := fileSize - 60
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, err
	}
	encryptedKey := make([]byte, 60)
	_, err = file.Read(encryptedKey)
	if err != nil {
		return nil, err
	}
	// encryptedKey, err := os.ReadFile(encryptedKeyFile)
	// if err != nil {
	// 	return nil, err
	// }
	decryptedKey, err := DecryptData(encryptedKey, kekKey)
	if err != nil {
		return nil, err
	}
	err = file.Truncate(offset)
	if err != nil {
		return nil, err
	}

	return decryptedKey, nil
}

func DecryptFile(ciphertext []byte, dekKey []byte) ([]byte, error) {
	// ciphertext, err := os.ReadFile(cipherFile)
	// if err != nil {
	// 	return fmt.Errorf("error reading ciphertext file: %v", err)
	// }

	blockSize := 1024 + 4 + 16 + 12
	var decryptedData []byte

	for i := 0; i < len(ciphertext); i += blockSize {
		end := i + blockSize
		if end > len(ciphertext) {
			end = len(ciphertext)
		}
		block := ciphertext[i:end]

		decryptedBlock, err := DecryptData(block, dekKey)
		if err != nil {
			return nil, err
		}
		originalData, err := VerifyChecksum(decryptedBlock)
		if err != nil {
			return nil, err
		}
		decryptedData = append(decryptedData, originalData...)
	}
	// if err := os.Remove(cipherFile); err != nil {
	// 	return err
	// }
	// Remove the ".enc" extension from the file name
	// decryptedFilePath := cipherFile[:len(cipherFile)-4]

	// if err := os.WriteFile(decryptedFilePath, decryptedData, 0644); err != nil {
	// 	return err
	// }

	return decryptedData, nil
}

func DecryptData(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext is too short")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting: %v", err)
	}

	return plaintext, nil
}

func VerifyChecksum(data []byte) ([]byte, error) {
	blockSize := 1024
	var originalData []byte
	for i := 0; i < len(data); i += blockSize + 4 {
		end := i + blockSize + 4
		if end > len(data) {
			end = len(data)
		}
		blockWithChecksum := data[i:end]
		checksumBytes := blockWithChecksum[:4]
		block := blockWithChecksum[4:]

		checksum := crc32.ChecksumIEEE(block)
		if binary.BigEndian.Uint32(checksumBytes) != checksum {
			return nil, fmt.Errorf("checksum verification failed")
		}
		originalData = append(originalData, block...)
	}
	return originalData, nil
}
