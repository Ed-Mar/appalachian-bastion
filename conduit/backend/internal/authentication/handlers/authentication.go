package handlers

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

const KEY_IN_HEADER_WITH_TOKEN = "Authentication"
const LOG_MODE = false

var ErrAuthBaseError = fmt.Errorf("[ERROR] [AUTHENTICATION]")
var WarnAuthBaseWarning = fmt.Errorf("[WARNING-ERROR] [AUTHENTICATION]")
var ErrAuthLocalTokenVerificationBase = fmt.Errorf(" %v [LOCAL-JWT-VERIFICATION] ", ErrAuthBaseError)

var WarnAuthNoAccessTokenInHeader = fmt.Errorf("%v: No Acess Token Found in Header", WarnAuthBaseWarning)

var WarnAuthUnsupportedJWTIssuer = fmt.Errorf(" %v: Non Supported Issuer in Submitted JWT", WarnAuthBaseWarning)
var WarnAuthUnsupportedJWTAuthParty = fmt.Errorf(" %v: Non Supported Authorization Party in Submitted JWT", WarnAuthBaseWarning)
var WarnAuthJWTTimesBeforeNow = fmt.Errorf(" %v: AuthTime is after now somehow in Submitted JWT", WarnAuthBaseWarning)
var WarnAuthJWTTimeHasComeAndGone = fmt.Errorf(" %v: JWT has expired", WarnAuthBaseWarning)

var WarnAuthLocalTokenVerificationInvalidTokenBase = fmt.Errorf(" %v [LOCAL-JWT-VERIFICATION] [INVALID-TOKEN]", WarnAuthBaseWarning)
var WarnAuthInvalidToken = fmt.Errorf(" %v: Token was found to be Invalid upon local verifcation", WarnAuthLocalTokenVerificationInvalidTokenBase)
var WarnAuthInvalidTokenMalformed = fmt.Errorf(" %v: [MALFORMED-TOKEN]: Token was found to be Invalid upon local verifcation due to being Malformed", WarnAuthLocalTokenVerificationInvalidTokenBase)
var WarnAuthInvalidTokenExpired = fmt.Errorf(" %v: [EXPIRED-TOKEN]: Token was found to be Invalid due to the Token being expired", WarnAuthLocalTokenVerificationInvalidTokenBase)
var WarnAuthInvalidTokenPreDated = fmt.Errorf(" %v: [PREDATED-TOKEN]: Token was found to be Invalid due to the Token being predated", WarnAuthLocalTokenVerificationInvalidTokenBase)

var WarnAuthTokenSignedWithNonSupportedAlg = fmt.Errorf(" %v: Token Was Signed with a non Supported alg", WarnAuthLocalTokenVerificationInvalidTokenBase)

var ErrAuthPublicRSAKeyLoading = fmt.Errorf(" %v [LOCAL-JWT-VERIFICATION] [RSA PUBLIC]: Error encounter when trying to retive RSA Public Key from env file", ErrAuthBaseError)
var ErrAuthParsingPublicRSA = fmt.Errorf(" %v [LOCAL-JWT-VERIFICATION] [RSA PUBLIC]: Error encounter when trying to parse RSA Public Key", ErrAuthBaseError)

var ErrAuthTokenIntrospectiveBase = fmt.Errorf("%v [TOKEN-INTROSPECTIVE]", ErrAuthBaseError)
var WarnAuthTokenIntrospectiveBase = fmt.Errorf("%v [TOKEN-INTROSPECTIVE]", WarnAuthBaseWarning)
var WarnAuthTokenIntrospectiveTokenInvalidBase = fmt.Errorf("%v [INVALID-TOKEN] ", WarnAuthTokenIntrospectiveBase)

var ErrAuthLoadingTokenIntrospectiveClientCreds = fmt.Errorf("%v: Error Loading Token Introspective Client Credientaials from ENV File", ErrAuthTokenIntrospectiveBase)
var ErrAuthTokenIntrospectiveRequestCreation = fmt.Errorf("%v: Error Creating Request for Token Intospective", ErrAuthTokenIntrospectiveBase)
var ErrAuthTokenIntrospectiveSendingRequest = fmt.Errorf("%v: Error Sending Request for Token Intospective", ErrAuthTokenIntrospectiveBase)

var WarnAuthTokenIntrospectiveTokenInvalidNotActive = fmt.Errorf("%v: Token is not active", WarnAuthTokenIntrospectiveTokenInvalidBase)
var WarnAuthTokenIntrospectiveRequestStatusNotOK = fmt.Errorf("%v: Token intorspective response status NOT OK(200)", WarnAuthTokenIntrospectiveTokenInvalidBase)

