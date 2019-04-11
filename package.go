package pdiutil

import (
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

// DefaultPackageCheckInterval for set the default interval for status checking
// unit is second
const DefaultPackageCheckInterval = 10

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
					err = fmt.Errorf("Activation failed, please check at PDI UI.")
				}
				break
			}
		}

	} else {
		err = fmt.Errorf(gjson.Get(res, "EXPORTING.ET_MESSAGES").String())
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
					err = fmt.Errorf("Assemble failed, please check at PDI UI.")
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
// return base64 banary zip file
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
