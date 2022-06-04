package svelte
    
import (
	"github.com/gofiber/fiber/v2"
	"github.com/kokizzu/gotro/M"
)

var handlers = map[string]func(c *fiber.Ctx, path string) (res M.SX, err error){
	"index": index,
	"page1/subpage": page1_subpage,
	"page1/subpage3/index": page1_subpage3_index,
	"page2/index": page2_index,
}