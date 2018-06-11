package Config
import "testing"
import "reflect"
import "fmt"

func TestGetConfigHandler(t *testing.T) {
  str := "PORT = 8888"+
    "\nMAX_NUM_files= 16" +
   "\nmax_filesize_MB =8" +
   "\nmax_user_uploads_per_minute =       4"

  c1,err := ParseServerConfig(str)

  if err != nil {
    t.Error(err)
  }

  expected := ServerConfig{8888,16,8,4}
  if !reflect.DeepEqual(*c1,expected) {
    t.Error(fmt.Sprintf("Expected: |%s| Actual: |%s|",expected,c1))
  }
}

