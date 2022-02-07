package util

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const YES = 1
const NO = 0

var env string

type Enviroment struct {
	ENV string `json:"ENV"`
}

var TimeLayoutFull = "2006-01-02 15:04:05"
var TimeLayoutWithZone = "2006-01-02T15:04:05.000+08:00"
var TimeLayoutWithZoneSimple = "2006-01-02T15:04:05Z"
var TimeLayoutYmdHi = "2006-01-02 15:04"
var TimeLayoutYmd = "2006-01-02"

var location, _ = time.LoadLocation("Asia/Shanghai")

func IsExist(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() != 0
	case reflect.Bool:
		return value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return value.Float() != 0
	case reflect.Interface, reflect.Ptr:
		return !value.IsNil()
	}
	return false
}

// 字节的单位转换 保留两位小数
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzQWERTYUIOPASDFGHJKLZXCVBNM"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func GetRandomInt(l int)string {
	str := "0123456789"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func GetUserNum ()string{
	base := GetRandomInt(6)
	final := fmt.Sprintf("%d%s",121,base)
	return final
}

func GetGuestNum ()string{
	base := GetRandomInt(6)
	final := fmt.Sprintf("%d%s",000,base)
	return final
}
func CheckPhone(phone string) bool {
	if phone == ""{
		return true
	}
	regular := "^1([358][0-9]|4[579]|66|7[0135678]|9[89])[0-9]{8}$"
	// CheckMobileNum 手机号码的验证
	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}

func CheckMailFormat(email string) bool {

	mailCompile := regexp.MustCompile("^(.*)@(.*)\\.(.*)$")

	r := mailCompile.FindSubmatch([]byte(email))
	if len(r) != 4 {
		return false
	}
	return true
}

func NowUnixMs() int64 {
	return time.Now().In(location).UnixNano() / 1e6
}

func GetTodayStr() string {
	return time.Now().In(location).Format(TimeLayoutYmd)
}

func GetTodayTimeDetail()string{
	return time.Now().In(location).Format(TimeLayoutFull)
}

func TimeStampToString(timestampInt int64) string {
	return time.Unix(timestampInt, 0).In(location).Format(TimeLayoutFull)
}

func TimeStampToStringWithFormat(timestampInt int64, format string) string {
	return time.Unix(timestampInt, 0).In(location).Format(format)
}

func StringToTime(timeString string) int64 {
	timeStampInt := int64(0)
	timeLen := len(timeString)
	var timeStampTime time.Time
	var err error

	if timeLen == 29 {
		timeStampTime, err = time.ParseInLocation(TimeLayoutWithZone, timeString, time.Local)
	} else if timeLen == 20 {
		timeStampTime, err = time.ParseInLocation(TimeLayoutWithZoneSimple, timeString, time.Local)
	} else if timeLen == 16 {
		timeStampTime, err = time.ParseInLocation(TimeLayoutYmdHi, timeString, time.Local)
	} else if timeLen == 10 {
		timeStampTime, err = time.ParseInLocation(TimeLayoutYmd, timeString, time.Local)
	} else {
		timeStampTime, err = time.ParseInLocation(TimeLayoutFull, timeString, time.Local)
	}

	if err == nil {
		timeStampInt = timeStampTime.Unix() * 1e3
	}

	return timeStampInt
}

