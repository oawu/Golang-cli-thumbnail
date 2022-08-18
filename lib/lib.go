package lib

import (
  "os"
  "fmt"
  "strings"
)

func Args(arguments []string, keys ...string) string {
  keyword := []rune("-")[0]

  key := ""
  args := &map[string]string{}

  for _, argument := range arguments {
    char1 := []rune(argument)[0]
    if char1 == keyword {
      key = argument
      (*args)[key] = ""
    } else if key != "" {
      if (*args)[key] == "" {
        (*args)[key] = argument
      } else {
        (*args)[key] = (*args)[key] + " " + argument
      }
    }
  }

  for _, key := range keys {
    v, ok := (*args)[key]
    if ok {
      return v
    }
  }

  return ""
}

func StrWidth(str string) uint64 {
  w := uint64(0)
  for _, c := range []rune(str) {
    s := fmt.Sprintf("%c", c)
    if len(s) == 1 {
      w += 1
    } else {
      if (s == "➜") {
        w += 1
      } else {
        w += 2
      }
    }
  }
  return w
}

func inArray(strs []string, s string) bool {
  for _, str := range strs {
    if str == s {
      return true
    }
  }
  return false
}

func FileSize(path string) uint64 {
  reader, e := os.Open(path)
  if e != nil {
    return 0
  }
  defer reader.Close()

  stat, e := reader.Stat()
  if e != nil {
    return 0
  }
  return uint64(stat.Size())
}

func ShowHelpIosAppIcon(errs ...error) {
  if len(errs) > 0 && errs[0] != nil && errs[0].Error() != "" {
    str := errs[0].Error()
    padding := uint64(4)
    width := StrWidth(str) + padding * 2

    fmt.Println("")
    fmt.Printf("  \x1b[48;5;1m%s\x1b[0m\n", strings.Repeat(" ", int(width)))
    fmt.Printf("  \x1b[48;5;1m%s\x1b[0m\n", fmt.Sprintf("%s%s%s", strings.Repeat(" ", int(padding)), str, strings.Repeat(" ", int(padding))))
    fmt.Printf("  \x1b[48;5;1m%s\x1b[0m\n", strings.Repeat(" ", int(width)))
  }

  fmt.Println("")
  fmt.Println("  § 製作 iOS 的 APP 圖示")
  fmt.Println("")
  fmt.Println("    使用方式：")
  fmt.Println("")
  fmt.Println("        oago ios-app-icon -P <name>")
  fmt.Println("")
  fmt.Println("    參數說明：")
  fmt.Println("")
  fmt.Println("        1. -P <name>, --pic <name>")
  fmt.Println("          此參數代表要縮圖的檔案，請於參數後填寫檔案名稱或檔案位置，")
  fmt.Println("          檔案若不是絕對位置則會以終端機當時位置的相對檔名尋找。")
  fmt.Println("          目前僅能處理 .jpg 與 .png 類型的圖片。")
  fmt.Println("")
  fmt.Println("    注意事項：")
  fmt.Println("")
  fmt.Println("        縮圖功能必須要使用 [-P | --pic] 參數。")
  fmt.Println("")
}

func ShowHelpThumbnail(errs ...error) {
  if len(errs) > 0 && errs[0] != nil && errs[0].Error() != "" {
    str := errs[0].Error()
    padding := uint64(4)
    width := StrWidth(str) + padding * 2

    fmt.Println("")
    fmt.Printf("  \x1b[48;5;1m%s\x1b[0m\n", strings.Repeat(" ", int(width)))
    fmt.Printf("  \x1b[48;5;1m%s\x1b[0m\n", fmt.Sprintf("%s%s%s", strings.Repeat(" ", int(padding)), str, strings.Repeat(" ", int(padding))))
    fmt.Printf("  \x1b[48;5;1m%s\x1b[0m\n", strings.Repeat(" ", int(width)))
  }

  fmt.Println("")
  fmt.Println("  § 縮圖功能")
  fmt.Println("")
  fmt.Println("    使用方式：")
  fmt.Println("")
  fmt.Println("        oago thumbnail -P <name> [-W <width>] | [-H <height>] | [-S <size>] [-M | --multiple <1,2,3...>]")
  fmt.Println("")
  fmt.Println("    參數說明：")
  fmt.Println("")
  fmt.Println("        1. -P <name>, --pic <name>")
  fmt.Println("          此參數代表要縮圖的檔案，請於參數後填寫檔案名稱或檔案位置，")
  fmt.Println("          檔案若不是絕對位置則會以終端機當時位置的相對檔名尋找。")
  fmt.Println("          目前僅能處理 .jpg 與 .png 類型的圖片。")
  fmt.Println("")
  fmt.Println("        2. -W <width>, --width <width>")
  fmt.Println("          此參數代表縮圖後 @1x 尺寸下的寬度，")
  fmt.Println("          若只設定寬度而未設定高度，則會依據原圖寬高比例進行縮圖。")
  fmt.Println("")
  fmt.Println("        3. -H <height>, --height <height>")
  fmt.Println("          此參數代表縮圖後 @1x 尺寸下的高度，")
  fmt.Println("          若只設定高度而未設定寬度，則會依據原圖寬高比例進行縮圖。")
  fmt.Println("")
  fmt.Println("        4. -S <size>, --size <size>")
  fmt.Println("          此參數代表縮圖後 @1x 尺寸下的寬度與高度。")
  fmt.Println("          因為此參數是同時設定寬度與高度，故使用此參數所輸出後的縮圖皆為方形，")
  fmt.Println("          若指令中同時有 [-W | --width] 或 [-H | --height] 的有效參數設定，則此參數則為無效。")
  fmt.Println("")
  fmt.Println("        5. -M <1,2,3...> | --multiple <1,2,3...>")
  fmt.Println("          需要縮圖的倍數，可以不用填寫，預設為 1,2,3，")
  fmt.Println("          若使用 -M 1,2 則會分別產出 @1x 與 @2x 的縮圖。")
  fmt.Println("")
  fmt.Println("    注意事項：")
  fmt.Println("")
  fmt.Println("        縮圖功能必須要使用 [-P | --pic] 參數，")
  fmt.Println("        同時至少要有 [-W | --width] 或 [-H | --height] 或 [-S | --size] 一項參數。")
  fmt.Println("")
}

func ShowHelp() {
  fmt.Println("")
  fmt.Println("  § 歡迎使用 OAGO 工具")
  fmt.Println("")
  fmt.Println("    使用方式：")
  fmt.Println("")
  fmt.Println("        oago <指令> [參數]")
  fmt.Println("")
  fmt.Println("    指令如下：")
  fmt.Println("")
  fmt.Println("        thumbnail        縮圖")
  fmt.Println("        ios-app-icon     製作 iOS 的 APP 圖示")
  fmt.Println("")
  fmt.Println("    注意事項：")
  fmt.Println("")
  fmt.Println("        可以使用 \"oago help <指令>\" 來查詢其參數說明。")
  fmt.Println("") 
}