package handlers

import (
	"backend/internal/authentication/config"
	"backend/internal/helper"
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"time"
)

const KEY_IN_HEADER_WITH_TOKEN = "Authentication"
const LOG_MODE = true

var ErrAuthBaseError = fmt.Errorf("[ERROR] [AUTHENTICATION]")
var WarnAuthBaseWarning = fmt.Errorf("[WARNING-ERROR] [AUTHENTICATION]")

var WarnAuthNoAccessTokenInHeader = fmt.Errorf("%v: No Acess Token Found in Header", WarnAuthBaseWarning)

var WarnAuthUnsupportedJWTIssuer = fmt.Errorf(" %v: Non Supported Issuer in Submitted JWT", WarnAuthBaseWarning)
var WarnAuthUnsupportedJWTAuthParty = fmt.Errorf(" %v: Non Supported Authorization Party in Submitted JWT", WarnAuthBaseWarning)
var WarAuthJWTTimesBeforeNow = fmt.Errorf(" %v: AuthTime is after now somehow in Submitted JWT", WarnAuthBaseWarning)
var WarnAuthJWTTimeHasComeAndGone = fmt.Errorf(" %v: JWT has expired", WarnAuthBaseWarning)

var WarnAuthInvalidTokenBase = fmt.Errorf(" %v [LOCAL-JWT-VERIFICATION] [INVALID-TOKEN]", WarnAuthBaseWarning)
var WarnAuthInvalidToken = fmt.Errorf(" %v: Token was found to be Invalid upon local verifcation", WarnAuthInvalidTokenBase)
var WarnAuthInvalidTokenMalformed = fmt.Errorf(" %v: [MALFORMED-TOKEN]: Token was found to be Invalid upon local verifcation due to being Malformed", WarnAuthInvalidTokenBase)
var WarnAuthInvalidTokenExpired = fmt.Errorf(" %v: [EXPIRED-TOKEN]: Token was found to be Invalid due to the Token being expired", WarnAuthInvalidTokenBase)
var WarnAuthInvalidTokenPreDated = fmt.Errorf(" %v: [PREDATED-TOKEN]: Token was found to be Invalid due to the Token being predated", WarnAuthInvalidTokenBase)

var WarnAuthTokenSignedWithNonSupportedAlg = fmt.Errorf(" %v: Token Was Signed with a non Supported alg", WarnAuthInvalidTokenBase)

var ErrAuthPublicRSAKeyLoading = fmt.Errorf(" %v [LOCAL-JWT-VERIFICATION] [RSA PUBLIC]: Error encounter when trying to retive RSA Public Key from env file", ErrAuthBaseError)
var ErrAuthParsingPublicRSA = fmt.Errorf(" %v [LOCAL-JWT-VERIFICATION] [RSA PUBLIC]: Error encounter when trying to parse RSA Public Key", ErrAuthBaseError)

// This is to get a const list, cause you cant const string slice
func getListOfAcceptedSSOIssuers() []string {
	return []string{"http://keycloak.test/realms/gatehouse"}
}
func getListOfAcceptedAuthorzationParties() []string {
	return []string{"dev-conduit-rust"}
}

type KeyAuth struct{}

type auth struct {
	//    "upn": "test-user-####",
	//    "sid": "UUID",
	//    "groups": [],
	//    "auth_time": "Unix-Time",
	//    "exp": "Expiration Time",
	//    "azp": "authorization Party",
	//    "iss": "issuer"
	//Authentication.ExternalID the 'sid' in the JWT this will be used to id the user in this Application and the external SSO
	ExternalID uuid.UUID
	//Authentication.UserPrincipalName the 'upn' in the JWT
	UserPrincipalName string
	//Authentication.JWTIssuer the 'iss' in the JWT
	JWTIssuer string
	//Authentication.JWTAuthorizationParty the 'azp' in the JWT
	JWTAuthorizationParty string
	//Authentication.JWTIssuerGroups the 'groups' in teh JWT
	JWTIssuerGroups []string
	//Authentication.JWTAuthorizationTime the 'auth_time' in the JWT
	JWTAuthorizationTime *time.Time
	//Authentication.JWTExpiration the 'exp' in the JWT
	JWTExpiration *time.Time
	//Authentication.LastTokenIntrospective The time of the last token introspective if any
	LastTokenIntrospective *time.Time
	//Authentication.LastTokenLocalVerification The time of the last local Verification of the token
	LastTokenLocalVerification *time.Time
}

