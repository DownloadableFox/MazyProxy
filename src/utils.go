package main

import "math/rand"

func NetworkEncryptEx(buff []byte) []byte {
	// rand
	encryptor := GetRandomByte()
	offset := GetRandomByte()

	return NetworkEncrypt(buff, encryptor, offset)
}

func NetworkEncrypt(buff []byte, encryptor, offset uint8) []byte {
	encrypted := make([]byte, 0)

	// initialize
	encrypted = append(encrypted, encryptor)
	encrypted = append(encrypted, offset)

	for _, value := range buff {
		encrypted = append(encrypted, value^encryptor)
		encryptor = encryptor&0xff + offset&0xff
	}

	return encrypted
}

func NetworkDecrypt(buff []byte) []byte {
	// preparing
	decrypted := make([]byte, 0)

	// values
	encryptor := buff[0]
	offset := buff[1]

	for i := 2; i < len(buff); i++ {
		decrypted = append(decrypted, buff[i]^encryptor)
		encryptor = encryptor&0xff + offset&0xff
	}

	return decrypted
}

func GetRandomByte() uint8 {
	return uint8(rand.Intn(255))
}
