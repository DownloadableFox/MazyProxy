package main

import "math/rand"

func NetworkRandomEncrypt(buff []byte) []byte {
	// rand
	x := GetRandomByte()
	y := GetRandomByte()

	return NetworkEncrypt(buff, x, y)
}

func NetworkEncrypt(buff []byte, x, y byte) []byte {
	encrypted := make([]byte, 0)

	// initialize
	encrypted = append(encrypted, x)
	encrypted = append(encrypted, y)

	for _, value := range buff {
		encrypted = append(encrypted, value^x)
		x = x + y
	}

	return encrypted
}

func NetworkDecrypt(buff []byte) []byte {
	// preparing
	decrypted := make([]byte, 0)

	// values
	x := buff[0]
	y := buff[1]

	for i := 2; i < len(buff); i++ {
		decrypted = append(decrypted, buff[i]^x)
		x = x + y
	}

	return decrypted
}

func GetRandomByte() byte {
	return byte(rand.Intn(255))
}