func newAuthViaIntrospective(sid uuid.UUID, upn string, iss string, azp string, groups []string, authTime *time.Time, exp *time.Time, introspectiveTime *time.Time) (*auth, error) {

	//Note I don't really think any of the checks are necessary, but they make me feel better so that's
	// gotta be worth something.

	flag := false
	//Checks to see if the issuer is one of our accepted ones
	for _, issuer := range getListOfAcceptedSSOIssuers() {
		if iss == issuer {
			flag = true
			break
		}
	}
	if flag != true {
		return &auth{}, WarnAuthUnsupportedJWTIssuer
	}
	//Checks to see if the authorization party is one of our accepted ones
	for _, issuer := range getListOfAcceptedAuthorzationParties() {
		if azp == issuer {
			flag = true
			break
		}
	}
	if flag != true {
		return &auth{}, WarnAuthUnsupportedJWTAuthParty
	}

	//Sanctity check if I need these I'm already fucked. I really don't think I should ever need this in the current flow.
	if time.Now().Unix() < authTime.Unix() {
		return &auth{}, fmt.Errorf("%v | Now: %v < auth_time: %v \n", WarAuthJWTTimesBeforeNow, time.Now(), authTime)
	}
	if time.Now().Unix() > exp.Unix() {
		return &auth{}, fmt.Errorf("%v | Now: %v > exp: %v \n", WarnAuthJWTTimeHasComeAndGone, time.Now(), authTime)
	}
	if introspectiveTime.Unix() > exp.Unix() {
		return &auth{}, WarnAuthJWTTimeHasComeAndGone
	}

	a := auth{
		ExternalID:                 sid,
		UserPrincipalName:          upn,
		JWTIssuer:                  iss,
		JWTAuthorizationParty:      azp,
		JWTIssuerGroups:            groups,
		JWTAuthorizationTime:       authTime,
		JWTExpiration:              exp,
		LastTokenIntrospective:     introspectiveTime,
		LastTokenLocalVerification: nil,
	}
	return &a, nil

}
func newAuthViaLocalTokenLocalVerification(sid uuid.UUID, upn string, iss string, azp string, groups []string, authTime *time.Time, exp *time.Time, tokenLocalVerification *time.Time) (*auth, error) {

	//Note I don't really think any of the checks are necessary, but they make me feel better so that's
	// gotta be worth something.

	flag := false
	//Checks to see if the issuer is one of our accepted ones
	for _, issuer := range getListOfAcceptedSSOIssuers() {
		if iss == issuer {
			flag = true
			break
		}
	}
	if flag != true {
		return &auth{}, WarnAuthUnsupportedJWTIssuer
	}
	//Checks to see if the authorization party is one of our accepted ones
	for _, issuer := range getListOfAcceptedAuthorzationParties() {
		if azp == issuer {
			flag = true
			break
		}
	}
	if flag != true {
		return nil, WarnAuthUnsupportedJWTAuthParty
	}

	//Sanctity check if I need these I'm already fucked
	if time.Now().Unix() < authTime.Unix() {
		return nil, fmt.Errorf("%v | Now: %v < auth_time: %v \n", WarAuthJWTTimesBeforeNow, time.Now(), authTime)
	}
	if time.Now().Unix() > exp.Unix() {
		return nil, fmt.Errorf("%v | Now: %v > exp: %v \n", WarnAuthJWTTimeHasComeAndGone, time.Now(), authTime)
	}
	if tokenLocalVerification.Unix() > exp.Unix() {
		return &auth{}, WarnAuthJWTTimeHasComeAndGone
	}

	a := auth{
		ExternalID:                 sid,
		UserPrincipalName:          upn,
		JWTIssuer:                  iss,
		JWTAuthorizationParty:      azp,
		JWTIssuerGroups:            groups,
		JWTAuthorizationTime:       authTime,
		JWTExpiration:              exp,
		LastTokenIntrospective:     nil,
		LastTokenLocalVerification: tokenLocalVerification,
	}
	return &a, nil

}
func (a *auth) getExternalID() uuid.UUID {
	return a.ExternalID
}
func (a *auth) getUserPrincipalName() string {
	return a.UserPrincipalName
}
func (a *auth) getJWTIssuer() string {
	return a.JWTIssuer
}
func (a *auth) getJWTAuthorizationParty() string {
	return a.JWTAuthorizationParty
}
func (a *auth) getJWTIssuerGroups() []string {
	return a.JWTIssuerGroups
}
func (a *auth) getJWTAuthorizationTime() *time.Time {
	return a.JWTAuthorizationTime
}
func (a *auth) getJWTExpiration() *time.Time {
	return a.JWTExpiration
}
func (a *auth) getLastTokenIntrospective() *time.Time {
	return a.LastTokenIntrospective
}
func (a *auth) getLastTokenLocalVerification() *time.Time {
	return a.LastTokenLocalVerification
}

