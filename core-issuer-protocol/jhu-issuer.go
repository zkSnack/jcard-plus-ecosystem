package main

import (
	"math/big"
	"zkSnacks/walletsdk"
)

func generateAgeClaim(holderID string, token string) (*walletsdk.ClaimAPI, error) {

	studentInfo, err := getStudentInfoByToken(token)
	if err != nil {
		return nil, err
	}

	claimSchemaHashHex := walletsdk.GetHashFromClaimSchema(CLAIM_SCHEMA_ROOT_DIR+"student-age.json-ld", "AgeCredential")

	birthday := new(big.Int)
	birthday.SetString(studentInfo.BirthDate, 10)

	claimAPI := walletsdk.ClaimAPI{
		SubjectID:          holderID,
		ClaimSchemaHashHex: claimSchemaHashHex,
		IndexSlotA:         birthday,
	}

	return &claimAPI, nil
}
