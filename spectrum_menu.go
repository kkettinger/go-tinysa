package tinysa

import (
	"fmt"
	"strconv"
	"strings"
)

// TriggerMenu virtually clicks on the menu items on the device. First element starts with 1.
// Example: [6, 4, 2] enables the waterfall display.
func (d *Device) TriggerMenu(menuIDs []uint) error {
	d.logger.Info("triggering menu", "menu_ids", menuIDs)
	strs := make([]string, len(menuIDs))
	for i, v := range menuIDs {
		if v <= 0 || v > 255 {
			return fmt.Errorf("invalid menu id: %d", v)
		}
		strs[i] = strconv.Itoa(int(v))
	}
	menuStr := strings.Join(strs, " ")
	_, err := d.sendCommand(fmt.Sprintf("menu %s", menuStr))
	return err
}
