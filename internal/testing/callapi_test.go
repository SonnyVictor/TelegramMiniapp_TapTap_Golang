package testing

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

type Post struct {
}

func testcallapi() {
	client := resty.New()

	// URL của API
	url := "http://localhost:8080/clicktoearn"

	// Tạo dữ liệu POST
	post := Post{}

	const lol = "user=%7B%22id%22%3A5204989815%2C%22first_name%22%3A%22Sonny%22%2C%22last_name%22%3A%22Victok%22%2C%22username%22%3A%22SonnyVitok%22%2C%22language_code%22%3A%22en%22%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2FrUi7bmeMCZfIcid5iCSQZWMRUTPriPLJVNmdvts9Kc1umIhBNJ2c4_lJlVCeSlX0.svg%22%7D&chat_instance=-4534184936334365392&chat_type=private&auth_date=1740494970&signature=jszRckMd4AZn0o2OeYwzRuHfiKuAnUzwKUVtU_SOi-XhN3-Ml3fRu3TvvMmhx_PHqUmkfKktQcUmtyMH0u_tBQ&hash=17b32d1151055f447f3b75484bb48f753570ccd3bcd8f16fe46f20544a29483e"

	// Gọi API POST với custom header
	response, err := client.R().
		SetHeader("tma", lol).
		SetBody(post).
		Post(url)
	if err != nil {
		log.Fatalf("Error calling API: %v", err)
	}

	// In kết quả
	fmt.Println("Response:", response.String())
}
func TestCallApi(t *testing.T) {
	for i := 0; i < 1000; i++ {
		testcallapi()
		fmt.Println("Call api", i)
		time.Sleep(100 * time.Millisecond)
	}
}