func StringTimeFormatConversion(originTimeString string, format string) string {
	originTime := StringToTime(originTimeString)
	return time.Unix(originTime/1000, 0).In(location).Format(format)
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func RandomId(length int64, strong int64) string {
	arrChars := [...]string{
		"9350742186",
		"clz2q4xh7bm8ay1ev0idk9jn3op5sft6rwgu",
		"dehOPQ9EKLbciCjk3WX6noplNAHs78JMI45RzrVfgYZ1Dxy2qFGa0STmtuvwUB",
		"lwsCjkci4zrVfAHvUgY5N+Z1D/2qxyR7)F$a0S@8LbGTp(mJMI3WX6*tuB-_dehOPQ9EKno=",
	}
	retStr := ""
	if int64(len(arrChars)) < strong || 0 >= length {
		return ""
	}
	chars := arrChars[strong]
	runeChars := []rune(chars)
	charsLen := int64(len(runeChars))

	timeNow := time.Now().Unix()
	for timeNow > 0 {
		start := timeNow % charsLen
		retStr += string(runeChars[start:(start + 1)])
		timeNow = timeNow / charsLen
	}
	retStrLength := int64(len(retStr))
	if length < (retStrLength + 4) {
		return ""
	}
	randomSeed := time.Now().UnixNano()
	rand.Seed(randomSeed)
	for length > int64(len(retStr)) {
		start := rand.Int63n(charsLen)
		retStr += string(runeChars[start:(start + 1)])
	}
	return retStr
}

//反射 性能低 慎用～
func InArray(value interface{}, arr interface{}) bool {
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(arr)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return true
			}
		}

	}

	return false
}

func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122 {
		strArry[0] -= 32
	}
	return string(strArry)
}

func GetEnv() string {
	if env == "" {
		load()
	}

	return strings.ToLower(env)
}

func DoNothing() {
	return
}

func load() {
	path := "./settings/config.json"

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var envStruct Enviroment

	err = json.Unmarshal(bs, &envStruct)
	if err != nil {
		panic(err)
	}

	env = envStruct.ENV
}

func GetSSOQuery(platformId, openApiSecret string) map[string]interface{} {
	ret := make(map[string]interface{})
	now := time.Now().Unix()
	nowstr := strconv.FormatInt(now, 10)
	var signature = Md5(fmt.Sprintf("%s_%s", nowstr, openApiSecret))
	ret["signature"] = signature
	ret["timestamp"] = nowstr
	ret["platform_id"] = platformId
	return ret
}

func Md5(str string) string {
	data := []byte(str)
	sum := md5.Sum(data)
	return fmt.Sprintf("%x", sum)
}

func SwitchStruct(oldStruct, newStruct interface{}) error {
	oldStructByte, jsonErr := json.Marshal(oldStruct)
	if jsonErr != nil {
		return jsonErr
	}
	jsonErr = json.Unmarshal(oldStructByte, newStruct)
	if jsonErr != nil {
		return jsonErr
	}
	return nil
}

func String2StringArrWithSeparate(str string, separate string, filterZero bool) []string {
	var resultStr = []string{}

	str = strings.Trim(str, " "+separate)
	if str != "" {
		strArr := strings.Split(str, separate)
		for _, strItem := range strArr {
			//过滤0值
			if filterZero && strItem == "" {
				continue
			}

			resultStr = append(resultStr, strItem)
		}
	}

	return resultStr
}

func String2IntArrWithSeparate(str string, separate string, filterZero bool) []int {
	var resultInt = []int{}

	str = strings.Trim(str, " "+separate)
	if str != "" {
		strArr := strings.Split(str, separate)
		for _, strItem := range strArr {
			//过滤0值
			if filterZero && strItem == "" {
				continue
			}

			intItem, err := strconv.Atoi(strItem)
			if err == nil {
				resultInt = append(resultInt, intItem)
			}
		}
	}

	return resultInt
}

func StrArr2StringWithSeparate(originStrArr []string, separate string, filterZero bool) string {

	var strArr []string
	for _, item := range originStrArr {

		//跳过0值
		if filterZero && item == "" {
			continue
		}

		strArr = append(strArr, item)
	}

	if len(strArr) > 0 {
		return separate + strings.Join(strArr, separate) + separate
	}

	return ""
}

func IntArr2StringWithSeparate(intArr []int, separate string, filterZero bool) string {

	var strArr []string
	for _, item := range intArr {

		//跳过0值
		if filterZero && item == 0 {
			continue
		}

		strArr = append(strArr, strconv.Itoa(item))
	}

	if len(strArr) > 0 {
		return separate + strings.Join(strArr, separate) + separate
	}

	return ""
}

func Int64ArrToStringWithSeparate(intArr []int64, sep string) string {
	if len(intArr) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, intItem := range intArr {
		buffer.WriteString(sep + strconv.Itoa(int(intItem)))
	}
	buffer.WriteString(sep)
	return buffer.String()
}

