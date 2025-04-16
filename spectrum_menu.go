package tinysa

import (
	"fmt"
	"strconv"
	"strings"
)

// TriggerMenu virtually clicks on the menu items on the device. First element starts with 1.
// Example: [6, 4, 2] enables the waterfall display.
func (d *Device) TriggerMenu(menuIds []uint) error {
	d.logger.Info("triggering menu", "menu_ids", menuIds)
	strs := make([]string, len(menuIds))
	for i, v := range menuIds {
		if v <= 0 {
			return fmt.Errorf("%w: invalid menu id: %d", ErrCommandFailed, v)
		}
		strs[i] = strconv.Itoa(int(v))
	}
	menuStr := strings.Join(strs, " ")
	_, err := d.sendCommand(fmt.Sprintf("menu %s", menuStr))
	return err
}
