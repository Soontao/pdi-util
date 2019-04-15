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