func IsEnvDev() bool {
	return GetEnv() == "dev"
}

func IsEnvTest() bool {
	return GetEnv() == "test"
}

func IsEnvProduction() bool {
	return GetEnv() == "online"
}
func IsDigitString(str string) bool{
	for _,x:=range []rune(str){
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}

func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

func IntList2Map(list []int) map[int]bool {
	m := make(map[int]bool)
	for _, item := range list {
		m[item] = true
	}
	return m
}

func StringList2Map(list []string) map[string]bool {
	m := make(map[string]bool)
	for _, item := range list {
		m[item] = true
	}
	return m
}

func InIntegerArray(value int, arr []int) bool {
	for _, elem := range arr {
		if value == elem {
			return true
		}
	}
	return false
}

func Ternary(a bool, b, c interface{}) interface{} {
	if a {
		return b
	}
	return c
}

func UniqueIntArr(originArr []int) []int {
	targetArr := []int{}
	checkMap := map[int]interface{}{}

	for _, item := range originArr {
		if _, ok := checkMap[item]; !ok {
			targetArr = append(targetArr, item)
			checkMap[item] = 0
		}
	}
	return targetArr
}

func UniqueStrArr(originArr []string) []string {
	targetArr := []string{}
	checkMap := map[string]interface{}{}

	for _, item := range originArr {
		if _, ok := checkMap[item]; !ok {
			targetArr = append(targetArr, item)
			checkMap[item] = 0
		}
	}
	return targetArr
}

func InStringArray(value string, arr []string) bool {
	for _, elem := range arr {
		if value == elem {
			return true
		}
	}
	return false
}

func IntArrayIndex(value int, arr []int) int{
	for k,v := range arr{
		if value == v{
			return k
		}
	}
	return 0
}

func GetOffsetAndLimit(page, size int) (int, int) {
	limit := 15
	if size > 0 {
		limit = size
	}
	offset := 0
	if page > 0 {
		offset = (page - 1) * limit
	}
	return offset, limit
}

func GetPreciseFloatOf(num float64, round int) float64 {
	v, _ := decimal.NewFromFloat(num).Round(int32(round)).Float64()
	return v
}

//模拟PHP中的array_diff 返回在 array1中 中但是不在其他 array 里的值。
func ArrayDiffInt(array1 []int, othersParams ...[]int) []int {
	if len(array1) == 0 {
		return []int{}
	}
	if len(array1) > 0 && len(othersParams) == 0 {
		return array1
	}
	var tmp = make(map[int]int, len(array1))
	for _, v := range array1 {
		tmp[v] = 1
	}
	for _, param := range othersParams {
		for _, arg := range param {
			if tmp[arg] != 0 {
				tmp[arg]++
			}
		}
	}
	var res = make([]int, 0, len(tmp))
	for k, v := range tmp {
		if v == 1 {
			res = append(res, k)
		}
	}
	return res
}

//俩Int数组交集
func IntersectionIntArray(s1, s2 []int) []int {
	m := make(map[int]int)
	for k := range s1 {
		m[s1[k]] += 1
	}
	var a = make([]int, 0)
	for k := range s2 {
		for key, value := range m {
			if key == s2[k] && value > 0 {
				m[k] -= 1
				a = append(a, key)
			}
		}
	}
	return a
}

//俩String数组交集
func IntersectionStrArray(s1, s2 []string) []string {
	m := make(map[string]int)
	for _, k := range s1 {
		m[k] = 0
	}
	var a []string
	for _, i := range s2 {
		if _, ok := m[i]; ok {
			a = append(a, i)
		}
	}
	return a
}

func TimeStampToDate(stamp int64) string {
	t := time.Unix(stamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

func TimeStampToDateYMD(stamp int64) string {
	t := time.Unix(stamp, 0)
	return t.Format("2006-01-02")
}

func AddToSet(list []int) mapset.Set {
	set := mapset.NewSet()
	for _, item := range list {
		set.Add(item)
	}

	return set
}

