package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 密码哈希函数
func HashPassword(password string) (string, error) {
	// cost 参数控制哈希计算强度（4-31）
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 密码验证函数
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