var ErrAuthTokenIntrospectiveResponseToJSON = fmt.Errorf("%v: [JSON] [DECODEING]: There was a an error decoding the Token Introsepctive Response Body From JSON to struct. ", ErrAuthTokenIntrospectiveBase)
var WarnAuthIsNotInContext = fmt.Errorf("%v: Authenication is not in Context", WarnAuthBaseWarning)

func getListOfAcceptedSSOIssuers() []string {
	return []string{"http://keycloak.test/realms/gatehouse"}
}
func getListOfAcceptedAuthorizationParties() []string {
	return []string{"dev-conduit-rust"}
}

type KeyAuthentication struct{}

type authentication struct {
	//newAuthentication.ExternalID the 'sub' in the JWT this will be used to id the user in this Application and the external SSO
	ExternalID uuid.UUID
	//newAuthentication.UserPrincipalName the 'upn' in the JWT
	UserPrincipalName string
	//newAuthentication.JWTIssuer the 'iss' in the JWT
	JWTIssuer string
	//newAuthentication.JWTAuthorizationParty the 'azp' in the JWT
	JWTAuthorizationParty string
	//newAuthentication.JWTIssuerGroups the 'groups' in teh JWT
	JWTIssuerGroups []string
	//newAuthentication.JWTAuthorizationTime the 'auth_time' in the JWT
	JWTAuthorizationTime *time.Time
	//newAuthentication.JWTExpiration the 'exp' in the JWT
	JWTExpiration *time.Time

	//newAuthentication.LastTokenLocalVerification The time of the last local Verification of the token
	LastVerificationTime          *time.Time
	LastVerificationOperationType authenticationVerificationType
}
type authenticationVerificationType interface {
	GetAuthenticationVerificationType() string
}

type localAuthenticationVerificationType struct {
	authenticationVerificationType
}

func (t *localAuthenticationVerificationType) GetAuthenticationVerificationType() string {
	return "Local-Token-Verification"
}

type tokenIntrospectionAuthenticationVerificationType struct {
	authenticationVerificationType
}

func (t *tokenIntrospectionAuthenticationVerificationType) GetAuthenticationVerificationType() string {
	return "Token-Introspective-Verification"
}
func (a *authentication) getExternalID() uuid.UUID {
	return a.ExternalID
}
func (a *authentication) getUserPrincipalName() string {
	return a.UserPrincipalName
}
func (a *authentication) getJWTIssuer() string {
	return a.JWTIssuer
}
func (a *authentication) getJWTAuthorizationParty() string {
	return a.JWTAuthorizationParty
}
func (a *authentication) getJWTIssuerGroups() []string {
	return a.JWTIssuerGroups
}
func (a *authentication) getJWTAuthorizationTime() *time.Time {
	return a.JWTAuthorizationTime
}
func (a *authentication) getJWTExpiration() *time.Time {
	return a.JWTExpiration
}
func (a *authentication) getLastVerification() *time.Time {
	return a.LastVerificationTime
}
func (a *authentication) getLastTokenLocalVerificationType() string {
	return a.LastVerificationOperationType.GetAuthenticationVerificationType()
}
func newAuthentication(sub uuid.UUID, upn string, iss string, azp string, groups []string, authTime *time.Time, exp *time.Time, timeOfVerification *time.Time, typeOfVerification authenticationVerificationType) (*authentication, error) {
	flag := false
	//Checks to see if the issuer is one of our accepted ones
	for _, issuer := range getListOfAcceptedSSOIssuers() {
		if iss == issuer {
			flag = true
			break
		}
	}
	if flag != true {
		return nil, WarnAuthUnsupportedJWTIssuer
	}
	//Checks to see if the authorization party is one of our accepted ones
	for _, issuer := range getListOfAcceptedAuthorizationParties() {
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
		return nil, fmt.Errorf("%v | Now: %v < auth_time: %v \n", WarnAuthJWTTimesBeforeNow, time.Now(), authTime)
	}
	if time.Now().Unix() > exp.Unix() {
		return nil, fmt.Errorf("%v | Now: %v > exp: %v \n", WarnAuthJWTTimeHasComeAndGone, time.Now(), authTime)
	}
	if timeOfVerification.Unix() > exp.Unix() {
		return nil, WarnAuthJWTTimeHasComeAndGone
	}

	a := authentication{
		ExternalID:                    sub,
		UserPrincipalName:             upn,
		JWTIssuer:                     iss,
		JWTAuthorizationParty:         azp,
		JWTIssuerGroups:               groups,
		JWTAuthorizationTime:          authTime,
		JWTExpiration:                 exp,
		LastVerificationTime:          timeOfVerification,
		LastVerificationOperationType: typeOfVerification,
	}
	return &a, nil
}
func SetAuthenticationInContext(ctx context.Context, a *authentication) context.Context {
	ctx = context.WithValue(ctx, KeyAuthentication{}, a)
	return ctx
}

