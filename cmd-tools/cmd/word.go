package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/summerKK/go-code-snippet-library/cmd-tools/internal/word"
)

const (
	MODE_UPPER = iota + 1
	MODE_LOWER
	MODE_UNDERSCORE_TO_UPPER_CAMELCASE
	MODE_UNDERSCORE_TO_LOWER_CAMELCASE
	MODE_CAMELCASE_TO_UNDERSCORE
)

var desc = strings.Join([]string{
	"改子命令支持各种单词格式转换,模式如下: ",
	"1.全部单词转换为大写",
	"2.全部单词转换为小写",
	"3.下划线单词转换为大驼峰",
	"4.下划线单词转换为小驼峰",
	"5.驼峰单词转换为下划线",
}, "\n")

var str string
var mode int8

func init() {
	wordCmd.Flags().StringVarP(&str, "str", "s", "", "请输入单词内容")
	wordCmd.Flags().Int8VarP(&mode, "mode", "m", 0, "请输入单词转换模式")
}

var wordCmd = &cobra.Command{
	Use:   "word",
	Short: "单词格式转换",
	Long:  desc,
	Run: func(cmd *cobra.Command, args []string) {
		var content string
		switch mode {
		case MODE_UPPER:
			content = word.ToUpper(str)
		case MODE_LOWER:
			content = word.ToLower(str)
		case MODE_CAMELCASE_TO_UNDERSCORE:
			content = word.CamelCaseToUnderscore(str)
		case MODE_UNDERSCORE_TO_LOWER_CAMELCASE:
			content = word.UnderscoreToLowerCamelCase(str)
		case MODE_UNDERSCORE_TO_UPPER_CAMELCASE:
			content = word.UnderscoreToUpperCamelCase(str)
		default:
			log.Fatalf("暂时不支持改转换模式,请执行 help word 查看帮助文档")
		}

		log.Printf("输出结果: %s", content)
	},
}
