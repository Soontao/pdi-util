package pdiutil

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// DefaultPackageCheckInterval for set the default interval for status checking
// unit is second
const DefaultPackageCheckInterval = 20

// ActivationSolution sync
func (c *PDIClient) ActivationSolution(solution string) (err error) {

	if c.sessionID == "" {
		return fmt.Errorf("session id lost")
	}

	res := c.xrepRequest("00163E028DAE1EE1ABE2189C1FF64B07", JSONObject{
		"IMPORTING": JSONObject{
			"IV_ACT_SPLIT":    "X",
			"IV_CALLER":       "AssembleAndDownload",
			"IV_DELTA_MODE":   "X",
			"IV_MODE":         "D",
			"IV_PRODUCT_NAME": solution,
			"IV_SESSION_ID":   c.sessionID,
		},
	})

	success := gjson.Get(res, "EXPORTING.EV_SUCCESS").String() == "X"

	if success {

		// wait job finished
		for {

			solutionHeader := c.GetSolutionStatus(solution)

			if solutionHeader.IsRunningJob {
				// still in running
				// wait interval then check it.
				time.Sleep(DefaultPackageCheckInterval * time.Second)
			} else {
				// finished
				if !solutionHeader.CanAssemble {
					// finished but can not assemble
					// error happened
					err = fmt.Errorf("Activation failed, please check at PDI UI")
				}
				break
			}
		}

	} else {
		err = fmt.Errorf(gjson.Get(res, "EXPORTING.ET_MESSAGES").String())
	}

	return err

}

// FindBACAndActivateIt
func (c *PDIClient) FindBACAndActivateIt(solution string) (err error) {

	// NEED do more update
	if c.sessionID == "" {
		return fmt.Errorf("session id lost")
	}

	files := c.GetSolutionXrepFileList(solution)
	bacFile := ""

	for _, file := range files {
		if strings.HasSuffix(file, ".bac") {
			bacFile = file
			break
		}
	}

	if bacFile == "" {
		err = fmt.Errorf("Not found BAC from solution: %v", solution)
	} else {
		log.Printf("Updating BAC file: %s", bacFile)
		payload := JSONObject{
			"IMPORTING": JSONObject{
				"IV_CONTENT_TYPE": "BUSINESSCONFIGURATION",
				"IV_SESSION_ID":   c.sessionID,
				"IV_XREP_PATH":    bacFile,
			},
		}
		if _, err := c.xrepRequestE("00163E0115B01DDFB194EC88B8EE4C9B", payload); err != nil {
			log.Println("Updating BAC failed")
		}
	}

	return err
}

// AssembleSolution from server
func (c *PDIClient) AssembleSolution(solution string) (err error) {

	if c.sessionID == "" {
		return fmt.Errorf("session id lost")
	}

	res := c.xrepRequest("00163E0975CB1ED4B79AD60DC0D91314", JSONObject{
		"IMPORTING": JSONObject{
			"IV_PRODUCT_NAME": solution,
			"IV_SESSION_ID":   c.sessionID,
			"IV_USER":         nil,
		},
	})

	success := gjson.Get(res, "EXPORTING.EV_SUCCESS").String() == "X"

	if success {

		// wait job finished
		for {

			solutionHeader := c.GetSolutionStatus(solution)

			if solutionHeader.IsRunningJob {
				// still in running
				// wait interval then check it.
				time.Sleep(DefaultPackageCheckInterval * time.Second)
			} else {
				// finished
				if !solutionHeader.CanDownload {
					// finished but can not assemble
					// error happened
					err = fmt.Errorf("Assemble failed, please check at PDI UI")
				}
				break
			}
		}

	} else {
		err = fmt.Errorf(gjson.Get(res, "EXPORTING.ET_MESSAGES").String())
	}

	return err
}

// DownloadSolution from tenant
// return base64 binary zip file
func (c *PDIClient) DownloadSolution(solution, version string) (err error, output string) {

	res := c.xrepRequest("00163E0975CB1ED4B79AD6AC1C161314", JSONObject{
		"IMPORTING": JSONObject{
			"IV_MINOR_VERSION": version,
			"IV_PROJECT_NAME":  solution,
			"IV_PROJECT_TYPE":  "ONE_OFF",
		},
	})

	success := gjson.Get(res, "EXPORTING.EV_SUCCESS").String() == "X"

	if success {
		output = gjson.Get(res, "EXPORTING.EV_SOLUTION_ASSEMBLE").String()
	} else {
		err = fmt.Errorf(gjson.Get(res, "EXPORTING.ET_MESSAGES").String())
	}

	return err, output
}

// CreatePatch solution
func (c *PDIClient) CreatePatch(solution string) (err error) {
	checkInterval := DefaultPackageCheckInterval * time.Second

	payload := JSONObject{
		"IMPORTING": JSONObject{
			"IV_DELETION_PATCH": false,
			"IV_PRODUCT_NAME":   solution,
			"IV_SESSION_ID":     c.sessionID,
			"IV_USER":           nil,
		},
	}

	if _, err = c.xrepRequestE("00163E1267F91EE5B7D7285EE2C105CE", payload); err != nil {
		m := "Create patch solution failed."
		log.Printf(m)
		err = fmt.Errorf(m)
		return err
	}

	time.Sleep(checkInterval)

	// wait patch solution created
	for {

		// retrieve status
		solutionHeader := c.GetSolutionStatus(solution)

		if solutionHeader.IsCreatingPatch {
			// still in running
			// wait interval then check it.
			time.Sleep(checkInterval)
		} else {
			// finished
			// in development now
			if solutionHeader.Status == S_STATUS_IN_DEV {
				break
			} else {
				err = fmt.Errorf("Patch created, but not in development")
			}
		}
	}

	return err
}
