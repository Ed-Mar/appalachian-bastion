package handlers

import (
	"backend/internal"
	"backend/internal/authentication/config"
	"backend/internal/helper"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (ha *ServiceHandler) AuthenticationMiddlewareViaTokenIntrospective(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Header[KEY_IN_HEADER_WITH_TOKEN] == nil {
			ha.ServiceLogger.Println(WarnAuthNoAccessTokenInHeader)
			http.Error(rw, WarnAuthNoAccessTokenInHeader.Error(), http.StatusBadRequest)
			return
		}

		accessToken := r.Header.Get(KEY_IN_HEADER_WITH_TOKEN)
		const introspectURL = "http://keycloak.test/realms/gatehouse/protocol/openid-connect/token/introspect"

		const filePath = "../internal/authentication/config"
		const fileName = "devKeycloakTokenIntrospectiveClientInfo"
		const fileType = "env"

		type Payload struct {
			Token string `json:"token"`
		}

		formData := &Payload{
			Token: accessToken,
		}

		payload := url.Values{
			"token": {formData.Token},
		}
		req, err := http.NewRequest("POST", introspectURL, strings.NewReader(payload.Encode()))
		if err != nil {
			err = fmt.Errorf("%v | %v", ErrAuthTokenIntrospectiveRequestCreation, err)
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			err = internal.ToJSON(GenericError{Message: err.Error()}, rw)
			return
		}
		clientCreds, err := config.GetTokenIntrospectiveClientBasicInfo(filePath, fileName, fileType)
		if err != nil {
			err = fmt.Errorf("%v | %v", ErrAuthLoadingTokenIntrospectiveClientCreds, err)
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			err = internal.ToJSON(GenericError{Message: err.Error()}, rw)
			return
		}
		req.Header.Add("Authorization", *clientCreds)
		req.Header.Set("Accept-Encoding", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(payload.Encode())))
		client := &http.Client{
			Timeout: time.Second * 11,
		}
		if LOG_MODE {
			reqDump, err := httputil.DumpRequestOut(req, true)
			if err != nil {
				ha.ServiceLogger.Printf("[VERBOSE LOGGING] [ERROR] ReqDump Error: %#v", err)
			}
			ha.ServiceLogger.Printf("[VERBOSE LOGGING] TOKEN Introspect REQUEST:\n%s", string(reqDump))
		}
		// Sending the request for the introspective.
		response, err := client.Do(req)
		if err != nil {
			err = fmt.Errorf("%v | %v", ErrAuthTokenIntrospectiveSendingRequest, err)
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				ha.ServiceLogger.Printf("[ERROR] Closing Body?! ", err)
				return
			}
		}(response.Body)

		if LOG_MODE {
			respDump, err := httputil.DumpResponse(response, true)
			if err != nil {
				ha.ServiceLogger.Printf("[VERBOSE LOGGING] [ERROR] Response Error: %#v", err)
			}
			ha.ServiceLogger.Printf("[VERBOSE LOGGING] TOKEN Introspect RESPONSE:\n%s", string(respDump))
		}
		if response.Status != "200 OK" {
			err = fmt.Errorf("%v %v", WarnAuthTokenIntrospectiveRequestStatusNotOK, fmt.Errorf(" | Status: %v Body: %v", response.Status, response.Body))
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		tokenIntrospectiveRespObj := &authenticationJSON{}

		err = internal.FromJSON(tokenIntrospectiveRespObj, response.Body)
		if err != nil {
			err = fmt.Errorf("%v | %v ", ErrAuthTokenIntrospectiveResponseToJSON, err)
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if !tokenIntrospectiveRespObj.Active {
			ha.ServiceLogger.Println(WarnAuthTokenIntrospectiveTokenInvalidNotActive)
			http.Error(rw, WarnAuthTokenIntrospectiveTokenInvalidNotActive.Error(), http.StatusForbidden)
			return
		}
		ha.ServiceLogger.Println("-----------------------------------------")
		ha.ServiceLogger.Println(tokenIntrospectiveRespObj)
		ha.ServiceLogger.Println("-----------------------------------------")

		subString := tokenIntrospectiveRespObj.Sub
		//sid Conversion from string into uuid
		subUUID, err := uuid.FromString(subString)
		if err != nil {
			err = fmt.Errorf("%v | [MAPPING] [\"sub\"] Error Converting String to UUID | %v", ErrAuthTokenIntrospectiveResponseToJSON, err)
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// upn---
		upn := tokenIntrospectiveRespObj.Upn
		// iss--
		iss := tokenIntrospectiveRespObj.Iss
		// azp--
		azp := tokenIntrospectiveRespObj.Azp
		// groups--
		groups := tokenIntrospectiveRespObj.Groups
		// auth_time--
		authTimeInt64 := tokenIntrospectiveRespObj.AuthTime
		// auth_time Conversion to *time
		authTime := time.Unix(authTimeInt64, 0)
		//Sanctity Check I am checking if the time could make any sense by comparing it to the current time +/- a month (604800 seconds in a month)
		//I also just noticed this will catch if the time is Zero which is nice.
		if !(authTime.Unix() > (time.Now().Unix()-604800) && (authTime.Unix() < (time.Now().Unix() + 604800))) {
			err = fmt.Errorf("%v | [MAPPING] [\"auth_time\"] Error while converting it to time. The time doesn't make sense so the conversion is most likly off | Recived: %v Current Time: %v", ErrAuthTokenIntrospectiveResponseToJSON, authTime, time.Now())
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// expFloat64 --
		expInt64 := tokenIntrospectiveRespObj.Exp

		// auth_time Conversion to *time
		expTime := time.Unix(expInt64, 0)
		//Sanctity Check I am checking if the time could make any sense by comparing it to the current time +/- a month (604800 seconds in a month)
		if !(expTime.Unix() > (time.Now().Unix()-604800) && (expTime.Unix() < (time.Now().Unix() + 604800))) {
			err = fmt.Errorf("%v | [MAPPING] [\"exp\"] Error while converting it to time. The time doesn't make sense so the conversion is most likly off  | Recived: %v Current Time: %v", ErrAuthTokenIntrospectiveResponseToJSON, authTime, time.Now())
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		//Getting the Time of the token introspective
		timeOfTokenIntrospection := time.Time{}
		//Checks the head to see if the response has the Date Key
		if response.Header["Date"] == nil {
			//If not then just uses the the current time(should be at least close)
			timeOfTokenIntrospection = time.Now()
		} else {
			//IF the date key is in the header then try to convert that value into a time.Time for go useage
			timeOfTokenIntrospection, err = time.Parse(time.RFC1123, response.Header.Get("Date"))
			// If that fails then just fall back to use the current time.
			if err != nil {
				timeOfTokenIntrospection = time.Now()
			}
		}
		auth, err := newAuthentication(subUUID, upn, iss, azp, groups, &authTime, &expTime, &timeOfTokenIntrospection, &tokenIntrospectionAuthenticationVerificationType{})
		if err != nil {
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		//log.Println(auth)
		ctx := SetAuthenticationInContext(r.Context(), auth)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})
}

func (ha *ServiceHandler) AuthenticationMiddlewareViaLocalVerification(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Header[KEY_IN_HEADER_WITH_TOKEN] == nil {
			ha.ServiceLogger.Println(WarnAuthNoAccessTokenInHeader)
			http.Error(rw, WarnAuthNoAccessTokenInHeader.Error(), http.StatusBadRequest)
			return
		}
		accessToken := r.Header.Get(KEY_IN_HEADER_WITH_TOKEN)

		const filePath = "../internal/authentication/config"
		const fileName = "localDevPublicRSASecret"
		const fileType = "env"

		rsaSecret, err := config.GetPublicRSASecret(filePath, fileName, fileType)
		if err != nil {
			err = fmt.Errorf("%v | %v", ErrAuthPublicRSAKeyLoading, err)
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		// PEM formatting
		rsaSecret = "-----BEGIN CERTIFICATE-----\n" +
			rsaSecret +
			"\n-----END CERTIFICATE-----"
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(rsaSecret))
		if err != nil {
			err = fmt.Errorf("%v | %v", ErrAuthParsingPublicRSA, err)
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		timeOfVerification := time.Time{}

		jwtToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			//Validate the security algorithm:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("%v | Token was Signed with unexpected signing method [%v] while the expected being [rsa]", WarnAuthTokenSignedWithNonSupportedAlg, token.Header["alg"])
			}
			timeOfVerification = time.Now()
			return key, nil
		})
		if !jwtToken.Valid {
			http.Error(rw, WarnAuthInvalidToken.Error(), http.StatusForbidden)
			return
		}

		if errors.Is(err, jwt.ErrTokenMalformed) {
			err = fmt.Errorf("%v | %v", WarnAuthInvalidTokenMalformed, err)
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusForbidden)
			return
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			err = fmt.Errorf("%v | %v", WarnAuthInvalidTokenExpired, err)
			ha.ServiceLogger.Println("[WARNING] [WARNING] [WARNING] Post-Date JWT Attempt:")
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusForbidden)
			return
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			err = fmt.Errorf("%v | %v", WarnAuthInvalidTokenPreDated, err)
			ha.ServiceLogger.Println("[WARNING] [WARNING] [WARNING]  Pre-Date JWT Attempt:")
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusForbidden)
			return
		} else if err != nil {
			err = fmt.Errorf("%v | %v", WarnAuthInvalidToken, err)
			ha.ServiceLogger.Println(err)
			http.Error(rw, WarnAuthInvalidToken.Error(), http.StatusForbidden)
			return
		}
		jwtMappedClaims, _ := jwtToken.Claims.(jwt.MapClaims)

		subString, ok := jwtMappedClaims["sub"].(string)
		if !ok {
			err = fmt.Errorf("%v %v", ErrAuthLocalTokenVerificationBase, mappingErrorOuputFormatting("sub"))
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return

		}
		//sid Conversion from string into uuid
		subUUID, err := uuid.FromString(subString)
		if err != nil {
			err = fmt.Errorf("%v [MAPPING] [\"sub\"] Error Converting String to UUID | %v", ErrAuthLocalTokenVerificationBase, err)
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// upn---
		upn, ok := jwtMappedClaims["upn"].(string)
		if !ok {
			err = fmt.Errorf("%v %v", ErrAuthLocalTokenVerificationBase, mappingErrorOuputFormatting("upn"))
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		// iss--
		iss, ok := jwtMappedClaims["iss"].(string)
		if !ok {
			err = fmt.Errorf("%v %v", ErrAuthLocalTokenVerificationBase, mappingErrorOuputFormatting("iss"))
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		// azp--
		azp, ok := jwtMappedClaims["azp"].(string)
		if !ok {
			err = fmt.Errorf("%v %v", ErrAuthLocalTokenVerificationBase, mappingErrorOuputFormatting("azp"))
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		// groups--
		groupsInterface, ok := jwtMappedClaims["groups"].([]interface{})
		if !ok {
			err = fmt.Errorf("%v %v", ErrAuthLocalTokenVerificationBase, mappingErrorOuputFormatting("groups"))
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		//groups Conversion in to String Slice
		groups := make([]string, len(groupsInterface))
		for i, v := range groupsInterface {
			groups[i], ok = v.(string)
			if !ok {
				//(I would tell what it is but if it wont become a string there is not hope in logging it)
				err = fmt.Errorf("%v [MAPPING] [\"sid\"] Error Converting items in groups to string  | at index %v | %#v", ErrAuthLocalTokenVerificationBase, i, groups)
				ha.ServiceLogger.Println(err)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		// auth_time--
		authTimeFloat64, ok := jwtMappedClaims["auth_time"].(float64)
		if !ok {
			err = fmt.Errorf("%v %v", ErrAuthLocalTokenVerificationBase, mappingErrorOuputFormatting("auth_time"))
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		// auth_time Conversion to *time
		authTime := helper.FloatToUnixTime(authTimeFloat64)
		//Sanctity Check I am checking if the time could make any sense by comparing it to the current time +/- a month (604800 seconds in a month)
		if !(authTime.Unix() > (time.Now().Unix()-604800) && (authTime.Unix() < (time.Now().Unix() + 604800))) {
			err = fmt.Errorf("%v [MAPPING] [\"auth_time\"] Error while converting it to time. The time doesn't make sense so the conversion is most likly off | %#v ", ErrAuthLocalTokenVerificationBase, authTime)
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// expFloat64 --
		expFloat64, ok := jwtMappedClaims["exp"].(float64)
		if !ok {
			err = fmt.Errorf("%v %v", ErrAuthLocalTokenVerificationBase, mappingErrorOuputFormatting("exp"))
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		// auth_time Conversion to *time
		expTime := helper.FloatToUnixTime(expFloat64)
		//Sanctity Check I am checking if the time could make any sense by comparing it to the current time +/- a month (604800 seconds in a month)
		if !(expTime.Unix() > (time.Now().Unix()-604800) && (expTime.Unix() < (time.Now().Unix() + 604800))) {
			err = fmt.Errorf("%v [MAPPING] [\"exp\"] Error while converting it to time. The time doesn't make sense so the conversion is most likly off | %#v", ErrAuthLocalTokenVerificationBase, expTime)
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		//--Encase of shitty Debugging
		//log.Printf("%#v \n %#v \n %#v \n %#v \n %#v \n %#v \n %#v \n %#v \n", subString, upn, iss, azp, groups, authTime.String(), expTime.String(), TimeOfVerification.String())

		auth, err := newAuthentication(subUUID, upn, iss, azp, groups, &authTime, &expTime, &timeOfVerification, &localAuthenticationVerificationType{})
		if err != nil {
			ha.ServiceLogger.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx := SetAuthenticationInContext(r.Context(), auth)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})
}
