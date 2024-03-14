package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"hash/crc32"
	"io"
	"os"
)

func EncryptFile(filePath string, dek []byte, kek []byte) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}

	srcFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	plaintext, err := io.ReadAll(srcFile)
	if err != nil {
		return err
	}

	blockSize := 1024
	var encryptedBlocks [][]byte
	for i := 0; i < len(plaintext); i += blockSize {
		end := i + blockSize
		if end > len(plaintext) {
			end = len(plaintext)
		}
		block := plaintext[i:end]

		checksum := crc32.ChecksumIEEE(block)
		checksumBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(checksumBytes, checksum)
		checksumBlock := append(checksumBytes, block...)

		encryptedBlock, err := EncryptData(checksumBlock, dek)
		if err != nil {
			return err
		}
		encryptedBlocks = append(encryptedBlocks, encryptedBlock)
	}

	encryptedData := bytes.Join(encryptedBlocks, nil)

	if err := os.Remove(filePath); err != nil {
		return err
	}

	outputFilePath := filePath + ".enc"
	dstFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.Write(encryptedData)
	if err != nil {
		return err
	}
	encryptedKey, err := EncryptData(dek, kek)
	if err != nil {
		return err
	}
	_, err = dstFile.Write(encryptedKey)
	if err != nil {
		return err
	}
	return nil
}

func EncryptData(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}
