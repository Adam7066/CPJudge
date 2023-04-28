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

## 功能展示

### 啟動設定

可設定執行時間、平行處理數量、是否自動批改

![main](https://i.imgur.com/BsKayU2.png)

### 檢視編譯狀況

![out](https://i.imgur.com/9n0sg6t.png)
![error](https://i.imgur.com/KevG8AL.png)

### 顯示執行結果與正確答案比對

可以切換顯示學生執行結果及與正確答案的差異

![hw0201/1](https://i.imgur.com/QIpwmlf.png)
![hw0201/1/diff](https://i.imgur.com/VVTqGO5.png)

### 超時終止

![timeout](https://i.imgur.com/6OA8q3Z.png)

### 全螢幕顯示

套用 less 指令顯示詳細檔案內容

![full screen](https://i.imgur.com/fqrY6Uw.png)
