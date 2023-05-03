package controllers

import (
	"ex/part2/metrics"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// var logger2 *zap.Logger
// var sugar_logger2 *zap.SugaredLogger
// var atomic_level = zap.NewAtomicLevel()

// type score_struct struct {
// 	Url                          string  `json:"Url"`
// 	NetScore                     float64 `json:"NetScore"`
// 	RampUpScore                  float64 `json:"RampUpScore"`
// 	CorrectnessScore             float64 `json:"CorrectnessScore"`
// 	BusFactorScore               float32 `json:"BusFactorScore"`
// 	ResponsivnessMaintainerScore float64 `json:"ResponsivnessMaintainerScore"`
// 	LicenseScore                 float64 `json:"LicenseScore"`
// 	CodeReviewScore              float64 `json:"CodeReviewScore"`
// 	DependencyScore              float64 `json:"DependencyScore"`
// }

func get_score(url string) score_struct {
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

func CreatePackageRate(github_url string) score_struct {
	var result score_struct
	result = get_score(github_url)
	return result
}
