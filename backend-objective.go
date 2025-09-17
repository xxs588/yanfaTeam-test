package main

import( 
    "fmt"
	"errors"
	"math/rand"
	"regexp"//AI说是为了用正则表达式处理数据
	"time"
)

type verifyInfo struct {
	code       string    
	sendTime   time.Time // 发送时间
	dailyCount int       // 当日发送次数
	isUsed     bool      // 是否已使用
}

var phoneVerifyMap = make(map[string]*verifyInfo) //值为结构体的指针为了这一步		info = &verifyInfo{}phoneVerifyMap[phone] = info 存到info这个新指针变量
var chinaPhoneRegex  = regexp.MustCompile(`^1[345789]\d{9}$`) //声明正则 验证是否为中国手机号
func init() {
	rand.Seed(time.Now().UnixNano())
}//这个是初始化验证码种子，方便每次都会刷新验证码

func verificationCode(phone string)(string, error){
	if len(phone) != 11{
		return "" , errors.New("请输入正确的电话号码格式，是11位哦")
	}
	chinaPhoneRegex := `^1[345789]\d{9}$`
	if !regexp.MustCompile(chinaPhoneRegex).MatchString(phone) {
		return "", errors.New("请输入中国的电话号码哦")//运用正则表达式，用MatchString检验
	}
	// 生成6位随机验证码(包含前导0)
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code, nil //格式化输出如果数不足6位的话用0补成六位
}

func getCode(phone string) (string, error) {
	// 生成验证码(先验证手机号)
	code, err := verificationCode(phone)
	if err != nil {
		return "", err
	}

	now := time.Now()
	info, exists := phoneVerifyMap[phone]

	if !exists {
		info = &verifyInfo{}
		phoneVerifyMap[phone] = info
	}

	if info.dailyCount >= 5 {
		return "", errors.New("今日验证码获取次数已达上限(5次)")
	}

	if  now.Sub(info.sendTime) < time.Minute {
		remaining := time.Minute - now.Sub(info.sendTime)
		return "", fmt.Errorf("请%d秒后再获取验证码", int(remaining.Seconds()))
	}

	info.code = code
	info.sendTime = now
	info.dailyCount++
	info.isUsed = false

	return code, nil
}

func login(phone, inputCode string) error {
	info := phoneVerifyMap[phone]

	if time.Now().Sub(info.sendTime) > 5*time.Minute {
		return errors.New("验证码已过期(有效期5分钟)，请重新获取")
	}

	if inputCode != info.code {
		return errors.New("验证码错误")
	}

	// 验证码立即失效哦
	info.isUsed = true
	return nil
}


func main(){
	var phone string
	fmt.Println("请输入你的手机号码登录")
	fmt.Scanln(&phone)
	code, err := getCode(phone)
	    if err != nil{
		    fmt.Println("错误：",err)
		    return
 	    }
		fmt.Printf("验证码已发送：%s(有效期5分钟)\n", code)
	for{
		fmt.Println("1:输入验证码进行登录  2:获取验证码")
		var count string
		fmt.Scanln(&count)
        switch count{
			case "1" : 
				fmt.Print("请输入验证码：")
			    var verifyCode string 
            	fmt.Scanln(&verifyCode)
				err := login(phone,verifyCode)
			if err != nil {
				fmt.Println("登录失败：", err)
				continue
			}
			fmt.Println("登录成功！")
			return
		
		case "2" :
			newCode, err := getCode(phone)
			if err != nil {
				fmt.Println("获取失败：", err)
				continue
			}
			fmt.Printf("新验证码已发送：%s(有效期5分钟)\n", newCode)
		default:
			fmt.Println("未知操作，请输入1或2")
		}
	}
}
