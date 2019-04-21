package pdiutil

// DeploySolution to CURRENT tenant
//
// maybe you need two PDIClient instance to do CD operation
func (c *PDIClient) DeploySolution(solutionBase64String string) (err error) {
	payload := JSONObject{
		"IMPORTING": JSONObject{
			"IV_INSTALLATION_KEY":         "",
			"IV_IS_UPLOAD_INTO_PATCH_SOL": " ",
			"IV_SESSION_ID":               c.sessionID,
			"IV_SOLUTION_ASSEMBLE":        solutionBase64String,
		},
	}
	_, err = c.xrepRequestE("00163E0115B01ED09084CBD556024087", payload)
	return err
}

// ActivateDeployedSolution func
// after deploy a solution to target tenant, when uploading progress finished, please active it.
func (c *PDIClient) ActivateDeployedSolution(solution string) (err error) {
	payload := JSONObject{
		"IMPORTING": JSONObject{
			"IV_ACTIVATE_SYNC": nil,
			"IV_PHASE":         "ACT",
			"IV_SESSION_ID":    c.sessionID,
			"IV_SOLUTION_NAME": solution,
		},
	}
	_, err = c.xrepRequestE("00163E0115B01ED09084CF0DBA3A4087", payload)
	return err
}
