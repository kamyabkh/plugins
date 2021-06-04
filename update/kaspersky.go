package update

import (
	"context"
	"fmt"
	"malice-new/scan"
	"time"

	"github.com/labstack/echo"
)

func kasperskyupdate(key string) interface{} {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	// defer cancel()

	if key != "" {
		err := updateLicense(key)
		return echo.Map{"kaspersky": err.Error()}

	}

	fmt.Println("Updating Kaspersky...")
	output, err := scan.RunCommand(nil, "/opt/kaspersky/kav4fs/bin/kav4fs-control", "-T", "--start-task", "6", "--progress")

	if err != nil {
		return echo.Map{"kaspersky": err.Error()}
	}

	return echo.Map{"kaspersky": output}

}

func updateLicense(key string) error {
	// kaspersky needs to have the daemon started first

	// check for exec context timeout

	if len(key) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
		defer cancel()
		_, err := scan.RunCommand(ctx, "/opt/kaspersky/kav4fs/bin/kav4fs-control", "--revoke-active-key", key)
		if err != nil {
			return err
		}
		_, err = scan.RunCommand(ctx, "/opt/kaspersky/kav4fs/bin/kav4fs-control", "--install-active-key", key)
		if err != nil {
			return err
		}
	}

	return nil
}
