# CPJudge

- 用於程式設計課程批改學生繳交檔案的程式

## 使用方法

### 前置作業

1. 將 moodle 下載的 zip 檔放在 `根目錄下`
2. 修改 `.env` 中的 `HWZip` 為 zip 檔名
3. 將要複製進去的檔案放在：`judgeEnv/share/copy` 下
4. 將測資放在：`judgeEnv/share/testcase` 下
   - 並以執行檔的名稱包一層資料夾 => `judgeEnv/share/testcase/執行檔名/測資`
   - e.g. `judgeEnv/share/testcase/hw0101/1`

### 執行

- `go run main.go`

## 輸出檔案結構

- zip 檔名資料夾下：
  - extract：解壓後學生繳交的檔案
  - output：自動批改後的檔案
