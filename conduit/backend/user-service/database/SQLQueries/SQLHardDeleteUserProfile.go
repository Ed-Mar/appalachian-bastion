package SQLQueries

func SQLHardDeleteUserProfile() string {
	return sqlHardDeleteUserProfile
}

const sqlHardDeleteUserProfile = `DELETE FROM user_profiles
	WHERE external_auth_id = $1;
`