//func getAuthenticationFromContext(ctx context.Context) (*authentication, error) {
//	a := ctx.Value(KeyAuthentication{}).(*authentication)
//	if a == nil {
//		return nil, fmt.Errorf("[WARNING] Token Not in Context")
//	}
//	return a, nil
//}

// authenticationJSON used to easily decode from Json with the necessary fields
type authenticationJSON struct {
	Sub      string   `json:"sub"`
	Upn      string   `json:"upn"`
	Iss      string   `json:"iss"`
	Azp      string   `json:"azp"`
	Groups   []string `json:"groups"`
	AuthTime int64    `json:"auth_time"`
	Exp      int64    `json:"exp"`
	//authenticationJSON.ClientId is useful for token introspective.
	ClientId string `json:"client_id"`
	Active   bool   `json:"active"`
}

func AuthGetExternalIDFromContext(ctx context.Context) (uuid.UUID, error) {
	authUUID := ctx.Value(KeyAuthentication{}).(*authentication).getExternalID()
	if authUUID == uuid.Nil {
		return uuid.Nil, WarnAuthIsNotInContext
	}
	return authUUID, nil

}
func AuthGetUserPrincipalName(ctx context.Context) (string, error) {
	authUPN := ctx.Value(KeyAuthentication{}).(*authentication).getUserPrincipalName()
	if authUPN == "" {
		return "", WarnAuthIsNotInContext
	}
	return authUPN, nil
}
func AuthGetJWTIssuer(ctx context.Context) (string, error) {
	authIssuer := ctx.Value(KeyAuthentication{}).(*authentication).getJWTIssuer()
	if authIssuer == "" {
		return "", WarnAuthIsNotInContext
	}
	return authIssuer, nil
}
func AuthGetJWTAuthorizationParty(ctx context.Context) (string, error) {
	authAuthParty := ctx.Value(KeyAuthentication{}).(*authentication).getJWTAuthorizationParty()
	if authAuthParty == "" {
		return "", WarnAuthIsNotInContext
	}
	return authAuthParty, nil
}
func AuthGetJWTIssuerGroups(ctx context.Context) ([]string, error) {
	authGroups := ctx.Value(KeyAuthentication{}).(*authentication).getJWTIssuerGroups()
	if len(authGroups) <= 0 {
		return nil, WarnAuthIsNotInContext
	}
	return authGroups, nil
}
func AuthGetJWTAuthorizationTime(ctx context.Context) (*time.Time, error) {
	authAuthTime := ctx.Value(KeyAuthentication{}).(*authentication).getJWTAuthorizationTime()
	if authAuthTime.IsZero() || authAuthTime == nil {
		return nil, WarnAuthIsNotInContext
	}
	return authAuthTime, nil
}

func AuthGetJWTExpiration(ctx context.Context) (*time.Time, error) {
	authExp := ctx.Value(KeyAuthentication{}).(*authentication).getJWTExpiration()
	if authExp.IsZero() || authExp == nil {
		return nil, WarnAuthIsNotInContext
	}
	return authExp, nil
}
func AuthGetLastTokenIntrospective(ctx context.Context) (*time.Time, error) {
	authLastVerification := ctx.Value(KeyAuthentication{}).(*authentication).getLastVerification()
	if authLastVerification.IsZero() || authLastVerification == nil {
		return nil, WarnAuthIsNotInContext
	}
	return authLastVerification, nil
}
func AuthGetLastTokenVerificationType(ctx context.Context) (string, error) {
	authLastTokenVerificationType := ctx.Value(KeyAuthentication{}).(*authentication).getLastTokenLocalVerificationType()
	if authLastTokenVerificationType == "" {
		return "", WarnAuthIsNotInContext
	}
	return authLastTokenVerificationType, nil
}
func mappingErrorOuputFormatting(key string) string {
	return fmt.Sprintf("[JWT-CLAIMS] [MAPPING]: [\"%v\"] Error while pulling from map from claims on JWT", key)

}
