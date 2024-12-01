package main

import (
	"fmt"
	"github.com/logxxx/utils/fileutil"
	"image"
	"image/color"
	_ "image/jpeg" // 导入jpeg支持
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	_, err := run("test3.jpg")
	if err != nil {
		panic(err)
	}
}

func run(filePath string) (respObj interface{}, e error) {
	// 示例图片路径、要获取颜的x、y坐标

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()

	maxX := bounds.Max.X
	maxY := bounds.Max.Y

	round := 0

	step := 10

	colorMatrix := [][]string{}
	for i := 0; i < maxX; i += step {
		colorColumn := []string{}
		for j := 0; j < maxY; j += step {
			rectColors := map[string]int{}
			for k := i; k <= i+step; k++ {
				for l := j; l <= j+step; l++ {
					orig, respCode, resp, err := GetPixelColor3(img, k, l)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}
					_ = orig
					_ = resp
					_ = respCode
					//fmt.Printf("The color at (%d, %d) is: %v => [%v]%v\n", k, l, orig, respCode, resp)
					rectColors[respCode]++
				}
			}
			rectMainColorCount := 0
			rectMainColor := ""
			for c, count := range rectColors {
				if count > rectMainColorCount {
					rectMainColorCount = count
					rectMainColor = c
				}
			}
			round++
			//log.Printf("rect%v color:%v", round, rectMainColor)
			colorColumn = append(colorColumn, rectMainColor)
		}
		colorMatrix = append(colorMatrix, colorColumn)
	}

	//for i := 0; i < len(colorMatrix); i++ {
	//	col := colorMatrix[i]
	//	for j := 0; j < len(col); j++ {
	//		log.Printf("rect(%v,%v):%v", i+1, j+1, ConvColorCode2Name(col[j]))
	//	}
	//}

	//draw start
	reportImg := image.NewNRGBA(image.Rect(0, 0, maxX, maxY))

	for i := 0; i < len(colorMatrix); i++ {
		for j := 0; j < len(colorMatrix[i]); j++ {

			r, g, b := ConvColorCode2RGB(colorMatrix[i][j])

			// 定义一个颜色
			fillColor := color.RGBA{r, g, b, 255}

			log.Printf("(%v,%v)~(%v,%v): %v", i*step, j*step, i*step+step, j*step+step, fillColor)

			// 给一个矩形区域填充颜色
			drawRectangle(reportImg, image.Rect(i*step, j*step, i*step+step, j*step+step), fillColor)
		}
	}

	// 保存图像为PNG格式
	outFile, _ := os.Create("test1_output.png")
	defer outFile.Close()
	png.Encode(outFile, reportImg)
	//draw end

	reportCode := []string{}

	for i := 0; i < len(colorMatrix); i++ {
		reportCode = append(reportCode, "")
		for j := 0; j < len(colorMatrix[i]); j++ {
			reportCode[i] = reportCode[i] + " " + colorMatrix[i][j]
		}
	}
	fileutil.WriteToFile([]byte(strings.Join(reportCode, "\n")), "reportCode.txt")

	return
}

// 画一个填充指定颜色的矩形
func drawRectangle(img *image.NRGBA, rect image.Rectangle, col color.Color) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			img.Set(x, y, col)
		}
	}
}

func GetPixelColor3(img image.Image, x, y int) (orig string, respCode, resp string, err error) {
	color := img.At(x, y)
	r, g, b, _ := color.RGBA()

	R := uint8(r)
	G := uint8(g)
	B := uint8(b)

	orig = fmt.Sprintf("%02X%02X%02X", R, G, B)

	R = fixColor(R)
	G = fixColor(G)
	B = fixColor(B)

	respCode = fmt.Sprintf("%02X%02X%02X", R, G, B)

	resp = getColorName(R, G, B)

	return
}

var (
	colorModel = map[string]string{
		"00FF00": "绿",
		"00FFFF": "青",
		"0000FF": "蓝",
		"800000": "棕",
		"808000": "橄榄",
		"800080": "紫",
		"80FFFF": "浅青",
		"808080": "灰",
		"FFFF00": "黄",
		"FFFFFF": "白",
		"FF0000": "红",
		"FF00FF": "品红",
		"000000": "黑",
		"000080": "深蓝",
		"FFFF80": "黄",
		"008080": "水绿",
		"FF8080": "胭脂粉",
		"008000": "绿",
		"FF8000": "橙",
		"8080FF": "蓝紫",
		"80FF80": "翠绿",
		"FF80FF": "紫粉",
		"80FF00": "亮绿",
		"0080FF": "蓝",

		"C00000": "品红",
		"C0C000": "黄绿",
		"C0C0C0": "银",
		"C000C0": "大紫",
		"0000C0": "深蓝",
		"00C0C0": "绿蓝",
		"00C000": "大绿",

		"80C000": "草绿",
		"8000C0": "紫",
		"C08000": "屎黄",
		"C0C080": "土绿",
		"80C0C0": "绿蓝",
		"C08080": "胭脂",
		"8080C0": "紫",
		"C080C0": "紫",
		"0080C0": "蓝",
		"80C080": "绿",
		"C00080": "紫",

		"C0FF00": "翠绿",
		"FFC000": "橙黄",
		"FF00C0": "艳紫",
		"C0FFFF": "淡天蓝",
		"C000FF": "深紫",
		"C0C0FF": "浅紫",
		"FFFFC0": "鹅黄",
		"FFC0C0": "胭脂2",
		"C0FFC0": "淡绿",
		"FFC0FF": "淡粉",

		"C080FF": "浅紫",
		"C0FF80": "草绿",
		"80FFC0": "绿",
		"FF80C0": "粉红",
		"FFC080": "黄",
		"80C0FF": "蓝",
	}
)

func ConvColorCode2RGB(req string) (r uint8, g uint8, b uint8) {
	//log.Printf("ConvColorCode2RGB r:%v g:%v b:%v", req[:2], req[2:4], req[4:])
	r = convColorChannelCode2Num(req[:2])
	g = convColorChannelCode2Num(req[2:4])
	b = convColorChannelCode2Num(req[4:])
	return
}

func convColorChannelCode2Num(req string) (resp uint8) {
	num, err := strconv.ParseInt(req, 16, 64)
	if err != nil {
		return
	}
	if num < 0 || num > 255 {
		panic(fmt.Sprintf("数值超出uint8范围:%v", num))
		return
	}

	return uint8(num)

}

func ConvColorCode2Name(req string) string {
	return colorModel[req]
}

func getColorName(r, g, b uint8) (name string) {

	code := fmt.Sprintf("%02X%02X%02X", r, g, b)
	name = colorModel[code]
	if name == "" {
		panic(fmt.Sprintf("code not found:%v", code))
	}
	return name
}

func fixColor(i uint8) (resp uint8) {

	//levels := []uint8{32, 64, 96, 128, 160, 192}

	if i < 64 {
		return 0
	}

	if i < 128 {
		return 128
	}

	if i < 192 {
		return 192
	}

	return 255

}

func getPixelColor(filePath string, x, y int) (color.Color, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	if x < bounds.Min.X || x > bounds.Max.X || y < bounds.Min.Y || y > bounds.Max.Y {
		return nil, fmt.Errorf("pixel coordinates out of bounds")
	}

	return img.At(x, y), nil
}
