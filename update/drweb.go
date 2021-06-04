package update

import (
	"fmt"
	"malice-new/scan"

	"github.com/labstack/echo"
)

func drwebupdate(key string) interface{} {

	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	// defer cancel()

	if key != "" {
		err := updatedrwebLicense(key)
		if err != nil {
			return echo.Map{"drweb": err.Error()}
		}
	}

	fmt.Println("Updating Dr.WEB...")
	output, err := scan.RunCommand(nil, "/opt/drweb.com/bin/drweb-ctl", "update")

	if err != nil {
		return echo.Map{"drweb": err.Error()}
	}

	return echo.Map{"drweb": output}
}

func updatedrwebLicense(key string) error {

	if len(key) > 0 {
		_, err := scan.RunCommand(nil, "/opt/drweb.com/bin/drweb-ctl", "license", "--GetRegistered", key)
		if err != nil {
			return err
		}
	} else {
		_, err := scan.RunCommand(nil, "/opt/drweb.com/bin/drweb-ctl", "license", "--GetDemo")
		if err != nil {
			return err
		}
	}

	return nil
}
