package main

import (
	"math/big"
	"zkSnacks/walletsdk"
)

func generateAgeClaim(holderID string, token string) *walletsdk.ClaimAPI {

	studentInfo := getStudentInfoByToken(token)

	claimSchemaHashHex := walletsdk.GetHashFromClaimSchema(CLAIM_SCHEMA_ROOT_DIR+"student-age.json-ld", "AgeCredential")

	birthday := new(big.Int)
	birthday.SetString(studentInfo.BirthDate, 10)

	claimAPI := walletsdk.ClaimAPI{
		SubjectID:          holderID,
		ClaimSchemaHashHex: claimSchemaHashHex,
		IndexSlotA:         birthday,
	}

	return &claimAPI
}
