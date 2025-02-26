package sec

import "golang.org/x/crypto/bcrypt"

/*
Convert our password string in hash
 [pass] string
*/
func Hash(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}

/*
Compare if pass and hadhed pass is equal
 [pass] string
 [hashed_pass] string
*/
func ComparePassHash(pass, hashed_pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed_pass), []byte(pass))
}