func newAuthViaLocalVerification(rw *http.ResponseWriter, r *http.Request, logger *log.Logger) (*auth, error) {
	if r.Header[KEY_IN_HEADER_WITH_TOKEN] == nil {
		http.Error(*rw, WarnAuthNoAccessTokenInHeader.Error(), http.StatusBadRequest)
		return &auth{}, WarnAuthNoAccessTokenInHeader
	}
	tokenString := r.Header.Get(KEY_IN_HEADER_WITH_TOKEN)
	jwtToken, TimeOfVerification, err := localJWTVerification(&tokenString)
	unwrapped := errors.Unwrap(err)
	if err == WarnAuthInvalidToken || unwrapped == WarnAuthInvalidToken {
		http.Error(*rw, WarnAuthInvalidToken.Error(), http.StatusForbidden)
		return nil, err
	} else if unwrapped == WarnAuthInvalidTokenBase || err == WarnAuthInvalidTokenBase {
		http.Error(*rw, err.Error(), http.StatusBadRequest)
		return nil, err
	} else if unwrapped == ErrAuthBaseError || err == ErrAuthPublicRSAKeyLoading || err == ErrAuthParsingPublicRSA {
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil, err
	} else if err != nil {
		http.Error(*rw, WarnAuthInvalidToken.Error(), http.StatusForbidden)
		return nil, err
	}
	jwtMappedClaims, err := localJWTVerificationClaims(jwtToken)
	if err != nil {
		return nil, err
	}

	//log.Printf("%v", jwtMappedClaims)

	//ExternalID:                 sid,
	//UserPrincipalName:          upn,
	//JWTIssuer:                  iss,
	//JWTAuthorizationParty:      azp,
	//JWTIssuerGroups:            groups,
	//JWTAuthorizationTime:       authTimeFloat64,
	//JWTExpiration:              expFloat64,
	//LastTokenIntrospective:     nil,
	//LastTokenLocalVerification: tokenLocalVerification,

	// sid ---
	sidString, ok := jwtMappedClaims["sid"].(string)
	if !ok {
		err = fmt.Errorf("%v %v", ErrAuthBaseError, mappingErrorOuputFormatting("sid"))
		logger.Println(err)
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil, err

	}
	//sid Conversion from string into uuid
	sidUUID, err := uuid.FromString(sidString)
	if err != nil {
		logger.Println(err)
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil, fmt.Errorf("%v [MAPPING] [\"sid\"] Error Converting String to UUID | %v", ErrAuthBaseError, err)
	}

	// upn---
	upn, ok := jwtMappedClaims["upn"].(string)
	if !ok {
		err = fmt.Errorf("%v %v", ErrAuthBaseError, mappingErrorOuputFormatting("upn"))
		logger.Println(err)
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	// iss--
	iss, ok := jwtMappedClaims["iss"].(string)
	if !ok {
		err = fmt.Errorf("%v %v", ErrAuthBaseError, mappingErrorOuputFormatting("iss"))
		logger.Println(err)
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	// azp--
	azp, ok := jwtMappedClaims["azp"].(string)
	if !ok {
		err = fmt.Errorf("%v %v", ErrAuthBaseError, mappingErrorOuputFormatting("azp"))
		logger.Println(err)
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	// groups--
	groupsInterface, ok := jwtMappedClaims["groups"].([]interface{})
	if !ok {
		err = fmt.Errorf("%v %v", ErrAuthBaseError, mappingErrorOuputFormatting("groups"))
		logger.Println(err)
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	//groups Conversion in to String Slice
	groups := make([]string, len(groupsInterface))
	for i, v := range groupsInterface {
		groups[i], ok = v.(string)
		if !ok {
			//(I would tell what it is but if it wont become a string there is not hope in logging it)
			err = fmt.Errorf("%v [MAPPING] [\"sid\"] Error Converting items in groups to string  | at index %v | %#v", ErrAuthBaseError, i, groups)
			logger.Println(err)
			http.Error(*rw, err.Error(), http.StatusInternalServerError)
			return nil, err
		}
	}
	// auth_time--
	authTimeFloat64, ok := jwtMappedClaims["auth_time"].(float64)
	if !ok {
		err = fmt.Errorf("%v %v", ErrAuthBaseError, mappingErrorOuputFormatting("auth_time"))
		logger.Println(err)
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	// auth_time Conversion to *time
	authTime := helper.FloatToUnixTime(authTimeFloat64)
	//Sanctity Check I am checking if the time could make any sense by comparing it to the current time +/- a month (604800 seconds in a month)
	if !(authTime.Unix() > (time.Now().Unix()-604800) && (authTime.Unix() < (time.Now().Unix() + 604800))) {
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil, fmt.Errorf("%v [MAPPING] [\"auth_time\"] Error while converting it to time. The time doesn't make sense so the conversion is most likly off | %#v ", ErrAuthBaseError, authTime)
	}

	// expFloat64 --
	expFloat64, ok := jwtMappedClaims["exp"].(float64)
	if !ok {
		err = fmt.Errorf("%v %v", ErrAuthBaseError, mappingErrorOuputFormatting("exp"))
		logger.Println(err)
		http.Error(*rw, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	// auth_time Conversion to *time
	expTime := helper.FloatToUnixTime(expFloat64)
	//Sanctity Check I am checking if the time could make any sense by comparing it to the current time +/- a month (604800 seconds in a month)
	if !(expTime.Unix() > (time.Now().Unix()-604800) && (expTime.Unix() < (time.Now().Unix() + 604800))) {
		return nil, fmt.Errorf("%v [MAPPING] [\"exp\"] Error while converting it to time. The time doesn't make sense so the conversion is most likly off | %#v", ErrAuthBaseError, expTime)
	}

	//--Encase of shitty Debugging
	//log.Printf("%#v \n %#v \n %#v \n %#v \n %#v \n %#v \n %#v \n %#v \n", sidString, upn, iss, azp, groups, authTime.String(), expTime.String(), TimeOfVerification.String())

	localVerifiedAuth, err := newAuthViaLocalTokenLocalVerification(sidUUID, upn, iss, azp, groups, &authTime, &expTime, TimeOfVerification)
	if err != nil {
		http.Error(*rw, err.Error(), http.StatusForbidden)
		logger.Println(err)
		return nil, err
	}

	return localVerifiedAuth, nil
}
func NewAuthViaLocalVerificationInContext(r *http.Request, rw *http.ResponseWriter, logger *log.Logger) (*http.Request, http.ResponseWriter, bool) {
	auth, err := newAuthViaLocalVerification(rw, r, logger)
	if err != nil {
		logger.Println(err)
		return r, *rw, false
	}
	r = r.WithContext(SetAuthInContext(r.Context(), auth))
	return r, *rw, true

}

func SetAuthInContext(ctx context.Context, a *auth) context.Context {
	ctx = context.WithValue(ctx, KeyAuth{}, a)
	return ctx
}
func GetAuthFromContext(ctx context.Context) (a *auth, err error) {
	a, ok := ctx.Value(KeyAuth{}).(*auth)
	if !ok {
		temp := &auth{}
		return temp, fmt.Errorf("[ERROR] [AUTH] error loading Auth Obj from the context | %v", ok)
	}

	return a, nil
}

func localJWTVerification(accessToken *string) (*jwt.Token, *time.Time, error) {

	const filePath = "../internal/authentication/config"
	const fileName = "localDevPublicRSASecret"
	const fileType = "env"

	rsaSecret, err := config.GetPublicRSASecret(filePath, fileName, fileType)
	if err != nil {
		return nil, nil, fmt.Errorf("%v | %v", ErrAuthPublicRSAKeyLoading, err)
	}
	// PEM formatting
	rsaSecret = "-----BEGIN CERTIFICATE-----\n" +
		rsaSecret +
		"\n-----END CERTIFICATE-----"
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(rsaSecret))
	if err != nil {
		return nil, nil, fmt.Errorf("%v | %v", ErrAuthParsingPublicRSA, err)
	}
	timeOfVerification := time.Time{}

	token, err := jwt.Parse(*accessToken, func(token *jwt.Token) (interface{}, error) {
		//Validate the security alg-orithm:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("%v | Token was Signed with unexpected signing method [%v] while the expected being [rsa]", WarnAuthTokenSignedWithNonSupportedAlg, token.Header["alg"])
		}
		timeOfVerification = time.Now()
		return key, nil
	})

	if token.Valid {
		return token, &timeOfVerification, err
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, nil, WarnAuthInvalidTokenMalformed
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		log.Println("[WARNING] [WARNING] [WARNING] Post-Date JWT Attempt: ", err)
		return nil, nil, WarnAuthInvalidTokenExpired
	} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
		log.Println("[WARNING] [WARNING] [WARNING] Pre-Date JWT Attempt:", err)
		return nil, nil, WarnAuthInvalidTokenPreDated
	} else {
		return nil, nil, WarnAuthInvalidToken
	}
}

func localJWTVerificationClaims(jwtToken *jwt.Token) (jwt.MapClaims, error) {

	//Note-- This also returns a bool I'm Not postive what I need to do with cause I am not using any of its Verify function
	if jwtToken.Valid {
		jwtMappedClaims, _ := jwtToken.Claims.(jwt.MapClaims)
		return jwtMappedClaims, nil
	}
	return nil, WarnAuthInvalidToken

}

// I hate doing things over and over again.
func mappingErrorOuputFormatting(key string) string {
	return fmt.Sprintf("[JWT-CLAIMS] [MAPPING]: [\"%v\"] Error while pulling from map from claims on JWT", key)

}
