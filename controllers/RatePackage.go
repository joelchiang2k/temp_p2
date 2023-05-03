package controllers

import (
	"encoding/json"
	"ex/part2/metrics"
	"ex/part2/models"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sugar_logger *zap.SugaredLogger
var atomic_level = zap.NewAtomicLevel()

type score_struct struct {
	Url                          string  `json:"Url"`
	NetScore                     float64 `json:"NetScore"`
	RampUpScore                  float64 `json:"RampUpScore"`
	CorrectnessScore             float64 `json:"CorrectnessScore"`
	BusFactorScore               float32 `json:"BusFactorScore"`
	ResponsivnessMaintainerScore float64 `json:"ResponsivnessMaintainerScore"`
	LicenseScore                 float64 `json:"LicenseScore"`
	CodeReviewScore              float64 `json:"CodeReviewScore"`
	DependencyScore              float64 `json:"DependencyScore"`
}

type package_info struct {
	Repository struct {
		URL string `json:"url"`
	} `json:"repository"`
}

func get_git_url(npm_url string) string {
	re_npm_url, _ := regexp.Compile("/\\w+")
	raw_module_name := re_npm_url.FindAllString(npm_url, -1)
	if len(raw_module_name) == 0 {
		log.Println("Error: The npmjs url provided is invalid!")
		return ""
	}
	module_name := raw_module_name[len(raw_module_name)-1]
	url := fmt.Sprintf("https://registry.npmjs.org/%s", module_name)

	res, err := http.Get(url)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var info package_info
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}
	re_git_url, _ := regexp.Compile("github.com/\\w+/\\w+")
	match_url := "https://" + re_git_url.FindString(info.Repository.URL)
	return match_url
}

func convert_byte_to_string(b []byte) string {
	str := ""
	for _, v := range b {
		if string(v) == "," {
			str += string(v) + " "
		} else {
			str += string(v)
		}
	}
	str += "\n"
	return str
}

func analyze_git(url string) score_struct {
	var result score_struct
	result.Url = url
	result.NetScore = 0.0
	result.RampUpScore = 0.0
	result.CorrectnessScore = 0.0
	result.BusFactorScore = 0.0
	result.ResponsivnessMaintainerScore = 0.0
	result.LicenseScore = 0.0
	result.DependencyScore = 0.0
	result.CodeReviewScore = 0.0
	if url == "" {
		log.Println("Error: The git url provided is invalid!")
		return result
	}

	var personal_token string
	godotenv.Load()
	personal_token = os.Getenv("GITHUB_TOKEN")

	sugar_logger.Info("Getting ramp-up score...")
	ramp_up_score_num, owner, repo := metrics.Ramp_up_score(personal_token, url)
	repo = strings.TrimSuffix(repo, ".git")
	sugar_logger.Info("Completed getting ramp-up score!")

	sugar_logger.Info("Getting correctness score...")
	correctness_score_num := metrics.Correctness_score(personal_token, owner, repo)
	sugar_logger.Info("Completed correctness score!")

	sugar_logger.Info("Getting responseviness score...")
	responseviness_score_num := metrics.Responseviness_score(personal_token, owner, repo)
	sugar_logger.Info("Completed getting responseviness score!")

	sugar_logger.Info("Getting bus factor score...")
	bus_factor_score_num := metrics.Bus_factor_score(personal_token, owner, repo)
	sugar_logger.Info("Completed getting bus factor score!")

	sugar_logger.Info("Getting license compatibility score...")
	license_compatibility_score_num := metrics.License_score(personal_token, owner, repo)
	sugar_logger.Info("Completed getting license compatibility score!")

	sugar_logger.Info("Getting code review score...")
	code_review_score_num := metrics.Code_review_metric(personal_token, owner, repo)
	sugar_logger.Info("Completed getting code review score!")

	sugar_logger.Info("Getting code review score...")
	dependency_score_num := metrics.Dependency_score(owner, repo)
	sugar_logger.Info("Completed getting code review score!")

	// Calculate net score
	net_score_raw := 0.15*ramp_up_score_num + 0.15*correctness_score_num + 0.15*float64(bus_factor_score_num) + 0.2*responseviness_score_num + 0.1*license_compatibility_score_num + 0.1*code_review_score_num + .15*dependency_score_num
	net_score, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", net_score_raw), 64)

	result.NetScore = net_score
	result.RampUpScore = ramp_up_score_num
	result.CorrectnessScore = correctness_score_num
	result.BusFactorScore = bus_factor_score_num
	result.ResponsivnessMaintainerScore = responseviness_score_num
	result.LicenseScore = license_compatibility_score_num
	result.CodeReviewScore = code_review_score_num
	result.DependencyScore = dependency_score_num
	return result
}

func Calc_score(url_file string) score_struct {
	/*file, _ := os.Open(url_file)
	defer file.Close()*/

	//var scores score_struct
	// Process inputted URLs
	/*scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "npmjs") {
			new_line := get_git_url(line)
			sugar_logger.Infof("URL: " + new_line)
			result := analyze_git(line, new_line)
			scores = result
		} else {
			if strings.Contains(line, ".git") {
				sugar_logger.Infof("URL: " + line)*/
	result := analyze_git(url_file)
	//scores = result
			/*} else {
				new_line := line + ".git"
				sugar_logger.Infof("URL: " + new_line)
				result := analyze_git(line, new_line)
				scores = result
			}
		}*/

	//}

	// sort URLs based on decending order of net score and output as NDJSON format
	/*b, err := json.Marshal(scores)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	fmt.Print(convert_byte_to_string(b))*/
	return result
}

func init() {
	// Set up logger
	encode_config := zap.NewProductionEncoderConfig()
	encode_config.EncodeTime = zapcore.ISO8601TimeEncoder
	log_file, _ := os.Create(os.Getenv("LOG_FILE"))
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encode_config),
		zapcore.AddSync(log_file), atomic_level)
	logger = zap.New(core, zap.AddCaller())
	sugar_logger = logger.Sugar()
	defer sugar_logger.Sync()
	log_level, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	atomic_level.SetLevel(zap.FatalLevel)
	switch log_level {
	case 1:
		atomic_level.SetLevel(zap.InfoLevel)
	case 2:
		atomic_level.SetLevel(zap.DebugLevel)
	default:
		atomic_level.SetLevel(zap.FatalLevel)
	}
}

/*func CalcScores(url string) score_struct {
	result := analyze_git(url)
	return result
}*/

func RatePackage(c *gin.Context) {
	var packageToRate models.PackageCreate
	if c.Param("{id}") == "/" {
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
	} else if err := models.DB.Where("id = ?", c.Param("{id}")).First(&packageToRate).Error; err != nil {
		c.JSON(404, "Package does not exist.")
	}
	var scores score_struct
	scores = Calc_score(packageToRate.URL) //return analyze git

	c.JSON(200, gin.H{

		"BusFactor":            scores.BusFactorScore,
		"Correctness":          scores.CorrectnessScore,
		"RampUp":               scores.RampUpScore,
		"ResponsiveMaintainer": scores.ResponsivnessMaintainerScore,
		"LicenseScore":         scores.LicenseScore,
		"GoodPinningPractice":  scores.DependencyScore,
	})
}
