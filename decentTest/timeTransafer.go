package main
import(
	"fmt"
	"time"
	"os"
	"strconv"
)

func main(){
	unix, _ := strconv.ParseInt(os.Args[1], 10, 64)
	t := time.Unix(0, unix)
	fmt.Println(t)
}
