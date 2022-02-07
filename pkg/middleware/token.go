package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)
const (
	ExpireTime = 720 * time.Hour
	GuestExpire = 2*time.Hour
	secret     = "morethanthis"
	RoleUser   = 1
	RoleGuest  = 2
	RoleAdmin  = int8(127)
)
var jwtSecret = []byte(secret)
type Claims struct {
	UserId   int
	UserName string
	UserRole int8
	UserNum string
	jwt.StandardClaims
}

func GenerateToken(uid int, role int8,userName ,userNum,id string,expireSec time.Duration) (string, error) {

	nowTime := time.Now()
	expireTime := nowTime.Add(expireSec)

	claims := Claims{
		uid,
		userName,
		role,
		userNum,
		jwt.StandardClaims{
			Id: id,
			ExpiresAt: expireTime.Unix(),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)


	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}


//手工jwt
//func NewToken(uid string, role int) (string, error) {
//	token := NewJWT(uid, role, time.Now().Add(expireTime).Unix())
//	s, err := token.SignedString(secret)
//	if err != nil {
//		return "", err
//	}
//	return s, nil
//}
//var (
//	jwtHeader = JWTHeader{"HS256", "JWT"}
//)
//
//// JWT JWT
//type JWT struct {
//	Header  *JWTHeader
//	Payload *JWTPayload
//}
//
//// JWTHeader token的header部分
//type JWTHeader struct {
//	Algorithm string `json:"alg"`
//	TokenType string `json:"typ"`
//}
//
//// JWTPayload token的payload部分,由用户id和过期时间组成
//type JWTPayload struct {
//	UserID     string `json:"uid"`
//	ExpireTime int64  `json:"exp"`
//	Role       int    `json:"role"`
//}
//
//// NewJWT NewJWT
//func NewJWT(uid string, role int, exp int64) *JWT {
//	return &JWT{
//		Header: &jwtHeader,
//		Payload: &JWTPayload{
//			UserID:     uid,
//			ExpireTime: exp,
//			Role:       role,
//		},
//	}
//}
//
//// SignedString 使用HS256加密标头返回token字符串
//func (t *JWT) SignedString(secret string) (string, error) {
//	header, err := encoding.Base64EncodeJSON(base64.RawURLEncoding, t.Header)
//	if err != nil {
//		return "", err
//	}
//	payload, err := encoding.Base64EncodeJSON(base64.RawURLEncoding, t.Payload)
//	if err != nil {
//		return "", err
//	}
//	signature := encoding.Base64Encode(base64.RawURLEncoding,
//		encoding.HMAC(sha256.New, []byte(string(header)+"."+string(payload)), []byte(secret)))
//
//	return strings.Join([]string{string(header), string(payload), string(signature)}, "."), nil
//}
//
////JWTParse 解析token字符串并验证签名，如果一切正常，将设置有效负载
//func (t *JWT) JWTParse(token, secret string) error {
//	s := strings.Split(token, ".")
//	if len(s) != 3 {
//		return fmt.Errorf("token format error")
//	}
//	header, payload, signature := s[0], s[1], s[2]
//	// 验证签名
//	if !encoding.ConstTimeEqual([]byte(signature), encoding.Base64Encode(base64.RawURLEncoding,
//		encoding.HMAC(sha256.New, []byte(header+"."+payload), []byte(secret)))) {
//		return fmt.Errorf("token signature error")
//	}
//	// decode header
//	if err := encoding.Base64DecodeJSON(base64.RawURLEncoding, []byte(header), &t.Header); err != nil {
//		return err
//	}
//	// decode payload
//	return encoding.Base64DecodeJSON(base64.RawURLEncoding, []byte(payload), &t.Payload)
//}
//
////Expired checks if access token is expired
//func (t *JWT) Expired() bool {
//	return time.Now().After(time.Unix(t.Payload.ExpireTime, 0))
//}
//
//
//func GetTokenClaims(claims string) (res *JWTPayload){
//	binary ,_:= json.Marshal(claims)
//	json.Unmarshal(binary,&res)
//	return
//}