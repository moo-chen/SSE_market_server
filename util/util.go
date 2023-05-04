package util
import(
	"math/rand"
	"time"
)
// 工具函数，若注册时没有传名字，生成随机字符串
func RandomString(n int) string {
	var letters = []byte("fajvhfaufvafbiauAABIFWFIIudfuwfuwdagcqiuetoeh")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}